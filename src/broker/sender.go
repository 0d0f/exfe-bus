package broker

import (
	"gobus"
	"model"
)

type Sender struct {
	dispatcher *gobus.Dispatcher
}

func NewSender(config *model.Config) (*Sender, error) {
	table := gobus.NewTable(config.Dispatcher)
	return &Sender{
		dispatcher: gobus.NewDispatcher(table),
	}, nil
}

func (s Sender) Send(to model.Recipient, private, public string, info *model.InfoData) (string, error) {
	arg := model.ThirdpartSend{
		PrivateMessage: private,
		PublicMessage:  public,
		Info:           info,
	}
	arg.To = to

	var ids string
	err := s.dispatcher.DoWithIdentity(to.Provider, "bus://Thirdpart", "Send", &arg, &ids)

	if err != nil {
		return "", err
	}
	return ids, nil
}
