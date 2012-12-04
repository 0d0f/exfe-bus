package main

import (
	"broker"
	"bytes"
	"encoding/json"
	"fmt"
	gcms "github.com/googollee/go-gcm"
	"github.com/googollee/go-logger"
	"github.com/virushuo/Go-Apns"
	"gobus"
	"model"
	"net/http"
	"ringcache"
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
	sendCache *ringcache.RingCache
}

func NewThirdpart(config *model.Config) (*Thirdpart, error) {
	if config.Thirdpart.MaxStateCache == 0 {
		return nil, fmt.Errorf("config.Thirdpart.MaxStateCache should be bigger than 0")
	}

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
		sendCache: ringcache.New(int(config.Thirdpart.MaxStateCache)),
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
	if arg.To.ExternalID == "" {
		go func() {
			err := t.thirdpart.UpdateIdentity(&arg.To)
			if err != nil {
				t.config.Log.Crit("update %s identity error: %s", arg.To, err)
			}
		}()
	}
	*id, err = t.thirdpart.Send(&arg.To, arg.PrivateMessage, arg.PublicMessage, arg.Info)

	key := fmt.Sprintf("%s(%s)@%s", arg.To.ExternalID, arg.To.ExternalUsername, arg.To.Provider)
	lastErr := t.sendCache.Get(key)
	if lastErr == nil {
		t.sendCallback(arg.To, err)
	} else {
		lastError := lastErr.(string)
		if (lastError == "" && err != nil) || (lastError != "" && err == nil) {
			t.sendCallback(arg.To, err)
		}
	}

	if err != nil {
		t.sendCache.Push(key, err.Error())
	} else {
		t.sendCache.Push(key, "")
	}
	return err
}

// 同步更新to在第三方网站的个人信息（头像，bio之类）
//
// 例子：
//
//   > curl http://127.0.0.1:23333/Thirdpart?method=UpdateIdentity -d '{"external_id":"123","external_username":"name","auth_data":"","provider":"twitter","identity_id":789,"user_id":1}'
//
func (t *Thirdpart) UpdateIdentity(meta *gobus.HTTPMeta, to model.ThirdpartTo, i *int) error {
	return t.thirdpart.UpdateIdentity(&to.To)
}

// 同步更新to在第三方网站的好友信息
//
// 例子：
//
//   > curl http://127.0.0.1:23333/Thirdpart?method=UpdateFriends -d '{"external_id":"123","external_username":"name","auth_data":"","provider":"twitter","identity_id":789,"user_id":1}'
//
func (t *Thirdpart) UpdateFriends(meta *gobus.HTTPMeta, to model.ThirdpartTo, i *int) error {
	return t.thirdpart.UpdateFriends(&to.To)
}

func (t *Thirdpart) sendCallback(recipient model.Recipient, err error) {
	arg := callbackArg{
		Recipient: recipient,
	}
	if err != nil {
		arg.Error = err.Error()
	}

	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err = encoder.Encode(arg)
	if err != nil {
		t.log.Crit("encoding json error: %s", err)
		return
	}

	url := fmt.Sprintf("%s/v2/gobus/NotificationCallback", t.config.SiteApi)
	resp, err := http.Post(url, "application/json", buf)
	if err != nil {
		t.log.Crit("send callback(%s) to %s error: %s", buf.String(), url, err)
	}
	if resp.StatusCode != 200 {
		t.log.Crit("send callback(%s) to %s failed: %s", buf.String(), url, resp.Status)
	}
}

func getApnErrorHandler(log *logger.SubLogger) apn.ErrorHandler {
	return func(err apns.NotificationError) {
		log.Err("%s", err)
	}
}

type callbackArg struct {
	Recipient model.Recipient `json:"recipient"`
	Error     string          `json:"error"`
}
