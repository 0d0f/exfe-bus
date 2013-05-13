package main

import (
	"broker"
	"fmt"
	"formatter"
	"github.com/googollee/go-rest"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"model"
	"net/http"
	"notifier"
	"os"
)

type V3Notifier struct {
	rest.Service      `prefix:"/v3/notifier"`
	CrossDigest       rest.Processor `path:"/cross/digest" method:"POST"`
	CrossRemind       rest.Processor `path:"/cross/remind" method:"POST"`
	CrossInvitation   rest.Processor `path:"/cross/invitation" method:"POST"`
	CrossSummary      rest.Processor `path:"/cross/summary" method:"POST"`
	CrossConversation rest.Processor `path:"/cross/conversation" method:"POST"`
	UserWelcome       rest.Processor `path:"/user/welcome" method:"POST"`
	UserVerify        rest.Processor `path:"/user/verify" method:"POST"`
	UserReset         rest.Processor `path:"/user/reset" method:"POST"`

	cross *notifier.Cross
	user  *notifier.User
}

func NewV3Notifier(local *formatter.LocalTemplate, config *model.Config, platform *broker.Platform) (*V3Notifier, error) {
	pin, err := os.Open(fmt.Sprintf("%s/image_data/map_pin_blue.png", config.TemplatePath))
	if err != nil {
		return nil, err
	}
	defer pin.Close()
	config.Pin, _, err = image.Decode(pin)
	if err != nil {
		return nil, err
	}
	ribbon, err := os.Open(fmt.Sprintf("%s/image_data/ribbon_280.png", config.TemplatePath))
	if err != nil {
		return nil, err
	}
	defer ribbon.Close()
	config.Ribbon, _, err = image.Decode(ribbon)
	if err != nil {
		return nil, err
	}

	return &V3Notifier{
		cross: notifier.NewCross(local, config, platform),
		user:  notifier.NewUser(local, config, platform),
	}, nil
}

func (n V3Notifier) HandleCrossDigest(requests []model.CrossDigestRequest) {
	if len(requests) == 0 {
		n.Error(http.StatusBadRequest, n.GetError(1, "no request"))
		return
	}
	err := n.cross.V3Digest(requests)
	if err != nil {
		n.Error(http.StatusInternalServerError, n.GetError(2, err.Error()))
		return
	}
}

func (n V3Notifier) HandleCrossRemind(requests []model.CrossDigestRequest) {
	if len(requests) == 0 {
		n.Error(http.StatusBadRequest, n.GetError(1, "no request"))
		return
	}
	err := n.cross.V3Remind(requests)
	if err != nil {
		n.Error(http.StatusInternalServerError, n.GetError(2, err.Error()))
		return
	}
}

func (n V3Notifier) HandleCrossInvitation(invitation model.CrossInvitation) {
	err := n.cross.V3Invitation(invitation)
	if err != nil {
		n.Error(http.StatusInternalServerError, n.GetError(3, err.Error()))
		return
	}
}

func (n V3Notifier) HandleCrossSummary(updates []model.CrossUpdate) {
	err := n.cross.V3Summary(updates)
	if err != nil {
		n.Error(http.StatusInternalServerError, n.GetError(7, err.Error()))
		return
	}
}

func (n V3Notifier) HandleCrossConversation(updates []model.ConversationUpdate) {
	err := n.cross.V3Conversation(updates)
	if err != nil {
		n.Error(http.StatusInternalServerError, n.GetError(8, err.Error()))
		return
	}
}

func (n V3Notifier) HandleUserWelcome(arg model.UserWelcome) {
	err := n.user.V3Welcome(arg)
	if err != nil {
		n.Error(http.StatusInternalServerError, n.GetError(4, err.Error()))
		return
	}
}

func (n V3Notifier) HandleUserVerify(arg model.UserVerify) {
	err := n.user.V3Verify(arg)
	if err != nil {
		n.Error(http.StatusInternalServerError, n.GetError(5, err.Error()))
		return
	}
}

func (n V3Notifier) HandleUserReset(arg model.UserVerify) {
	err := n.user.V3ResetPassword(arg)
	if err != nil {
		n.Error(http.StatusInternalServerError, n.GetError(6, err.Error()))
		return
	}
}
