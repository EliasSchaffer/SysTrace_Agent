package transport

import (
	"SysTrace_Agent/internal/data/static"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ENVLoader struct {
}

func (l *ENVLoader) GetSettings() static.Settings {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	sendGPS, _ := strconv.ParseBool(os.Getenv("SENDGPS"))
	staticGPS, _ := strconv.ParseBool(os.Getenv("STATICGPS"))
	latitude, _ := strconv.ParseFloat(os.Getenv("GPS_LATITUDE"), 64)
	longitude, _ := strconv.ParseFloat(os.Getenv("GPS_LONGITUDE"), 64)

	return static.Settings{
		GEOLOCATION_API_KEY: os.Getenv("GEOLOCATION_API_KEY"),
		MASTER_SERVER_URL:   os.Getenv("MASTER_SERVER_URL"),
		LOGFILE_PATH:        os.Getenv("LOGFILE_PATH"),
		SENDGPS:             sendGPS,
		STATICGPS:           staticGPS,
		GPS_LATITUDE:        latitude,
		GPS_LONGITUDE:       longitude,
		GPS_CITY:            os.Getenv("GPS_CITY"),
		GPS_REGION:          os.Getenv("GPS_REGION"),
		GPS_COUNTRY:         os.Getenv("GPS_COUNTRY"),
	}
}

func (e *ENVLoader) GetGeoLocationAPIKey() string {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	return os.Getenv("GEOLOCATION_API_KEY")
}

func (e *ENVLoader) GetMasterServerURL() string {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	return os.Getenv("MASTER_SERVER_URL")

}

func (e *ENVLoader) SetMasterServerURL(url string) {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	err = godotenv.Write(map[string]string{"MASTER_SERVER_URL": url}, ".env")
	if err != nil {
		panic("Error writing to .env file")
	}
}

func (e *ENVLoader) SetGeoLocationAPIKey(apiKey string) {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	err = godotenv.Write(map[string]string{"GEOLOCATION_API_KEY": apiKey}, ".env")
	if err != nil {
		panic("Error writing to .env file")
	}
}

func (e *ENVLoader) GetSendGPSData() bool {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	return os.Getenv("SENDGPS") == "true"
}

func (e *ENVLoader) SetSendGPSData(data bool) {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	err = godotenv.Write(map[string]string{"SENDGPS": strconv.FormatBool(data)}, ".env")

	if err != nil {
		panic("Error writing to .env file")
	}
}

func (e *ENVLoader) GetStaticGPSData() static.GPS {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	latitude, _ := strconv.ParseFloat(os.Getenv("GPS_LATITUDE"), 64)
	longitude, _ := strconv.ParseFloat(os.Getenv("GPS_LONGITUDE"), 64)
	altitude, _ := strconv.ParseFloat(os.Getenv("GPS_ALTITUDE"), 64)
	accuracy, _ := strconv.ParseFloat(os.Getenv("GPS_ACCURACY"), 64)

	return static.GPS{
		Latitude:  latitude,
		Longitude: longitude,
		Altitude:  altitude,
		Accuracy:  accuracy,
		City:      os.Getenv("GPS_CITY"),
		Country:   os.Getenv("GPS_COUNTRY"),
		Region:    os.Getenv("GPS_REGION"),
	}

}

func (e *ENVLoader) SetStaticGPSData(data static.GPS) {
	envMap, err := godotenv.Read(".env")
	if err != nil {
		if !os.IsNotExist(err) {
			panic("Error loading .env file")
		}
		envMap = map[string]string{}
	}

	envMap["GPS_LATITUDE"] = strconv.FormatFloat(data.Latitude, 'f', -1, 64)
	envMap["GPS_LONGITUDE"] = strconv.FormatFloat(data.Longitude, 'f', -1, 64)
	envMap["GPS_ALTITUDE"] = strconv.FormatFloat(data.Altitude, 'f', -1, 64)
	envMap["GPS_ACCURACY"] = strconv.FormatFloat(data.Accuracy, 'f', -1, 64)
	envMap["GPS_CITY"] = data.City
	envMap["GPS_COUNTRY"] = data.Country
	envMap["GPS_REGION"] = data.Region

	err = godotenv.Write(envMap, ".env")
	if err != nil {
		panic("Error writing to .env file")
	}
}
