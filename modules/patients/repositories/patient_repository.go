package repositories

import (
	"errors"

	"github.com/Teemo4621/Hospital-Api/modules/entities"
	"gorm.io/gorm"
)

type PatientRepo struct {
	Db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) entities.PatientRepository {
	return &PatientRepo{Db: db}
}

func (r *PatientRepo) Create(patient *entities.Patient) (*entities.Patient, error) {
	if err := r.Db.Create(&patient).Error; err != nil {
		return nil, err
	}

	return patient, nil
}

func (r *PatientRepo) Update(patient *entities.Patient) (*entities.Patient, error) {
	if err := r.Db.Save(&patient).Error; err != nil {
		return nil, err
	}

	return patient, nil
}

func (r *PatientRepo) Delete(id uint) (*entities.Patient, error) {
	patient, err := r.FindById(id)
	if err != nil {
		return nil, err
	}

	if patient == nil {
		return nil, errors.New("patient not found")
	}

	if err := r.Db.Delete(&patient).Error; err != nil {
		return nil, err
	}

	return patient, nil
}

func (r *PatientRepo) FindPatientCount() (int64, error) {
	var count int64
	if err := r.Db.Model(&entities.Patient{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *PatientRepo) FindAll(page int, limit int) ([]entities.Patient, error) {
	var patients []entities.Patient
	if err := r.Db.Preload("Hospital").Offset((page - 1) * limit).Limit(limit).Find(&patients).Error; err != nil {
		return nil, err
	}
	return patients, nil
}

func (r *PatientRepo) FindById(id uint) (*entities.Patient, error) {
	var patient entities.Patient
	if err := r.Db.Preload("Hospital").First(&patient, id).Error; err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r *PatientRepo) FindByIdNationalOrPassport(id string) (*entities.Patient, error) {
	var patient entities.Patient
	if err := r.Db.Preload("Hospital").Where("national_id = ? OR passport_id = ?", id, id).First(&patient).Error; err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r *PatientRepo) FindByName(firstName string, lastName string) ([]entities.Patient, error) {
	var patients []entities.Patient
	if err := r.Db.Preload("Hospital").Where("first_name_th = ?", firstName).Where("last_name_th = ?", lastName).Find(&patients).Error; err != nil {
		return nil, err
	}
	return patients, nil
}

func (r *PatientRepo) FindByAdvanceSearch(input entities.PatientSearchInput, page int, limit int) ([]entities.Patient, int, error) {
	var patients []entities.Patient
	var totalCount int64

	query := r.Db.Model(&entities.Patient{}).Where("hospital_id = ?", input.HospitalID)

	if input.NationalID != "" {
		query = query.Where("national_id = ?", input.NationalID)
	}
	if input.PassportID != "" {
		query = query.Where("passport_id = ?", input.PassportID)
	}
	if input.FirstName != "" {
		query = query.Where("first_name_th ILIKE ? OR first_name_en ILIKE ?", "%"+input.FirstName+"%", "%"+input.FirstName+"%")
	}
	if input.MiddleName != "" {
		query = query.Where("middle_name_th ILIKE ? OR middle_name_en ILIKE ?", "%"+input.MiddleName+"%", "%"+input.MiddleName+"%")
	}
	if input.LastName != "" {
		query = query.Where("last_name_th ILIKE ? OR last_name_en ILIKE ?", "%"+input.LastName+"%", "%"+input.LastName+"%")
	}
	if input.DateOfBirth != nil {
		query = query.Where("date_of_birth = ?", input.DateOfBirth)
	}
	if input.PhoneNumber != "" {
		query = query.Where("phone_number ILIKE ?", "%"+input.PhoneNumber+"%")
	}
	if input.Email != "" {
		query = query.Where("email ILIKE ?", "%"+input.Email+"%")
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * limit).Limit(limit).Find(&patients).Error; err != nil {
		return nil, 0, err
	}

	return patients, int(totalCount), nil
}
