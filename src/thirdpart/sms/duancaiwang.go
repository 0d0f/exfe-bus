package sms

import (
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
	"model"
	"net/http"
	"net/url"
)

type DuanCaiWang struct {
	url string
	log *logger.SubLogger
}

func NewDuanCaiWang(config *model.Config) *DuanCaiWang {
	return &DuanCaiWang{
		url: config.Thirdpart.Sms.DuanCaiWang.Url,
		log: config.Log.SubPrefix("sms-DuanCaiWang"),
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

func (t *DuanCaiWang) Send(phone string, contents []string) (string, error) {
	params := make(url.Values)
	params.Add("mobile", phone)
	ret := ""
	for _, content := range contents {
		params.Add("content", content)
		resp, err := http.PostForm(t.url, params)
		if err != nil {
			return "", fmt.Errorf("send to %s failed: %s", phone, err)
		}
		decoder := json.NewDecoder(resp.Body)
		var reply duancaiwangReply
		err = decoder.Decode(&reply)
		if err != nil {
			t.log.Err("send to %s reply decode failed: %s", phone, err)
		}
		if resp.StatusCode != 200 || !reply.Result {
			if reply.Msg == nil {
				return "", fmt.Errorf("send to %s response: %s(%+v)", phone, resp.Status, reply)
			} else {
				return "", fmt.Errorf("send to %s response: %s(%+v)", phone, *reply.Msg, reply)
			}
		}
		if reply.Active != nil {
			ret += fmt.Sprintf(",%d-%d", reply.Active, reply.ID)
		} else {
			ret += fmt.Sprintf(",%d", reply.ID)
		}
	}
	if len(ret) > 0 {
		ret = ret[1:]
	}
	return ret, nil
}
