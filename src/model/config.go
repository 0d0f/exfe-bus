package model

type Config struct {
	SiteUrl            string  `json:"site_url"`
	SiteApi            string  `json:"site_api"`
	SiteImg            string  `json:"site_img"`
	AppUrl             string  `json:"app_url"`
	TemplatePath       string  `json:"template_path"`
	DefaultLang        string  `json:"default_lang"`
	AccessDomain       string  `json:"access_domain"`
	Proxy              string  `json:"proxy"`
	TutorialBotUserIds []int64 `json:"tutorial_bot_user_ids"`

	Debug bool `json:"debug"`

	DB struct {
		Addr              string `json:"addr"`
		Port              uint   `json:"port"`
		Username          string `json:"username"`
		Password          string `json:"password"`
		DbName            string `json:"db_name"`
		MaxConnections    uint   `json:"max_connections"`
		HeartBeatInSecond uint   `json:"heart_beat_in_second"`
	} `json:"db"`
	Redis struct {
		Netaddr  string `json:"netaddr"`
		Db       int    `json:"db"`
		Password string `json:"password"`
	} `json:"redis"`
	RedisCache struct {
		Netaddr  string `json:"netaddr"`
		Db       int    `json:"db"`
		Password string `json:"password"`
	} `json:"redis_cache"`
	Email struct {
		Host             string `json:"host"`
		Username         string `json:"username"`
		Password         string `json:"password"`
		Name             string `json:"name"`
		Prefix           string `json:"prefix"`
		Domain           string `json:"domain"`
		IdleTimeoutInSec uint   `json:"idle_timeout_in_sec"`
		IntervalInSec    uint   `json:"interval_in_sec"`
	} `json:"email"`
	AWS struct {
		S3 struct {
			Key          string `json:"key"`
			Secret       string `json:"secret"`
			Domain       string `json:"domain"`
			BucketPrefix string `json:"bucket_prefix"`
		} `json:"s3"`
	}
	Google struct {
		Key string `json:"key"`
	} `json:"google"`

	Dispatcher map[string]map[string]string `json:"dispatcher"`

	ExfeService struct {
		Addr     string `json:"addr"`
		Port     uint   `json:"port"`
		Services struct {
			Token     bool `json:"token"`
			Iom       bool `json:"iom"`
			Thirdpart bool `json:"thirdpart"`
			Notifier  bool `json:"notifier"`
			Live      bool `json:"live"`
			Splitter  bool `json:"splitter"`
			Routex    bool `json:"routex"`
		} `json:"services"`
	} `json:"exfe_service"`
	ExfeQueue struct {
		Addr string `json:"addr"`
		Port uint   `json:"port"`
	} `json:"exfe_queue"`
	Wechat map[string]struct {
		Addr     string `json:"addr"`
		Port     uint   `json:"port"`
		PingId   string `json:"ping_id"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"wechat"`

	Splitter struct {
		SpeedOn map[string]int64 `json:"speed_on"`
	}
	Here struct {
		Threshold       float64 `json:"threshold"`
		SignThreshold   float64 `json:"sign_threshold"`
		TimeoutInSecond int     `json:"timeout_in_second"`
	} `json:"here"`
	Routex struct {
		TutorialCreator  int64             `json:"tutorial_creator"`
		TutorialDataFile map[string]string `json:"tutorial_data_file"`
	}
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
			Apps map[string]struct {
				Cert   string `json:"cert"`
				Key    string `json:"key"`
				Server string `json:"server"`
				RootCA string `json:"rootca"`
			} `json:"apps"`
			Default string `json:"default"`
		} `json:"apn"`
		Gcm struct {
			Key string `json:"key"`
		} `json:"gcm"`
		Sms struct {
			AllToiMsg bool `json:"all_to_imsg"`
			Twilio    struct {
				Url       string `json:"url"`
				FromPhone string `json:"from_phone"`
			} `json:"twilio"`
			DuanCaiWang struct {
				Url string `json:"url"`
			} `json:"duancaiwang"`
		}
		IMessage struct {
			Address        string   `json:"address"`
			Origin         string   `json:"origin"`
			QueueDepth     int      `json:"queue_depth"`
			Channels       []string `json:"channels"`
			PeriodInSecond int      `json:"period_in_second"`
		} `json:"imessage"`
		Dropbox struct {
			Key    string `json:"key"`
			Secret string `json:"secret"`
		} `json:"dropbox"`
		Photostream struct {
			Domain string `json:"domain"`
		} `json:"photostream"`
	} `json:"thirdpart"`
	Bot struct {
		Email struct {
			IMAPHost        string `json:"imap_host"`
			IMAPUser        string `json:"imap_user"`
			IMAPPassword    string `json:"imap_password"`
			TimeoutInSecond uint   `json:"timeout_in_second"`
		} `json:"email"`
	} `json:"bot"`
}
