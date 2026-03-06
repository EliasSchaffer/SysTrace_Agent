package services

import (
	"SysTrace_Agent/data"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func (a *Agent) CollectGPSData() {
	gpsData := a.GetGPSDataByLocationAPI()
	if gpsData != nil {
		a.device.SetGPS(*gpsData)
	} else {
		fmt.Println("Failed to get GPS data from Location API, trying IP-based geolocation...")
		gpsData = a.GetGPSDataByIP()
		if gpsData != nil {
			fmt.Println("")
			a.device.SetGPS(*gpsData)
		} else {
			fmt.Println("Failed to get GPS data from both methods.")
		}
	}
}

// GetGPSDataByIP retrieves GPS data based on the agent's IP address.
//
// It constructs a URL using the GeoLocation API key, makes an HTTP GET request to fetch the GPS data,
// and decodes the JSON response into a map. The function extracts city, region, country, latitude, and
// longitude from the response and sets them in a data.GPS struct. If any errors occur during the HTTP
// request or JSON decoding, it logs the error and returns nil. The function returns a pointer to the
// populated data.GPS struct or nil if an error occurred.
func (a *Agent) GetGPSDataByIP() *data.GPS {
	apiKey := a.envLoader.GetGeoLocationAPIKey()
	url := fmt.Sprintf("https://api.ipgeolocation.io/ipgeo?apiKey=%s", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching GPS data: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Printf("Error decoding GPS data: %v\n", err)
		return nil
	}

	gps := &data.GPS{}

	if city, ok := result["city"].(string); ok {
		gps.SetCity(city)
	}
	if region, ok := result["state_prov"].(string); ok {
		gps.SetRegion(region)
	}
	if country, ok := result["country_name"].(string); ok {
		gps.SetCountry(country)
	}

	if latStr, ok := result["latitude"].(string); ok {
		if lat, err := strconv.ParseFloat(latStr, 64); err == nil {
			gps.SetLatitude(lat)
		}
	} else if latFloat, ok := result["latitude"].(float64); ok {
		gps.SetLatitude(latFloat)
	}

	if lonStr, ok := result["longitude"].(string); ok {
		if lon, err := strconv.ParseFloat(lonStr, 64); err == nil {
			gps.SetLongitude(lon)
		}
	} else if lonFloat, ok := result["longitude"].(float64); ok {
		gps.SetLongitude(lonFloat)
	}

	return gps
}

func getGpsHelperPath() string {
	cmd := exec.Command("powershell", "-Command",
		"(Get-AppxPackage -Name 'GpsHelper').InstallLocation")
	out, err := cmd.Output()
	if err != nil || strings.TrimSpace(string(out)) == "" {
		return ""
	}
	return filepath.Join(strings.TrimSpace(string(out)), "gpshelper.exe")
}

// GetGPSDataByLocationAPI retrieves GPS data using a helper application.
func (a *Agent) GetGPSDataByLocationAPI() *data.GPS {
	gpsHelperPath := getGpsHelperPath()
	if gpsHelperPath == "" {
		fmt.Println("Error: GpsHelper App nicht installiert!")
		return nil
	}

	cmd := exec.Command(gpsHelperPath)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing GPS helper: %v\n", err)
		return nil
	}

	gps := &data.GPS{}
	if err = json.Unmarshal(output, gps); err != nil {
		fmt.Printf("Error unmarshaling GPS data: %v\n", err)
		return nil
	}

	//	if err = enrichGPSData(gps); err != nil {
	//		fmt.Printf("Warning: Could not enrich GPS data: %v\n", err)
	//	}

	return gps
}
