package twitter_job

import (
	"time"
)

type Config struct {
	Site_url string
	Site_api string
	Twitter struct {
		Client_token string
		Client_secret string
		Access_token string
		Access_secret string
		Screen_name string
	}
	Redis struct {
		Netaddr string
		Db int
		Password string
	}
	Service struct {
		Time_out time.Duration
	}
}
