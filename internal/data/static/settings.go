package static

type Settings struct {
	GEOLOCATION_API_KEY string
	MASTER_SERVER_URL   string
	LOGFILE_PATH        string
	SENDGPS             bool
	STATICGPS           bool
	GPS_LATITUDE        float64
	GPS_LONGITUDE       float64
	GPS_CITY            string
	GPS_REGION          string
	GPS_COUNTRY         string
}
