package configs

import (
	"awesomeProject/util/logger"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		logger.Sugar.Fatal("Error loading .env file")
	}
}
