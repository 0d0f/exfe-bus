package main

import (
	"bytes"
	"crypto/tls"
	"daemon"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-aws/s3"
	"io/ioutil"
	"logger"
	"model"
	"net/http"
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
	Count       int
	List        []ContactRequest
	SyncKey     SyncKey
	Msg         Message
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
}

func New() (*WeChat, error) {
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
	logger.NOTICE("login: https://login.weixin.qq.com/qrcode/%s?t=webwx", uuid)

	login := fmt.Sprintf("https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?uuid=%s&tip=1", uuid)
	var tokenUrl string
	for {
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
	return &WeChat{
		client:      client,
		baseRequest: baseRequest,
		syncKey:     ret.SyncKey,
		userName:    ret.User.UserName,
	}, nil
}

func (wc *WeChat) SendMessage(to, content string) error {
	req := Request{
		BaseRequest: wc.baseRequest,
		Msg: Message{
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
	return nil
}

func (wc *WeChat) GetLast() (*Response, error) {
	req := Request{
		BaseRequest: wc.baseRequest,
		SyncKey:     wc.syncKey,
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
		SyncKey:     wc.syncKey,
	}
	var resp Response
	err := wc.postJson("https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxbatchgetcontact?type=ex", req, &resp)
	if err != nil {
		return nil, err
	}
	return resp.ContactList, nil
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

func main() {
	var config model.Config
	_, quit := daemon.Init("exfe.json", &config)

	aws := s3.New(config.AWS.S3.Domain, config.AWS.S3.Key, config.AWS.S3.Secret)
	aws.SetACL(s3.ACLPublicRead)
	aws.SetLocationConstraint(s3.LC_AP_SINGAPORE)
	bucket, err := aws.GetBucket(fmt.Sprintf("%s-3rdpart-photos", config.AWS.S3.BucketPrefix))
	if err != nil {
		panic(err)
	}

	wc, err := New()
	if err != nil {
		panic(err)
	}
	defer func() {
		logger.NOTICE("quit")
	}()
	logger.NOTICE("login as %s", wc.userName)

	i := 0
	msgs := []string{"早", "hi", "喂", "what", "艹", "fuck", "呸", "lol"}
	last := time.Now()

	for {
		select {
		case <-quit:
			return
		default:
		}
		// err := wc.Ping()
		// if err != nil {
		// 	panic(err)
		// }
		fmt.Print(".")
		resp, err := wc.GetLast()
		if err != nil {
			panic(err)
		}
		if len(resp.AddMsgList) > 0 {
			var contactReq []ContactRequest
			for _, msg := range resp.AddMsgList {
				if msg.MsgType != 10000 {
					continue
				}
				contactReq = append(contactReq, ContactRequest{
					UserName: msg.FromUserName,
				})
			}
			if len(contactReq) > 0 {
				contacts, err := wc.GetContact(contactReq)
				if err != nil {
					panic(err)
				}
				for i, c := range contacts {
					if c.MemberCount == 0 {
						fmt.Printf("%d: person %+v\n", i, c)
						continue
					}
					var contactReq []ContactRequest
					for _, m := range c.MemberList {
						contactReq = append(contactReq, ContactRequest{
							UserName:   m.UserName,
							ChatRoomId: c.Uin,
						})
					}
					cs, err := wc.GetContact(contactReq)
					if err != nil {
						panic(err)
					}
					cross := model.Cross{}
					cross.Title = "Cross with "
					cross.Exfee.Invitations = make([]model.Invitation, len(cs))
					var host *model.Identity
					for i, member := range cs {
						if i < 3 {
							cross.Title += member.NickName + ", "
						}
						headerPath := fmt.Sprintf("/thirdpart/weichat/%d.jpg", member.Uin)
						resp, err := wc.get("https://wx.qq.com" + member.HeadImgUrl)
						if err != nil {
							panic(err)
						}
						obj, err := bucket.CreateObject(headerPath, resp.Header.Get("Content-Type"))
						if err != nil {
							panic(err)
						}
						obj.SetDate(time.Now())
						length, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
						if err != nil {
							panic(err)
						}
						err = obj.SaveReader(resp.Body, length)
						if err != nil {
							panic(err)
						}
						cross.Exfee.Invitations[i].Identity = model.Identity{
							ExternalID:       fmt.Sprintf("%d", member.Uin),
							Provider:         "wechat",
							ExternalUsername: member.UserName,
							Nickname:         member.NickName,
							Avatar:           obj.URL(),
						}
						if member.Uin == c.OwnerUin {
							cross.Exfee.Invitations[i].Host = true
							host = &cross.Exfee.Invitations[i].Identity
						}
					}
					cross.Title = cross.Title[:len(cross.Title)-2]
					cross.By = *host
					for i := range cross.Exfee.Invitations {
						cross.Exfee.Invitations[i].By = *host
						cross.Exfee.Invitations[i].UpdatedBy = *host
					}
					b, _ := json.Marshal(cross)
					fmt.Printf("%d: cross %s\n", i, string(b))
				}
			}
		}
		time.Sleep(time.Second * 30)
		if time.Now().Sub(last) > time.Minute*30 {
			last = time.Now()
			err = wc.SendMessage("leaskh", msgs[i])
			if err != nil {
				panic(err)
			}
			i++
			i = i % len(msgs)
		}
	}
}
