package controllers

import (
	"strconv"

	"github.com/Teemo4621/Hospital-Api/configs"
	"github.com/Teemo4621/Hospital-Api/modules/entities"
	"github.com/Teemo4621/Hospital-Api/pkgs/middlewares"
	"github.com/Teemo4621/Hospital-Api/pkgs/utils"
	"github.com/gin-gonic/gin"
)

type StaffCon struct {
	Cfg            configs.Config
	StaffUsecase   entities.StaffUseCase
	AuthMiddleware middlewares.AuthMiddleware
}

func NewStaffController(c *gin.RouterGroup, cfg configs.Config, staffUsecase entities.StaffUseCase, authMiddleware middlewares.AuthMiddleware) {
	controller := &StaffCon{
		Cfg:            cfg,
		StaffUsecase:   staffUsecase,
		AuthMiddleware: authMiddleware,
	}
	c.GET("/", controller.FindAll)
	c.GET("/:id", controller.FindById)
	c.POST("/create", controller.Create)
	c.POST("/login", controller.Login)
	c.POST("/update", controller.AuthMiddleware.JwtAuthentication(), controller.Update)
	c.GET("/me", controller.AuthMiddleware.JwtAuthentication(), controller.Me)
}

func (a *StaffCon) FindAll(c *gin.Context) {
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
		utils.BadRequestResponse(c, "limit is required and must be an integer")
		return
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		utils.BadRequestResponse(c, "page is required and must be an integer")
		return
	}

	if pageInt < 1 {
		pageInt = 1
	}

	if limitInt < 1 {
		limitInt = 10
	}

	staffs, totalPage, err := a.StaffUsecase.FindAll(pageInt, limitInt)
	if err != nil {
		utils.ErrorResponse(c, err.Error())
		return
	}

	var staffFindResponse []entities.StaffResponse
	for _, staff := range staffs {
		staffFindResponse = append(staffFindResponse, entities.StaffResponse{
			ID:           staff.ID,
			FirstNameTH:  staff.FirstNameTH,
			MiddleNameTH: staff.MiddleNameTH,
			LastNameTH:   staff.LastNameTH,
			FirstNameEN:  staff.FirstNameEN,
			MiddleNameEN: staff.MiddleNameEN,
			LastNameEN:   staff.LastNameEN,
			Gender:       staff.Gender,
			Hospital:     staff.Hospital,
		})
	}

	if len(staffFindResponse) == 0 {
		staffFindResponse = []entities.StaffResponse{}
	}

	response := gin.H{
		"staffs": staffFindResponse,
		"meta": gin.H{
			"page":       pageInt,
			"limit":      limitInt,
			"page_total": totalPage,
		},
	}

	utils.OkResponse(c, response)
}

func (a *StaffCon) FindById(c *gin.Context) {
	id := c.Param("id")

	staffID, err := strconv.Atoi(id)
	if err != nil {
		utils.BadRequestResponse(c, "id is required and must be an integer")
		return
	}

	staff, err := a.StaffUsecase.FindById(uint(staffID))
	if err != nil {
		utils.NotFoundResponse(c, "staff not found")
		return
	}

	staffFindResponse := entities.StaffResponse{
		ID:           staff.ID,
		FirstNameTH:  staff.FirstNameTH,
		MiddleNameTH: staff.MiddleNameTH,
		LastNameTH:   staff.LastNameTH,
		FirstNameEN:  staff.FirstNameEN,
		MiddleNameEN: staff.MiddleNameEN,
		LastNameEN:   staff.LastNameEN,
		Gender:       staff.Gender,
		Hospital:     staff.Hospital,
	}

	utils.OkResponse(c, staffFindResponse)
}

func (a *StaffCon) Create(c *gin.Context) {
	var staffCreateRequest entities.StaffCreateRequest
	if err := c.ShouldBindJSON(&staffCreateRequest); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	staff, err := a.StaffUsecase.Create(&staffCreateRequest)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.OkResponse(c, staff)
}

func (a *StaffCon) Login(c *gin.Context) {
	var loginRequest entities.StaffLoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	staff, err := a.StaffUsecase.Login(&a.Cfg, &loginRequest)
	if err != nil {
		utils.NotFoundResponse(c, "username, password or hospital is invalid")
		return
	}

	c.SetCookie("access_token", staff.AccessToken, 60*60, "/", "localhost", false, true)

	utils.OkResponse(c, staff)
}

func (a *StaffCon) Update(c *gin.Context) {
	userData, exists := c.Get("user_data")
	if !exists {
		utils.UnauthorizedResponse(c, "Unauthorized")
		return
	}

	var staffReq entities.StaffUpdateRequest
	if err := c.ShouldBindJSON(&staffReq); err != nil {
		utils.BadRequestResponse(c, "bad request")
		return
	}

	staffID := userData.(*entities.JwtClaim).Id
	staffReq.ID = staffID

	updatedStaff, err := a.StaffUsecase.Update(&staffReq)
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	updatedStaffFindResponse := entities.StaffMeResponse{
		ID:           updatedStaff.ID,
		Username:     updatedStaff.Username,
		FirstNameTH:  updatedStaff.FirstNameTH,
		MiddleNameTH: updatedStaff.MiddleNameTH,
		LastNameTH:   updatedStaff.LastNameTH,
		FirstNameEN:  updatedStaff.FirstNameEN,
		MiddleNameEN: updatedStaff.MiddleNameEN,
		LastNameEN:   updatedStaff.LastNameEN,
		Gender:       updatedStaff.Gender,
		Hospital:     updatedStaff.Hospital,
	}

	utils.OkResponse(c, updatedStaffFindResponse)
}

func (a *StaffCon) Me(c *gin.Context) {
	userData, exists := c.Get("user_data")

	if !exists {
		utils.UnauthorizedResponse(c, "Unauthorized")
		return
	}

	staffID := userData.(*entities.JwtClaim).Id

	staff, err := a.StaffUsecase.FindById(staffID)
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	staffFindResponse := entities.StaffMeResponse{
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
	}

	utils.OkResponse(c, staffFindResponse)
}
