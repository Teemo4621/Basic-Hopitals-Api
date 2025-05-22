package mocks

import (
	"github.com/Teemo4621/Hospital-Api/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockHospitalUseCase struct {
	mock.Mock
}

func NewMockHospitalUseCase() *MockHospitalUseCase {
	return &MockHospitalUseCase{}
}

func (m *MockHospitalUseCase) Create(hospital *entities.Hospital) (*entities.Hospital, error) {
	args := m.Called(hospital)
	return args.Get(0).(*entities.Hospital), args.Error(1)
}

func (m *MockHospitalUseCase) Update(hospital *entities.Hospital) (*entities.Hospital, error) {
	args := m.Called(hospital)
	return args.Get(0).(*entities.Hospital), args.Error(1)
}

func (m *MockHospitalUseCase) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockHospitalUseCase) FindAll(page int, limit int) ([]entities.Hospital, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]entities.Hospital), args.Int(1), args.Error(2)
}

func (m *MockHospitalUseCase) FindById(id uint) (*entities.Hospital, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Hospital), args.Error(1)
}

func (m *MockHospitalUseCase) FindByName(name string) (*entities.Hospital, error) {
	args := m.Called(name)
	return args.Get(0).(*entities.Hospital), args.Error(1)
}
