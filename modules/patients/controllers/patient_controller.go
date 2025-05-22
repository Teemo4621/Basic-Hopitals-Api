package controllers

import (
	"strconv"

	"github.com/Teemo4621/Hospital-Api/configs"
	"github.com/Teemo4621/Hospital-Api/modules/entities"
	"github.com/Teemo4621/Hospital-Api/pkgs/middlewares"
	"github.com/Teemo4621/Hospital-Api/pkgs/utils"
	"github.com/gin-gonic/gin"
)

type PatientCon struct {
	Cfg            configs.Config
	PatientUsecase entities.PatientUseCase
	AuthMiddleware middlewares.AuthMiddleware
}

func NewPatientController(c *gin.RouterGroup, cfg configs.Config, patientUsecase entities.PatientUseCase, authMiddleware middlewares.AuthMiddleware) {
	controller := &PatientCon{
		Cfg:            cfg,
		PatientUsecase: patientUsecase,
		AuthMiddleware: authMiddleware,
	}

	c.GET("/search/:id", controller.AuthMiddleware.JwtAuthentication(), controller.FindById)
	c.POST("/create", controller.AuthMiddleware.JwtAuthentication(), controller.Create)
	c.POST("/update", controller.AuthMiddleware.JwtAuthentication(), controller.Update)
	c.DELETE("/:id", controller.AuthMiddleware.JwtAuthentication(), controller.Delete)
	c.POST("/search", controller.AuthMiddleware.JwtAuthentication(), controller.FindByAdvanceSearch)
}

func (a *PatientCon) Create(c *gin.Context) {
	userData, exists := c.Get("user_data")
	if !exists {
		utils.UnauthorizedResponse(c, "Unauthorized")
		return
	}

	HospitalID := userData.(*entities.JwtClaim).HospitalID

	var patient entities.Patient

	if err := c.ShouldBindJSON(&patient); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	if patient.NationalID == "" && patient.PassportID == "" {
		utils.BadRequestResponse(c, "national_id or passport_id is required")
		return
	}

	if len(patient.NationalID) > 13 {
		utils.BadRequestResponse(c, "national_id must be less than 13 characters")
		return
	}

	patient.HospitalID = HospitalID

	createdPatient, err := a.PatientUsecase.Create(&patient)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.OkResponse(c, createdPatient)
}

func (a *PatientCon) Update(c *gin.Context) {
	userData, exists := c.Get("user_data")
	if !exists {
		utils.UnauthorizedResponse(c, "Unauthorized")
		return
	}

	HospitalID := userData.(*entities.JwtClaim).HospitalID

	var patient entities.Patient

	if err := c.ShouldBindJSON(&patient); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	updatedPatient, err := a.PatientUsecase.Update(&patient, HospitalID)
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	utils.OkResponse(c, updatedPatient)
}

func (a *PatientCon) FindById(c *gin.Context) {
	userData, exists := c.Get("user_data")
	if !exists {
		utils.UnauthorizedResponse(c, "Unauthorized")
		return
	}

	HospitalID := userData.(*entities.JwtClaim).HospitalID

	id := c.Param("id")

	patient, err := a.PatientUsecase.FindByIdNationalOrPassport(id, HospitalID)
	if err != nil {
		utils.NotFoundResponse(c, "patient not found")
		return
	}

	utils.OkResponse(c, patient)
}

func (a *PatientCon) Delete(c *gin.Context) {
	userData, exists := c.Get("user_data")
	if !exists {
		utils.UnauthorizedResponse(c, "Unauthorized")
		return
	}

	HospitalID := userData.(*entities.JwtClaim).HospitalID

	id := c.Param("id")

	patientID, err := strconv.Atoi(id)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	_, err = a.PatientUsecase.Delete(uint(patientID), HospitalID)
	if err != nil {
		utils.NotFoundResponse(c, "patient not found")
		return
	}

	utils.OkResponse(c, "deleted successfully")
}

func (a *PatientCon) FindByAdvanceSearch(c *gin.Context) {
	userData, exists := c.Get("user_data")
	if !exists {
		utils.UnauthorizedResponse(c, "Unauthorized")
		return
	}

	HospitalID := userData.(*entities.JwtClaim).HospitalID

	var input entities.PatientSearchInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	input.HospitalID = HospitalID

	page := c.Query("page")
	limit := c.Query("limit")

	if page == "" {
		page = "1"
	}

	if limit == "" {
		limit = "10"
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	patient, totalPage, err := a.PatientUsecase.FindByAdvanceSearch(input, pageInt, limitInt)
	if err != nil {
		utils.ErrorResponse(c, err.Error())
		return
	}

	utils.OkResponse(c, gin.H{
		"patients": patient,
		"meta": gin.H{
			"page":       pageInt,
			"limit":      limitInt,
			"page_total": totalPage,
		},
	})
}
