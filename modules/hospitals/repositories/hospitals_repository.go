package repositories

import (
	"github.com/Teemo4621/Hospital-Api/modules/entities"
	"gorm.io/gorm"
)

type HospitalRepo struct {
	Db *gorm.DB
}

func NewHospitalRepository(db *gorm.DB) entities.HospitalRepository {
	return &HospitalRepo{Db: db}
}

func (r *HospitalRepo) Create(hospital *entities.Hospital) (*entities.Hospital, error) {
	if err := r.Db.Create(&hospital).Error; err != nil {
		return nil, err
	}

	return hospital, nil
}

func (r *HospitalRepo) Update(hospital *entities.Hospital) (*entities.Hospital, error) {
	if err := r.Db.Save(&hospital).Error; err != nil {
		return nil, err
	}

	return hospital, nil
}

func (r *HospitalRepo) Delete(id uint) error {
	return r.Db.Delete(&entities.Hospital{}, id).Error
}

func (r *HospitalRepo) FindHospitalCount() (int64, error) {
	var count int64
	if err := r.Db.Model(&entities.Hospital{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *HospitalRepo) FindAll(page int, limit int) ([]entities.Hospital, error) {
	var hospitals []entities.Hospital
	if err := r.Db.Offset((page - 1) * limit).Limit(limit).Find(&hospitals).Error; err != nil {
		return nil, err
	}
	return hospitals, nil
}

func (r *HospitalRepo) FindById(id uint) (*entities.Hospital, error) {
	var hospital entities.Hospital
	if err := r.Db.First(&hospital, id).Error; err != nil {
		return nil, err
	}
	return &hospital, nil
}

func (r *HospitalRepo) FindByName(name string) (*entities.Hospital, error) {
	var hospital entities.Hospital
	if err := r.Db.Where("hospital_name = ?", name).First(&hospital).Error; err != nil {
		return nil, err
	}
	return &hospital, nil
}
