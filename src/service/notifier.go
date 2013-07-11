package main

import (
	"broker"
	"formatter"
	"github.com/googollee/go-rest"
	"model"
	"net/http"
	"notifier"
)

type V3Notifier struct {
	rest.Service      `prefix:"/v3/notifier"`
	CrossDigest       rest.Processor `path:"/cross/digest" method:"POST"`
	CrossRemind       rest.Processor `path:"/cross/remind" method:"POST"`
	CrossInvitation   rest.Processor `path:"/cross/invitation" method:"POST"`
	CrossPreview      rest.Processor `path:"/cross/preview" method:"POST"`
	CrossUpdate       rest.Processor `path:"/cross/update" method:"POST"`
	CrossConversation rest.Processor `path:"/cross/conversation" method:"POST"`
	UserWelcome       rest.Processor `path:"/user/welcome" method:"POST"`
	UserVerify        rest.Processor `path:"/user/verify" method:"POST"`
	UserReset         rest.Processor `path:"/user/reset" method:"POST"`
	WechatRoutex      rest.Processor `path:"/wechat/routex" method:"POST"`

	cross  *notifier.Cross
	user   *notifier.User
	wechat *notifier.Wechat
}

func NewV3Notifier(local *formatter.LocalTemplate, config *model.Config, platform *broker.Platform) (*V3Notifier, error) {
	return &V3Notifier{
		cross:  notifier.NewCross(local, config, platform),
		user:   notifier.NewUser(local, config, platform),
		wechat: notifier.NewWechat(local, config, platform),
	}, nil
}

func (n V3Notifier) HandleCrossDigest(requests []model.CrossDigestRequest) {
	if len(requests) == 0 {
		n.Error(http.StatusBadRequest, n.DetailError(1, "no request"))
		return
	}
	err := n.cross.V3Digest(requests)
	if err != nil {
		n.Error(http.StatusInternalServerError, n.DetailError(2, "%s", err))
		return
	}
}

func (n V3Notifier) HandleCrossRemind(requests []model.CrossDigestRequest) {
	if len(requests) == 0 {
		n.Error(http.StatusBadRequest, n.DetailError(1, "no request"))
		return
	}
	err := n.cross.V3Remind(requests)
	if err != nil {
		n.Error(http.StatusInternalServerError, n.DetailError(2, "%s", err))
		return
	}
}

func (n V3Notifier) HandleCrossInvitation(invitation model.CrossInvitation) {
	err := n.cross.V3Invitation(invitation)
	if err != nil {
		n.Error(http.StatusInternalServerError, n.DetailError(3, "%s", err))
		return
	}
}

func (n V3Notifier) HandleCrossPreview(invitation model.CrossInvitation) {
	err := n.cross.V3Preview(invitation)
	if err != nil {
		n.Error(http.StatusInternalServerError, n.DetailError(3, "%s", err))
		return
	}
}

func (n V3Notifier) HandleCrossUpdate(updates []model.CrossUpdate) {
	err := n.cross.V3Update(updates)
	if err != nil {
		n.Error(http.StatusInternalServerError, n.DetailError(7, "%s", err))
		return
	}
}

func (n V3Notifier) HandleCrossConversation(updates []model.ConversationUpdate) {
	err := n.cross.V3Conversation(updates)
	if err != nil {
		n.Error(http.StatusInternalServerError, n.DetailError(8, "%s", err))
		return
	}
}

func (n V3Notifier) HandleUserWelcome(arg model.UserWelcome) {
	err := n.user.V3Welcome(arg)
	if err != nil {
		n.Error(http.StatusInternalServerError, n.DetailError(4, "%s", err))
		return
	}
}

func (n V3Notifier) HandleUserVerify(arg model.UserVerify) {
	err := n.user.V3Verify(arg)
	if err != nil {
		n.Error(http.StatusInternalServerError, n.DetailError(5, "%s", err))
		return
	}
}

func (n V3Notifier) HandleUserReset(arg model.UserVerify) {
	err := n.user.V3ResetPassword(arg)
	if err != nil {
		n.Error(http.StatusInternalServerError, n.DetailError(6, "%s", err))
		return
	}
}

func (n V3Notifier) HandleWechatRoutex(to model.Recipient) {
	err := n.wechat.RoutexNotice(to)
	if err != nil {
		n.Error(http.StatusInternalServerError, n.DetailError(6, "%s", err))
		return
	}
}
