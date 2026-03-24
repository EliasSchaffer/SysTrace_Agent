package transport

import (
	"SysTrace_Agent/internal/data/static"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type ENVLoader struct {
}

func resolveEnvPath() string {
	if wd, err := os.Getwd(); err == nil && wd != "" {
		wdPath := filepath.Join(wd, ".env")
		if _, err := os.Stat(wdPath); err == nil {
			return wdPath
		}
	}

	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		if exeDir != "" {
			exeEnvPath := filepath.Join(exeDir, ".env")
			if _, err := os.Stat(exeEnvPath); err == nil {
				return exeEnvPath
			}
		}
	}

	if wd, err := os.Getwd(); err == nil && wd != "" {
		return filepath.Join(wd, ".env")
	}

	return ".env"
}

func loadEnv() string {
	envPath := resolveEnvPath()
	if err := godotenv.Overload(envPath); err != nil {
		panic(fmt.Sprintf("Error loading .env file (%s): %v", envPath, err))
	}
	return envPath
}

func readEnvMap(envPath string) map[string]string {
	envMap, err := godotenv.Read(envPath)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}
		}
		panic(fmt.Sprintf("Error reading .env file (%s): %v", envPath, err))
	}
	return envMap
}

func writeEnvMap(envPath string, envMap map[string]string) {
	if err := godotenv.Write(envMap, envPath); err != nil {
		panic(fmt.Sprintf("Error writing to .env file (%s): %v", envPath, err))
	}
}

func (l *ENVLoader) GetSettings() static.Settings {
	loadEnv()

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
	loadEnv()
	return os.Getenv("GEOLOCATION_API_KEY")
}

func (e *ENVLoader) GetMasterServerURL() string {
	loadEnv()
	return os.Getenv("MASTER_SERVER_URL")
}

func (e *ENVLoader) SetMasterServerURL(url string) {
	envPath := loadEnv()
	envMap := readEnvMap(envPath)
	envMap["MASTER_SERVER_URL"] = url
	writeEnvMap(envPath, envMap)
}

func (e *ENVLoader) SetGeoLocationAPIKey(apiKey string) {
	envPath := loadEnv()
	envMap := readEnvMap(envPath)
	envMap["GEOLOCATION_API_KEY"] = apiKey
	writeEnvMap(envPath, envMap)
}

func (e *ENVLoader) GetSendGPSData() bool {
	loadEnv()
	return os.Getenv("SENDGPS") == "true"
}

func (e *ENVLoader) SetSendGPSData(data bool) {
	envPath := loadEnv()
	envMap := readEnvMap(envPath)
	envMap["SENDGPS"] = strconv.FormatBool(data)
	writeEnvMap(envPath, envMap)
}

func (e *ENVLoader) GetStaticGPSData() static.GPS {
	loadEnv()

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
	envPath := loadEnv()
	envMap := readEnvMap(envPath)
	envMap["GPS_LATITUDE"] = strconv.FormatFloat(data.Latitude, 'f', -1, 64)
	envMap["GPS_LONGITUDE"] = strconv.FormatFloat(data.Longitude, 'f', -1, 64)
	envMap["GPS_ALTITUDE"] = strconv.FormatFloat(data.Altitude, 'f', -1, 64)
	envMap["GPS_ACCURACY"] = strconv.FormatFloat(data.Accuracy, 'f', -1, 64)
	envMap["GPS_CITY"] = data.City
	envMap["GPS_COUNTRY"] = data.Country
	envMap["GPS_REGION"] = data.Region
	writeEnvMap(envPath, envMap)
}

func (e *ENVLoader) SetSettings(settings static.Settings) error {
	envPath := loadEnv()
	envMap := readEnvMap(envPath)
	envMap["GEOLOCATION_API_KEY"] = settings.GEOLOCATION_API_KEY
	envMap["MASTER_SERVER_URL"] = settings.MASTER_SERVER_URL
	envMap["LOGFILE_PATH"] = settings.LOGFILE_PATH
	envMap["SENDGPS"] = strconv.FormatBool(settings.SENDGPS)
	envMap["STATICGPS"] = strconv.FormatBool(settings.STATICGPS)
	envMap["GPS_LATITUDE"] = strconv.FormatFloat(settings.GPS_LATITUDE, 'f', -1, 64)
	envMap["GPS_LONGITUDE"] = strconv.FormatFloat(settings.GPS_LONGITUDE, 'f', -1, 64)
	envMap["GPS_CITY"] = settings.GPS_CITY
	envMap["GPS_REGION"] = settings.GPS_REGION
	envMap["GPS_COUNTRY"] = settings.GPS_COUNTRY

	if err := godotenv.Write(envMap, envPath); err != nil {
		return err
	}
	return nil
}
