package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
	"time"
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
	var user *model.User
	err := u.DB.Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, err
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
        SELECT DISTINCT scope FROM (
            SELECT CONCAT(m.menu_name, ':', p.permission_name) AS scope
            FROM user u
            JOIN position_role pr 
                ON u.position_id = pr.position_id
            JOIN role_menu_permission rmp 
                ON pr.role_id = rmp.role_id
            JOIN menu m 
                ON rmp.menu_id = m.menu_id
            JOIN permissions p 
                ON rmp.permission_id = p.permission_id
            WHERE u.user_id = ?
            UNION
            SELECT CONCAT(m.menu_name, ':', p.permission_name) AS scope
            FROM user_permission up
            JOIN menu m 
                ON up.menu_id = m.menu_id
            JOIN permissions p 
                ON up.permission_id = p.permission_id
            WHERE up.user_id = ?
        ) t
    `, userId, userId).Pluck("scope", &authorities).Error

	return authorities, err
}
func (u *UserRepository) GetById(id int) (*model.User, error) {
	var user *model.User
	err := u.Instance.Where("user_id = ?", id).First(&user).Error
	return user, err
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
