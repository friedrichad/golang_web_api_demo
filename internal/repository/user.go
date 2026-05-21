package repository

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/shared"
	"gorm.io/gorm"
)

type IUserRepository interface {
	IBaseRepository[model.User, int]
	GetByUsername(username string) (*model.User, error)
	GetAuthorities(userId int) ([]string, error)
	GetAllByCondition(query model.UserFilter) ([]model.User, int, error)
	Delete(ids []int) error
	GetById(id int) (*model.User, error)
	Save(user *model.User) error
	Update(user *model.User) error
	AddUserRole(userId int, roleId int) error
	GetUserPermissionScopes(userId int) ([]shared.UserPermissionScope, error)
}

type UserRepository struct {
	BaseRepository[model.User, int]
	DB *gorm.DB
}

var userRepository IUserRepository

func NewUserRepository() IUserRepository {
	if userRepository == nil {
		userRepository = &UserRepository{DB: db.Instance}
		userRepository.SetInstance(db.Instance)
	}
	return userRepository
}

func (u *UserRepository) GetByUsername(username string) (*model.User, error) {

	var user model.User

	err := u.DB.
		Joins("LEFT JOIN position p ON p.position_id = user.position_id").
		Where("user.username = ?", username).
		Select("user.*, p.position_name, p.position_level").
		First(&user).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &user, err
}
func (u *UserRepository) GetAllByCondition(query model.UserFilter) ([]model.User, int, error) {
	return u.GetPage("SELECT * FROM user "+
		"where (? is null OR username like ?) "+
		" and (? is null OR display_name like ?) "+
		" and (? is null OR status_int = ?)"+
		" and (? is null or created_at >= ?) "+
		" and (? is null or created_at < ?) ", query.Page, query.Size, query.Username, query.Username, query.DisplayName, query.DisplayName, query.StatusInt, query.StatusInt, query.GetDateFrom(), query.GetDateFrom(), query.GetDateTo(), query.GetDateTo())
}

func (u *UserRepository) GetAuthorities(userId int) ([]string, error) {
	var authorities []string

	err := u.DB.Raw(`
        SELECT DISTINCT CONCAT(m.menu_name, ':', p.permission_name) AS scope
		FROM user u
		JOIN position_role pr 
			ON u.position_id = pr.position_id
		JOIN role_menu_permission rmp 
			ON pr.role_id = rmp.role_id
		JOIN menu_permission mp 
			ON rmp.menu_permission_id = mp.menu_permission_id
		JOIN menu m 
			ON mp.menu_id = m.menu_id
		JOIN permissions p 
			ON mp.permission_id = p.permission_id
		WHERE u.user_id = ?;
    `, userId).Pluck("scope", &authorities).Error

	return authorities, err
}

func (u *UserRepository) GetById(id int) (*model.User, error) {
	var user model.User

	err := u.DB.
		Joins("LEFT JOIN position p ON p.position_id = user.position_id").
		Where("user.user_id = ?", id).
		Select("user.*, p.position_name, p.position_level, position_id").
		First(&user).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &user, err
}

func (u *UserRepository) Save(user *model.User) error {
	return u.BaseRepository.Create(user)
}

func (u *UserRepository) Update(user *model.User) error {
	return u.BaseRepository.Update(user)
}

func (u *UserRepository) AddUserRole(userId int, roleId int) error {
	userRole := model.UserRole{
		UserID:    userId,
		RoleID:    roleId,
		CreatedAt: time.Now(), // Assuming we want to set created_at
	}
	return u.DB.Create(&userRole).Error
}

func (u *UserRepository) Delete(ids []int) error {
	err := u.DB.Where("user_id IN ?", ids).Delete(&model.User{}).Error
	return err
}

func (u *UserRepository) WithTx(tx *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: BaseRepository[model.User, int]{Instance: tx},
		DB:             tx,
	}
}

func (u *UserRepository) GetUserPermissionScopes(userId int) ([]shared.UserPermissionScope, error) {

	var result []shared.UserPermissionScope

	err := u.DB.Raw(`
        SELECT 
            CONCAT(m.menu_name, ':', p.permission_name) AS scope,
            UNIX_TIMESTAMP(up.expired_date) AS expired_date
        FROM user_permission up
        JOIN menu_permission mp
            ON up.menu_permission_id = mp.menu_permission_id
        JOIN menu m 
            ON mp.menu_id = m.menu_id
        JOIN permissions p 
            ON mp.permission_id = p.permission_id
        WHERE up.user_id = ? 
        AND (up.expired_date IS NULL OR up.expired_date >= NOW())
    `, userId).Scan(&result).Error

	return result, err
}
