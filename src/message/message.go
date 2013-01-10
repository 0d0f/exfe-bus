package message

import (
	"github.com/googollee/go-logger"
	"model"
)

type Dispatcher interface {
	DoWithTicket(ticket, url, method string, arg, reply interface{}) error
}

type Message struct {
	dispatcher Dispatcher
	log        *logger.SubLogger
}

func New(config *model.Config, dispatcher Dispatcher) (*Message, error) {
	return &Message{
		dispatcher: dispatcher,
		log:        config.Log.SubPrefix("message"),
	}, nil
}

func (m *Message) Send(service, ticket string, recipients []model.Recipient, data interface{}) error {
	head10Recipients := make([]model.Recipient, 0)
	instantRecipients := make([]model.Recipient, 0)
	for _, to := range recipients {
		switch to.Provider {
		case "Android":
			fallthrough
		case "iOS":
			instantRecipients = append(instantRecipients, to)
		default:
			head10Recipients = append(head10Recipients, to)
		}
	}
	service, method := m.workaround(service)
	push := model.QueuePush{
		Service:  service,
		Method:   method,
		MergeKey: ticket,
		Data:     data,
	}
	if len(head10Recipients) > 0 {
		push.Tos = head10Recipients
		var i int
		err := m.dispatcher.DoWithTicket(ticket, "bus://exfe_queue/head10", "POST", push, &i)
		if err != nil {
			m.log.Err("can't send to head10: %s", err)
		}
	}
	if len(instantRecipients) > 0 {
		push.Tos = instantRecipients
		var i int
		err := m.dispatcher.DoWithTicket(ticket, "bus://exfe_queue/instant", "POST", push, &i)
		if err != nil {
			m.log.Err("can't send to instant: %s", err)
		}
	}
	return nil
}

func (m *Message) workaround(url string) (string, string) {
	switch url {
	case "bus://exfe_service/notifier/conversation":
		return "Conversation", "Update"
	case "bus://exfe_service/notifier/cross/invitation":
		return "Cross", "Invite"
	case "bus://exfe_service/notifier/cross/summary":
		return "Cross", "Summary"
	case "bus://exfe_service/notifier/user/welcome":
		return "User", "Welcome"
	case "bus://exfe_service/notifier/user/verify":
		return "User", "Verify"
	case "bus://exfe_service/notifier/user/password":
		return "User", "ResetPassword"
	}
	return "", ""
}
