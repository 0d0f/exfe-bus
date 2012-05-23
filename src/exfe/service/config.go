package exfe_service

import (
	"time"
)

type Config struct {
	Site_url string
	Site_api string

	Redis struct {
		Netaddr string
		Db int
		Password string
	}
	Cross struct {
		Time_out time.Duration
		Delay map[string]int
	}

	Twitter struct {
		Client_token string
		Client_secret string
		Access_token string
		Access_secret string
		Screen_name string
		Time_out time.Duration
	}
	Apn struct {
		Time_out time.Duration
		Cert string
		Key string
		Server string
	}
	Email struct {
		Time_out time.Duration
		Host string
		Port uint
		User string
		Password string
	}
}
