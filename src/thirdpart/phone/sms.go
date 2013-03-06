package sms

import (
	"fmt"
	"formatter"
	"model"
	"strings"
	"thirdpart"
	"unicode/utf8"
)

type Sms struct {
	senders map[string]Sender
	config  *model.Config
	imsg    thirdpart.Sender
}

func New(config *model.Config, imsg thirdpart.Sender) (*Sms, error) {
	ret := &Sms{
		senders: make(map[string]Sender),
		config:  config,
		imsg:    imsg,
	}

	senders := [...]Sender{NewTwilio(config), NewDuanCaiWang(config)}
	for _, sender := range senders {
		for _, code := range sender.Codes() {
			ret.senders[code] = sender
		}
	}

	return ret, nil
}

func (s *Sms) Provider() string {
	return "phone"
}

func (s *Sms) Send(to *model.Recipient, text string) (id string, err error) {
	if s.imsg != nil {
		id, err = s.imsg.Send(to, text)
		if err == nil {
			return
		}
	}

	phone := to.ExternalID
	var sender Sender
	for i := 3; i > 0; i-- {
		if len(phone) < i {
			continue
		}
		code := phone[:i]
		var ok bool
		sender, ok = s.senders[code]
		if ok {
			break
		}
	}

	if sender == nil {
		return "", fmt.Errorf("invalid recipient %s", to)
	}
	lines := strings.Split(text, "\n")
	contents := make([]string, 0)
	for _, line := range lines {
		line = strings.Trim(line, " \n\r\t")
		if line == "" {
			continue
		}

		cutter, err := formatter.CutterParse(line, smsLen)
		if err != nil {
			return "", fmt.Errorf("parse cutter error: %s", err)
		}

		for _, content := range cutter.Limit(140) {
			contents = append(contents, content)
		}
	}
	return sender.Send(phone, contents)
}

func smsLen(content string) int {
	allAsc := true
	for _, r := range content {
		if r > 127 {
			allAsc = false
			break
		}
	}
	if allAsc {
		return len([]byte(content))
	}
	return utf8.RuneCountInString(content) * 2
}
