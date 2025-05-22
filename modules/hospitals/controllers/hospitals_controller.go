package controllers

import (
	"strconv"

	"github.com/Teemo4621/Hospital-Api/configs"
	"github.com/Teemo4621/Hospital-Api/modules/entities"
	"github.com/Teemo4621/Hospital-Api/pkgs/utils"
	"github.com/gin-gonic/gin"
)

type HospitalCon struct {
	Cfg             configs.Config
	HospitalUsecase entities.HospitalUseCase
}

func NewHospitalController(c *gin.RouterGroup, cfg configs.Config, hospitalUsecase entities.HospitalUseCase) {
	controller := &HospitalCon{
		Cfg:             cfg,
		HospitalUsecase: hospitalUsecase,
	}
	c.GET("/", controller.FindAll)
	c.GET("/:id", controller.FindById)
	c.POST("/", controller.Create)
	c.DELETE("/:id", controller.Delete)
}

func (a *HospitalCon) FindAll(c *gin.Context) {
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

	hospital, totalPage, err := a.HospitalUsecase.FindAll(pageInt, limitInt)
	if err != nil {
		utils.ErrorResponse(c, err.Error())
		return
	}

	if len(hospital) == 0 {
		hospital = []entities.Hospital{}
	}

	utils.OkResponse(c, gin.H{
		"hospitals": hospital,
		"meta": gin.H{
			"page":       pageInt,
			"limit":      limitInt,
			"page_total": totalPage,
		},
	})
}

func (a *HospitalCon) FindById(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.BadRequestResponse(c, "id is not a number")
		return
	}

	hospital, err := a.HospitalUsecase.FindById(uint(idInt))
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	utils.OkResponse(c, hospital)
}

func (a *HospitalCon) Create(c *gin.Context) {
	var reqHospital entities.HospitalCreateRequest
	if err := c.ShouldBindJSON(&reqHospital); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	hospital := &entities.Hospital{
		HospitalName: reqHospital.HospitalName,
		Address:      reqHospital.Address,
	}

	hospital, err := a.HospitalUsecase.Create(hospital)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.OkResponse(c, hospital)
}

func (a *HospitalCon) Delete(c *gin.Context) {
	id := c.Param("id")

	hospitalID, err := strconv.Atoi(id)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	err = a.HospitalUsecase.Delete(uint(hospitalID))
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	utils.OkResponse(c, nil)
}
