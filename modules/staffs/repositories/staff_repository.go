package repositories

import (
	"github.com/Teemo4621/Hospital-Api/modules/entities"
	"gorm.io/gorm"
)

type StaffRepo struct {
	Db *gorm.DB
}

func NewStaffRepository(db *gorm.DB) entities.StaffRepository {
	return &StaffRepo{Db: db}
}

func (r *StaffRepo) Create(staff *entities.Staff) (*entities.Staff, error) {
	if err := r.Db.Create(&staff).Error; err != nil {
		return nil, err
	}

	return staff, nil
}

func (r *StaffRepo) Update(staff *entities.Staff) (*entities.Staff, error) {
	if err := r.Db.Save(&staff).Error; err != nil {
		return nil, err
	}

	return staff, nil
}

func (r *StaffRepo) Delete(id uint) error {
	return r.Db.Delete(&entities.Staff{}, id).Error
}

func (r *StaffRepo) FindStaffCount() (int64, error) {
	var count int64
	if err := r.Db.Model(&entities.Staff{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *StaffRepo) FindAll(page int, limit int) ([]entities.Staff, error) {
	var staffs []entities.Staff
	if err := r.Db.Preload("Hospital").Offset((page - 1) * limit).Limit(limit).Find(&staffs).Error; err != nil {
		return nil, err
	}
	return staffs, nil
}

func (r *StaffRepo) FindById(id uint) (*entities.Staff, error) {
	var staff entities.Staff
	if err := r.Db.Preload("Hospital").First(&staff, id).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

func (r *StaffRepo) FindByUsername(username string) (*entities.Staff, error) {
	var staff entities.Staff
	if err := r.Db.Preload("Hospital").Where("username = ?", username).First(&staff).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}
