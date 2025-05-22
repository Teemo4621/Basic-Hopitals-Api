package mocks

import (
	"github.com/Teemo4621/Hospital-Api/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockPatientUseCase struct {
	mock.Mock
}

func NewMockPatientUseCase() *MockPatientUseCase {
	return &MockPatientUseCase{}
}

func (m *MockPatientUseCase) Create(input *entities.Patient) (*entities.Patient, error) {
	args := m.Called(input)
	return args.Get(0).(*entities.Patient), args.Error(1)
}

func (m *MockPatientUseCase) Update(patient *entities.Patient, staffHospitalId uint) (*entities.Patient, error) {
	args := m.Called(patient, staffHospitalId)
	return args.Get(0).(*entities.Patient), args.Error(1)
}

func (m *MockPatientUseCase) Delete(id uint, staffHospitalId uint) (*entities.Patient, error) {
	args := m.Called(id, staffHospitalId)
	return args.Get(0).(*entities.Patient), args.Error(1)
}

func (m *MockPatientUseCase) FindByIdNationalOrPassport(id string, staffHospitalId uint) (*entities.Patient, error) {
	args := m.Called(id, staffHospitalId)
	return args.Get(0).(*entities.Patient), args.Error(1)
}

func (m *MockPatientUseCase) FindByAdvanceSearch(input entities.PatientSearchInput, page int, limit int) ([]entities.Patient, int, error) {
	args := m.Called(input, page, limit)
	return args.Get(0).([]entities.Patient), args.Get(1).(int), args.Error(2)
}
