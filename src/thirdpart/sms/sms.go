package sms

import (
	"fmt"
	"model"
	"strings"
)

type Sms struct {
	senders map[string]Sender
	config  *model.Config
}

func New(config *model.Config) *Sms {
	ret := &Sms{
		senders: make(map[string]Sender),
		config:  config,
	}

	senders := [...]Sender{NewTwilio(config), NewDuanCaiWang(config)}
	for _, sender := range senders {
		for _, code := range sender.Codes() {
			ret.senders[code] = sender
		}
	}

	return ret
}

func (s *Sms) Provider() string {
	return "sms"
}

func (s *Sms) Send(to *model.Recipient, privateMessage string, publicMessage string, data *model.InfoData) (id string, err error) {
	phone := to.ExternalID
	var sender Sender
	for i := 3; i > 0; i-- {
		code := phone[0:i]
		var ok bool
		sender, ok = s.senders[code]
		if ok {
			break
		}
	}
	if sender == nil {
		return "", fmt.Errorf("can't send to %s, no support code", to)
	}
	lines := strings.Split(privateMessage, "\n")
	contents := make([]string, 0)
	for _, line := range lines {
		if line != "" {
			contents = append(contents, line)
		}
	}
	return sender.Send(phone, contents)
}
