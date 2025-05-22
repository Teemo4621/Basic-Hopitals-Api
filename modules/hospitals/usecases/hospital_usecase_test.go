package usecases_test

import (
	"errors"
	"testing"

	"github.com/Teemo4621/Hospital-Api/modules/entities"
	"github.com/Teemo4621/Hospital-Api/modules/hospitals/usecases"
	"github.com/Teemo4621/Hospital-Api/modules/mocks"
	"github.com/stretchr/testify/assert"
)

// ---------- TEST CASES ---------- //

func TestCreateHospital(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewHospitalUseCase(mockRepo)
		input := &entities.Hospital{HospitalName: "Test Hospital", Address: "Bangkok"}

		mockRepo.On("FindByName", "Test Hospital").Return((*entities.Hospital)(nil), nil)
		mockRepo.On("Create", input).Return(input, nil)

		result, err := usecase.Create(input)
		assert.NoError(t, err)
		assert.Equal(t, "Test Hospital", result.HospitalName)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Hospital already exists", func(t *testing.T) {
		mockRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewHospitalUseCase(mockRepo)
		input := &entities.Hospital{HospitalName: "Test Hospital", Address: "Bangkok"}

		mockRepo.On("FindByName", "Test Hospital").Return(input, nil)

		_, err := usecase.Create(input)
		assert.EqualError(t, err, "hospital name already exists")
	})
}

func TestUpdateHospital(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewHospitalUseCase(mockRepo)

		existing := &entities.Hospital{ID: 1, HospitalName: "Old", Address: "Old Addr"}
		updated := &entities.Hospital{ID: 1, HospitalName: "New", Address: "New Addr"}

		mockRepo.On("FindById", uint(1)).Return(existing, nil)
		mockRepo.On("Update", existing).Return(updated, nil)

		result, err := usecase.Update(updated)
		assert.NoError(t, err)
		assert.Equal(t, "New", result.HospitalName)
	})

	t.Run("Hospital not found", func(t *testing.T) {
		mockRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewHospitalUseCase(mockRepo)

		mockRepo.On("FindById", uint(1)).Return((*entities.Hospital)(nil), errors.New("hospital not found"))

		_, err := usecase.Update(&entities.Hospital{ID: 1})
		assert.Equal(t, "hospital not found", err.Error())
	})
}

func TestDeleteHospital(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewHospitalUseCase(mockRepo)

		h := &entities.Hospital{ID: 1}
		mockRepo.On("FindById", uint(1)).Return(h, nil)
		mockRepo.On("Delete", uint(1)).Return(nil)

		err := usecase.Delete(1)
		assert.NoError(t, err)
	})

	t.Run("Hospital not found", func(t *testing.T) {
		mockRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewHospitalUseCase(mockRepo)

		mockRepo.On("FindById", uint(1)).Return((*entities.Hospital)(nil), errors.New("hospital not found"))

		err := usecase.Delete(1)
		assert.EqualError(t, err, "hospital not found")
	})
}

func TestFindAllHospitals(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewHospitalUseCase(mockRepo)

		hospitals := []entities.Hospital{{ID: 1}, {ID: 2}}

		mockRepo.On("FindHospitalCount").Return(int64(1), nil)
		mockRepo.On("FindAll", 1, 10).Return(hospitals, nil)

		result, _, err := usecase.FindAll(1, 10)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("Failed", func(t *testing.T) {
		mockRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewHospitalUseCase(mockRepo)

		mockRepo.On("FindHospitalCount").Return(int64(1), errors.New("failed to find hospitals"))
		mockRepo.On("FindAll", 1, 10).Return(nil, errors.New("failed to find hospitals"))

		_, _, err := usecase.FindAll(1, 10)
		assert.EqualError(t, err, "failed to find hospitals")
	})
}

func TestFindById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewHospitalUseCase(mockRepo)

		hospital := &entities.Hospital{ID: 1, HospitalName: "H1"}
		mockRepo.On("FindById", uint(1)).Return(hospital, nil)

		result, err := usecase.FindById(1)
		assert.NoError(t, err)
		assert.Equal(t, "H1", result.HospitalName)
	})

	t.Run("Hospital not found", func(t *testing.T) {
		mockRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewHospitalUseCase(mockRepo)

		mockRepo.On("FindById", uint(1)).Return((*entities.Hospital)(nil), errors.New("hospital not found"))

		_, err := usecase.FindById(1)
		assert.EqualError(t, err, "hospital not found")
	})
}

func TestFindByName(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewHospitalUseCase(mockRepo)

		hospital := &entities.Hospital{HospitalName: "Test"}
		mockRepo.On("FindByName", "Test").Return(hospital, nil)

		result, err := usecase.FindByName("Test")
		assert.NoError(t, err)
		assert.Equal(t, "Test", result.HospitalName)
	})

	t.Run("Hospital not found", func(t *testing.T) {
		mockRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewHospitalUseCase(mockRepo)

		mockRepo.On("FindByName", "Unknown").Return((*entities.Hospital)(nil), errors.New("hospital not found"))

		_, err := usecase.FindByName("Unknown")
		assert.EqualError(t, err, "hospital not found")
	})
}
