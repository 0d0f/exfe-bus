package main

import (
	"broker"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"logger"
	"math/rand"
	"model"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type WeChat struct {
	grabSelector *regexp.Regexp
	client       *http.Client
	baseRequest  BaseRequest
	syncKey      SyncKey
	userName     string
	lastPing     time.Time
	pingId       string
	pingIndex    int
}

func New(username, password, pingId string, config *model.Config) (*WeChat, error) {
	jar, err := cookiejar.New(new(cookiejar.Options))
	if err != nil {
		return nil, err
	}
	ret := &WeChat{
		grabSelector: regexp.MustCompile(`selector:"(.*?)"`),
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Jar: jar,
		},
	}
	ret.baseRequest.DeviceID = ret.getDeviceId()

	query := make(url.Values)
	query.Set("appid", "wx782c26e4c19acffb")
	query.Set("redirect_uri", "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage")
	query.Set("fun", "new")
	query.Set("lang", "en_US")
	query.Set("_", fmt.Sprintf("%d", timestamp()))
	b, err := resp(ret.request("GET", "https://login.weixin.qq.com/jslogin", query, nil))
	if err != nil {
		return nil, err
	}

	grabId := regexp.MustCompile(`window.QRLogin.uuid *= *"(.*?)"`)
	ids := grabId.FindAllStringSubmatch(string(b), -1)
	if len(ids) == 0 || len(ids[0]) == 0 {
		return nil, fmt.Errorf("can't find key uuid in %s", string(b))
	}
	uuid := ids[0][1]
	loginQrUrl := fmt.Sprintf("https://login.weixin.qq.com/qrcode/%s?t=webwx", uuid)
	logger.NOTICE("login: %s", loginQrUrl)

	mail := fmt.Sprintf(`Content-Type: text/plain
To: srv-op@exfe.com
From: =?utf-8?B?U2VydmljZSBOb3RpZmljYXRpb24=?= <x@exfe.com>
Subject: =?utf-8?B?V2VjaGF0IFNlcnZpY2UgTm90aWZpY2F0aW9uCg==?=

WeChat need login!!! Help!!!!
QR: %s
Username: %s
Password: %s`, loginQrUrl, username, password)
	sendmail(config, mail)

	grabLoginUri := regexp.MustCompile(`window.redirect_uri="(.*?)";`)
	query = make(url.Values)
	query.Set("uuid", uuid)
	query.Set("tip", "0")
	var tokenUrl string
	startTime := time.Now()
	for {
		now := time.Now()
		if now.Sub(startTime) > 5*time.Minute {
			return nil, fmt.Errorf("login timeout, need restart")
		}
		query.Set("_", fmt.Sprintf("%d", timestamp()))
		r, err := resp(ret.request("GET", "https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login", query, nil))
		if err != nil {
			return nil, err
		}
		urls := grabLoginUri.FindAllStringSubmatch(string(r), -1)
		if len(urls) > 0 && len(urls[0]) > 0 {
			tokenUrl = urls[0][1]
			break
		}
	}

	_, err = resp(ret.request("GET", tokenUrl, nil, nil))
	if err != nil {
		return nil, err
	}
	u, err := url.Parse("https://wx.qq.com/")
	if err != nil {
		return nil, err
	}
	cookies := ret.client.Jar.Cookies(u)
	for _, c := range cookies {
		switch c.Name {
		case "wxsid":
			ret.baseRequest.Sid = c.Value
		case "wxuin":
			ret.baseRequest.Uin, err = strconv.ParseUint(c.Value, 10, 64)
			if err != nil {
				return nil, err
			}
		}
	}

	req := Request{
		BaseRequest: ret.baseRequest,
	}
	query = make(url.Values)
	query.Set("r", fmt.Sprintf("%d", timestamp()))
	var resp Response
	err = ret.postJson("https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxinit", query, req, &resp)
	if err != nil {
		return nil, err
	}
	if resp.BaseResponse.Ret != 0 {
		return nil, fmt.Errorf("webwxinit error: %s", resp.BaseResponse.ErrMsg)
	}

	ret.baseRequest.Skey = resp.Skey
	ret.syncKey = resp.SyncKey
	ret.userName = resp.User.UserName
	ret.lastPing = time.Now()
	ret.pingId = pingId

	sendmail(config, `Content-Type: text/plain
To: srv-op@exfe.com
From: =?utf-8?B?U2VydmljZSBOb3RpZmljYXRpb24=?= <x@exfe.com>
Subject: =?utf-8?B?V2VjaGF0IFNlcnZpY2UgTm90aWZpY2F0aW9uCg==?=

WeChat logined as `+ret.userName)

	return ret, nil
}

func (wc *WeChat) SendMessage(to, content string) error {
	req := Request{
		BaseRequest: wc.baseRequest,
		Msg: &Message{
			FromUserName: wc.userName,
			ToUserName:   to,
			Type:         1,
			Content:      content,
			ClientMsgId:  timestamp(),
			LocalID:      timestamp(),
		},
	}
	var resp Response
	params := make(url.Values)
	params.Set("sid", wc.baseRequest.Sid)
	params.Set("r", fmt.Sprintf("%d", timestamp()))
	err := wc.postJson("https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxsendmsg", params, req, &resp)
	if err != nil {
		return err
	}
	if resp.BaseResponse.Ret != 0 {
		return fmt.Errorf("(%d)%s", resp.BaseResponse.Ret, resp.BaseResponse.ErrMsg)
	}
	wc.baseRequest.Skey = resp.Skey
	wc.lastPing = time.Now()
	return nil
}

func (wc *WeChat) Ping(timeout time.Duration) error {
	if time.Now().Sub(wc.lastPing) < timeout {
		return nil
	}
	msgs := []string{"早", "hi", "喂", "what", "敲", "lol"}
	err := wc.SendMessage(wc.pingId, msgs[wc.pingIndex])
	if err != nil {
		return err
	}
	wc.pingIndex = (wc.pingIndex + 1) % len(msgs)
	return nil
}

func (wc *WeChat) Verify(user, ticket string) error {
	req := VerifyRequest{
		BaseRequest:        wc.baseRequest,
		Opcode:             3,
		SceneList:          []int{0},
		SceneListCount:     1,
		VerifyUserList:     []VerifyUser{VerifyUser{user, ticket}},
		VerifyUserListSize: 1,
	}
	query := make(url.Values)
	query.Set("r", fmt.Sprintf("%d", timestamp()))
	var resp Response
	err := wc.postJson("https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxverifyuser", query, req, &resp)
	if err != nil {
		return err
	}
	if resp.BaseResponse.Ret != 0 {
		return fmt.Errorf("(%d)%s", resp.BaseResponse.Ret, resp.BaseResponse.ErrMsg)
	}
	return nil
}

func (wc *WeChat) Check() (string, error) {
	params := make(url.Values)
	params.Set("callback", "_")
	params.Set("sid", wc.baseRequest.Sid)
	params.Set("uin", fmt.Sprintf("%d", wc.baseRequest.Uin))
	params.Set("deviceid", wc.baseRequest.DeviceID)
	params.Set("synckey", makeSyncQuery(wc.syncKey.List))
	params.Set("_", fmt.Sprintf("%d", timestamp()))
	resp, err := resp(wc.request("GET", "https://webpush.weixin.qq.com/cgi-bin/mmwebwx-bin/synccheck", params, nil))
	if err != nil {
		return "", err
	}

	s := wc.grabSelector.FindAllStringSubmatch(string(resp), -1)
	if len(s) == 0 || len(s[0]) == 0 {
		return "", fmt.Errorf("sync error, no retcode: %s", string(resp))
	}

	return s[0][1], nil
}

func (wc *WeChat) GetLast() (*Response, error) {
	req := Request{
		BaseRequest: wc.baseRequest,
		SyncKey:     &wc.syncKey,
	}
	query := make(url.Values)
	query.Set("sid", wc.baseRequest.Sid)
	query.Set("r", fmt.Sprintf("%d", timestamp()))
	var resp Response
	err := wc.postJson("https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxsync", query, req, &resp)
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
	query := make(url.Values)
	query.Set("type", "ex")
	var resp Response
	err := wc.postJson("https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxbatchgetcontact", query, req, &resp)
	if err != nil {
		return nil, err
	}
	return resp.ContactList, nil
}

func (wc *WeChat) GetChatroomHeader(username, chatroomId string) (*http.Response, error) {
	query := make(url.Values)
	query.Set("username", username)
	query.Set("chatroomid", chatroomId)
	resp, err := wc.request("GET", "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxgeticon", query, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("%s", resp.Status)
	}
	return resp, nil
}

func (wc *WeChat) postJson(urlStr string, query url.Values, data interface{}, reply interface{}) error {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(data)
	if err != nil {
		return err
	}
	// fmt.Println("post json req:", buf.String())
	resp, err := wc.request("POST", urlStr, query, buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", resp.Status)
	}

	// b, err := ioutil.ReadAll(resp.Body)
	// fmt.Println("post json reply:", string(b))
	// return json.Unmarshal(b, &reply)

	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(reply)
}

func (wc *WeChat) request(method, urlStr string, query url.Values, body io.Reader) (*http.Response, error) {
	if query != nil {
		urlStr = urlStr + "?" + query.Encode()
	}
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.116 Safari/537.36")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Origin", "https://wx.qq.com")
	req.Header.Set("Referer", "https://wx.qq.com/")
	return wc.client.Do(req)
}

func (wc *WeChat) getDeviceId() string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	d := r.Intn(1e16)
	return fmt.Sprintf("e%d", d)
}

func timestamp() int64 {
	now := time.Now().UTC()
	return now.UnixNano() / int64(time.Millisecond)
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
	post := fmt.Sprintf("http://%s:%d/v3/poster/message/email/srv-op@exfe.com", config.ExfeService.Addr, config.ExfeService.Port)
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
