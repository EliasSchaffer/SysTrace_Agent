package services

import (
	"os"

	"github.com/joho/godotenv"
)

type ENVLoader struct {
}

func (e *ENVLoader) GetGeoLocationAPIKey() string {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	return os.Getenv("GEOLOCATION_API_KEY")
}
