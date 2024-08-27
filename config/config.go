package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB       string
	User     string
	Password string
	DBHost   string
	DBPort   string
}

func Init() (Config, error) {
	err := godotenv.Load("./config/config.env")
	if err != nil {
		return Config{}, err
	}
	return Config{
		DB:       os.Getenv("POSTGRES_DB"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBHost:   os.Getenv("POSTGRES_HOST"),
		DBPort:   os.Getenv("POSTGRES_PORT"),
	}, nil
}
