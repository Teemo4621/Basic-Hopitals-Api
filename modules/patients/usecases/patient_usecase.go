package usecases

import (
	"errors"

	"github.com/Teemo4621/Hospital-Api/modules/entities"
)

type PatientUseCase struct {
	repo entities.PatientRepository
}

func NewPatientUseCase(repo entities.PatientRepository) entities.PatientUseCase {
	return &PatientUseCase{repo: repo}
}

func (u *PatientUseCase) Create(patient *entities.Patient) (*entities.Patient, error) {
	if exist, err := u.repo.FindByName(patient.FirstNameTH, patient.LastNameTH); err != nil {
		return nil, err
	} else if len(exist) > 0 {
		return nil, errors.New("patient already exists")
	}

	createdPatient, err := u.repo.Create(patient)
	if err != nil {
		return nil, err
	}

	return createdPatient, nil
}

func (u *PatientUseCase) Update(patient *entities.Patient, staffHospitalId uint) (*entities.Patient, error) {
	if patient.HospitalID != staffHospitalId {
		return nil, errors.New("patient not found")
	}

	return u.repo.Update(patient)
}

func (u *PatientUseCase) Delete(id uint, staffHospitalId uint) (*entities.Patient, error) {
	exist, err := u.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	if exist == nil {
		return nil, errors.New("patient not found")
	}

	if exist.HospitalID != staffHospitalId {
		return nil, errors.New("patient not found")
	}

	return u.repo.Delete(id)
}

func (u *PatientUseCase) FindByIdNationalOrPassport(id string, staffHospitalId uint) (*entities.Patient, error) {
	exist, err := u.repo.FindByIdNationalOrPassport(id)
	if err != nil {
		return nil, err
	}
	if exist == nil {
		return nil, errors.New("patient not found")
	}

	if exist.HospitalID != staffHospitalId {
		return nil, errors.New("patient not found")
	}

	return exist, nil
}

func (u *PatientUseCase) FindByAdvanceSearch(input entities.PatientSearchInput, page int, limit int) ([]entities.Patient, int, error) {
	patients, totalPage, err := u.repo.FindByAdvanceSearch(input, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return patients, totalPage, nil
}
