package phone

import (
	"broker"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-encoding"
	"logger"
	"model"
	"net/url"
	"unicode/utf8"
)

type DuanCaiWang struct {
	url string
}

func NewDuanCaiWang(config *model.Config) *DuanCaiWang {
	return &DuanCaiWang{
		url: config.Thirdpart.Sms.DuanCaiWang.Url,
	}
}

func (t *DuanCaiWang) Codes() []string {
	return []string{"+86"}
}

type duancaiwangReply struct {
	Result  bool    `json:"result"`
	Msg     *string `json:"msg"`
	Active  *int    `json:"active"`
	ErrCode int     `json:"errcode"`
	ID      int     `json:"msg_id"`
}

func (t *DuanCaiWang) Send(phone string, content string) (string, error) {
	var err error
	if content, err = filter("gb2312", content, ""); err != nil {
		return "", err
	}
	phone = phone[3:]
	params := make(url.Values)
	params.Add("mobile", phone)
	params.Add("content", content)
	resp, err := broker.HttpForm(t.url, params)
	if err != nil {
		return "", fmt.Errorf("send to %s failed: %s", phone, err)
	}
	defer resp.Close()
	decoder := json.NewDecoder(resp)
	var reply duancaiwangReply
	err = decoder.Decode(&reply)
	if err != nil {
		logger.ERROR("send to %s reply decode failed: %s with %s", t.url, err, params.Encode())
	}
	if !reply.Result {
		if reply.Msg == nil {
			return "", fmt.Errorf("send to %s response: %+v", phone, reply)
		} else {
			return "", fmt.Errorf("send to %s response: %s(%+v)", phone, *reply.Msg, reply)
		}
	}
	if reply.Active != nil {
		return fmt.Sprintf("%d-%d", *reply.Active, reply.ID), nil
	}
	return fmt.Sprintf("%d", reply.ID), nil
}

func filter(codec, s, c string) (string, error) {
	iconv, err := encoding.NewIconv(codec, "utf-8")
	if err != nil {
		return "", err
	}
	from := []byte(s)
	ret := make([]byte, 0, len(from))
	b := make([]byte, len(from))
	for len(from) > 0 {
		inlen, _, err := iconv.Conv(from, b)
		ret = append(ret, from[:inlen]...)
		from = from[inlen:]
		if err == nil {
			continue
		}
		_, size := utf8.DecodeRune(from)
		from = from[size:]
		ret = append(ret, []byte(c)...)
	}
	return string(ret), nil
}
