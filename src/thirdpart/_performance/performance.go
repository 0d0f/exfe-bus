package _performance

import (
	"fmt"
	"model"
	"strconv"
	"time"
)

type Performance struct {
	delay         time.Duration
	sendCount     int
	friendCount   int
	identityCount int
}

func New() *Performance {
	return new(Performance)
}

func (p Performance) Provider() string {
	return "_performance"
}

func (p *Performance) Send(to *model.Recipient, privateMessage string, publicMessage string, data *model.InfoData) (string, error) {
	switch privateMessage {
	case "set delay":
		i, err := strconv.Atoi(publicMessage)
		if err != nil {
			return "", err
		}
		p.delay = time.Duration(i) * time.Second
		return "", nil
	case "reset":
		p.sendCount = 0
		p.friendCount = 0
		p.identityCount = 0
		return "", nil
	case "get":
		return fmt.Sprintf("send:%d friend:%d identity:%d", p.sendCount, p.friendCount, p.identityCount), nil
	}
	p.sendCount++
	time.Sleep(p.delay)
	return "", nil
}

func (p Performance) UpdateFriends(to *model.Recipient) error {
	p.friendCount++
	time.Sleep(p.delay)
	return nil
}

func (p Performance) UpdateIdentity(to *model.Recipient) error {
	p.identityCount++
	time.Sleep(p.delay)
	return nil
}
