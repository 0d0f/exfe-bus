package _performance

import (
	"fmt"
	"model"
	"strconv"
	"strings"
	"thirdpart"
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

func (p *Performance) Provider() string {
	return "_performance"
}

func (p *Performance) SetPosterCallback(callback thirdpart.Callback) (time.Duration, bool) {
	return 10 * time.Second, false
}

func (p *Performance) Post(from, id, text string) (string, error) {
	lines := strings.Split(text, "\n")
	if len(lines) > 1 {
		key, value := lines[0], lines[1]
		switch key {
		case "set delay":
			i, err := strconv.Atoi(value)
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
	}
	fmt.Println("send to", id, "message:", text)
	ret := p.sendCount
	p.sendCount++
	time.Sleep(p.delay)
	return fmt.Sprintf("%d", ret), nil
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
