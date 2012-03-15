package twitter_job

import (
	"time"
	"config"
)

type Config struct {
	Site_url string
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

func Load(filename string) *Config {
	config := config.LoadFile(filename)
	ret := &Config{}
	ret.Site_url = config.String("site_url")
	ret.Twitter.Client_token = config.String("twitter.client_token")
	ret.Twitter.Client_secret = config.String("twitter.client_secret")
	ret.Twitter.Access_token = config.String("twitter.access_token")
	ret.Twitter.Access_secret = config.String("twitter.access_secret")
	ret.Twitter.Screen_name = config.String("twitter.screen_name")
	ret.Redis.Netaddr = config.String("redis.netaddr")
	ret.Redis.Db = config.Int("redis.db")
	ret.Redis.Password = config.String("redis.password")
	ret.Service.Time_out = time.Duration(config.Int("service.time_out"))
	return ret
}
