package config

import (
	"os"
	"strings"
)

type IConfig struct {
	Server struct {
		Port		string
		GO_ENV	string
		BaseUrl	string
		Version	string
		Username string
		Password string
	}
	SSL struct {
		ChallengeBaseUrl string
		Email string
	}
}

var Config IConfig

func Load() IConfig {
	Config.Server.Port = os.Getenv("PORT")
	if Config.Server.Port == "" {
		Config.Server.Port = "80"
	}

	Config.Server.GO_ENV = os.Getenv("GO_ENV")
	if Config.Server.GO_ENV == "" {
		Config.Server.GO_ENV = "development"
	}


	Config.Server.Version = "0.1.0"

	Config.Server.BaseUrl = strings.TrimSuffix(os.Getenv("APP_BASE_URL"), "/")

	Config.Server.Username = os.Getenv("APP_BASIC_AUTH_USER")
	Config.Server.Password = os.Getenv("APP_BASIC_AUTH_PASSWORD")

	Config.SSL.ChallengeBaseUrl = strings.TrimSuffix(os.Getenv("APP_SSL_CHALLENGE_BASE_URL"), "/")
	Config.SSL.Email = os.Getenv("APP_SSL_EMAIL")
	if Config.SSL.Email == "" {
		Config.SSL.Email = "admin@example.com"
	}

	return Config
}