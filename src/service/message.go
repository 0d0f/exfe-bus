package main

import (
	"gobus"
	"message"
	"model"
)

type Message struct {
	message *message.Message
}

func NewMessage(config *model.Config, dispatcher *gobus.Dispatcher, platform message.Platform) (*Message, error) {
	msg, err := message.New(config, dispatcher, platform)
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

// 通过Message发送一条消息arg
//
// 例子：
//
// > curl 'http://127.0.0.1:23333/message' -d '{"service":"bus://exfe_service/conversation",
//       "ticket":"email_cross123",
//       "recipients":[{"identity_id":12,"user_id":3,"name":"email1 name","auth_data":"","timezone":"+0800","token":"recipient_email1_token","language":"en_US","provider":"email","external_id":"sender1@gmail.com","external_username":"sender1@gmail.com"},{"identity_id":12,"user_id":3,"name":"email1 name","auth_data":"","timezone":"+0800","token":"recipient_email1_token","language":"en_US","provider":"email","external_id":"sender2@hotmail.com","external_username":"sender2@hotmail.com"}],
//       "data":{"cross":{"id":123,"by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"sender2@hotmail.com","external_username":"sender2@hotmail.com"},"title":"Test Cross","description":"test cross description","time":{"begin_at":{"date_word":"","date":"","time_word":"","time":"","timezone":""},"origin":"","output_format":0},"place":{"id":0,"title":"","description":"","lng":"","lat":"","provider":"","external_id":""},"exfee":{"id":0,"name":"","invitations":null}},"post":{"id":1,"by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"sender2@hotmail.com","external_username":"sender2@hotmail.com"},"content":"email1 post sth","via":"abc","created_at":"2012-10-24 16:31:00"}}}'
//
func (m *Message) Send(params map[string]string, arg model.Message) (int, error) {
	return 0, m.message.Send(arg.Service, arg.Ticket, arg.Recipients, arg.Data)
}
