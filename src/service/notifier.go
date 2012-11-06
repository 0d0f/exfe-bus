package main

import (
	"formatter"
	"gobus"
	"model"
	"notifier"
)

type Conversation struct {
	conversation *notifier.Conversation
}

func NewConversation(localTemplate *formatter.LocalTemplate, config *model.Config) *Conversation {
	return &Conversation{
		conversation: notifier.NewConversation(localTemplate, config),
	}
}

// 发送Conversation的更新消息updates
//
// 例子：
//
// > curl 'http://127.0.0.1:23333/Conversation?method=Update' -d '[{"to":{"identity_id":11,"user_id":1,"name":"email1 name","auth_data":"","timezone":"+0800","token":"recipient_email1_token","language":"en_US","provider":"email","external_id":"sender1@gmail.com","external_username":"sender1@gmail.com"},"cross":{"id":123,"by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"sender1@gmail.com","external_username":"sender1@gmail.com"},"title":"Test Cross","description":"test cross description","time":{"begin_at":{"date_word":"","date":"","time_word":"","time":"","timezone":""},"origin":"","output_format":0},"place":{"id":0,"title":"","description":"","lng":"","lat":"","provider":"","external_id":""},"exfee":{"id":0,"name":"","invitations":null}},"post":{"id":1,"by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"sender1@gmail.com","external_username":"sender1@gmail.com"},"content":"email1 post sth","via":"abc","created_at":"2012-10-24 16:31:00"}},{"recipient":{"identity_id":11,"user_id":1,"name":"email1 name","auth_data":"","timezone":"+0800","token":"recipient_email1_token","language":"en_US","provider":"email","external_id":"sender1@gmail.com","external_username":"sender1@gmail.com"},"cross":{"id":123,"by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"sender1@gmail.com","external_username":"sender1@gmail.com"},"title":"Test Cross","description":"test cross description","time":{"begin_at":{"date_word":"","date":"","time_word":"","time":"","timezone":""},"origin":"","output_format":0},"place":{"id":0,"title":"","description":"","lng":"","lat":"","provider":"","external_id":""},"exfee":{"id":0,"name":"","invitations":null}},"post":{"id":2,"by_identity":{"id":22,"name":"twitter3 name","nickname":"twitter3 nick","bio":"twitter3 bio","timezone":"+0800","connected_user_id":3,"avatar_filename":"http://path/to/twitter3.avatar","provider":"twitter","external_id":"twitter3@domain.com","external_username":"twitter3@domain.com"},"content":"twitter3 post sth","via":"abc","created_at":"2012-10-24 16:40:00"}}]'
//
func (c *Conversation) Update(meta *gobus.HTTPMeta, updates model.ConversationUpdates, i *int) error {
	*i = 0
	err := c.conversation.Update(updates)
	if err == notifier.SendSelfError {
		meta.Log.Info("send to %s: %s", updates[0].To, err)
		return nil
	}
	return err
}

type Cross struct {
	cross *notifier.Cross
}

func NewCross(localTemplate *formatter.LocalTemplate, config *model.Config) *Cross {
	return &Cross{
		cross: notifier.NewCross(localTemplate, config),
	}
}

// 发送Cross的邀请消息invitations
//
func (c Cross) Invite(meta *gobus.HTTPMeta, invitations model.CrossInvitations, i *int) error {
	for *i = range invitations {
		err := c.cross.Invite(invitations[*i])
		if err != nil {
			return err
		}
	}
	return nil
}

// 发送Cross的更新汇总消息updates
//
func (c Cross) Summary(meta *gobus.HTTPMeta, updates model.CrossUpdates, i *int) error {
	*i = 0
	err := c.cross.Summary(updates)
	return err
}

type User struct {
	user *notifier.User
}

func NewUser(localTemplate *formatter.LocalTemplate, config *model.Config) *User {
	return &User{
		user: notifier.NewUser(localTemplate, config),
	}
}

// 发送给用户的邀请
//
func (u User) Welcome(meta *gobus.HTTPMeta, welcomes model.UserWelcomes, i *int) error {
	for *i = range welcomes {
		err := u.user.Welcome(welcomes[*i])
		if err != nil {
			return err
		}
	}
	return nil
}

// 发送给用户的确认请求
//
func (u User) Confirm(meta *gobus.HTTPMeta, confirmations model.UserConfirms, i *int) error {
	for *i = range confirmations {
		err := u.user.Confirm(confirmations[*i])
		if err != nil {
			return err
		}
	}
	return nil
}

// 发送给用户的重置密码请求
//
func (u User) ResetPassword(meta *gobus.HTTPMeta, tos model.ThirdpartTos, i *int) error {
	for *i = range tos {
		err := u.user.ResetPassword(tos[*i])
		if err != nil {
			return err
		}
	}
	return nil
}
