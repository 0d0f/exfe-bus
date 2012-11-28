package broker

import (
	"fmt"
	"gobus"
	"model"
)

type Sender struct {
	dispatcher    map[string]*gobus.Client
	defaultClient *gobus.Client
}

func NewSender(config *model.Config) (*Sender, error) {
	ret := &Sender{
		dispatcher: make(map[string]*gobus.Client),
	}
	for k, v := range config.Dispatcher.Sender {
		var err error
		url := fmt.Sprintf("http://%s/Thirdpart", v)
		if k == "_default" {
			ret.defaultClient, err = gobus.NewClient(url)
		} else {
			ret.dispatcher[k], err = gobus.NewClient(url)
		}
		if err != nil {
			return nil, err
		}
	}
	if ret.defaultClient == nil {
		url := fmt.Sprintf("http://%s:%d/Thirdpart", config.ExfeService.Addr, config.ExfeService.Port)
		var err error
		ret.defaultClient, err = gobus.NewClient(url)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (s Sender) Send(to model.Recipient, private, public string, info *model.InfoData) (string, error) {
	client, ok := s.dispatcher[to.Provider]
	if !ok {
		client = s.defaultClient
	}

	arg := model.ThirdpartSend{
		PrivateMessage: private,
		PublicMessage:  public,
		Info:           info,
	}
	arg.To = to

	var ids string
	err := client.Do("Send", &arg, &ids)

	if err != nil {
		return "", err
	}
	return ids, nil
}
