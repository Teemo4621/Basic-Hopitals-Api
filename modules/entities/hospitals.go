package entities

import "time"

type (
	Hospital struct {
		ID           uint      `gorm:"primaryKey" json:"id"`
		HospitalName string    `gorm:"unique;not null" json:"hospital_name"`
		Address      string    `gorm:"not null" json:"address"`
		CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
		UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

		// Relations
		Staffs   []Staff   `gorm:"foreignKey:HospitalID" json:"-"`
		Patients []Patient `gorm:"foreignKey:HospitalID" json:"-"`
	}

	HospitalRepository interface {
		Create(hospital *Hospital) (*Hospital, error)
		Update(hospital *Hospital) (*Hospital, error)
		Delete(id uint) error
		FindHospitalCount() (int64, error)
		FindAll(page int, limit int) ([]Hospital, error)
		FindById(id uint) (*Hospital, error)
		FindByName(name string) (*Hospital, error)
	}

	HospitalUseCase interface {
		Create(hospital *Hospital) (*Hospital, error)
		Update(hospital *Hospital) (*Hospital, error)
		Delete(id uint) error
		FindAll(page int, limit int) ([]Hospital, int, error)
		FindById(id uint) (*Hospital, error)
		FindByName(name string) (*Hospital, error)
	}

	HospitalCreateRequest struct {
		HospitalName string `json:"hospital_name" binding:"required"`
		Address      string `json:"address" binding:"required"`
	}
)
