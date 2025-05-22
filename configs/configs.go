package configs

type (
	Config struct {
		PostgreSQL PostgreSQLConfig
		App        Gin
		JWT        JWT
	}

	PostgreSQLConfig struct {
		Host     string
		User     string
		Password string
		DbName   string
		Port     string
		SSLMode  string
	}

	Gin struct {
		Host string
		Port string
	}

	JWT struct {
		Secret string
		Expire int
	}
)
