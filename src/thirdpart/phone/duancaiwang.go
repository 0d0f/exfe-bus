package phone

import (
	"broker"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
	"model"
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
	ErrCode string  `json:"errcode"`
	ID      int     `json:"msg_id"`
}

func (t *DuanCaiWang) Send(phone string, content string) (string, error) {
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
		t.log.Err("send to %s reply decode failed: %s", phone, err)
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
