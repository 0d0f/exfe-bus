package thirdpart

import (
	"github.com/googollee/go-rest"
	"model"
	"net/http"
	"thirdpart/_performance"
	"thirdpart/apn"
	"thirdpart/email"
	"thirdpart/facebook"
	"thirdpart/gcm"
	"thirdpart/imsg"
	"thirdpart/phone"
	"thirdpart/twitter"
)

type Message struct {
	rest.Service `prefix:"/v3/message"`

	Send rest.Processor `path:"/:channel/:id" method:"POST"`

	senders map[string]Sender
}

func NewMessage(config *model.Config, platform *broker.Platform) (*Message, error) {
	twitterBroker := broker.NewTwitter(config.Thirdpart.Twitter.ClientToken, config.Thirdpart.Twitter.ClientSecret)
	apns_, err := apns.New(config.Thirdpart.Apn.Cert, config.Thirdpart.Apn.Key, config.Thirdpart.Apn.Server, time.Duration(config.Thirdpart.Apn.TimeoutInMinutes)*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("can't connect apn: %s", err)
	}
	gcms_ := gcms.New(config.Thirdpart.Gcm.Key)
	helper := NewHelper(config)

	ret := &Message{
		senders: make(map[string]Sender),
	}

	twitter_ := twitter.New(config, twitterBroker, helper)
	ret.AddSender(twitter_)

	facebook_ := facebook.New(helper)
	ret.AddSender(facebook_)

	email_ := email.New(helper)
	ret.AddSender(email_)

	apn_ := apn.New(apns_, getApnErrorHandler(config.Log.SubPrefix("apn error")))
	ret.AddSender(apn_)

	gcm_ := gcm.New(gcms_)
	ret.AddSender(gcm_)

	imsg_, err := imessage.New(config)
	if err != nil {
		return nil, fmt.Errorf("can't connect imessage: %s", err)
	}
	ret.AddSender(imsg_)

	sms_, err := sms.New(config, imsg_)
	if err != nil {
		return nil, fmt.Errorf("can't create sms: %s", err)
	}
	ret.AddSender(sms_)

	return ret, nil
}

func (m Message) HandleSend(text string) string {
	channel := m.Vars()["channel"]
	id := m.Vars()["id"]
	sender, ok := m.senders[channel]
	if !ok {
		m.Error(http.StatusBadRequest, m.GetError(1, "invalid channel: %s", channel))
		return
	}
	to := model.Recipient{
		ExternalID:       id,
		ExternalUsername: id,
		Provider:         channel,
	}
	ret, err := sender.Send(&to, text)
	if err != nil {
		m.Error(http.StatusInternalServerError, m.GetError(2, "%s", err))
	}
}

func (m *Message) AddSender(sender Sender) {
	t.senders[sender.Provider()] = sender
}
