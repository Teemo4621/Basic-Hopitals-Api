package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Teemo4621/Hospital-Api/configs"
	"github.com/Teemo4621/Hospital-Api/modules/entities"
	"github.com/Teemo4621/Hospital-Api/modules/mocks"
	"github.com/Teemo4621/Hospital-Api/modules/staffs/controllers"
	"github.com/Teemo4621/Hospital-Api/pkgs/middlewares"
	"github.com/Teemo4621/Hospital-Api/pkgs/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ----------- Test Setup ----------- //

func setupRouter(usecase entities.StaffUseCase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	Cfg := &configs.Config{}
	Cfg.JWT.Secret = "test"
	Cfg.JWT.Expire = 1
	authMiddleware := middlewares.NewAuthMiddleware(Cfg)

	group := r.Group("/staff")
	controllers.NewStaffController(group, *Cfg, usecase, *authMiddleware)
	return r
}

// ----------- Tests ----------- //

func TestFindAllStaffHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUsecase := mocks.NewMockStaffUseCase()
		r := setupRouter(mockUsecase)
		staffs := []entities.Staff{
			{ID: 1, Username: "Test A", Password: "password", FirstNameTH: "Test", LastNameTH: "A", FirstNameEN: "Test", LastNameEN: "A", Gender: "M", HospitalID: 1},
		}
		mockUsecase.On("FindAll", 1, 10).Return(staffs, 1, nil)
		req, _ := http.NewRequest(http.MethodGet, "/staff/", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Pagination Page > Page Total", func(t *testing.T) {
		mockUsecase := mocks.NewMockStaffUseCase()
		r := setupRouter(mockUsecase)
		staffs := []entities.Staff{
			{ID: 1, Username: "", Password: "", FirstNameTH: "Test", LastNameTH: "A", FirstNameEN: "Test", LastNameEN: "A", Gender: "M", HospitalID: 1},
		}
		mockUsecase.On("FindAll", 2, 10).Return(staffs, 1, nil)

		req, _ := http.NewRequest(http.MethodGet, "/staff/?page=2", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var body map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &body)
		assert.NoError(t, err)

		assert.Equal(t, "success", body["message"])
		assert.Equal(t, "ok", body["status"])

		data := body["data"].(map[string]interface{})
		assert.NotNil(t, data["staffs"])
		assert.Equal(t, 1, len(data["staffs"].([]interface{})))

		meta := data["meta"].(map[string]interface{})
		assert.Equal(t, float64(2), meta["page"])
		assert.Equal(t, float64(10), meta["limit"])
		assert.Equal(t, float64(1), meta["page_total"])

		mockUsecase.AssertExpectations(t)
	})

	t.Run("Page is not a number", func(t *testing.T) {
		mockUsecase := mocks.NewMockStaffUseCase()
		r := setupRouter(mockUsecase)

		req, _ := http.NewRequest(http.MethodGet, "/staff/?page=dsawd", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Limit is not a number", func(t *testing.T) {
		mockUsecase := mocks.NewMockStaffUseCase()
		r := setupRouter(mockUsecase)

		req, _ := http.NewRequest(http.MethodGet, "/staff/?limit=dsawd", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Pagination Staff Count = 0", func(t *testing.T) {
		mockUsecase := mocks.NewMockStaffUseCase()
		r := setupRouter(mockUsecase)
		staffs := []entities.Staff{}
		mockUsecase.On("FindAll", 2, 10).Return(staffs, 0, nil)

		req, _ := http.NewRequest(http.MethodGet, "/staff/?page=2", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var body map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &body)
		assert.NoError(t, err)

		assert.Equal(t, "success", body["message"])
		assert.Equal(t, "ok", body["status"])

		data := body["data"].(map[string]interface{})
		assert.NotNil(t, data["staffs"])
		assert.Equal(t, 0, len(data["staffs"].([]interface{})))

		meta := data["meta"].(map[string]interface{})
		assert.Equal(t, float64(2), meta["page"])
		assert.Equal(t, float64(10), meta["limit"])
		assert.Equal(t, float64(0), meta["page_total"])

		mockUsecase.AssertExpectations(t)
	})
}

func TestFindById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUsecase := mocks.NewMockStaffUseCase()
		r := setupRouter(mockUsecase)

		staff := &entities.Staff{ID: 1, Username: "Test A", Password: "password", FirstNameTH: "Test", LastNameTH: "A", FirstNameEN: "Test", LastNameEN: "A", Gender: "M", HospitalID: 1}
		mockUsecase.On("FindById", uint(1)).Return(staff, nil)

		req, _ := http.NewRequest(http.MethodGet, "/staff/1", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockUsecase := mocks.NewMockStaffUseCase()
		mockUsecase.On("FindById", uint(99)).Return((*entities.Staff)(nil), errors.New("staff not found"))

		r := setupRouter(mockUsecase)
		req, _ := http.NewRequest(http.MethodGet, "/staff/99", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Id is not a number", func(t *testing.T) {
		mockUsecase := mocks.NewMockStaffUseCase()

		r := setupRouter(mockUsecase)
		req, _ := http.NewRequest(http.MethodGet, "/staff/dsawd", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		var body map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &body)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, "id is required and must be an integer", body["message"])
		mockUsecase.AssertExpectations(t)
	})
}

func TestMe(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUsecase := mocks.NewMockStaffUseCase()
		r := setupRouter(mockUsecase)

		staff := &entities.Staff{ID: 1, Username: "Test A", Password: "password", FirstNameTH: "Test", LastNameTH: "A", FirstNameEN: "Test", LastNameEN: "A", Gender: "M", HospitalID: 1}
		mockUsecase.On("FindById", uint(1)).Return(staff, nil)

		Cfg := &configs.Config{}
		Cfg.JWT.Secret = "test"
		Cfg.JWT.Expire = 1

		accessToken, err := utils.GenerateAccessToken(Cfg, &entities.Jwtpassport{
			Id:         staff.ID,
			Username:   staff.Username,
			Hospital:   staff.Hospital.HospitalName,
			HospitalID: staff.Hospital.ID,
		})

		if err != nil {
			t.Fatal(err)
		}

		req, _ := http.NewRequest(http.MethodGet, "/staff/me", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: accessToken,
		})
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var body map[string]interface{}
		err = json.Unmarshal(resp.Body.Bytes(), &body)
		assert.NoError(t, err)

		assert.Equal(t, "success", body["message"])
		assert.Equal(t, "ok", body["status"])

		data := body["data"].(map[string]interface{})

		staffID := uint(data["id"].(float64))

		assert.Equal(t, staff.ID, staffID)
		assert.Equal(t, staff.Username, data["username"])
		assert.Equal(t, staff.FirstNameTH, data["first_name_th"])
		assert.Equal(t, staff.MiddleNameTH, data["middle_name_th"])
		assert.Equal(t, staff.LastNameTH, data["last_name_th"])
		assert.Equal(t, staff.FirstNameEN, data["first_name_en"])
		assert.Equal(t, staff.MiddleNameEN, data["middle_name_en"])
		assert.Equal(t, staff.LastNameEN, data["last_name_en"])
		assert.Equal(t, staff.Gender, data["gender"])
		assert.Equal(t, staff.Hospital.HospitalName, data["hospital"].(map[string]interface{})["hospital_name"])
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockUsecase := mocks.NewMockStaffUseCase()

		r := setupRouter(mockUsecase)
		req, _ := http.NewRequest(http.MethodGet, "/staff/me", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		mockUsecase.AssertExpectations(t)
	})
}

func TestCreate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUsecase := mocks.NewMockStaffUseCase()
		r := setupRouter(mockUsecase)

		staff := &entities.Staff{ID: 1, Username: "Test A", FirstNameTH: "Test", LastNameTH: "A", FirstNameEN: "Test", LastNameEN: "A", Gender: "M", HospitalID: 1}

		newStaff := &entities.StaffCreateResponse{
			ID:          staff.ID,
			Username:    staff.Username,
			FirstNameTH: staff.FirstNameTH,
			LastNameTH:  staff.LastNameTH,
			FirstNameEN: staff.FirstNameEN,
			LastNameEN:  staff.LastNameEN,
			Gender:      staff.Gender,
		}

		mockUsecase.On("Create", &entities.StaffCreateRequest{
			Username: "Test A",
			Password: "password",
			Hospital: "Hospital A",
		}).Return(newStaff, nil)

		req, _ := http.NewRequest(http.MethodPost, "/staff/create", bytes.NewBufferString(`{
			"username": "Test A",
			"password": "password",
			"hospital": "Hospital A"
		}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var body map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &body)
		assert.NoError(t, err)

		assert.Equal(t, "ok", body["status"])
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Username is already exists", func(t *testing.T) {
		mockUsecase := mocks.NewMockStaffUseCase()
		r := setupRouter(mockUsecase)
		mockUsecase.On("Create", &entities.StaffCreateRequest{
			Username: "Test A",
			Password: "password",
			Hospital: "Hospital A",
		}).Return((*entities.StaffCreateResponse)(nil), errors.New("username is already exists"))

		req, _ := http.NewRequest(http.MethodPost, "/staff/create", bytes.NewBufferString(`{
			"username": "Test A",
			"password": "password",
			"hospital": "Hospital A"
		}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		var body map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &body)
		assert.NoError(t, err)

		assert.Equal(t, "username is already exists", body["message"])
		assert.Equal(t, "error", body["status"])
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Username is required", func(t *testing.T) {
		mockUsecase := mocks.NewMockStaffUseCase()
		r := setupRouter(mockUsecase)

		req, _ := http.NewRequest(http.MethodPost, "/staff/create", bytes.NewBufferString(`{
			"password": "password",
			"hospital": "Hospital A"
		}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Password is required", func(t *testing.T) {
		mockUsecase := mocks.NewMockStaffUseCase()
		r := setupRouter(mockUsecase)

		req, _ := http.NewRequest(http.MethodPost, "/staff/create", bytes.NewBufferString(`{
			"username": "Test A",
			"hospital": "Hospital A"
		}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Hospital is required", func(t *testing.T) {
		mockUsecase := mocks.NewMockStaffUseCase()
		r := setupRouter(mockUsecase)

		req, _ := http.NewRequest(http.MethodPost, "/staff/create", bytes.NewBufferString(`{
			"username": "Test A",
			"password": "password"
		}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Hospital is not found", func(t *testing.T) {
		mockUsecase := mocks.NewMockStaffUseCase()
		r := setupRouter(mockUsecase)

		mockUsecase.On("Create", &entities.StaffCreateRequest{
			Username: "Test A",
			Password: "password",
			Hospital: "Hospital A",
		}).Return((*entities.StaffCreateResponse)(nil), errors.New("hospital not found"))

		req, _ := http.NewRequest(http.MethodPost, "/staff/create", bytes.NewBufferString(`{
			"username": "Test A",
			"password": "password",
			"hospital": "Hospital A"
		}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		var body map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &body)
		assert.NoError(t, err)

		assert.Equal(t, "hospital not found", body["message"])
		assert.Equal(t, "error", body["status"])
		mockUsecase.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUsecase := mocks.NewMockStaffUseCase()
		Cfg := &configs.Config{}
		Cfg.JWT.Secret = "test"
		Cfg.JWT.Expire = 1
		r := setupRouter(mockUsecase)

		staff := &entities.Staff{ID: 1, Username: "Test A", FirstNameTH: "Test", LastNameTH: "A", FirstNameEN: "Test", LastNameEN: "A", Gender: "M", HospitalID: 1}

		mockUsecase.On("Login", Cfg, &entities.StaffLoginRequest{
			Username: "test",
			Password: "password",
			Hospital: "Hospital A",
		}).Return(&entities.StaffLoginResponse{
			Staff: &entities.StaffMeResponse{
				ID:           staff.ID,
				Username:     staff.Username,
				FirstNameTH:  staff.FirstNameTH,
				MiddleNameTH: staff.MiddleNameTH,
				LastNameTH:   staff.LastNameTH,
				FirstNameEN:  staff.FirstNameEN,
				MiddleNameEN: staff.MiddleNameEN,
				LastNameEN:   staff.LastNameEN,
				Gender:       staff.Gender,
				Hospital:     staff.Hospital,
			},
			AccessToken: "test",
		}, nil)
		req, _ := http.NewRequest(http.MethodPost, "/staff/login", bytes.NewBufferString(`{
			"username": "test",
			"password": "password",
			"hospital": "Hospital A"	
		}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		var body map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &body)
		assert.NoError(t, err)

		assert.Equal(t, "success", body["message"])
		assert.Equal(t, "ok", body["status"])
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Username is required", func(t *testing.T) {
		mockUsecase := mocks.NewMockStaffUseCase()
		r := setupRouter(mockUsecase)

		req, _ := http.NewRequest(http.MethodPost, "/staff/login", bytes.NewBufferString(`{
			"password": "password",
			"hospital": "Hospital A"
		}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Password is required", func(t *testing.T) {
		mockUsecase := mocks.NewMockStaffUseCase()
		r := setupRouter(mockUsecase)

		req, _ := http.NewRequest(http.MethodPost, "/staff/login", bytes.NewBufferString(`{
			"username": "test",
			"hospital": "Hospital A"
		}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Hospital is required", func(t *testing.T) {
		mockUsecase := mocks.NewMockStaffUseCase()
		r := setupRouter(mockUsecase)

		req, _ := http.NewRequest(http.MethodPost, "/staff/login", bytes.NewBufferString(`{
			"username": "test",
			"password": "password"
		}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Hospital not found", func(t *testing.T) {
		mockUsecase := mocks.NewMockStaffUseCase()
		r := setupRouter(mockUsecase)

		Cfg := &configs.Config{}
		Cfg.JWT.Secret = "test"
		Cfg.JWT.Expire = 1

		mockUsecase.On("Login", Cfg, &entities.StaffLoginRequest{
			Username: "test",
			Password: "password",
			Hospital: "Hospital A",
		}).Return((*entities.StaffLoginResponse)(nil), errors.New("hospital not found"))

		req, _ := http.NewRequest(http.MethodPost, "/staff/login", bytes.NewBufferString(`{
			"username": "test",
			"password": "password",
			"hospital": "Hospital A"
		}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		var body map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &body)
		assert.NoError(t, err)

		assert.Equal(t, "username, password or hospital is invalid", body["message"])
		assert.Equal(t, "error", body["status"])
		mockUsecase.AssertExpectations(t)
	})
}
