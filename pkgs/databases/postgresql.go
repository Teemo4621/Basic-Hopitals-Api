package databases

import (
	"log"

	"github.com/Teemo4621/Hospital-Api/configs"
	"github.com/Teemo4621/Hospital-Api/modules/entities"
	"github.com/Teemo4621/Hospital-Api/pkgs/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnection(config configs.Config) (*gorm.DB, error) {
	url, err := utils.BuildConnectionUrl("postgres", config)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("PostgreSQL database has been connected 📦")

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&entities.Staff{}, &entities.Patient{}, &entities.Hospital{})
}
