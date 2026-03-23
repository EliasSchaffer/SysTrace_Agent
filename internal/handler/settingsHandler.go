package handler

import "SysTrace_Agent/internal/data/static"

type SettingsHandler struct {
	settings static.Settings
}

func (s *SettingsHandler) HandleSettingsChange(settings static.Settings) {
	if s.settings == settings {
		return
	}

	if s.settings.MASTER_SERVER_URL != settings.MASTER_SERVER_URL {
		// TODO: reconnect agent with new master server URL
	}

	if s.settings.GEOLOCATION_API_KEY != settings.GEOLOCATION_API_KEY {
		// TODO: refresh geolocation fallback configuration
	}

	gpsChanged := false

	if s.settings.SENDGPS != settings.SENDGPS {
		gpsChanged = true
	}

	if s.settings.STATICGPS != settings.STATICGPS {
		gpsChanged = true
	}

	if s.settings.GPS_LATITUDE != settings.GPS_LATITUDE {
		gpsChanged = true
	}

	if s.settings.GPS_LONGITUDE != settings.GPS_LONGITUDE {
		gpsChanged = true
	}

	if s.settings.GPS_CITY != settings.GPS_CITY {
		gpsChanged = true
	}

	if s.settings.GPS_REGION != settings.GPS_REGION {
		gpsChanged = true
	}

	if s.settings.GPS_COUNTRY != settings.GPS_COUNTRY {
		gpsChanged = true
	}

	if gpsChanged {
		// TODO: apply GPS collection mode/values update
	}

	s.settings = settings
}
