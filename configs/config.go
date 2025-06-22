package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DBConfig  DBConfig
	JWTConfig JWTConfig
}

func NewConfig() *Config {
	return &Config{
		DBConfig: DBConfig{
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
		},
		JWTConfig: JWTConfig{
			Secret: os.Getenv("JWT_SECRET"),
		},
	}
}

func InitConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
