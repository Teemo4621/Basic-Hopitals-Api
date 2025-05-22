package usecases_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Teemo4621/Hospital-Api/modules/entities"
	"github.com/Teemo4621/Hospital-Api/modules/mocks"
	"github.com/Teemo4621/Hospital-Api/modules/patients/usecases"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		input := &entities.Patient{FirstNameTH: "Test", LastNameTH: "A", FirstNameEN: "Test", LastNameEN: "A", Gender: "M", HospitalID: 1}
		exist := []entities.Patient{}
		mockRepo.On("FindByName", "Test", "A").Return(exist, nil)
		mockRepo.On("Create", input).Return(input, nil)

		result, err := usecase.Create(input)
		assert.NoError(t, err)
		assert.Equal(t, "Test", result.FirstNameTH)
		assert.Equal(t, "A", result.LastNameTH)
		assert.Equal(t, "Test", result.FirstNameEN)
		assert.Equal(t, "A", result.LastNameEN)
		assert.Equal(t, "M", result.Gender)
		assert.Equal(t, uint(1), result.HospitalID)
	})
}

func TestFindAllPatientUseCase(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		patients := []entities.Patient{
			{ID: 1, FirstNameTH: "Test", LastNameTH: "A", FirstNameEN: "Test", LastNameEN: "A", Gender: "M", HospitalID: 1},
		}
		mockRepo.On("FindByAdvanceSearch", entities.PatientSearchInput{}, 1, 10).Return(patients, 1, nil)

		result, total, err := usecase.FindByAdvanceSearch(entities.PatientSearchInput{}, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, 1, total)
		assert.Equal(t, patients, result)
	})

	t.Run("Failed", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)

		expectedErr := errors.New("failed to find patients")
		mockRepo.On("FindByAdvanceSearch", entities.PatientSearchInput{}, 1, 10).Return([]entities.Patient{}, 0, expectedErr)

		result, total, err := usecase.FindByAdvanceSearch(entities.PatientSearchInput{}, 1, 10)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, 0, total)
		assert.Nil(t, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByIdNational", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		patient := &entities.Patient{ID: 1, FirstNameTH: "Test", LastNameTH: "A", FirstNameEN: "Test", LastNameEN: "A", Gender: "M", HospitalID: 1, NationalID: "11231231241231"}
		input := entities.PatientSearchInput{NationalID: "11231231241231"}
		mockRepo.On("FindByAdvanceSearch", input, 1, 10).Return([]entities.Patient{*patient}, 1, nil)

		result, _, err := usecase.FindByAdvanceSearch(input, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, "Test", result[0].FirstNameTH)
		assert.Equal(t, "A", result[0].LastNameTH)
		assert.Equal(t, "Test", result[0].FirstNameEN)
		assert.Equal(t, "A", result[0].LastNameEN)
		assert.Equal(t, "M", result[0].Gender)
		assert.Equal(t, uint(1), result[0].HospitalID)
	})

	t.Run("FindByIdNationalFailed", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		input := entities.PatientSearchInput{NationalID: "11231231241231"}
		mockRepo.On("FindByAdvanceSearch", input, 1, 10).Return([]entities.Patient{}, 0, errors.New("failed to find patients"))

		result, _, err := usecase.FindByAdvanceSearch(input, 1, 10)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})

	t.Run("FindByIdPassport", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		patient := &entities.Patient{ID: 1, FirstNameTH: "Test", LastNameTH: "A", FirstNameEN: "Test", LastNameEN: "A", Gender: "M", HospitalID: 1, PassportID: "11231231241231"}
		input := entities.PatientSearchInput{PassportID: "11231231241231"}
		mockRepo.On("FindByAdvanceSearch", input, 1, 10).Return([]entities.Patient{*patient}, 1, nil)

		result, _, err := usecase.FindByAdvanceSearch(input, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, "Test", result[0].FirstNameTH)
		assert.Equal(t, "A", result[0].LastNameTH)
		assert.Equal(t, "Test", result[0].FirstNameEN)
		assert.Equal(t, "A", result[0].LastNameEN)
		assert.Equal(t, "M", result[0].Gender)
		assert.Equal(t, uint(1), result[0].HospitalID)
	})

	t.Run("FindByIdPassportFailed", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		input := entities.PatientSearchInput{PassportID: "11231231241231"}
		mockRepo.On("FindByAdvanceSearch", input, 1, 10).Return([]entities.Patient{}, 0, errors.New("failed to find patients"))

		result, _, err := usecase.FindByAdvanceSearch(input, 1, 10)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})

	t.Run("FindByFirstName", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		patient := &entities.Patient{ID: 1, FirstNameTH: "Test", LastNameTH: "A", FirstNameEN: "Test", LastNameEN: "A", Gender: "M", HospitalID: 1, NationalID: "11231231241231"}
		input := entities.PatientSearchInput{FirstName: "Test"}
		mockRepo.On("FindByAdvanceSearch", input, 1, 10).Return([]entities.Patient{*patient}, 1, nil)

		result, _, err := usecase.FindByAdvanceSearch(input, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, "Test", result[0].FirstNameTH)
		assert.Equal(t, "A", result[0].LastNameTH)
		assert.Equal(t, "Test", result[0].FirstNameEN)
		assert.Equal(t, "A", result[0].LastNameEN)
		assert.Equal(t, "M", result[0].Gender)
		assert.Equal(t, uint(1), result[0].HospitalID)
	})

	t.Run("FindByFirstNameFailed", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		input := entities.PatientSearchInput{FirstName: "Test"}
		mockRepo.On("FindByAdvanceSearch", input, 1, 10).Return([]entities.Patient{}, 0, errors.New("failed to find patients"))

		result, _, err := usecase.FindByAdvanceSearch(input, 1, 10)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})

	t.Run("FindByMiddleName", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		patient := &entities.Patient{ID: 1, FirstNameTH: "Test", MiddleNameTH: "TestMid", LastNameTH: "A", FirstNameEN: "Test", LastNameEN: "A", Gender: "M", HospitalID: 1, NationalID: "11231231241231"}
		input := entities.PatientSearchInput{MiddleName: "TestMid"}
		mockRepo.On("FindByAdvanceSearch", input, 1, 10).Return([]entities.Patient{*patient}, 1, nil)

		result, _, err := usecase.FindByAdvanceSearch(input, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, "Test", result[0].FirstNameTH)
		assert.Equal(t, "A", result[0].LastNameTH)
		assert.Equal(t, "Test", result[0].FirstNameEN)
		assert.Equal(t, "A", result[0].LastNameEN)
		assert.Equal(t, "M", result[0].Gender)
		assert.Equal(t, uint(1), result[0].HospitalID)
	})

	t.Run("FindByMiddleNameFailed", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		input := entities.PatientSearchInput{MiddleName: "TestMid"}
		mockRepo.On("FindByAdvanceSearch", input, 1, 10).Return([]entities.Patient{}, 0, errors.New("failed to find patients"))

		result, _, err := usecase.FindByAdvanceSearch(input, 1, 10)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})

	t.Run("FindByLastName", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		patient := &entities.Patient{ID: 1, FirstNameTH: "Test", LastNameTH: "A", FirstNameEN: "Test", LastNameEN: "A", Gender: "M", HospitalID: 1, NationalID: "11231231241231"}
		input := entities.PatientSearchInput{LastName: "A"}
		mockRepo.On("FindByAdvanceSearch", input, 1, 10).Return([]entities.Patient{*patient}, 1, nil)

		result, _, err := usecase.FindByAdvanceSearch(input, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, "Test", result[0].FirstNameTH)
		assert.Equal(t, "A", result[0].LastNameTH)
		assert.Equal(t, "Test", result[0].FirstNameEN)
		assert.Equal(t, "A", result[0].LastNameEN)
		assert.Equal(t, "M", result[0].Gender)
		assert.Equal(t, uint(1), result[0].HospitalID)
	})

	t.Run("FindByLastNameFailed", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		input := entities.PatientSearchInput{LastName: "A"}
		mockRepo.On("FindByAdvanceSearch", input, 1, 10).Return([]entities.Patient{}, 0, errors.New("failed to find patients"))

		result, _, err := usecase.FindByAdvanceSearch(input, 1, 10)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})

	t.Run("FindByBirthDate", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		date := time.Now()
		patient := &entities.Patient{ID: 1, FirstNameTH: "Test", LastNameTH: "A", FirstNameEN: "Test", LastNameEN: "A", Gender: "M", HospitalID: 1, NationalID: "11231231241231", DateOfBirth: &date}
		input := entities.PatientSearchInput{DateOfBirth: &date}
		mockRepo.On("FindByAdvanceSearch", input, 1, 10).Return([]entities.Patient{*patient}, 1, nil)

		result, _, err := usecase.FindByAdvanceSearch(input, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, "Test", result[0].FirstNameTH)
		assert.Equal(t, "A", result[0].LastNameTH)
		assert.Equal(t, "Test", result[0].FirstNameEN)
		assert.Equal(t, "A", result[0].LastNameEN)
		assert.Equal(t, "M", result[0].Gender)
		assert.Equal(t, uint(1), result[0].HospitalID)
	})

	t.Run("FindByBirthDateFailed", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		date := time.Now()
		input := entities.PatientSearchInput{DateOfBirth: &date}
		mockRepo.On("FindByAdvanceSearch", input, 1, 10).Return([]entities.Patient{}, 0, errors.New("failed to find patients"))

		result, _, err := usecase.FindByAdvanceSearch(input, 1, 10)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})

	t.Run("FindByPhoneNumber", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		patient := &entities.Patient{ID: 1, FirstNameTH: "Test", LastNameTH: "A", FirstNameEN: "Test", LastNameEN: "A", Gender: "M", HospitalID: 1, NationalID: "11231231241231", PhoneNumber: "0812345678"}
		input := entities.PatientSearchInput{PhoneNumber: "0812345678"}
		mockRepo.On("FindByAdvanceSearch", input, 1, 10).Return([]entities.Patient{*patient}, 1, nil)

		result, _, err := usecase.FindByAdvanceSearch(input, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, "Test", result[0].FirstNameTH)
		assert.Equal(t, "A", result[0].LastNameTH)
		assert.Equal(t, "Test", result[0].FirstNameEN)
		assert.Equal(t, "A", result[0].LastNameEN)
		assert.Equal(t, "M", result[0].Gender)
		assert.Equal(t, uint(1), result[0].HospitalID)
	})

	t.Run("FindByPhoneNumberFailed", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		input := entities.PatientSearchInput{PhoneNumber: "0812345678"}
		mockRepo.On("FindByAdvanceSearch", input, 1, 10).Return([]entities.Patient{}, 0, errors.New("failed to find patients"))

		result, _, err := usecase.FindByAdvanceSearch(input, 1, 10)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})

	t.Run("FindByEmail", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		patient := &entities.Patient{ID: 1, FirstNameTH: "Test", LastNameTH: "A", FirstNameEN: "Test", LastNameEN: "A", Gender: "M", HospitalID: 1, NationalID: "11231231241231", Email: "test@gmail.com"}
		input := entities.PatientSearchInput{Email: "test@gmail.com"}
		mockRepo.On("FindByAdvanceSearch", input, 1, 10).Return([]entities.Patient{*patient}, 1, nil)

		result, _, err := usecase.FindByAdvanceSearch(input, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, "Test", result[0].FirstNameTH)
		assert.Equal(t, "A", result[0].LastNameTH)
		assert.Equal(t, "Test", result[0].FirstNameEN)
		assert.Equal(t, "A", result[0].LastNameEN)
		assert.Equal(t, "M", result[0].Gender)
		assert.Equal(t, uint(1), result[0].HospitalID)
	})

	t.Run("FindByEmailFailed", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		input := entities.PatientSearchInput{Email: "test@gmail.com"}
		mockRepo.On("FindByAdvanceSearch", input, 1, 10).Return([]entities.Patient{}, 0, errors.New("failed to find patients"))

		result, _, err := usecase.FindByAdvanceSearch(input, 1, 10)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})

	t.Run("FindByAllData", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		date := time.Now()
		patient := &entities.Patient{ID: 1, FirstNameTH: "Test", LastNameTH: "A", FirstNameEN: "Test", LastNameEN: "A", Gender: "M", HospitalID: 1, NationalID: "11231231241231", Email: "test@gmail.com", PhoneNumber: "0812345678", DateOfBirth: &date}
		input := entities.PatientSearchInput{}
		input.NationalID = patient.NationalID
		input.PassportID = patient.PassportID
		input.FirstName = patient.FirstNameTH
		input.MiddleName = patient.MiddleNameTH
		input.LastName = patient.LastNameTH
		input.DateOfBirth = patient.DateOfBirth
		input.PhoneNumber = patient.PhoneNumber
		input.Email = patient.Email
		mockRepo.On("FindByAdvanceSearch", input, 1, 10).Return([]entities.Patient{*patient}, 1, nil)

		result, _, err := usecase.FindByAdvanceSearch(input, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, "Test", result[0].FirstNameTH)
		assert.Equal(t, "A", result[0].LastNameTH)
		assert.Equal(t, "Test", result[0].FirstNameEN)
		assert.Equal(t, "A", result[0].LastNameEN)
		assert.Equal(t, "M", result[0].Gender)
		assert.Equal(t, uint(1), result[0].HospitalID)
	})

	t.Run("FindByAllDataFailed", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		input := entities.PatientSearchInput{}
		date := time.Now()
		input.NationalID = "11231231241231"
		input.PassportID = "11231231241231"
		input.FirstName = "Test"
		input.MiddleName = "Test"
		input.LastName = "A"
		input.DateOfBirth = &date
		input.PhoneNumber = "0812345678"
		input.Email = "test@gmail.com"
		mockRepo.On("FindByAdvanceSearch", input, 1, 10).Return([]entities.Patient{}, 0, errors.New("failed to find patients"))

		result, _, err := usecase.FindByAdvanceSearch(input, 1, 10)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})
}

func TestFindByIdNationalOrPassport(t *testing.T) {
	t.Run("FindByNationalId", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		mockHospitalRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		mockHospitalRepo.On("FindById", uint(1)).Return(&entities.Hospital{ID: 1}, nil)
		patient := &entities.Patient{ID: 1, FirstNameTH: "Test", LastNameTH: "A", FirstNameEN: "Test", LastNameEN: "A", Gender: "M", HospitalID: 1, NationalID: "11231231241231"}
		mockRepo.On("FindByIdNationalOrPassport", "11231231241231").Return(patient, nil)

		patient, err := usecase.FindByIdNationalOrPassport("11231231241231", uint(1))

		assert.NoError(t, err)
		assert.Equal(t, "Test", patient.FirstNameTH)
		assert.Equal(t, "A", patient.LastNameTH)
		assert.Equal(t, "Test", patient.FirstNameEN)
		assert.Equal(t, "A", patient.LastNameEN)
		assert.Equal(t, "M", patient.Gender)
		assert.Equal(t, uint(1), patient.HospitalID)
		assert.Equal(t, "11231231241231", patient.NationalID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByNationalIdFailed", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		mockRepo.On("FindByIdNationalOrPassport", "11231231241231").Return((*entities.Patient)(nil), errors.New("failed to find patient"))

		patient, err := usecase.FindByIdNationalOrPassport("11231231241231", uint(1))

		assert.Nil(t, patient)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByPassportId", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		mockHospitalRepo := mocks.NewMockHospitalRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		mockHospitalRepo.On("FindById", uint(1)).Return(&entities.Hospital{ID: 1}, nil)
		patient := &entities.Patient{ID: 1, FirstNameTH: "Test", LastNameTH: "A", FirstNameEN: "Test", LastNameEN: "A", Gender: "M", HospitalID: 1, PassportID: "DB11241231"}
		mockRepo.On("FindByIdNationalOrPassport", "DB11241231").Return(patient, nil)

		patient, err := usecase.FindByIdNationalOrPassport("DB11241231", uint(1))

		assert.NoError(t, err)
		assert.Equal(t, "Test", patient.FirstNameTH)
		assert.Equal(t, "A", patient.LastNameTH)
		assert.Equal(t, "Test", patient.FirstNameEN)
		assert.Equal(t, "A", patient.LastNameEN)
		assert.Equal(t, "M", patient.Gender)
		assert.Equal(t, uint(1), patient.HospitalID)
		assert.Equal(t, "DB11241231", patient.PassportID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByPassportIdFailed", func(t *testing.T) {
		mockRepo := mocks.NewMockPatientRepository()
		usecase := usecases.NewPatientUseCase(mockRepo)
		mockRepo.On("FindByIdNationalOrPassport", "DB11241231").Return((*entities.Patient)(nil), errors.New("failed to find patient"))

		patient, err := usecase.FindByIdNationalOrPassport("DB11241231", uint(1))

		assert.Nil(t, patient)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

}
