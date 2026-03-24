package handler

import (
	"SysTrace_Agent/internal/data/static"
)

type SettingsHandler struct {
	settings static.Settings
}

func (s *SettingsHandler) HandleSettingsChange(settings static.Settings, setNewMasterServer func(string), changeStaticGPS func(static.GPS), changeSendStaticGPS func(bool), changeSendGPS func(bool)) {
	if s.settings == settings {
		return
	}

	if s.settings.MASTER_SERVER_URL != settings.MASTER_SERVER_URL {
		if setNewMasterServer != nil {
			setNewMasterServer(settings.MASTER_SERVER_URL)
		}
	}

	if s.settings.SENDGPS != settings.SENDGPS {
		changeSendGPS(settings.SENDGPS)
	}

	if s.settings.STATICGPS != settings.STATICGPS {
		changeSendGPS(settings.STATICGPS)
	}

	if s.settings.GPS_LATITUDE != settings.GPS_LATITUDE ||
		s.settings.GPS_LONGITUDE != settings.GPS_LONGITUDE ||
		s.settings.GPS_CITY != settings.GPS_CITY ||
		s.settings.GPS_REGION != settings.GPS_REGION ||
		s.settings.GPS_COUNTRY != settings.GPS_COUNTRY {
		gps := static.GPS{
			Latitude:  settings.GPS_LATITUDE,
			Longitude: settings.GPS_LONGITUDE,
			City:      settings.GPS_CITY,
			Region:    settings.GPS_REGION,
			Country:   settings.GPS_COUNTRY,
		}
		changeStaticGPS(gps)

	}

	s.settings = settings
}
