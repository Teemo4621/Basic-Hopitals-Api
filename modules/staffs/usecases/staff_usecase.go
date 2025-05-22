package usecases

import (
	"errors"

	"github.com/Teemo4621/Hospital-Api/configs"
	"github.com/Teemo4621/Hospital-Api/modules/entities"
	"github.com/Teemo4621/Hospital-Api/pkgs/utils"
)

type StaffUseCase struct {
	repo         entities.StaffRepository
	hospitalRepo entities.HospitalRepository
}

func NewStaffUseCase(repo entities.StaffRepository, hospitalRepo entities.HospitalRepository) entities.StaffUseCase {
	return &StaffUseCase{repo: repo, hospitalRepo: hospitalRepo}
}

func (u *StaffUseCase) Create(staff *entities.StaffCreateRequest) (*entities.StaffCreateResponse, error) {
	exist, _ := u.repo.FindByUsername(staff.Username)

	if exist != nil {
		return nil, errors.New("staff name already exists")
	}

	hashedPassword, err := utils.HashPassword(staff.Password)
	if err != nil {
		return nil, err
	}
	staff.Password = hashedPassword

	hospital, _ := u.hospitalRepo.FindByName(staff.Hospital)
	if hospital == nil {
		return nil, errors.New("hospital not found")
	}

	createdStaff := entities.Staff{
		Username:   staff.Username,
		Password:   staff.Password,
		HospitalID: hospital.ID,
	}

	newStaff, err := u.repo.Create(&createdStaff)
	if err != nil {
		return nil, err
	}

	return &entities.StaffCreateResponse{
		ID:           newStaff.ID,
		Username:     newStaff.Username,
		FirstNameTH:  newStaff.FirstNameTH,
		MiddleNameTH: newStaff.MiddleNameTH,
		LastNameTH:   newStaff.LastNameTH,
		FirstNameEN:  newStaff.FirstNameEN,
		MiddleNameEN: newStaff.MiddleNameEN,
		LastNameEN:   newStaff.LastNameEN,
		Gender:       newStaff.Gender,
	}, nil
}

func (u *StaffUseCase) Update(staff *entities.StaffUpdateRequest) (*entities.Staff, error) {
	exist, err := u.repo.FindById(staff.ID)
	if err != nil {
		return nil, err
	}
	if exist == nil {
		return nil, errors.New("staff not found")
	}

	exist.FirstNameTH = staff.FirstNameTH
	exist.MiddleNameTH = staff.MiddleNameTH
	exist.LastNameTH = staff.LastNameTH
	exist.FirstNameEN = staff.FirstNameEN
	exist.MiddleNameEN = staff.MiddleNameEN
	exist.LastNameEN = staff.LastNameEN
	exist.Gender = staff.Gender

	data, err := u.repo.Update(exist)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (u *StaffUseCase) Delete(id uint) error {
	exist, err := u.repo.FindById(id)
	if err != nil {
		return err
	}
	if exist == nil {
		return errors.New("staff not found")
	}

	return u.repo.Delete(id)
}

func (u *StaffUseCase) FindAll(page int, limit int) ([]entities.Staff, int, error) {
	totalCount, err := u.repo.FindStaffCount()
	if err != nil {
		return nil, 0, err
	}

	totalPage := int((totalCount + int64(limit) - 1) / int64(limit))

	staffs, err := u.repo.FindAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	return staffs, totalPage, nil
}

func (u *StaffUseCase) FindById(id uint) (*entities.Staff, error) {
	exist, err := u.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	if exist == nil {
		return nil, errors.New("staff not found")
	}

	return exist, nil
}

func (u *StaffUseCase) FindByUsername(username string) (*entities.Staff, error) {
	exist, err := u.repo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if exist == nil {
		return nil, errors.New("staff not found")
	}

	return exist, nil
}

func (u *StaffUseCase) Login(cfg *configs.Config, loginRequest *entities.StaffLoginRequest) (*entities.StaffLoginResponse, error) {
	exist, err := u.repo.FindByUsername(loginRequest.Username)
	if err != nil {
		return nil, err
	}

	if exist == nil {
		return nil, errors.New("staff not found")
	}

	if !utils.CheckPassword(loginRequest.Password, exist.Password) {
		return nil, errors.New("invalid password")
	}

	if exist.Hospital.HospitalName != loginRequest.Hospital {
		return nil, errors.New("hospital name not match")
	}

	accessToken, err := utils.GenerateAccessToken(cfg, &entities.Jwtpassport{
		Id:         exist.ID,
		Username:   exist.Username,
		Hospital:   exist.Hospital.HospitalName,
		HospitalID: exist.Hospital.ID,
	})

	if err != nil {
		return nil, err
	}

	loginResponse := entities.StaffLoginResponse{
		Staff: &entities.StaffMeResponse{
			ID:           exist.ID,
			Username:     exist.Username,
			FirstNameTH:  exist.FirstNameTH,
			MiddleNameTH: exist.MiddleNameTH,
			LastNameTH:   exist.LastNameTH,
			FirstNameEN:  exist.FirstNameEN,
			MiddleNameEN: exist.MiddleNameEN,
			LastNameEN:   exist.LastNameEN,
			Gender:       exist.Gender,
			Hospital:     exist.Hospital,
		},
		AccessToken: accessToken,
	}

	return &loginResponse, nil
}
