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
	"strings"
	"thirdpart"
	"thirdpart/_performance"
	"thirdpart/apn"
	"thirdpart/dropbox"
	"thirdpart/email"
	"thirdpart/facebook"
	"thirdpart/gcm"
	"thirdpart/imessage"
	"thirdpart/phone"
	"thirdpart/photostream"
	"thirdpart/twitter"
	"time"
)

func registerThirdpart(config *model.Config, platform *broker.Platform) (*thirdpart.Poster, error) {
	poster, err := thirdpart.NewPoster()
	if err != nil {
		return nil, err
	}

	if config.Thirdpart.MaxStateCache == 0 {
		return nil, fmt.Errorf("config.Thirdpart.MaxStateCache should be bigger than 0")
	}

	apns_, err := apns.New(config.Thirdpart.Apn.Cert, config.Thirdpart.Apn.Key, config.Thirdpart.Apn.Server, time.Duration(config.Thirdpart.Apn.TimeoutInMinutes)*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("can't connect apn: %s", err)
	}
	gcms_ := gcms.New(config.Thirdpart.Gcm.Key)
	helper := thirdpart.NewHelper(config)

	twitter_ := twitter.New(config, helper)
	poster.Add(twitter_)

	facebook_ := facebook.New(helper)
	poster.Add(facebook_)

	email_ := email.New(helper)
	poster.Add(email_)

	apn_ := apn.New(apns_, getApnErrorHandler(config.Log.SubPrefix("apn error")))
	poster.Add(apn_)

	gcm_ := gcm.New(gcms_)
	poster.Add(gcm_)

	imsg_, err := imessage.New(config)
	if err != nil {
		return nil, fmt.Errorf("can't connect imessage: %s", err)
	}
	poster.Add(imsg_)

	phone_, err := phone.New(config)
	if err != nil {
		return nil, fmt.Errorf("can't create phone: %s", err)
	}
	poster.Add(phone_)

	imsgPhone := phone.NewIMsgPhone(phone_, imsg_)
	poster.Add(imsgPhone)

	if config.Debug {
		performance := _performance.New()
		poster.Add(performance)
	}

	return poster, nil
}

type Thirdpart struct {
	thirdpart *thirdpart.Thirdpart
	log       *logger.SubLogger
	config    *model.Config
	platform  *broker.Platform
}

func NewThirdpart(config *model.Config, platform *broker.Platform) (*Thirdpart, error) {
	if config.Thirdpart.MaxStateCache == 0 {
		return nil, fmt.Errorf("config.Thirdpart.MaxStateCache should be bigger than 0")
	}

	helper := thirdpart.NewHelper(config)

	t := thirdpart.New(config)

	twitter_ := twitter.New(config, helper)
	t.AddUpdater(twitter_)

	facebook_ := facebook.New(helper)
	t.AddUpdater(facebook_)

	if config.Debug {
		performance := _performance.New()
		t.AddUpdater(performance)
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
	}, nil
}

func (t *Thirdpart) SetRoute(route gobus.RouteCreater) error {
	json := new(gobus.JSON)
	route().Methods("POST").Path("/thirdpart/identity").HandlerMethod(json, t, "UpdateIdentity")
	route().Methods("POST").Path("/thirdpart/friends").HandlerMethod(json, t, "UpdateFriends")
	route().Methods("POST").Path("/thirdpart/photographers").HandlerMethod(json, t, "GrabPhotos")
	route().Methods("POST").Path("/thirdpart/photographers/photos").HandlerMethod(json, t, "GetPhotos")

	// old
	route().Methods("POST").Path("/Thirdpart").Queries("method", "UpdateIdentity").HandlerMethod(json, t, "UpdateIdentity")
	route().Methods("POST").Path("/Thirdpart").Queries("method", "UpdateFriends").HandlerMethod(json, t, "UpdateFriends")

	return nil
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

// 抓取渠道to上图片pictureIDs的图片。bus地质：bus://exfe_service/thirdpart/photographers
//
// 例子：
//
//   > curl "http://127.0.0.1:23333/thirdpart/photographers/photos?picture_id=/Photos/underwater/001.jpg,/Photos/underwater/002.jpg" -d '{"external_id":"123","external_username":"name","auth_data":"{\"oauth_token\":\"key\",\"oauth_token_secret\":\"secret\"}","provider":"dropbox","identity_id":789,"user_id":1}'
//
func (t *Thirdpart) GetPhotos(params map[string]string, to model.Recipient) ([]string, error) {
	pictureIDs := strings.Split(params["picture_id"], ",")
	if len(pictureIDs) == 0 {
		return nil, fmt.Errorf("must give picture_id")
	}
	datas, err := t.thirdpart.GetPhotos(to, pictureIDs)
	if err != nil {
		return nil, err
	}
	return datas, nil
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
	return func(err error) {
		log.Err("%s", err)
	}
}

type callbackArg struct {
	Recipient model.Recipient `json:"recipient"`
	Error     string          `json:"error"`
}
