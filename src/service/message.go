package main

import (
	"gobus"
	"message"
	"model"
)

type Message struct {
	message *message.Message
}

func NewMessage(config *model.Config, dispatcher *gobus.Dispatcher) (*Message, error) {
	msg, err := message.New(config, dispatcher)
	if err != nil {
		return nil, err
	}
	return &Message{
		message: msg,
	}, nil
}

func (m *Message) SetRoute(r gobus.RouteCreater) error {
	json := new(gobus.JSON)
	return r().Methods("POST").Path("/message").HandlerMethod(json, m, "Send")
}

func (m *Message) Send(params map[string]string, arg model.Message) (int, error) {
	return 0, m.message.Send(arg.Service, arg.Ticket, arg.Recipients, arg.Data)
}
