package main

import (
	"broker"
	"encoding/json"
	"fmt"
	gcms "github.com/googollee/go-gcm"
	"github.com/googollee/go-rest"
	"logger"
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
	// "thirdpart/imessage"
	"thirdpart/phone"
	"thirdpart/photostream"
	"thirdpart/twitter"
	"thirdpart/wechat"
)

func registerThirdpart(config *model.Config, platform *broker.Platform) (*thirdpart.Poster, error) {
	poster, err := thirdpart.NewPoster()
	if err != nil {
		return nil, err
	}

	if config.Thirdpart.MaxStateCache == 0 {
		return nil, fmt.Errorf("config.Thirdpart.MaxStateCache should be bigger than 0")
	}

	gcms_ := gcms.New(config.Thirdpart.Gcm.Key)
	helper := thirdpart.NewHelper(config)

	twitter_ := twitter.New(config, helper)
	poster.Add(twitter_)

	facebook_ := facebook.New(helper)
	poster.Add(facebook_)

	email_ := email.New(helper)
	poster.Add(email_)

	wechat := wechat.New(config)
	poster.Add(wechat)

	apn_, err := apn.New(config)
	if err != nil {
		return nil, fmt.Errorf("can't connect apn: %s", err)
	}
	poster.Add(apn_)

	gcm_ := gcm.New(gcms_)
	poster.Add(gcm_)

	// imsg_, err := imessage.New(config)
	// if err != nil {
	// 	return nil, fmt.Errorf("can't connect imessage: %s", err)
	// }
	// poster.Add(imsg_)

	phone_, err := phone.New(config)
	if err != nil {
		return nil, fmt.Errorf("can't create phone: %s", err)
	}
	poster.Add(phone_)

	if config.Debug {
		performance := _performance.New()
		poster.Add(performance)
	}

	return poster, nil
}

type Thirdpart struct {
	rest.Service `prefix:"/thirdpart"`

	identity      rest.SimpleNode `route:"/identity" method:"POST"`
	friends       rest.SimpleNode `route:"/friends" method:"POST"`
	photographers rest.SimpleNode `route:"/photographers" method:"POST"`
	photos        rest.SimpleNode `route:"/photographers/photos" method:"POST"`

	updateIdentity rest.SimpleNode `path:"/Thirdpart/UpdateIdentity" method:"POST"`
	updateFriends  rest.SimpleNode `path:"/Thirdpart/UpdateFriends" method:"POST"`

	thirdpart *thirdpart.Thirdpart
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
		config:    config,
		platform:  platform,
	}, nil
}

// 同步更新to在第三方网站的个人信息（头像，bio之类）
//
// 例子：
//
//   > curl http://127.0.0.1:23333/thirdpart/identity -d '{"external_id":"123","external_username":"name","auth_data":"","provider":"twitter","identity_id":789,"user_id":1}'
//
func (t *Thirdpart) Identity(ctx rest.Context, to model.ThirdpartTo) {
	if err := t.thirdpart.UpdateIdentity(&to.To); err != nil {
		ctx.Return(http.StatusBadRequest, "%s", err)
		return
	}
}

// 同步更新to在第三方网站的好友信息
//
// 例子：
//
//   > curl http://127.0.0.1:23333/thirdpart/friends -d '{"external_id":"123","external_username":"name","auth_data":"","provider":"twitter","identity_id":789,"user_id":1}'
//
func (t *Thirdpart) Friends(ctx rest.Context, to model.ThirdpartTo) {
	if err := t.thirdpart.UpdateFriends(&to.To); err != nil {
		ctx.Return(http.StatusBadRequest, "%s", err)
		return
	}
}

// 抓取渠道to上图片库albumID里的图片，并加入crossID里。bus地质：bus://exfe_service/thirdpart/photographers
//
// 例子：
//
//   > curl "http://127.0.0.1:23333/thirdpart/photographers?album_id=/Photos/underwater&cross_id=100354" -d '{"external_id":"123","external_username":"name","auth_data":"{\"oauth_token\":\"key\",\"oauth_token_secret\":\"secret\"}","provider":"dropbox","identity_id":789,"user_id":1}'
//
func (t *Thirdpart) Photographers(ctx rest.Context, to model.Recipient) {
	var albumID, photoxID string
	ctx.Bind("album_id", &albumID)
	ctx.Bind("photox_id", &photoxID)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, "%s", err)
		return
	}
	photos, err := t.thirdpart.GrabPhotos(to, albumID)
	if err != nil {
		ctx.Return(http.StatusInternalServerError, "%s", err)
		return
	}
	if err := t.platform.UploadPhoto(photoxID, photos); err != nil {
		ctx.Return(http.StatusInternalServerError, "%s", err)
		return
	}
}

// 抓取渠道to上图片pictureIDs的图片。bus地质：bus://exfe_service/thirdpart/photographers
//
// 例子：
//
//   > curl "http://127.0.0.1:23333/thirdpart/photographers/photos?picture_id=/Photos/underwater/001.jpg,/Photos/underwater/002.jpg" -d '{"external_id":"123","external_username":"name","auth_data":"{\"oauth_token\":\"key\",\"oauth_token_secret\":\"secret\"}","provider":"dropbox","identity_id":789,"user_id":1}'
//
func (t *Thirdpart) Photos(ctx rest.Context, to model.Recipient) {
	var ids string
	ctx.Bind("picture_id", &ids)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, "%s", err)
		return
	}
	pictureIDs := strings.Split(ids, ",")
	if len(pictureIDs) == 0 {
		ctx.Return(http.StatusBadRequest, "must give picture_id")
		return
	}
	datas, err := t.thirdpart.GetPhotos(to, pictureIDs)
	if err != nil {
		ctx.Return(http.StatusInternalServerError, "%s", err)
		return
	}
	ctx.Render(datas)
}

func (t *Thirdpart) UpdateIdentity(ctx rest.Context, to model.ThirdpartTo) {
	t.Identity(ctx, to)
}

func (t *Thirdpart) UpdateFriends(ctx rest.Context, to model.ThirdpartTo) {
	t.Friends(ctx, to)
}

func (t *Thirdpart) sendCallback(recipient model.Recipient, err error) {
	arg := callbackArg{
		Recipient: recipient,
	}
	if err != nil {
		arg.Error = err.Error()
	}

	url := fmt.Sprintf("%s /v3/bus/notificationcallback", t.config.SiteApi)
	b, err := json.Marshal(arg)
	if err != nil {
		logger.ERROR("encode %s error: %s with %+v", url, err, arg)
		return
	}

	resp, err := broker.HttpResponse(broker.Http("Post", url, "application/json", b))
	if err != nil {
		logger.ERROR("post %s error: %s with %s", url, err, string(b))
		return
	}
	defer resp.Close()
}

type callbackArg struct {
	Recipient model.Recipient `json:"recipient"`
	Error     string          `json:"error"`
}
