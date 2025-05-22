package usecases_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Teemo4621/Hospital-Api/configs"
	"github.com/Teemo4621/Hospital-Api/modules/entities"
	"github.com/Teemo4621/Hospital-Api/modules/mocks"
	"github.com/Teemo4621/Hospital-Api/modules/staffs/usecases"
	"github.com/Teemo4621/Hospital-Api/pkgs/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---------- TEST CASES ---------- //

func TestCreateStaff(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := mocks.NewMockStaffRepository()
		mockHospitalRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewStaffUseCase(mockRepo, mockHospitalRepo)
		input := &entities.StaffCreateRequest{Username: "test", Password: "test", Hospital: "test"}
		mockHospitalRepo.On("FindByName", "test").Return(&entities.Hospital{ID: 1, HospitalName: "test"}, nil)
		mockRepo.On("FindByUsername", input.Username).Return((*entities.Staff)(nil), nil)
		mockRepo.On("Create", mock.Anything).Return(&entities.Staff{
			ID:           0,
			Username:     input.Username,
			FirstNameTH:  "test",
			MiddleNameTH: "test",
			LastNameTH:   "test",
			FirstNameEN:  "test",
			MiddleNameEN: "test",
			LastNameEN:   "test",
			Gender:       "M",
		}, nil)

		result, err := usecase.Create(input)
		assert.NoError(t, err)
		assert.Equal(t, "test", result.Username)

		mockRepo.AssertExpectations(t)
		mockHospitalRepo.AssertExpectations(t)
	})

	t.Run("Hospital not found", func(t *testing.T) {
		mockRepo := mocks.NewMockStaffRepository()
		mockHospitalRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewStaffUseCase(mockRepo, mockHospitalRepo)
		input := &entities.StaffCreateRequest{Username: "test", Password: "test", Hospital: "test"}
		mockHospitalRepo.On("FindByName", "test").Return((*entities.Hospital)(nil), nil)
		mockRepo.On("FindByUsername", input.Username).Return((*entities.Staff)(nil), nil)

		_, err := usecase.Create(input)
		assert.EqualError(t, err, "hospital not found")
	})

	t.Run("Staff already exists", func(t *testing.T) {
		mockRepo := mocks.NewMockStaffRepository()
		mockHospitalRepo := mocks.NewMockHospitalRepository()
		hospital := &entities.Hospital{ID: 1, HospitalName: "test"}
		mockHospitalRepo.On("FindByName", "test").Return(hospital, nil)

		usecase := usecases.NewStaffUseCase(mockRepo, mockHospitalRepo)
		input := &entities.StaffCreateRequest{Username: "test", Password: "test", Hospital: "test"}

		mockRepo.On("FindByUsername", "test").Return(&entities.Staff{Username: "test"}, nil)

		_, err := usecase.Create(input)
		assert.EqualError(t, err, "staff name already exists")
	})
}

func TestUpdateStaff(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := mocks.NewMockStaffRepository()
		mockHospitalRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewStaffUseCase(mockRepo, mockHospitalRepo)
		input := &entities.StaffUpdateRequest{ID: uint(1), FirstNameTH: "test11", FirstNameEN: "test11", Gender: "M"}

		OldStaff := &entities.Staff{
			ID:           uint(1),
			Username:     "test",
			FirstNameTH:  "test",
			MiddleNameTH: "test",
			LastNameTH:   "test",
			FirstNameEN:  "test",
			MiddleNameEN: "test",
			LastNameEN:   "test",
			Gender:       "M",
		}

		mockRepo.On("FindById", input.ID).Return(OldStaff, nil)
		OldStaff.FirstNameTH = input.FirstNameTH
		OldStaff.FirstNameEN = input.FirstNameEN
		OldStaff.Gender = input.Gender
		mockRepo.On("Update", OldStaff).Return(OldStaff, nil)

		result, err := usecase.Update(input)
		assert.NoError(t, err)
		assert.Equal(t, OldStaff.ID, result.ID)
		assert.Equal(t, "test", result.Username)
		assert.Equal(t, "test11", result.FirstNameTH)
		assert.Equal(t, "test11", result.FirstNameEN)
		assert.Equal(t, "", result.MiddleNameTH)
		assert.Equal(t, "", result.MiddleNameEN)
		assert.Equal(t, "", result.LastNameTH)
		assert.Equal(t, "", result.LastNameEN)
		assert.Equal(t, "M", result.Gender)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Staff not found", func(t *testing.T) {
		mockRepo := mocks.NewMockStaffRepository()
		mockHospitalRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewStaffUseCase(mockRepo, mockHospitalRepo)
		input := &entities.StaffUpdateRequest{ID: uint(1), FirstNameTH: "test", FirstNameEN: "test", Gender: "M"}
		mockRepo.On("FindById", input.ID).Return((*entities.Staff)(nil), nil)

		_, err := usecase.Update(input)
		assert.EqualError(t, err, "staff not found")
	})
}

func TestFindAllStaffs(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := mocks.NewMockStaffRepository()
		mockHospitalRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewStaffUseCase(mockRepo, mockHospitalRepo)

		staffs := []entities.Staff{{ID: 1}, {ID: 2}}

		mockRepo.On("FindStaffCount").Return(int64(1), nil)
		mockRepo.On("FindAll", 1, 10).Return(staffs, nil)

		result, _, err := usecase.FindAll(1, 10)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("Failed", func(t *testing.T) {
		mockRepo := mocks.NewMockStaffRepository()
		mockHospitalRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewStaffUseCase(mockRepo, mockHospitalRepo)

		mockRepo.On("FindStaffCount").Return(int64(1), errors.New("failed to find staffs"))
		mockRepo.On("FindAll", 1, 10).Return(nil, errors.New("record not found"))

		_, _, err := usecase.FindAll(1, 10)
		assert.EqualError(t, err, "failed to find staffs")
	})

	t.Run("Record not found", func(t *testing.T) {
		mockRepo := mocks.NewMockStaffRepository()
		mockHospitalRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewStaffUseCase(mockRepo, mockHospitalRepo)

		mockRepo.On("FindStaffCount").Return(int64(0), nil)
		mockRepo.On("FindAll", 1, 10).Return([]entities.Staff{}, nil)

		staffs, total, _ := usecase.FindAll(1, 10)
		assert.Equal(t, []entities.Staff{}, staffs)
		assert.Equal(t, 0, total)
	})

	t.Run("Page out of range", func(t *testing.T) {
		mockRepo := mocks.NewMockStaffRepository()
		mockHospitalRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewStaffUseCase(mockRepo, mockHospitalRepo)

		mockRepo.On("FindStaffCount").Return(int64(0), nil)
		mockRepo.On("FindAll", 2, 10).Return([]entities.Staff{}, nil)

		staffs, total, _ := usecase.FindAll(2, 10)
		assert.Equal(t, []entities.Staff{}, staffs)
		assert.Equal(t, 0, total)
		mockRepo.AssertExpectations(t)
		mockHospitalRepo.AssertExpectations(t)
	})
}

func TestFindById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := mocks.NewMockStaffRepository()
		mockHospitalRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewStaffUseCase(mockRepo, mockHospitalRepo)

		staff := &entities.Staff{ID: 1, Username: "test"}
		mockRepo.On("FindById", uint(1)).Return(staff, nil)

		result, err := usecase.FindById(1)
		assert.NoError(t, err)
		assert.Equal(t, "test", result.Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Staff not found", func(t *testing.T) {
		mockRepo := mocks.NewMockStaffRepository()
		mockHospitalRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewStaffUseCase(mockRepo, mockHospitalRepo)

		mockRepo.On("FindById", uint(1)).Return((*entities.Staff)(nil), errors.New("staff not found"))

		_, err := usecase.FindById(1)
		assert.EqualError(t, err, "staff not found")
	})
}

func TestFindByUsername(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := mocks.NewMockStaffRepository()
		mockHospitalRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewStaffUseCase(mockRepo, mockHospitalRepo)

		staff := &entities.Staff{ID: 1, Username: "test"}
		mockRepo.On("FindByUsername", "test").Return(staff, nil)

		result, err := usecase.FindByUsername("test")
		assert.NoError(t, err)
		assert.Equal(t, "test", result.Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Staff not found", func(t *testing.T) {
		mockRepo := mocks.NewMockStaffRepository()
		mockHospitalRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewStaffUseCase(mockRepo, mockHospitalRepo)

		mockRepo.On("FindByUsername", "test").Return((*entities.Staff)(nil), errors.New("staff not found"))

		_, err := usecase.FindByUsername("test")
		assert.EqualError(t, err, "staff not found")
	})
}

func TestLogin(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		cfg := &configs.Config{}
		cfg.JWT.Secret = "test"
		cfg.JWT.Expire = 1
		mockRepo := mocks.NewMockStaffRepository()
		mockHospitalRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewStaffUseCase(mockRepo, mockHospitalRepo)

		hashedPassword, _ := utils.HashPassword("test")
		staff := &entities.Staff{
			ID:       1,
			Username: "test",
			Password: hashedPassword,
			Hospital: entities.Hospital{
				ID:           1,
				HospitalName: "test",
				Address:      "test",
				UpdatedAt:    time.Now(),
				CreatedAt:    time.Now(),
			},
		}
		mockRepo.On("FindByUsername", "test").Return(staff, nil)

		result, err := usecase.Login(cfg, &entities.StaffLoginRequest{
			Username: "test",
			Password: "test",
			Hospital: "test",
		})

		assert.NoError(t, err)
		assert.Equal(t, "test", result.Staff.Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Staff not found", func(t *testing.T) {
		cfg := &configs.Config{}
		cfg.JWT.Secret = "test"
		cfg.JWT.Expire = 1
		mockRepo := mocks.NewMockStaffRepository()
		mockHospitalRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewStaffUseCase(mockRepo, mockHospitalRepo)

		mockRepo.On("FindByUsername", "test").Return((*entities.Staff)(nil), errors.New("staff not found"))

		_, err := usecase.Login(cfg, &entities.StaffLoginRequest{Username: "test", Password: "test"})
		assert.EqualError(t, err, "staff not found")
	})
}
