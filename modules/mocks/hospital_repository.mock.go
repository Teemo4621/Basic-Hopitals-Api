package mocks

import (
	"github.com/Teemo4621/Hospital-Api/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockHospitalRepository struct {
	mock.Mock
}

func NewMockHospitalRepository() *MockHospitalRepository {
	return &MockHospitalRepository{}
}

func (m *MockHospitalRepository) Create(hospital *entities.Hospital) (*entities.Hospital, error) {
	args := m.Called(hospital)
	return args.Get(0).(*entities.Hospital), args.Error(1)
}

func (m *MockHospitalRepository) Update(hospital *entities.Hospital) (*entities.Hospital, error) {
	args := m.Called(hospital)
	return args.Get(0).(*entities.Hospital), args.Error(1)
}

func (m *MockHospitalRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockHospitalRepository) FindHospitalCount() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockHospitalRepository) FindAll(page int, limit int) ([]entities.Hospital, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]entities.Hospital), args.Error(1)
}

func (m *MockHospitalRepository) FindById(id uint) (*entities.Hospital, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Hospital), args.Error(1)
}

func (m *MockHospitalRepository) FindByName(name string) (*entities.Hospital, error) {
	args := m.Called(name)
	return args.Get(0).(*entities.Hospital), args.Error(1)
}
