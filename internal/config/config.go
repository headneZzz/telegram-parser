package config

import (
	"github.com/joho/godotenv"
)

func Config() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
}
