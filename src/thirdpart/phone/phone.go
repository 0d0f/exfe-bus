package phone

import (
	"fmt"
	"model"
	"strings"
)

type Phone struct {
	senders map[string]Sender
	config  *model.Config
}

func New(config *model.Config) (*Phone, error) {
	ret := &Phone{
		senders: make(map[string]Sender),
		config:  config,
	}

	senders := [...]Sender{NewTwilio(config), NewDuanCaiWang(config)}
	for _, sender := range senders {
		for _, code := range sender.Codes() {
			ret.senders[code] = sender
		}
	}

	return ret, nil
}

func (s *Phone) Provider() string {
	return "phone"
}

func (s *Phone) Post(from, id, text string) (string, error) {
	text = strings.Trim(text, " \r\n")

	var sender Sender
	for i := 3; i > 0; i-- {
		if len(id) < i {
			continue
		}
		code := id[:i]
		var ok bool
		sender, ok = s.senders[code]
		if ok {
			break
		}
	}

	if sender == nil {
		return "", fmt.Errorf("invalid recipient %s", id)
	}
	return sender.Send(id, text)
}
