package servers

import (
	_hospitalHttp "github.com/Teemo4621/Hospital-Api/modules/hospitals/controllers"
	_hospitalRepo "github.com/Teemo4621/Hospital-Api/modules/hospitals/repositories"
	_hospitalUseCase "github.com/Teemo4621/Hospital-Api/modules/hospitals/usecases"
	_patientHttp "github.com/Teemo4621/Hospital-Api/modules/patients/controllers"
	_patientRepo "github.com/Teemo4621/Hospital-Api/modules/patients/repositories"
	_patientUseCase "github.com/Teemo4621/Hospital-Api/modules/patients/usecases"
	_staffHttp "github.com/Teemo4621/Hospital-Api/modules/staffs/controllers"
	_staffRepo "github.com/Teemo4621/Hospital-Api/modules/staffs/repositories"
	_staffUseCase "github.com/Teemo4621/Hospital-Api/modules/staffs/usecases"
	"github.com/Teemo4621/Hospital-Api/pkgs/middlewares"
	"github.com/Teemo4621/Hospital-Api/pkgs/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) MapHandlers() error {
	apiGroup := s.App.Group("/api")
	v1 := apiGroup.Group("/v1")
	authMiddleware := middlewares.NewAuthMiddleware(s.Cfg)
	hospitalGroup := v1.Group("/hospitals")

	hospitalRepository := _hospitalRepo.NewHospitalRepository(s.Db)
	hospitalUseCase := _hospitalUseCase.NewHospitalUseCase(hospitalRepository)
	_hospitalHttp.NewHospitalController(hospitalGroup, *s.Cfg, hospitalUseCase)

	staffGroup := v1.Group("/staff")
	staffRepository := _staffRepo.NewStaffRepository(s.Db)
	staffUseCase := _staffUseCase.NewStaffUseCase(staffRepository, hospitalRepository)
	_staffHttp.NewStaffController(staffGroup, *s.Cfg, staffUseCase, *authMiddleware)

	patientGroup := v1.Group("/patient")
	patientRepository := _patientRepo.NewPatientRepository(s.Db)
	patientUseCase := _patientUseCase.NewPatientUseCase(patientRepository)
	_patientHttp.NewPatientController(patientGroup, *s.Cfg, patientUseCase, *authMiddleware)

	s.App.Use(func(c *gin.Context) {
		utils.ErrorResponse(c, "end point not found")
	})

	return nil
}
