package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Teemo4621/Hospital-Api/configs"
	"github.com/Teemo4621/Hospital-Api/modules/entities"
	"github.com/Teemo4621/Hospital-Api/modules/mocks"
	"github.com/Teemo4621/Hospital-Api/modules/patients/controllers"
	"github.com/Teemo4621/Hospital-Api/pkgs/middlewares"
	"github.com/Teemo4621/Hospital-Api/pkgs/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Placeholder for utils.APIResponse; replace with actual struct from utils package
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ----------- Test Setup ----------- //

func setupRouter(mockUseCase *mocks.MockPatientUseCase) (*gin.Engine, *configs.Config, *middlewares.AuthMiddleware) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	cfg := &configs.Config{}
	cfg.JWT.Secret = "test"
	cfg.JWT.Expire = 1
	authMiddleware := middlewares.NewAuthMiddleware(cfg)
	group := r.Group("/patient")
	controllers.NewPatientController(group, *cfg, mockUseCase, *authMiddleware)
	return r, cfg, authMiddleware
}

func createValidToken(cfg *configs.Config, entity *entities.Jwtpassport) string {
	token, _ := utils.GenerateAccessToken(cfg, entity)
	return token
}

func AddAccessTokenCookie(req *http.Request, token string) {
	req.AddCookie(&http.Cookie{
		Name:     "access_token",
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})
}

// ----------- Tests ----------- //

func TestCreatePatientController(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := mocks.NewMockPatientUseCase()
		r, cfg, _ := setupRouter(mockUseCase)

		date := time.Now().Truncate(time.Second)

		expectedPatient := &entities.Patient{
			FirstNameTH: "Test",
			LastNameTH:  "A",
			FirstNameEN: "Test",
			LastNameEN:  "A",
			DateOfBirth: &date,
			PatientHN:   "HN123",
			Gender:      "M",
			HospitalID:  1,
			NationalID:  "1234567890123",
		}
		mockUseCase.On("Create", expectedPatient).Return(expectedPatient, nil)

		reqBody := fmt.Sprintf(`{
            "first_name_th":"Test",
            "last_name_th":"A",
            "first_name_en":"Test",
            "last_name_en":"A",
            "date_of_birth":"%s",
            "patient_hn":"HN123",
            "gender":"M",
            "national_id":"1234567890123"
        }`, date.Format(time.RFC3339))
		req, _ := http.NewRequest(http.MethodPost, "/patient/create", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		AddAccessTokenCookie(req, createValidToken(cfg, &entities.Jwtpassport{HospitalID: 1}))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var response APIResponse
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "success", response.Message)
		assert.Equal(t, "ok", response.Status)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockUseCase := mocks.NewMockPatientUseCase()
		r, _, _ := setupRouter(mockUseCase)

		reqBody := `{
            "first_name_th":"Test",
            "last_name_th":"A",
            "first_name_en":"Test",
            "last_name_en":"A",
            "patient_hn":"HN123",
            "gender":"M",
            "national_id":"1234567890123"
        }`
		req, _ := http.NewRequest(http.MethodPost, "/patient/create", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		var response APIResponse
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "Unauthorized", response.Message)
		mockUseCase.AssertNotCalled(t, "Create")
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		mockUseCase := mocks.NewMockPatientUseCase()
		r, cfg, _ := setupRouter(mockUseCase)

		reqBody := `{"first_name_th":"Test","last_name_th":"A",`
		req, _ := http.NewRequest(http.MethodPost, "/patient/create", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		AddAccessTokenCookie(req, createValidToken(cfg, &entities.Jwtpassport{HospitalID: 1}))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var response APIResponse
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "error", response.Status)
		mockUseCase.AssertNotCalled(t, "Create")
	})

	t.Run("Missing Required Fields", func(t *testing.T) {
		mockUseCase := mocks.NewMockPatientUseCase()
		r, cfg, _ := setupRouter(mockUseCase)

		// Missing first_name_th, date_of_birth
		reqBody := `{"last_name_th":"A","national_id":"1234567890123"}`
		mockUseCase.On("Create", mock.Anything).Return((*entities.Patient)(nil), errors.New("missing required fields"))
		req, _ := http.NewRequest(http.MethodPost, "/patient/create", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		AddAccessTokenCookie(req, createValidToken(cfg, &entities.Jwtpassport{HospitalID: 1}))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Missing NationalID and PassportID", func(t *testing.T) {
		mockUseCase := mocks.NewMockPatientUseCase()
		r, cfg, _ := setupRouter(mockUseCase)

		date := time.Now().Truncate(time.Second)

		reqBody := fmt.Sprintf(`{
            "first_name_th":"Test",
            "last_name_th":"A",
            "first_name_en":"Test",
            "last_name_en":"A",
            "date_of_birth":"%s",
            "patient_hn":"HN123",
            "gender":"M"
        }`, date.Format(time.RFC3339))

		req, _ := http.NewRequest(http.MethodPost, "/patient/create", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		AddAccessTokenCookie(req, createValidToken(cfg, &entities.Jwtpassport{HospitalID: 1}))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var response APIResponse
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "national_id or passport_id is required", response.Message)
		mockUseCase.AssertNotCalled(t, "Create")
	})

	t.Run("NationalID Too Long", func(t *testing.T) {
		mockUseCase := mocks.NewMockPatientUseCase()
		r, cfg, _ := setupRouter(mockUseCase)

		date := time.Now().Truncate(time.Second)

		reqBody := fmt.Sprintf(`{
            "first_name_th":"Test",
            "last_name_th":"A",
            "first_name_en":"Test",
            "last_name_en":"A",
            "date_of_birth":"%s",
            "patient_hn":"HN123",
            "gender":"M",
            "national_id":"12345678901234"
        }`, date.Format(time.RFC3339))

		req, _ := http.NewRequest(http.MethodPost, "/patient/create", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		AddAccessTokenCookie(req, createValidToken(cfg, &entities.Jwtpassport{HospitalID: 1}))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var response APIResponse
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "national_id must be less than 13 characters", response.Message)
		mockUseCase.AssertNotCalled(t, "Create")
	})

	t.Run("Patient Already Exists", func(t *testing.T) {
		mockUseCase := mocks.NewMockPatientUseCase()
		r, cfg, _ := setupRouter(mockUseCase)

		mockUseCase.On("Create", mock.Anything).Return((*entities.Patient)(nil), errors.New("patient already exists"))

		date := time.Now().Truncate(time.Second)

		reqBody := fmt.Sprintf(`{
            "first_name_th":"Test",
            "last_name_th":"A",
            "first_name_en":"Test",
            "last_name_en":"A",
            "date_of_birth":"%s",
            "patient_hn":"HN123",
            "gender":"M",
            "national_id":"1234567890123"
        }`, date.Format(time.RFC3339))
		req, _ := http.NewRequest(http.MethodPost, "/patient/create", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		AddAccessTokenCookie(req, createValidToken(cfg, &entities.Jwtpassport{HospitalID: 1}))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var response APIResponse
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "patient already exists", response.Message)
		mockUseCase.AssertExpectations(t)
	})
}

func TestFindByIdPatientController(t *testing.T) {
	t.Run("Find By National ID", func(t *testing.T) {
		mockUseCase := mocks.NewMockPatientUseCase()
		r, cfg, _ := setupRouter(mockUseCase)

		expectedPatient := &entities.Patient{
			ID:          1,
			FirstNameTH: "Test",
			HospitalID:  1,
			NationalID:  "1234567890123",
		}
		mockUseCase.On("FindByIdNationalOrPassport", "1234567890123", uint(1)).Return(expectedPatient, nil)

		req, _ := http.NewRequest(http.MethodGet, "/patient/search/1234567890123", nil)
		AddAccessTokenCookie(req, createValidToken(cfg, &entities.Jwtpassport{HospitalID: 1}))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var response APIResponse
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "ok", response.Status)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Find By Passport ID", func(t *testing.T) {
		mockUseCase := mocks.NewMockPatientUseCase()
		r, cfg, _ := setupRouter(mockUseCase)

		expectedPatient := &entities.Patient{
			ID:          1,
			FirstNameTH: "Test",
			HospitalID:  1,
			PassportID:  "1234567890123",
		}
		mockUseCase.On("FindByIdNationalOrPassport", "1234567890123", uint(1)).Return(expectedPatient, nil)

		req, _ := http.NewRequest(http.MethodGet, "/patient/search/1234567890123", nil)
		AddAccessTokenCookie(req, createValidToken(cfg, &entities.Jwtpassport{HospitalID: 1}))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var response APIResponse
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "ok", response.Status)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockUseCase := mocks.NewMockPatientUseCase()
		r, _, _ := setupRouter(mockUseCase)

		req, _ := http.NewRequest(http.MethodGet, "/patient/search/1234567890123", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		var response APIResponse
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "Unauthorized", response.Message)
		mockUseCase.AssertNotCalled(t, "FindByIdNationalOrPassport")
	})

	t.Run("Patient Not Found", func(t *testing.T) {
		mockUseCase := mocks.NewMockPatientUseCase()
		r, cfg, _ := setupRouter(mockUseCase)

		mockUseCase.On("FindByIdNationalOrPassport", "1234567890123", uint(1)).Return((*entities.Patient)(nil), errors.New("patient not found"))

		req, _ := http.NewRequest(http.MethodGet, "/patient/search/1234567890123", nil)
		AddAccessTokenCookie(req, createValidToken(cfg, &entities.Jwtpassport{HospitalID: 1}))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		var response APIResponse
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "error", response.Status)
		mockUseCase.AssertExpectations(t)
	})
}

func TestDeletePatientController(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := mocks.NewMockPatientUseCase()
		r, cfg, _ := setupRouter(mockUseCase)

		expectedPatient := &entities.Patient{ID: 1, HospitalID: 1}
		mockUseCase.On("Delete", uint(1), uint(1)).Return(expectedPatient, nil)

		req, _ := http.NewRequest(http.MethodDelete, "/patient/1", nil)
		AddAccessTokenCookie(req, createValidToken(cfg, &entities.Jwtpassport{HospitalID: 1}))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var response APIResponse
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "ok", response.Status)
		assert.Equal(t, "deleted successfully", response.Data)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockUseCase := mocks.NewMockPatientUseCase()
		r, _, _ := setupRouter(mockUseCase)

		req, _ := http.NewRequest(http.MethodDelete, "/patient/1", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		var response APIResponse
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "Unauthorized", response.Message)
		mockUseCase.AssertNotCalled(t, "Delete")
	})

	t.Run("Invalid ID", func(t *testing.T) {
		mockUseCase := mocks.NewMockPatientUseCase()
		r, cfg, _ := setupRouter(mockUseCase)

		req, _ := http.NewRequest(http.MethodDelete, "/patient/dsad1231sad", nil)
		AddAccessTokenCookie(req, createValidToken(cfg, &entities.Jwtpassport{HospitalID: 1}))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var response APIResponse
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "error", response.Status)
		mockUseCase.AssertNotCalled(t, "Delete")
	})

	t.Run("Patient Not Found", func(t *testing.T) {
		mockUseCase := mocks.NewMockPatientUseCase()
		r, cfg, _ := setupRouter(mockUseCase)

		mockUseCase.On("Delete", uint(1), uint(1)).Return((*entities.Patient)(nil), errors.New("patient not found"))

		req, _ := http.NewRequest(http.MethodDelete, "/patient/1", nil)
		AddAccessTokenCookie(req, createValidToken(cfg, &entities.Jwtpassport{HospitalID: 1}))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		var response APIResponse
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "patient not found", response.Message)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Hospital ID Mismatch", func(t *testing.T) {
		mockUseCase := mocks.NewMockPatientUseCase()
		r, cfg, _ := setupRouter(mockUseCase)

		mockUseCase.On("Delete", uint(1), uint(1)).Return((*entities.Patient)(nil), errors.New("patient not found"))

		req, _ := http.NewRequest(http.MethodDelete, "/patient/1", nil)
		AddAccessTokenCookie(req, createValidToken(cfg, &entities.Jwtpassport{HospitalID: 1}))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		var response APIResponse
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "patient not found", response.Message)
		mockUseCase.AssertExpectations(t)
	})
}

func TestFindByAdvanceSearchPatientController(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := mocks.NewMockPatientUseCase()
		r, cfg, _ := setupRouter(mockUseCase)

		input := entities.PatientSearchInput{HospitalID: 1, NationalID: "1234567890123"}
		expectedPatients := []entities.Patient{{ID: 1, FirstNameTH: "Test", HospitalID: 1}}
		mockUseCase.On("FindByAdvanceSearch", input, 1, 10).Return(expectedPatients, 1, nil)

		reqBody := `{"national_id":"1234567890123"}`
		req, _ := http.NewRequest(http.MethodPost, "/patient/search?page=1&limit=10", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		AddAccessTokenCookie(req, createValidToken(cfg, &entities.Jwtpassport{HospitalID: 1}))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var body map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &body)
		assert.NoError(t, err)
		assert.Equal(t, "ok", body["status"])
		assert.Equal(t, "success", body["message"])

		data := body["data"].(map[string]interface{})
		assert.Equal(t, 1, len(data["patients"].([]interface{})))
		meta := data["meta"].(map[string]interface{})
		assert.Equal(t, float64(1), meta["page"])
		assert.Equal(t, float64(10), meta["limit"])
		assert.Equal(t, float64(1), meta["page_total"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockUseCase := mocks.NewMockPatientUseCase()
		r, _, _ := setupRouter(mockUseCase)

		reqBody := `{"national_id":"1234567890123"}`
		req, _ := http.NewRequest(http.MethodPost, "/patient/search?page=1&limit=10", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		var response APIResponse
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "Unauthorized", response.Message)
		mockUseCase.AssertNotCalled(t, "FindByAdvanceSearch")
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		mockUseCase := mocks.NewMockPatientUseCase()
		r, cfg, _ := setupRouter(mockUseCase)

		reqBody := `{"national_id":"1234567890123",`
		req, _ := http.NewRequest(http.MethodPost, "/patient/search?page=1&limit=10", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		AddAccessTokenCookie(req, createValidToken(cfg, &entities.Jwtpassport{HospitalID: 1}))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var response APIResponse
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "error", response.Status)
		mockUseCase.AssertNotCalled(t, "FindByAdvanceSearch")
	})

	t.Run("Invalid Page Param", func(t *testing.T) {
		mockUseCase := mocks.NewMockPatientUseCase()
		r, cfg, _ := setupRouter(mockUseCase)

		reqBody := `{"national_id":"1234567890123"}`
		req, _ := http.NewRequest(http.MethodPost, "/patient/search?page=invalid&limit=10", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		AddAccessTokenCookie(req, createValidToken(cfg, &entities.Jwtpassport{HospitalID: 1}))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var response APIResponse
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "error", response.Status)
		mockUseCase.AssertNotCalled(t, "FindByAdvanceSearch")
	})
}
