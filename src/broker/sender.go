package broker

import (
	"gobus"
	"model"
)

type Sender struct {
	dispatcher *gobus.Dispatcher
}

func NewSender(config *model.Config, dispatcher *gobus.Dispatcher) (*Sender, error) {
	return &Sender{
		dispatcher: dispatcher,
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
	err := s.dispatcher.DoWithIdentity(to.Provider, "bus://exfe_service/thirdpart/message", "POST", &arg, &ids)

	if err != nil {
		return "", err
	}
	return ids, nil
}
