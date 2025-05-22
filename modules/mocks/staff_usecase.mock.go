package mocks

import (
	"github.com/Teemo4621/Hospital-Api/configs"
	"github.com/Teemo4621/Hospital-Api/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockStaffUseCase struct {
	mock.Mock
}

func NewMockStaffUseCase() *MockStaffUseCase {
	return &MockStaffUseCase{}
}

func (m *MockStaffUseCase) Create(staff *entities.StaffCreateRequest) (*entities.StaffCreateResponse, error) {
	args := m.Called(staff)
	return args.Get(0).(*entities.StaffCreateResponse), args.Error(1)
}

func (m *MockStaffUseCase) Update(staff *entities.StaffUpdateRequest) (*entities.Staff, error) {
	args := m.Called(staff)
	return args.Get(0).(*entities.Staff), args.Error(1)
}

func (m *MockStaffUseCase) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStaffUseCase) FindAll(page, limit int) ([]entities.Staff, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]entities.Staff), args.Get(1).(int), args.Error(2)
}

func (m *MockStaffUseCase) FindById(id uint) (*entities.Staff, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Staff), args.Error(1)
}

func (m *MockStaffUseCase) FindByUsername(username string) (*entities.Staff, error) {
	args := m.Called(username)
	return args.Get(0).(*entities.Staff), args.Error(1)
}

func (m *MockStaffUseCase) Login(cfg *configs.Config, staff *entities.StaffLoginRequest) (*entities.StaffLoginResponse, error) {
	args := m.Called(cfg, staff)
	return args.Get(0).(*entities.StaffLoginResponse), args.Error(1)
}
