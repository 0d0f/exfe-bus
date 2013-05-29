package phone

import (
	"logger"
	"strings"
	"thirdpart/imessage"
)

type IMsgPhone struct {
	phone *Phone
	imsg  *imessage.IMessage
}

func NewIMsgPhone(phone *Phone, imsg *imessage.IMessage) *IMsgPhone {
	ret := &IMsgPhone{
		phone: phone,
		imsg:  imsg,
	}

	return ret
}

func (s *IMsgPhone) Provider() string {
	return "imessage,phone"
}

func (s *IMsgPhone) Post(id, text string) (string, error) {
	text = strings.Trim(text, " \r\n")
	if ret, err := s.imsg.Post(id, text); err == nil {
		logger.NOTICE("%s@phone sent from imessage", id)
		return ret, nil
	}
	return s.phone.Post(id, text)
}
