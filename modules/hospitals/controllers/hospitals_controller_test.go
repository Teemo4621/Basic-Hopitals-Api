package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Teemo4621/Hospital-Api/configs"
	"github.com/Teemo4621/Hospital-Api/modules/entities"
	"github.com/Teemo4621/Hospital-Api/modules/hospitals/controllers"
	"github.com/Teemo4621/Hospital-Api/modules/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ----------- Test Setup ----------- //

func setupRouter(usecase entities.HospitalUseCase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	group := r.Group("/hospitals")
	controllers.NewHospitalController(group, configs.Config{}, usecase)
	return r
}

// ----------- Tests ----------- //

func TestFindAllHospitalHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUsecase := mocks.NewMockHospitalUseCase()
		r := setupRouter(mockUsecase)
		hospitals := []entities.Hospital{
			{ID: 1, HospitalName: "Test A", Address: "Bangkok"},
		}
		mockUsecase.On("FindAll", 1, 10).Return(hospitals, 1, nil)
		req, _ := http.NewRequest(http.MethodGet, "/hospitals/", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Pagination Page > Page Total", func(t *testing.T) {
		mockUsecase := mocks.NewMockHospitalUseCase()
		r := setupRouter(mockUsecase)
		hospitals := []entities.Hospital{
			{ID: 1, HospitalName: "Test A", Address: "Bangkok"},
		}
		mockUsecase.On("FindAll", 2, 10).Return(hospitals, 1, nil)

		req, _ := http.NewRequest(http.MethodGet, "/hospitals/?page=2", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)

		var body map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &body)
		assert.NoError(t, err)

		assert.Equal(t, "success", body["message"])
		assert.Equal(t, "ok", body["status"])

		data := body["data"].(map[string]interface{})
		assert.NotNil(t, data["hospitals"])
		assert.Equal(t, 1, len(data["hospitals"].([]interface{})))

		meta := data["meta"].(map[string]interface{})
		assert.Equal(t, float64(2), meta["page"])
		assert.Equal(t, float64(10), meta["limit"])
		assert.Equal(t, float64(1), meta["page_total"])

		mockUsecase.AssertExpectations(t)
	})

	t.Run("Page is not a number", func(t *testing.T) {
		mockUsecase := mocks.NewMockHospitalUseCase()
		r := setupRouter(mockUsecase)

		req, _ := http.NewRequest(http.MethodGet, "/hospitals/?page=dsawd", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Limit is not a number", func(t *testing.T) {
		mockUsecase := mocks.NewMockHospitalUseCase()
		r := setupRouter(mockUsecase)

		req, _ := http.NewRequest(http.MethodGet, "/hospitals/?limit=dsawd", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Pagination Hospital Count = 0", func(t *testing.T) {
		mockUsecase := mocks.NewMockHospitalUseCase()
		r := setupRouter(mockUsecase)
		hospitals := []entities.Hospital{}
		mockUsecase.On("FindAll", 2, 10).Return(hospitals, 0, nil)

		req, _ := http.NewRequest(http.MethodGet, "/hospitals/?page=2", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)

		var body map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &body)
		assert.NoError(t, err)

		assert.Equal(t, "success", body["message"])
		assert.Equal(t, "ok", body["status"])

		data := body["data"].(map[string]interface{})
		assert.NotNil(t, data["hospitals"])
		log.Println(data["hospitals"])
		assert.Equal(t, 0, len(data["hospitals"].([]interface{})))

		meta := data["meta"].(map[string]interface{})
		assert.Equal(t, float64(2), meta["page"])
		assert.Equal(t, float64(10), meta["limit"])
		assert.Equal(t, float64(0), meta["page_total"])

		mockUsecase.AssertExpectations(t)
	})

}

func TestFindByIdHospitalHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUsecase := mocks.NewMockHospitalUseCase()
		hospital := &entities.Hospital{ID: 1, HospitalName: "Test A", Address: "Bangkok"}

		mockUsecase.On("FindById", uint(1)).Return(hospital, nil)

		r := setupRouter(mockUsecase)
		req, _ := http.NewRequest(http.MethodGet, "/hospitals/1", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockUsecase := mocks.NewMockHospitalUseCase()
		mockUsecase.On("FindById", uint(99)).Return((*entities.Hospital)(nil), errors.New("hospital not found"))

		r := setupRouter(mockUsecase)
		req, _ := http.NewRequest(http.MethodGet, "/hospitals/99", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Id is not a number", func(t *testing.T) {
		mockUsecase := mocks.NewMockHospitalUseCase()

		r := setupRouter(mockUsecase)
		req, _ := http.NewRequest(http.MethodGet, "/hospitals/dsawd", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		var body map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &body)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, "id is not a number", body["message"])
		mockUsecase.AssertExpectations(t)
	})
}

func TestCreate_Success(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUsecase := mocks.NewMockHospitalUseCase()

		hospital := &entities.Hospital{HospitalName: "Test", Address: "Bangkok"}
		mockUsecase.On("Create", mock.AnythingOfType("*entities.Hospital")).Return(hospital, nil)

		r := setupRouter(mockUsecase)

		body := `{"hospital_name":"Test","address":"Bangkok"}`
		req, _ := http.NewRequest(http.MethodPost, "/hospitals/", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Same Hospital Name", func(t *testing.T) {
		mockUsecase := mocks.NewMockHospitalUseCase()

		hospital := &entities.Hospital{HospitalName: "Test", Address: "Bangkok"}
		mockUsecase.On("Create", mock.AnythingOfType("*entities.Hospital")).Return(hospital, errors.New("hospital name already exists"))

		r := setupRouter(mockUsecase)

		body := `{"hospital_name":"Test","address":"Bangkok"}`
		req, _ := http.NewRequest(http.MethodPost, "/hospitals/", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		mockUsecase := mocks.NewMockHospitalUseCase()

		r := setupRouter(mockUsecase)

		body := `{"hospital_name":"Test"}`
		req, _ := http.NewRequest(http.MethodPost, "/hospitals/", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		mockUsecase.AssertExpectations(t)
	})
}

func TestDelete_Success(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUsecase := mocks.NewMockHospitalUseCase()
		mockUsecase.On("Delete", uint(1)).Return(nil)

		r := setupRouter(mockUsecase)

		req, _ := http.NewRequest(http.MethodDelete, "/hospitals/1", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockUsecase := mocks.NewMockHospitalUseCase()
		mockUsecase.On("Delete", uint(99)).Return(errors.New("hospital not found"))

		r := setupRouter(mockUsecase)

		req, _ := http.NewRequest(http.MethodDelete, "/hospitals/99", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Id is not a number", func(t *testing.T) {
		mockUsecase := mocks.NewMockHospitalUseCase()

		r := setupRouter(mockUsecase)

		req, _ := http.NewRequest(http.MethodDelete, "/hospitals/dsawd", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		mockUsecase.AssertExpectations(t)
	})
}
