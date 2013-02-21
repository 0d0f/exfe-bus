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
	"thirdpart/dropbox"
	"thirdpart/email"
	"thirdpart/facebook"
	"thirdpart/gcm"
	"thirdpart/imsg"
	"thirdpart/phone"
	"thirdpart/photostream"
	"thirdpart/twitter"
	"time"
)

type Thirdpart struct {
	thirdpart *thirdpart.Thirdpart
	log       *logger.SubLogger
	config    *model.Config
	platform  *Platform
	sendCache *ringcache.RingCache
}

func NewThirdpart(config *model.Config, streaming *Streaming, platform *Platform) (*Thirdpart, error) {
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

	sms_ := sms.New(config)
	t.AddSender(sms_)

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

	if streaming != nil {
		t.AddSender(streaming)
	}

	dropbox_, err := dropbox.New(config)
	if err != nil {
		return nil, fmt.Errorf("can't create dropbox: %s", err)
	}
	t.AddPhotographer(dropbox_)

	photostream_, err := photostream.New(config)
	if err != nil {
		return nil, fmt.Errorf("can't create photostream: %s", err)
	}
	t.AddPhotographer(photostream_)

	return &Thirdpart{
		thirdpart: t,
		log:       config.Log.SubPrefix("thirdpart"),
		config:    config,
		platform:  platform,
		sendCache: ringcache.New(int(config.Thirdpart.MaxStateCache)),
	}, nil
}

func (t *Thirdpart) SetRoute(route gobus.RouteCreater) error {
	json := new(gobus.JSON)
	route().Methods("POST").Path("/thirdpart/message").HandlerMethod(json, t, "Send")
	route().Methods("POST").Path("/thirdpart/identity").HandlerMethod(json, t, "UpdateIdentity")
	route().Methods("POST").Path("/thirdpart/friends").HandlerMethod(json, t, "UpdateFriends")
	route().Methods("POST").Path("/thirdpart/photographers").HandlerMethod(json, t, "GrabPhotos")

	// old
	route().Methods("POST").Path("/Thirdpart").Queries("method", "Send").HandlerMethod(json, t, "Send")
	route().Methods("POST").Path("/Thirdpart").Queries("method", "UpdateIdentity").HandlerMethod(json, t, "UpdateIdentity")
	route().Methods("POST").Path("/Thirdpart").Queries("method", "UpdateFriends").HandlerMethod(json, t, "UpdateFriends")

	return nil
}

// 发信息给to，如果是私人信息，就发送private的内容，如果是公开信息，就发送public的内容。info内是相关的应用信息。
//
// 例子：
//
//   > curl http://127.0.0.1:23333/thirdpart/message -d '{"to":{"external_id":"123","external_username":"name","auth_data":"","provider":"twitter","identity_id":789,"user_id":1},"private":"private","public":"public","info":null}'
//
func (t *Thirdpart) Send(params map[string]string, arg model.ThirdpartSend) (string, error) {
	if arg.To.ExternalID == "" {
		go func() {
			err := t.thirdpart.UpdateIdentity(&arg.To)
			if err != nil {
				t.config.Log.Crit("update %s identity error: %s", arg.To, err)
			}
		}()
	}
	id, err := t.thirdpart.Send(&arg.To, arg.PrivateMessage, arg.PublicMessage, arg.Info)

	key := fmt.Sprintf("%s(%s)@%s", arg.To.ExternalID, arg.To.ExternalUsername, arg.To.Provider)
	lastErr := t.sendCache.Get(key)
	if lastErr == nil {
		if err == thirdpart.Unreachable {
			t.sendCallback(arg.To, err)
		}
		if err != thirdpart.Unreachable {
			t.sendCallback(arg.To, nil)
		}
	} else {
		lastError := lastErr.(string)
		if lastError != "Unreachable" && err == thirdpart.Unreachable {
			t.sendCallback(arg.To, err)
		}
		if lastError == "Unreachable" && err != thirdpart.Unreachable {
			t.sendCallback(arg.To, nil)
		}
	}
	return id, err
}

// 同步更新to在第三方网站的个人信息（头像，bio之类）
//
// 例子：
//
//   > curl http://127.0.0.1:23333/thirdpart/identity -d '{"external_id":"123","external_username":"name","auth_data":"","provider":"twitter","identity_id":789,"user_id":1}'
//
func (t *Thirdpart) UpdateIdentity(params map[string]string, to model.ThirdpartTo) (int, error) {
	return 0, t.thirdpart.UpdateIdentity(&to.To)
}

// 同步更新to在第三方网站的好友信息
//
// 例子：
//
//   > curl http://127.0.0.1:23333/thirdpart/friends -d '{"external_id":"123","external_username":"name","auth_data":"","provider":"twitter","identity_id":789,"user_id":1}'
//
func (t *Thirdpart) UpdateFriends(params map[string]string, to model.ThirdpartTo) (int, error) {
	return 0, t.thirdpart.UpdateFriends(&to.To)
}

// 抓取渠道to上图片库albumID里的图片，并加入crossID里。bus地质：bus://exfe_service/thirdpart/photographers
//
// 例子：
//
//   > curl "http://127.0.0.1:23333/thirdpart/photographers?album_id=/Photos/underwater&cross_id=100354" -d '{"external_id":"123","external_username":"name","auth_data":"{\"oauth_token\":\"key\",\"oauth_token_secret\":\"secret\"}","provider":"dropbox","identity_id":789,"user_id":1}'
//
func (t *Thirdpart) GrabPhotos(params map[string]string, to model.Recipient) (int, error) {
	albumID := params["album_id"]
	photoxID := params["photox_id"]
	if albumID == "" || photoxID == "" {
		return 0, fmt.Errorf("must give album_id and photox_id")
	}
	photos, err := t.thirdpart.GrabPhotos(to, albumID)
	if err != nil {
		return 0, err
	}
	err = t.platform.UploadPhoto(photoxID, photos)
	if err != nil {
		return 0, err
	}
	return len(photos), nil
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
	bufString := buf.String()

	url := fmt.Sprintf("%s/v2/gobus/NotificationCallback", t.config.SiteApi)
	resp, err := http.Post(url, "application/json", buf)
	if err != nil {
		t.log.Crit("send callback(%s) to %s error: %s", bufString, url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.log.Crit("send callback(%s) to %s failed: %s", bufString, url, resp.Status)
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
