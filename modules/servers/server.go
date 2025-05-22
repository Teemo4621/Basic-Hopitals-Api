package servers

import (
	"log"

	"github.com/Teemo4621/Hospital-Api/configs"
	"github.com/Teemo4621/Hospital-Api/pkgs/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	App *gin.Engine
	Cfg *configs.Config
	Db  *gorm.DB
}

func NewServer(cfg *configs.Config, db *gorm.DB) *Server {
	return &Server{
		App: gin.Default(),
		Cfg: cfg,
		Db:  db,
	}
}

func (s *Server) Start() {
	if err := s.MapHandlers(); err != nil {
		log.Fatalln(err.Error())
		panic(err)
	}

	ginConnURL, err := utils.BuildConnectionUrl("gin", *s.Cfg)
	if err != nil {
		log.Fatalln(err.Error())
		panic(err)
	}
	host := s.Cfg.App.Host
	port := s.Cfg.App.Port
	log.Printf("Server is running on %s:%s 🥤", host, port)
	if err := s.App.Run(ginConnURL); err != nil {
		log.Fatalln(err.Error())
		panic(err)
	}
}
