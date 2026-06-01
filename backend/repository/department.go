package repository

import (
	"github.com/friedrichad/golang_web_api_demo/backend/configs/db"
	"github.com/friedrichad/golang_web_api_demo/backend/model"
	"gorm.io/gorm"
)

type IDepartmentRepository interface {
	IBaseRepository[model.Department, int]
	GetAllByCondition(filter *model.DepartmentFilter) ([]model.Department, int, error)
	Delete(departmentID []int) error
	GetById(id int) (*model.Department, error)
	Save(department *model.Department) error
	Update(department *model.Department) error
	CreateBatch(departments []*model.Department, batchSize int) error
	UpdateBatch(departments []*model.Department, batchSize int) error
}

type DepartmentRepository struct {
	BaseRepository[model.Department, int]
	DB *gorm.DB
}

var departmentRepository IDepartmentRepository

func NewDepartmentRepository() IDepartmentRepository {
	if departmentRepository == nil {
		departmentRepository = &DepartmentRepository{DB: db.Instance}
		departmentRepository.SetInstance(db.Instance)
	}
	return departmentRepository
}

func (d *DepartmentRepository) Create(department *model.Department) error {
	return d.BaseRepository.Create(department)
}

func (d *DepartmentRepository) Update(department *model.Department) error {
	return d.BaseRepository.Update(department)
}

func (d *DepartmentRepository) Delete(departmentID []int) error {
	return d.BaseRepository.Delete(departmentID)
}

func (d *DepartmentRepository) Save(department *model.Department) error {
	return d.DB.Save(department).Error
}

func (d *DepartmentRepository) GetById(id int) (*model.Department, error) {
	var department *model.Department
	err := d.DB.Where("department_id = ?", id).First(&department).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return department, err
}

func (d *DepartmentRepository) GetAllByCondition(filter *model.DepartmentFilter) ([]model.Department, int, error) {
	return d.GetPage("select d.* from department d "+
		" where (? is null or d.department_name like ?)"+
		" and (? is null or d.created_at >= ?)"+
		" and (? is null or d.created_at <= ?)",
		filter.Page, filter.Size,
		filter.DepartmentName, filter.DepartmentName,
		filter.GetDateFrom(), filter.GetDateFrom(),
		filter.GetDateTo(), filter.GetDateTo(),
	)
}

func (d *DepartmentRepository) CreateBatch(departments []*model.Department, batchSize int) error {
	return d.BaseRepository.CreateBatch(departments, batchSize)
}
func (d *DepartmentRepository) UpdateBatch(departments []*model.Department, batchSize int) error {
	return d.BaseRepository.UpdateBatch(departments, batchSize)
}
