package model

import (
	"github.com/googollee/go-logger"
)

type Config struct {
	SiteUrl      string `json:"site_url"`
	SiteApi      string `json:"site_api"`
	SiteImg      string `json:"site_img"`
	AppUrl       string `json:"app_url"`
	TemplatePath string `json:"template_path"`
	DefaultLang  string `json:"default_lang"`

	Test bool `json:"test"`
	DB   struct {
		Addr     string `json:"addr"`
		Port     uint   `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		DbName   string `json:"db_name"`
	} `json:"db"`
	Redis struct {
		Netaddr           string `json:"netaddr"`
		Db                int    `json:"db"`
		Password          string `json:"password"`
		MaxConnections    uint   `json:"max_connections"`
		HeartBeatInSecond uint   `json:"heart_beat_in_second"`
	} `json:"redis"`
	Email struct {
		Host             string `json:"host"`
		Port             uint   `json:"port"`
		Username         string `json:"username"`
		Password         string `json:"password"`
		Name             string `json:"name"`
		Domain           string `json:"domain"`
		IdleTimeoutInMin uint   `json:"idle_timeout_in_min"`
		IntervalInMin    uint   `json:"interval_in_min"`
	} `json:"email"`

	Dispatcher struct {
		Sender map[string]string `json:"sender"`
	} `json:"dispatcher"`

	ExfeService struct {
		Addr     string `json:"addr"`
		Port     uint   `json:"port"`
		Services struct {
			TokenManager bool `json:"tokenmanager"`
			Iom          bool `json:"iom"`
			Thirdpart    bool `json:"thirdpart"`
			Notifier     bool `json:"notifier"`
		} `json:"services"`
	} `json:"exfe_service"`
	ExfeQueue struct {
		Addr string          `json:"addr"`
		Port uint            `json:"port"`
		Head map[string]uint `json:"head"`
		Tail map[string]uint `json:"tail"`
	} `json:"exfe_queue"`

	TokenManager struct {
		TableName string `json:"table_name"`
	} `json:"token_manager"`
	Thirdpart struct {
		MaxStateCache uint `json:"max_state_cache"`
		Twitter       struct {
			ClientToken  string `json:"client_token"`
			ClientSecret string `json:"client_secret"`
			AccessToken  string `json:"access_token"`
			AccessSecret string `json:"access_secret"`
			ScreenName   string `json:"screen_name"`
		} `json:"twitter"`
		Apn struct {
			Cert             string `json:"cert"`
			Key              string `json:"key"`
			Server           string `json:"server"`
			RootCA           string `json:"rootca"`
			TimeoutInMinutes uint   `json:"timeout_in_minutes"`
		} `json:"apn"`
		Gcm struct {
			Key string `json:"key"`
		} `json:"gcm"`
	} `json:"thirdpart"`
	Bot struct {
		Email struct {
			IMAPHost        string `json:"imap_host"`
			IMAPUser        string `json:"imap_user"`
			IMAPPassword    string `json:"imap_password"`
			TimeoutInSecond uint   `json:"timeout_in_second"`
		} `json:"email"`
	} `json:"bot"`

	Log *logger.Logger `json:"-"`
}
