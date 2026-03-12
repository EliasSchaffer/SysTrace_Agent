package collector

import (
	"SysTrace_Agent/internal/data/static"
	"SysTrace_Agent/internal/transport"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type GPSCollector struct {
}

func (G GPSCollector) Collect() static.Data {
	gpsData := GetGPSDataByLocationAPI()
	if gpsData != nil {
		return *gpsData
	} else {
		fmt.Println("Failed to get GPS data from Location API, trying IP-based geolocation...")
		gpsData = GetGPSDataByIP()
		if gpsData != nil {
			fmt.Println("")
			return *gpsData
		} else {
			fmt.Println("Failed to get GPS data from both methods.")
		}
	}
	return nil
}

func GetGPSDataByIP() *static.GPS {
	envLoader := transport.ENVLoader{}
	apiKey := envLoader.GetGeoLocationAPIKey()
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

	gps := &static.GPS{}

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

func GetGPSDataByLocationAPI() *static.GPS {
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

	gps := &static.GPS{}
	if err = json.Unmarshal(output, gps); err != nil {
		fmt.Printf("Error unmarshaling GPS data: %v\n", err)
		return nil
	}

	enrichGPSData(gps)

	return gps
}

func enrichGPSData(gps *static.GPS) {
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=jsonv2&lat=%f&lon=%f",
		gps.GetLatitude(), gps.GetLongitude())

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	req.Header.Set("User-Agent", "SysTrace-Agent/1.0 (Windows)")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error fetching enriched GPS data: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var geoResult struct {
		Address struct {
			City    string `json:"city"`
			Town    string `json:"town"`
			Village string `json:"village"`
			State   string `json:"state"`
			Country string `json:"country"`
		} `json:"address"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&geoResult); err != nil {
		fmt.Printf("Error decoding enriched GPS data: %v\n", err)
		return
	}

	city := geoResult.Address.City
	if city == "" {
		city = geoResult.Address.Town
	}
	if city == "" {
		city = geoResult.Address.Village
	}

	gps.SetCity(city)
	gps.SetCountry(geoResult.Address.Country)
	gps.SetRegion(geoResult.Address.State)
}
