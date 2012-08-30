package exfe_service

import (
	"apn/service"
	"c2dm/service"
	"gobus"
)

type CrossPush struct {
	CrossProviderBase
	android *gobus.Client
}

func NewCrossPush(config *Config) (ret *CrossPush) {
	ret = &CrossPush{
		CrossProviderBase: NewCrossProviderBase("push", config),
	}
	ret.client = gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, "iOSAPN")
	ret.android = gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, "Android")
	ret.handler = ret
	ret.log.SetPrefix("exfe.cross.push")
	return
}

func (s *CrossPush) Handle(arg *ProviderArg) {
	s.sendNewCross(arg)
	s.sendCrossChange(arg)
	s.sendExfeeChange(arg)
}

func (s *CrossPush) sendNewCross(arg *ProviderArg) {
	if arg.Old_cross != nil {
		return
	}

	str, _ := arg.TextPrivateInvitation()
	s.push(arg, str, "default", "i", 1)
}

func (s *CrossPush) sendCrossChange(arg *ProviderArg) {
	if arg.Old_cross == nil {
		return
	}

	if arg.Old_cross.Title != arg.Cross.Title {
		msg, _ := arg.TextTitleChange()
		s.push(arg, msg, "default", "u", 1)
	}
	msg, _ := arg.TextCrossChange()
	s.push(arg, msg, "default", "u", 1)
}

func (s *CrossPush) sendExfeeChange(arg *ProviderArg) {
	if arg.Old_cross == nil {
		return
	}
	accepted, declined, newlyInvited, removed := arg.Diff(s.log)

	if len(accepted) > 0 {
		msg, _ := arg.TextAccepted()
		s.push(arg, msg, "default", "u", 1)
	}
	if len(declined) > 0 {
		msg, _ := arg.TextDeclined()
		s.push(arg, msg, "default", "u", 1)
	}
	if len(newlyInvited) > 0 {
		if _, ok := newlyInvited[arg.To_identity.DiffId()]; ok {
			msg, _ := arg.TextPrivateInvitation()
			s.push(arg, msg, "default", "i", 1)
		} else {
			msg, _ := arg.TextNewlyInvited()
			s.push(arg, msg, "default", "u", 1)
		}
	}
	if len(removed) > 0 {
		if _, ok := removed[arg.To_identity.DiffId()]; ok {
			msg, _ := arg.TextQuit()
			s.push(arg, msg, "default", "r", 1)
		} else {
			msg, _ := arg.TextRemoved()
			s.push(arg, msg, "default", "u", 1)
		}
	}
}

func (s *CrossPush) push(arg *ProviderArg, message, sound, messageType string, badge uint) {
	switch arg.To_identity.Provider {
	case "iOS":
		arg := apn_service.ApnSendArg{
			DeviceToken: arg.To_identity.External_id,
			Alert:       message,
			Badge:       badge,
			Sound:       sound,
			Cid:         arg.Cross.Id,
			T:           messageType,
		}
		s.client.Send("ApnSend", &arg, 5)
	case "Android":
		arg := c2dm_service.C2DMSendArg{
			DeviceID: arg.To_identity.External_id,
			Message:  message,
			Cid:      arg.Cross.Id,
			T:        messageType,
			Badge:    badge,
			Sound:    sound,
		}
		s.android.Send("C2DMSend", &arg, 5)
	}
}
