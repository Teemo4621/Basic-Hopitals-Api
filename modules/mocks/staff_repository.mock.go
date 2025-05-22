package mocks

import (
	"github.com/Teemo4621/Hospital-Api/modules/entities"
	"github.com/stretchr/testify/mock"
)

type MockStaffRepository struct {
	mock.Mock
}

func NewMockStaffRepository() *MockStaffRepository {
	return &MockStaffRepository{}
}

func (m *MockStaffRepository) Create(staff *entities.Staff) (*entities.Staff, error) {
	args := m.Called(staff)
	return args.Get(0).(*entities.Staff), args.Error(1)
}

func (m *MockStaffRepository) Update(staff *entities.Staff) (*entities.Staff, error) {
	args := m.Called(staff)
	return args.Get(0).(*entities.Staff), args.Error(1)
}

func (m *MockStaffRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStaffRepository) FindStaffCount() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStaffRepository) FindAll(page int, limit int) ([]entities.Staff, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]entities.Staff), args.Error(1)
}

func (m *MockStaffRepository) FindById(id uint) (*entities.Staff, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Staff), args.Error(1)
}

func (m *MockStaffRepository) FindByUsername(username string) (*entities.Staff, error) {
	args := m.Called(username)
	return args.Get(0).(*entities.Staff), args.Error(1)
}
