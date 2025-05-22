package entities

import "time"

type (
	Patient struct {
		ID           uint       `gorm:"primaryKey autoIncrement" json:"id"`
		FirstNameTH  string     `gorm:"not null" json:"first_name_th"`
		MiddleNameTH string     `json:"middle_name_th"`
		LastNameTH   string     `gorm:"not null" json:"last_name_th"`
		FirstNameEN  string     `gorm:"not null" json:"first_name_en"`
		MiddleNameEN string     `json:"middle_name_en"`
		LastNameEN   string     `gorm:"not null" json:"last_name_en"`
		DateOfBirth  *time.Time `gorm:"not null" json:"date_of_birth"`
		PatientHN    string     `json:"patient_hn"`
		NationalID   string     `json:"national_id"`
		PassportID   string     `json:"passport_id"`
		PhoneNumber  string     `json:"phone_number"`
		Email        string     `json:"email"`
		Gender       string     `gorm:"type:char(1); default 'M'" json:"gender"`
		HospitalID   uint       `gorm:"not null" json:"hospital_id"`
		Hospital     Hospital   `gorm:"foreignKey:HospitalID" json:"-"`
		CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
		UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	}

	PatientRepository interface {
		Create(patient *Patient) (*Patient, error)
		Update(patient *Patient) (*Patient, error)
		Delete(id uint) (*Patient, error)
		FindAll(page int, limit int) ([]Patient, error)
		FindById(id uint) (*Patient, error)
		FindByIdNationalOrPassport(id string) (*Patient, error)
		FindByName(firstName string, lastName string) ([]Patient, error)
		FindByAdvanceSearch(input PatientSearchInput, page int, limit int) ([]Patient, int, error)
	}

	PatientUseCase interface {
		Create(patient *Patient) (*Patient, error)
		Update(patient *Patient, staffHospitalId uint) (*Patient, error)
		Delete(id uint, staffHospitalId uint) (*Patient, error)
		FindByIdNationalOrPassport(id string, staffHospitalId uint) (*Patient, error)
		FindByAdvanceSearch(input PatientSearchInput, page int, limit int) ([]Patient, int, error)
	}

	PatientCreateRequest struct {
		FirstNameTH  string     `json:"first_name_th" binding:"required"`
		MiddleNameTH string     `json:"middle_name_th,omitempty"`
		LastNameTH   string     `json:"last_name_th" binding:"required"`
		FirstNameEN  string     `json:"first_name_en" binding:"required"`
		MiddleNameEN string     `json:"middle_name_en,omitempty"`
		LastNameEN   string     `json:"last_name_en" binding:"required"`
		DateOfBirth  *time.Time `json:"date_of_birth" binding:"required"`
		PatientHN    string     `json:"patient_hn" binding:"required"`
		NationalID   string     `json:"national_id"`
		PassportID   string     `json:"passport_id"`
		PhoneNumber  string     `json:"phone_number,omitempty"`
		Email        string     `json:"email,omitempty"`
		Gender       string     `json:"gender" binding:"required"`
		HospitalID   uint
	}

	PatientSearchInput struct {
		HospitalID  uint
		NationalID  string     `json:"national_id"`
		PassportID  string     `json:"passport_id"`
		FirstName   string     `json:"first_name"`
		MiddleName  string     `json:"middle_name"`
		LastName    string     `json:"last_name"`
		DateOfBirth *time.Time `json:"date_of_birth"`
		PhoneNumber string     `json:"phone_number"`
		Email       string     `json:"email"`
	}
)
