package exfe_service

import (
	"apn/service"
	"c2dm/service"
	"fmt"
	"gobus"
	"log"
	"os"
)

type CrossPost struct {
	twitter *gobus.Client
	ios     *gobus.Client
	android *gobus.Client
	log     *log.Logger
	config  *Config
}

func NewCrossPost(config *Config) *CrossPost {
	return &CrossPost{
		ios:     gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, "iOS"),
		android: gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, "Android"),
		log:     log.New(os.Stderr, "exfe.cross.post", log.LstdFlags),
		config:  config,
	}
}

func (s *CrossPost) SendPost(arg *OneIdentityUpdateArg) {
	if arg.By_identity.Connected_user_id == arg.To_identity.Connected_user_id {
		return
	}
	by := arg.By_identity.Name
	switch arg.To_identity.Provider {
	case "iOS":
		s.SendApn(arg, fmt.Sprintf("%s %s", by, arg.Post.Content))
	case "Android":
		s.SendAndroid(arg, fmt.Sprintf("%s %s", by, arg.Post.Content))
	}
}

func (s *CrossPost) SendApn(arg *OneIdentityUpdateArg, msg string) {
	a := apn_service.ApnSendArg{
		DeviceToken: arg.To_identity.External_id,
		Alert:       msg,
		Badge:       1,
		Sound:       "default",
		Cid:         arg.Cross.Id,
		T:           "c",
	}
	s.ios.Send("ApnSend", &a, 5)
}

func (s *CrossPost) SendAndroid(arg *OneIdentityUpdateArg, msg string) {
	a := c2dm_service.C2DMSendArg{
		DeviceID: arg.To_identity.External_id,
		Message:  msg,
		Cid:      arg.Cross.Id,
		T:        "c",
		Badge:    0,
		Sound:    "default",
	}
	s.android.Send("C2DMSend", &a, 5)
}
