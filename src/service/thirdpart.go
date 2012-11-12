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
	"thirdpart/_performance"
	"thirdpart/apn"
	"thirdpart/email"
	"thirdpart/facebook"
	"thirdpart/gcm"
	"thirdpart/imsg"
	"thirdpart/twitter"
	"time"
)

type Thirdpart struct {
	thirdpart *thirdpart.Thirdpart
	log       *logger.SubLogger
	config    *model.Config
}

func NewThirdpart(config *model.Config) (*Thirdpart, error) {
	twitterBroker := broker.NewTwitter(config.Thirdpart.Twitter.ClientToken, config.Thirdpart.Twitter.ClientSecret)
	apns_, err := apns.New(config.Thirdpart.Apn.Cert, config.Thirdpart.Apn.Key, config.Thirdpart.Apn.Server, time.Duration(config.Thirdpart.Apn.TimeoutInMinutes)*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("can't connect apn: %s", err)
	}
	gcms_ := gcms.New(config.Thirdpart.Gcm.Key)
	helper := thirdpart.NewHelper(config)

	t := thirdpart.New(config)

	twitter_ := twitter.New(config, twitterBroker, helper)
	t.AddSender(twitter_)
	t.AddUpdater(twitter_)

	facebook_ := facebook.New(helper)
	t.AddSender(facebook_)
	t.AddUpdater(facebook_)

	email_ := email.New(helper)
	t.AddSender(email_)

	apn_ := apn.New(apns_, getApnErrorHandler(config.Log.SubPrefix("apn error")))
	t.AddSender(apn_)

	gcm_ := gcm.New(gcms_)
	t.AddSender(gcm_)

	if config.Test {
		performance := _performance.New()
		t.AddSender(performance)
		t.AddUpdater(performance)
	}

	imsg_, err := imsg.New(config)
	if err != nil {
		return nil, fmt.Errorf("can't connect imessage: %s", err)
	}
	t.AddSender(imsg_)

	return &Thirdpart{
		thirdpart: t,
		log:       config.Log.SubPrefix("thirdpart"),
		config:    config,
	}, nil
}

// 发信息给to，如果是私人信息，就发送private的内容，如果是公开信息，就发送public的内容。info内是相关的应用信息。
//
// 例子：
//
//   > curl http://127.0.0.1:23333/Thirdpart?method=Send -d '{"to":{"external_id":"123","external_username":"name","auth_data":"","provider":"twitter","identity_id":789,"user_id":1},"private":"private","public":"public","info":null}'
//
func (t *Thirdpart) Send(meta *gobus.HTTPMeta, arg model.ThirdpartSend, id *string) error {
	var err error
	*id, err = t.thirdpart.Send(&arg.To, arg.PrivateMessage, arg.PublicMessage, arg.Info)
	return err
}

// 同步更新to在第三方网站的个人信息（头像，bio之类）
//
// 例子：
//
//   > curl http://127.0.0.1:23333/Thirdpart?method=UpdateIdentity -d '{"external_id":"123","external_username":"name","auth_data":"","provider":"twitter","identity_id":789,"user_id":1}'
//
func (t *Thirdpart) UpdateIdentity(meta *gobus.HTTPMeta, to *model.Recipient, i *int) error {
	return t.thirdpart.UpdateIdentity(to)
}

// 同步更新to在第三方网站的好友信息
//
// 例子：
//
//   > curl http://127.0.0.1:23333/Thirdpart?method=UpdateFriends -d '{"external_id":"123","external_username":"name","auth_data":"","provider":"twitter","identity_id":789,"user_id":1}'
//
func (t *Thirdpart) UpdateFriends(meta *gobus.HTTPMeta, to *model.Recipient, i *int) error {
	return t.thirdpart.UpdateFriends(to)
}

func getApnErrorHandler(log *logger.SubLogger) apn.ErrorHandler {
	return func(err apns.NotificationError) {
		log.Err("%s", err)
	}
}
