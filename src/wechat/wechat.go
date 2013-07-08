package main

import (
	"broker"
	"bytes"
	"crypto/tls"
	"daemon"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/googollee/go-aws/s3"
	"io/ioutil"
	"logger"
	"model"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type BaseRequest struct {
	Uin      uint64
	Sid      string
	Skey     string
	DeviceID string
}

type SyncKey struct {
	Count int
	List  []map[string]int
}

type ContactRequest struct {
	UserName   string
	ChatRoomId uint64
}

type Request struct {
	BaseRequest BaseRequest
	Count       int              `json:"Count,omitempty"`
	List        []ContactRequest `json:"List,omitempty"`
	SyncKey     *SyncKey         `json:"SyncKey,omitempty"`
	Msg         *Message         `json:"Msg,omitempty"`
}

type Member struct {
	AttrStatus      int
	DisplayName     string
	MemberStatus    int
	NickName        string
	PYInitial       string
	PYQuanPin       string
	RemarkPYInitial string
	RemarkPYQuanPin string
	Uin             uint64
	UserName        string
}

type Contact struct {
	Alias            string
	AppAccountFlag   int
	AttrStatus       int
	City             string
	ContactFlag      int
	HeadImgUrl       string
	HideInputBarFlag int
	MemberCount      int
	MemberList       []Member
	NickName         string
	OwnerUin         uint64
	PYInitial        string
	PYQuanPin        string
	Province         string
	RemarkName       string
	RemarkPYInitial  string
	RemarkPYQuanPin  string
	Sex              int
	Signature        string
	SnsFlag          int
	StarFriend       int
	Statues          int
	Uin              uint64
	UniFriend        int
	UserName         string
	VerifyFlag       int
}

type Message struct {
	AppInfo struct {
		AddID string
		Type  int
	}
	AppMsgType           int
	Content              string
	CreateTime           int64
	ClientMsgId          int
	FileName             string
	FileSize             string
	ForwardFlag          int
	FromUserName         string
	ImgStatus            int
	LocalID              int
	MediaId              string
	MsgId                int64
	MsgType              int
	PlayLength           int
	Status               int
	StatusNotifyCode     int
	StatusNotifyUserName string
	ToUserName           string
	Type                 int
	Url                  string
	VoiceLength          int
}

func (m Message) IsChatRoom() bool {
	return m.MsgType == 10000
}

type Response struct {
	BaseResponse struct {
		ErrMsg string
		Ret    int
	}
	User Contact

	AddMsgCount  int
	AddMsgList   []Message
	ContinueFlag int

	Count       int
	ContactList []Contact

	Skey    string
	SyncKey SyncKey
}

type WeChat struct {
	client      *http.Client
	baseRequest BaseRequest
	syncKey     SyncKey
	userName    string
	lastPing    time.Time
	pingId      string
	pingIndex   int
}

func New(pingId string, config *model.Config) (*WeChat, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	b, err := resp(http.Get("https://login.weixin.qq.com/jslogin?appid=wx782c26e4c19acffb&redirect_uri=https%3A%2F%2Fwx.qq.com%2Fcgi-bin%2Fmmwebwx-bin%2Fwebwxnewloginpage&fun=new&lang=zh_CN"))
	if err != nil {
		return nil, err
	}
	key := "window.QRLogin.uuid = \""
	index := bytes.Index(b, []byte(key))
	if index < 0 {
		return nil, fmt.Errorf("can't find key uuid in %s", string(b))
	}
	b = b[index+len(key):]
	end := bytes.Index(b, []byte("\""))
	if index < 0 {
		return nil, fmt.Errorf("can't find key uuid in %s", string(b))
	}
	uuid := string(b[:end])
	loginUrl := fmt.Sprintf("https://login.weixin.qq.com/qrcode/%s?t=webwx", uuid)
	logger.NOTICE("login: %s", loginUrl)

	sendmail(config, `Content-Type: text/plain
To: srv-op@exfe.com
From: =?utf-8?B?U2VydmljZSBOb3RpZmljYXRpb24=?= <x@exfe.com>
Subject: =?utf-8?B?ISEhIVdlQ2hhdCBuZWVkIGxvZ2luISEhIQo=?=

WeChat need login!!! Help!!!!
QR: `+loginUrl)

	login := fmt.Sprintf("https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?uuid=%s&tip=1", uuid)
	var tokenUrl string
	startTime := time.Now()
	for {
		now := time.Now()
		if now.Sub(startTime) > 5*time.Minute {
			return nil, fmt.Errorf("login timeout, need restart")
		}
		b, err = resp(client.Get(login))
		if err != nil {
			return nil, err
		}
		target := "window.redirect_uri=\""
		index := bytes.Index(b, []byte(target))
		if index < 0 {
			continue
		}
		b = b[index+len(target):]
		end := bytes.Index(b, []byte("\""))
		if end < 0 {
			return nil, fmt.Errorf("can't find token end in %s", string(b))
		}
		tokenUrl = string(b[:end]) + "&fun=new"
		break
	}

	re, err := client.Get(tokenUrl)
	b, err = resp(re, err)
	if err != nil {
		return nil, err
	}

	uin, sid, err := grabCookie(re)
	if err != nil {
		return nil, err
	}

	baseRequest := BaseRequest{
		Uin:      uin,
		Sid:      sid,
		DeviceID: uuid,
	}
	req := Request{
		BaseRequest: baseRequest,
	}

	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err = encoder.Encode(req)
	if err != nil {
		return nil, err
	}
	var ret Response
	re, err = client.Post("https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxinit", "application/json", buf)
	err = respJson(&ret, re, err)
	if err != nil {
		return nil, err
	}
	if ret.BaseResponse.Ret != 0 {
		return nil, fmt.Errorf("webwxinit error: %s", ret.BaseResponse.ErrMsg)
	}
	baseRequest.Skey = ret.Skey

	buf = bytes.NewBuffer(nil)
	encoder = json.NewEncoder(buf)
	err = encoder.Encode(map[string]interface{}{
		"BaseRequest":  baseRequest,
		"Code":         3,
		"FromUserName": ret.User.UserName,
		"ToUserName":   ret.User.UserName,
		"ClientMsgId":  timestamp(),
	})
	re, err = client.Post("https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxstatusnotify", "application/json", buf)
	if err != nil {
		return nil, err
	}
	re.Body.Close()

	sendmail(config, `Content-Type: text/plain
To: srv-op@exfe.com
From: =?utf-8?B?U2VydmljZSBOb3RpZmljYXRpb24=?= <x@exfe.com>
Subject: =?utf-8?B?V2VDaGF0IGxvZ2luZWQK?=

WeChat logined as `+ret.User.UserName)

	return &WeChat{
		client:      client,
		baseRequest: baseRequest,
		syncKey:     ret.SyncKey,
		userName:    ret.User.UserName,
		lastPing:    time.Now(),
		pingId:      pingId,
		pingIndex:   0,
	}, nil
}

func (wc *WeChat) SendMessage(to, content string) error {
	req := Request{
		BaseRequest: wc.baseRequest,
		Msg: &Message{
			FromUserName: wc.userName,
			ToUserName:   to,
			Type:         1,
			Content:      content,
			ClientMsgId:  1,
			LocalID:      1,
		},
	}
	var resp Response
	err := wc.postJson("https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxsendmsg?sid="+wc.baseRequest.Sid, req, &resp)
	if err != nil {
		return err
	}
	if resp.BaseResponse.Ret != 0 {
		return fmt.Errorf("%s", resp.BaseResponse.ErrMsg)
	}
	wc.baseRequest.Skey = resp.Skey
	wc.lastPing = time.Now()
	return nil
}

func (wc *WeChat) Ping(timeout time.Duration) error {
	if time.Now().Sub(wc.lastPing) < timeout {
		fmt.Print(".")
		return nil
	}
	fmt.Print("+")
	msgs := []string{"早", "hi", "喂", "what", "敲", "lol"}
	err := wc.SendMessage(wc.pingId, msgs[wc.pingIndex])
	if err != nil {
		return err
	}
	wc.pingIndex = (wc.pingIndex + 1) % len(msgs)
	return nil
}

func (wc *WeChat) Checkt() (string, error) {
	params := make(url.Values)
	params.Set("callback", "_")
	params.Set("sid", wc.baseRequest.Sid)
	params.Set("uin", fmt.Sprintf("%d", wc.baseRequest.Uin))
	params.Set("deviceid", wc.baseRequest.DeviceID)
	params.Set("synckey", makeSyncQuery(wc.syncKey.List))
	params.Set("_", timestamp())
	resp, err := wc.get("https://webpush.weixin.qq.com/cgi-bin/mmwebwx-bin/synccheck?" + params.Encode())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	const ret = "retcode:\""
	index := bytes.Index(b, []byte(ret))
	if index < 0 {
		return "", fmt.Errorf("sync error, no retcode: %s", string(b))
	}
	s := b[index+len(ret):]
	end := bytes.Index(s, []byte("\""))
	if end < 0 {
		return "", fmt.Errorf("sync error, no retcode: %s", string(b))
	}
	if retcode := string(s[:end]); retcode != "0" {
		return "", fmt.Errorf("sync error: %s", retcode)
	}

	const begin = "selector:\""
	index = bytes.Index(b, []byte(begin))
	if index < 0 {
		return "", fmt.Errorf("sync error, no selector: %s", string(b))
	}
	s = b[index+len(begin):]
	end = bytes.Index(s, []byte("\""))
	if end < 0 {
		return "", fmt.Errorf("sync error, no selector: %s", string(b))
	}
	return string(s[:end]), nil
}

func (wc *WeChat) GetLast() (*Response, error) {
	req := Request{
		BaseRequest: wc.baseRequest,
		SyncKey:     &wc.syncKey,
	}
	var resp Response
	err := wc.postJson("https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxsync?sid="+wc.baseRequest.Sid, req, &resp)
	if err != nil {
		return nil, err
	}
	if resp.BaseResponse.Ret != 0 {
		return nil, fmt.Errorf("%s", resp.BaseResponse.ErrMsg)
	}
	wc.baseRequest.Skey = resp.Skey
	wc.syncKey = resp.SyncKey
	return &resp, nil
}

func (wc *WeChat) GetContact(reqContacts []ContactRequest) ([]Contact, error) {
	req := Request{
		BaseRequest: wc.baseRequest,
		Count:       len(reqContacts),
		List:        reqContacts,
		SyncKey:     &wc.syncKey,
	}
	var resp Response
	err := wc.postJson("https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxbatchgetcontact?type=ex", req, &resp)
	if err != nil {
		return nil, err
	}
	return resp.ContactList, nil
}

func (wc *WeChat) ConvertCross(bucket *s3.Bucket, msg *Message) (uint64, model.Cross, error) {
	if !msg.IsChatRoom() {
		return 0, model.Cross{}, fmt.Errorf("%s", "not chat room")
	}
	chatroomReq := []ContactRequest{
		ContactRequest{
			UserName: msg.FromUserName,
		},
	}
	chatrooms, err := wc.GetContact(chatroomReq)
	if err != nil {
		return 0, model.Cross{}, err
	}
	var chatroom Contact
	for _, c := range chatrooms {
		if c.UserName == msg.FromUserName {
			chatroom = c
			break
		}
	}
	if chatroom.UserName != msg.FromUserName {
		return 0, model.Cross{}, fmt.Errorf("can't find chatroom %s", msg.FromUserName)
	}
	var contactsReq []ContactRequest
	for _, m := range chatroom.MemberList {
		contactsReq = append(contactsReq, ContactRequest{
			UserName:   m.UserName,
			ChatRoomId: chatroom.Uin,
		})
	}
	contacts, err := wc.GetContact(contactsReq)
	if err != nil {
		return 0, model.Cross{}, err
	}
	ret := model.Cross{}
	ret.Title = "Cross with "
	ret.Exfee.Invitations = make([]model.Invitation, len(contacts))
	var host *model.Identity
	for i, member := range contacts {
		if i < 3 {
			ret.Title += member.NickName + ", "
		}
		headerUrl := "https://wx.qq.com" + member.HeadImgUrl
		headerPath := fmt.Sprintf("/thirdpart/weichat/%d.jpg", member.Uin)
		resp, err := wc.get(headerUrl)
		if err == nil {
			obj, err := bucket.CreateObject(headerPath, resp.Header.Get("Content-Type"))
			if err == nil {
				obj.SetDate(time.Now())
				length, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
				if err == nil {
					err = obj.SaveReader(resp.Body, length)
					if err == nil {
						headerUrl = obj.URL()
					}
				}
			}
		}
		ret.Exfee.Invitations[i].Identity = model.Identity{
			ExternalID:       fmt.Sprintf("%d", member.Uin),
			Provider:         "wechat",
			ExternalUsername: member.UserName,
			Nickname:         member.NickName,
			Avatar:           headerUrl,
		}
		if member.Uin == chatroom.OwnerUin {
			ret.Exfee.Invitations[i].Host = true
			host = &ret.Exfee.Invitations[i].Identity
		}
	}
	ret.Title = ret.Title[:len(ret.Title)-2]
	ret.By = *host
	for i := range ret.Exfee.Invitations {
		ret.Exfee.Invitations[i].By = *host
		ret.Exfee.Invitations[i].UpdatedBy = *host
	}
	return chatroom.Uin, ret, nil
}

func (wc *WeChat) postJson(url string, data interface{}, reply interface{}) error {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Cookie", fmt.Sprintf("wxuin=%d; wxsid=%s", wc.baseRequest.Uin, wc.baseRequest.Sid))
	re, err := wc.client.Do(req)
	err = respJson(reply, re, err)
	if err != nil {
		return err
	}
	return nil
}

func (wc *WeChat) get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.116 Safari/537.36")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Origin", "https://wx.qq.com")
	req.Header.Set("Referer", "https://wx.qq.com/")
	req.Header.Set("Cookie", fmt.Sprintf("wxuin=%d; wxsid=%s", wc.baseRequest.Uin, wc.baseRequest.Sid))
	re, err := wc.client.Do(req)
	if err != nil {
		return nil, err
	}
	if re.StatusCode != http.StatusOK {
		re.Body.Close()
		return nil, fmt.Errorf("%s", re.Status)
	}
	return re, nil
}

func timestamp() string {
	now := time.Now().UTC()
	return fmt.Sprintf("%d\n", now.UnixNano()/int64(time.Millisecond))
}

func resp(r *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", r.Status)
	}
	return b, nil
}

func respJson(v interface{}, r *http.Response, err error) error {
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", r.Status)
	}
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(v)
}

func grabCookie(resp *http.Response) (uint64, string, error) {
	var uin uint64
	var sid string
	for _, v := range resp.Header["Set-Cookie"] {
		wxsid := "wxsid="
		if index := strings.Index(v, wxsid); index >= 0 {
			v = v[len(wxsid):]
			index = strings.Index(v, ";")
			if index < 0 {
				return 0, "", fmt.Errorf("invalid sid: %s", v)
			}
			sid = v[:index]
		}
		wxuin := "wxuin="
		if index := strings.Index(v, wxuin); index >= 0 {
			v = v[len(wxuin):]
			index = strings.Index(v, ";")
			if index < 0 {
				return 0, "", fmt.Errorf("invalid uin: %s", v)
			}
			var err error
			uin, err = strconv.ParseUint(v[:index], 10, 64)
			if err != nil {
				return 0, "", err
			}
		}
	}
	return uin, sid, nil
}

func makeSyncQuery(syncKey []map[string]int) string {
	ret := ""
	l := len(syncKey)
	for i, k := range syncKey {
		ret += fmt.Sprintf("%d_%d", k["Key"], k["Val"])
		if i != l-1 {
			ret += "|"
		}
	}
	return ret
}

func sendmail(config *model.Config, content string) {
	post := fmt.Sprintf("http://%s:%d/v3/poster/email/srv-op@exfe.com", config.ExfeService.Addr, config.ExfeService.Port)
	queue := fmt.Sprintf("http://%s:%d/v3/queue/-/POST/%s?ontime=%d&update=once", config.ExfeQueue.Addr, config.ExfeQueue.Port, base64.URLEncoding.EncodeToString([]byte(post)), time.Now().Unix())
	b, _ := json.Marshal(content)
	{
		resp, err := broker.Http("POST", queue, "plain/text", b)
		if err != nil {
			logger.ERROR("send notification %s failed: %s with %s", queue, err, string(b))
		} else {
			resp.Body.Close()
		}
	}
}

func main() {
	var config model.Config
	_, quit := daemon.Init("exfe.json", &config)

	workType := os.Args[len(os.Args)-1]
	work, ok := config.Wechat[workType]
	if !ok {
		logger.ERROR("unknow work type %s", workType)
		return
	}

	aws := s3.New(config.AWS.S3.Domain, config.AWS.S3.Key, config.AWS.S3.Secret)
	aws.SetACL(s3.ACLPublicRead)
	aws.SetLocationConstraint(s3.LC_AP_SINGAPORE)
	bucket, err := aws.GetBucket(fmt.Sprintf("%s-3rdpart-photos", config.AWS.S3.BucketPrefix))
	if err != nil {
		logger.ERROR("can't create bucket: %s", err)
		return
	}

	platform, err := broker.NewPlatform(&config)
	if err != nil {
		logger.ERROR("can't create platform: %s", err)
		return
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8&autocommit=true",
		config.DB.Username, config.DB.Password, config.DB.Addr, config.DB.Port, config.DB.DbName))
	if err != nil {
		logger.ERROR("mysql error:", err)
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		logger.ERROR("mysql error:", err)
		return
	}
	kvSaver := broker.NewKVSaver(db)

	wc, err := New(work.PingId, &config)
	if err != nil {
		logger.ERROR("can't create wechat: %s", err)
		return
	}
	defer func() {
		logger.NOTICE("quit")
	}()

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		to := r.URL.Query().Get("to")
		if strings.HasPrefix(to, "e") && strings.HasSuffix(to, "@exfe") {
			defer r.Body.Close()
			chatroomId, exist, err := kvSaver.Check([]string{to})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !exist {
				http.Error(w, fmt.Sprintf("can't find chatroom for %s", to), http.StatusBadRequest)
				return
			}
			to = chatroomId
		}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = wc.SendMessage(to, string(b))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	go http.ListenAndServe(fmt.Sprintf("%s:%d", work.Addr, work.Port), nil)

	logger.NOTICE("login as %s", wc.userName)

	for {
		select {
		case <-quit:
			return
		default:
		}
		ret, err := wc.Checkt()
		if err != nil {
			logger.ERROR("can't get last message: %s", err)
			return
		}
		if ret == "0" {
			err = wc.Ping(time.Minute * 30)
			if err != nil {
				logger.ERROR("ping error: %s", err)
				return
			}
			continue
		}
		resp, err := wc.GetLast()
		if err != nil {
			logger.ERROR("can't get last message: %s", err)
			return
		}
		for _, msg := range resp.AddMsgList {
			if !msg.IsChatRoom() {
				continue
			}
			uin, cross, err := wc.ConvertCross(bucket, &msg)
			if err != nil {
				logger.ERROR("can't convert to cross: %s", err)
				return
			}
			uinStr := fmt.Sprintf("%d", uin)
			idStr, exist, err := kvSaver.Check([]string{uinStr})
			if err != nil {
				logger.ERROR("can't check uin %s: %s", uinStr, err)
				continue
			}
			if exist {
				id, err := strconv.ParseInt(idStr, 10, 64)
				if err != nil {
					goto CREATE
				}
				oldCross, err := platform.FindCross(id, nil)
				if err != nil {
					goto CREATE
				}
				exfee := make(map[string]bool)
				host := cross.Exfee.Invitations[0].Identity
				for _, invitation := range cross.Exfee.Invitations {
					exfee[invitation.Identity.Id()] = true
					if invitation.Host {
						host = invitation.Identity
					}
				}
				for _, invitation := range oldCross.Exfee.Invitations {
					if exfee[invitation.Identity.Id()] {
						continue
					}
					invitation.Response = model.Removed
					cross.Exfee.Invitations = append(cross.Exfee.Invitations, invitation)
				}
				err = platform.BotCrossUpdate("cross_id", idStr, cross, host)
				if err != nil {
					logger.ERROR("can't update cross %s: %s", idStr, err)
					goto CREATE
				}
			}
		CREATE:
			cross, err = platform.BotCrossGather(cross)
			if err != nil {
				logger.ERROR("can't gather cross: %s", err)
				continue
			}
			err = kvSaver.Save([]string{fmt.Sprintf("%d", uin)}, fmt.Sprintf("%d", cross.ID))
			if err != nil {
				logger.ERROR("can't save cross id: %s", err)
			}
			err = kvSaver.Save([]string{fmt.Sprintf("e%d@exfe", cross.Exfee.ID)}, fmt.Sprintf("%d@chatroom", uin))
			if err != nil {
				logger.ERROR("can't save exfee id: %s", err)
			}
			logger.INFO("wechat_gather", msg.FromUserName, uin, cross.ID, cross.Exfee.ID, err)
		}
	}
}
