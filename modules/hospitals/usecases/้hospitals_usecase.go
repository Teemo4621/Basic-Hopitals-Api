package usecases

import (
	"errors"

	"github.com/Teemo4621/Hospital-Api/modules/entities"
)

type HospitalUseCase struct {
	repo entities.HospitalRepository
}

func NewHospitalUseCase(repo entities.HospitalRepository) entities.HospitalUseCase {
	return &HospitalUseCase{repo: repo}
}

func (u *HospitalUseCase) Create(hospital *entities.Hospital) (*entities.Hospital, error) {
	exist, _ := u.repo.FindByName(hospital.HospitalName)

	if exist != nil {
		return nil, errors.New("hospital name already exists")
	}

	hospital, err := u.repo.Create(hospital)
	if err != nil {
		return nil, err
	}

	return hospital, nil
}

func (u *HospitalUseCase) Update(hospital *entities.Hospital) (*entities.Hospital, error) {
	exist, err := u.repo.FindById(hospital.ID)
	if err != nil {
		return nil, err
	}
	if exist == nil {
		return nil, errors.New("hospital not found")
	}

	exist.HospitalName = hospital.HospitalName
	exist.Address = hospital.Address

	hospital, err = u.repo.Update(exist)
	if err != nil {
		return nil, err
	}

	return hospital, err
}

func (u *HospitalUseCase) Delete(id uint) error {
	exist, err := u.repo.FindById(id)
	if err != nil {
		return err
	}

	if exist == nil {
		return errors.New("hospital not found")
	}

	return u.repo.Delete(id)
}

func (u *HospitalUseCase) FindAll(page int, limit int) ([]entities.Hospital, int, error) {
	totalCount, err := u.repo.FindHospitalCount()
	if err != nil {
		return nil, 0, err
	}

	totalPage := int((totalCount + int64(limit) - 1) / int64(limit))

	hospitals, err := u.repo.FindAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	return hospitals, totalPage, nil
}

func (u *HospitalUseCase) FindById(id uint) (*entities.Hospital, error) {
	exist, err := u.repo.FindById(id)
	if err != nil {
		return nil, errors.New("hospital not found")
	}

	if exist == nil {
		return nil, errors.New("hospital not found")
	}

	return exist, nil
}

func (u *HospitalUseCase) FindByName(name string) (*entities.Hospital, error) {
	exist, err := u.repo.FindByName(name)
	if err != nil {
		return nil, err
	}

	if exist == nil {
		return nil, errors.New("hospital not found")
	}
	return exist, nil
}
