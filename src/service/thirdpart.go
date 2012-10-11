package main

import (
	"broker"
	"fmt"
	gcms "github.com/googollee/go-gcm"
	"github.com/googollee/go-logger"
	"github.com/virushuo/Go-Apns"
	"gobus"
	"model"
	"thirdpart"
	"thirdpart/apn"
	"thirdpart/facebook"
	"thirdpart/gcm"
	"thirdpart/twitter"
	"time"
)

type Thirdpart struct {
	thirdpart *thirdpart.Thirdpart
	log       *logger.SubLogger
	config    *model.Config
}

func NewThirdpart(config *model.Config) (*Thirdpart, error) {
	twitterBroker := broker.NewTwitter(config.Thirdpart.Twitter.ClientToken, config.Thirdpart.Twitter.ClientSecret, config.Thirdpart.Twitter.AccessToken, config.Thirdpart.Twitter.AccessSecret)
	apns_, err := apns.New(config.Thirdpart.Apn.Cert, config.Thirdpart.Apn.Key, config.Thirdpart.Apn.Server, time.Duration(config.Thirdpart.Apn.TimeoutInMinutes)*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("can't connect apn: %s", err)
	}
	gcms_ := gcms.New(config.Thirdpart.Gcm.Key)
	helper := thirdpart.NewHelper(config)

	t := thirdpart.New()

	twitter_ := twitter.New(config.Thirdpart.Twitter.ClientToken, config.Thirdpart.Twitter.ClientSecret, twitterBroker, helper)
	t.AddSender(twitter_)
	t.AddUpdater(twitter_)

	facebook_ := facebook.New(helper)
	t.AddSender(facebook_)
	t.AddUpdater(facebook_)

	apn_ := apn.New(apns_, getApnErrorHandler(config.Log.SubPrefix("apn error")))
	t.AddSender(apn_)

	gcm_ := gcm.New(gcms_)
	t.AddSender(gcm_)

	return &Thirdpart{
		thirdpart: t,
		log:       config.Log.SubPrefix("thirdpart"),
		config:    config,
	}, nil
}

type SendArg struct {
	To             *model.Recipient    `json:"to"`
	PrivateMessage string              `json:"private"`
	PublicMessage  string              `json:"public"`
	Info           *thirdpart.InfoData `json:"info"`
}

// 发信息给to，如果是私人信息，就发送private的内容，如果是公开信息，就发送public的内容。info内是相关的应用信息。
//
// 例子：
//
//   > curl http://127.0.0.1:23333/Thirdpart?method=Send -d '{"to":{"external_id":"123","external_username":"name","auth_data":"","provider":"twitter","identity_id":789,"user_id":1},"public":"public","private":"private","info":{"cross_id":234,"type":"u"}}'
//
func (t *Thirdpart) Send(meta *gobus.HTTPMeta, arg *SendArg, id *string) error {
	log := t.log.SubCode()
	log.Debug("send with %+v", arg)
	var err error
	*id, err = t.thirdpart.Send(arg.To, arg.PrivateMessage, arg.PublicMessage, arg.Info)
	if err != nil {
		log.Err("send with arg(%+v) fail: %s", arg, err)
		return err
	}
	log.Debug("success")
	return nil
}

func getApnErrorHandler(log *logger.SubLogger) apn.ErrorHandler {
	return func(err apns.NotificationError) {
		log.Err("%s", err)
	}
}
