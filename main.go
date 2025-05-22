package main

import (
	"os"
	"strconv"

	"github.com/Teemo4621/Hospital-Api/configs"
	"github.com/Teemo4621/Hospital-Api/modules/servers"
	"github.com/Teemo4621/Hospital-Api/pkgs/databases"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("error loading .env file")
	}

	cfg := new(configs.Config)
	cfg.App.Host = os.Getenv("GIN_HOST")
	cfg.App.Port = os.Getenv("GIN_PORT")

	cfg.PostgreSQL.Host = os.Getenv("POSTGRES_HOST")
	cfg.PostgreSQL.User = os.Getenv("POSTGRES_USER")
	cfg.PostgreSQL.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.PostgreSQL.DbName = os.Getenv("POSTGRES_DBNAME")
	cfg.PostgreSQL.Port = os.Getenv("POSTGRES_PORT")
	cfg.PostgreSQL.SSLMode = os.Getenv("POSTGRES_SSLMODE")

	cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	expire, err := strconv.Atoi(os.Getenv("JWT_EXPIRE"))
	if err != nil {
		panic(err)
	}
	cfg.JWT.Expire = expire

	db, err := databases.NewPostgresConnection(*cfg)
	if err != nil {
		panic(err)
	}

	if err := databases.Migrate(db); err != nil {
		panic(err)
	}

	server := servers.NewServer(cfg, db)
	server.Start()
}
