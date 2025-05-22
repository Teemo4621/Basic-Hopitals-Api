package utils

import (
	"errors"
	"fmt"

	"github.com/Teemo4621/Hospital-Api/configs"
)

func BuildConnectionUrl(stuff string, config configs.Config) (string, error) {
	var url string

	switch stuff {
	case "postgres":
		url = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			config.PostgreSQL.Host,
			config.PostgreSQL.User,
			config.PostgreSQL.Password,
			config.PostgreSQL.DbName,
			config.PostgreSQL.Port,
			config.PostgreSQL.SSLMode,
		)
	case "gin":
		url = config.App.Host + ":" + config.App.Port
	default:
		errMsg := fmt.Sprintf("error, connection url builder doesn't know the %s", stuff)
		return "", errors.New(errMsg)
	}

	return url, nil
}
