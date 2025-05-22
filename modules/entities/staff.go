package entities

import (
	"time"

	"github.com/Teemo4621/Hospital-Api/configs"
)

type (
	Staff struct {
		ID           uint      `gorm:"primaryKey autoIncrement" json:"id"`
		Username     string    `gorm:"unique;not null" json:"username"`
		Password     string    `gorm:"not null" json:"password"`
		FirstNameTH  string    `gorm:"not null" json:"first_name_th"`
		MiddleNameTH string    `json:"middle_name_th,omitempty"`
		LastNameTH   string    `gorm:"not null" json:"last_name_th"`
		FirstNameEN  string    `gorm:"not null" json:"first_name_en"`
		MiddleNameEN string    `json:"middle_name_en,omitempty"`
		LastNameEN   string    `gorm:"not null" json:"last_name_en"`
		Gender       string    `gorm:"type:char(1);default:'M'" json:"gender"`
		HospitalID   uint      `gorm:"not null" json:"hospital_id"`
		Hospital     Hospital  `gorm:"foreignKey:HospitalID" json:"-"`
		CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
		UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	}

	StaffRepository interface {
		Create(staff *Staff) (*Staff, error)
		Update(staff *Staff) (*Staff, error)
		Delete(id uint) error
		FindStaffCount() (int64, error)
		FindAll(page int, limit int) ([]Staff, error)
		FindById(id uint) (*Staff, error)
		FindByUsername(username string) (*Staff, error)
	}

	StaffUseCase interface {
		Create(staff *StaffCreateRequest) (*StaffCreateResponse, error)
		Update(staff *StaffUpdateRequest) (*Staff, error)
		Delete(id uint) error
		FindAll(page int, limit int) ([]Staff, int, error)
		FindById(id uint) (*Staff, error)
		FindByUsername(username string) (*Staff, error)
		Login(cfg *configs.Config, loginRequest *StaffLoginRequest) (*StaffLoginResponse, error)
	}

	StaffCreateRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Hospital string `json:"hospital" binding:"required"`
	}

	StaffCreateResponse struct {
		ID           uint   `json:"id"`
		Username     string `json:"username"`
		FirstNameTH  string `json:"first_name_th"`
		MiddleNameTH string `json:"middle_name_th,omitempty"`
		LastNameTH   string `json:"last_name_th"`
		FirstNameEN  string `json:"first_name_en"`
		MiddleNameEN string `json:"middle_name_en,omitempty"`
		LastNameEN   string `json:"last_name_en"`
		Gender       string `json:"gender"`
	}

	StaffLoginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Hospital string `json:"hospital" binding:"required"`
	}

	StaffLoginResponse struct {
		Staff       *StaffMeResponse `json:"staff"`
		AccessToken string           `json:"access_token"`
	}

	StaffUpdateRequest struct {
		ID           uint
		FirstNameTH  string `json:"first_name_th"`
		MiddleNameTH string `json:"middle_name_th"`
		LastNameTH   string `json:"last_name_th"`
		FirstNameEN  string `json:"first_name_en"`
		MiddleNameEN string `json:"middle_name_en"`
		LastNameEN   string `json:"last_name_en"`
		Gender       string `json:"gender"`
	}

	StaffResponse struct {
		ID           uint     `json:"id"`
		FirstNameTH  string   `json:"first_name_th"`
		MiddleNameTH string   `json:"middle_name_th"`
		LastNameTH   string   `json:"last_name_th"`
		FirstNameEN  string   `json:"first_name_en"`
		MiddleNameEN string   `json:"middle_name_en"`
		LastNameEN   string   `json:"last_name_en"`
		Gender       string   `json:"gender"`
		Hospital     Hospital `json:"hospital"`
	}

	StaffMeResponse struct {
		ID           uint     `json:"id"`
		Username     string   `json:"username"`
		FirstNameTH  string   `json:"first_name_th"`
		MiddleNameTH string   `json:"middle_name_th"`
		LastNameTH   string   `json:"last_name_th"`
		FirstNameEN  string   `json:"first_name_en"`
		MiddleNameEN string   `json:"middle_name_en"`
		LastNameEN   string   `json:"last_name_en"`
		Gender       string   `json:"gender"`
		Hospital     Hospital `json:"hospital"`
	}
)
