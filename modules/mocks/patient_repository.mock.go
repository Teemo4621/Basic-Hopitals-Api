package mocks

import (
	"github.com/Teemo4621/Hospital-Api/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockPatientRepository struct {
	mock.Mock
}

func NewMockPatientRepository() *MockPatientRepository {
	return &MockPatientRepository{}
}

func (m *MockPatientRepository) Create(patient *entities.Patient) (*entities.Patient, error) {
	args := m.Called(patient)
	return args.Get(0).(*entities.Patient), args.Error(1)
}

func (m *MockPatientRepository) Update(patient *entities.Patient) (*entities.Patient, error) {
	args := m.Called(patient)
	return args.Get(0).(*entities.Patient), args.Error(1)
}

func (m *MockPatientRepository) Delete(id uint) (*entities.Patient, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Patient), args.Error(1)
}

func (m *MockPatientRepository) FindAll(page int, limit int) ([]entities.Patient, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]entities.Patient), args.Error(1)
}

func (m *MockPatientRepository) FindById(id uint) (*entities.Patient, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Patient), args.Error(1)
}

func (m *MockPatientRepository) FindPatientCount() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPatientRepository) FindByIdNationalOrPassport(id string) (*entities.Patient, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Patient), args.Error(1)
}

func (m *MockPatientRepository) FindByName(firstName string, lastName string) ([]entities.Patient, error) {
	args := m.Called(firstName, lastName)
	return args.Get(0).([]entities.Patient), args.Error(1)
}

func (m *MockPatientRepository) FindByAdvanceSearch(input entities.PatientSearchInput, page int, limit int) ([]entities.Patient, int, error) {
	args := m.Called(input, page, limit)
	return args.Get(0).([]entities.Patient), args.Get(1).(int), args.Error(2)
}
