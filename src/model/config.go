package model

import (
	"github.com/googollee/go-logger"
)

type Config struct {
	SiteApi string `json:"site_api"`
	DB      struct {
		Addr     string `json:"addr"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		DbName   string `json:"db_name"`
	} `json:"db"`
	Redis struct {
		Netaddr  string `json:"netaddr"`
		Db       int    `json:"db"`
		Password string `json:"password"`
	} `json:"redis"`
	ThirdPart struct {
		Twitter struct {
			ClientToken  string `json:"client_token"`
			ClientSecret string `json:"client_secret"`
			AccessToken  string `json:"access_token"`
			AccessSecret string `json:"access_secret"`
			ScreenName   string `json:"screen_name"`
		} `json:"twitter"`
		Email struct {
			Host     string `json:"host"`
			Port     int    `json:"port"`
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"thirdpart"`
	}
	ExfeService struct {
		Addr string `json:"addr"`
		Port int    `json:"port"`
	} `json:"exfe_service"`
	TokenManager struct {
		TableName string `json:"table_name"`
	} `json:"token_manager"`

	Log *logger.Logger
}