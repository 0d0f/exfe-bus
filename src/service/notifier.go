package main

import (
	"broker"
	"formatter"
	"gobus"
	"model"
	"notifier"
)

type Notifier struct {
	conversation *notifier.Conversation
	cross        *notifier.Cross
	user         *notifier.User
}

func NewNotifier(local *formatter.LocalTemplate, config *model.Config, sender *broker.Sender) *Notifier {
	return &Notifier{
		conversation: notifier.NewConversation(local, config, sender),
		cross:        notifier.NewCross(local, config, sender),
		user:         notifier.NewUser(local, config, sender),
	}
}

func (n *Notifier) SetRoute(r gobus.RouteCreater) error {
	json := new(gobus.JSON)
	r().Methods("POST").Path("/notifier/conversation").HandlerMethod(json, n, "PostUpdate")
	r().Methods("POST").Path("/notifier/cross/invitation").HandlerMethod(json, n, "CrossInvite")
	r().Methods("POST").Path("/notifier/cross/summary").HandlerMethod(json, n, "CrossSummary")
	r().Methods("POST").Path("/notifier/user/welcome").HandlerMethod(json, n, "UserWelcome")
	r().Methods("POST").Path("/notifier/user/verify").HandlerMethod(json, n, "UserVerify")
	r().Methods("POST").Path("/notifier/user/password").HandlerMethod(json, n, "UserPassword")

	r().Queries("method", "Update").Path("/Conversation").HandlerMethod(json, n, "PostUpdate")
	r().Queries("method", "Invite").Path("/Cross").HandlerMethod(json, n, "CrossInvite")
	r().Queries("method", "Summary").Path("/Cross").HandlerMethod(json, n, "CrossSummary")
	r().Queries("method", "Welcome").Path("/User").HandlerMethod(json, n, "UserWelcome")
	r().Queries("method", "Verify").Path("/User").HandlerMethod(json, n, "UserVerify")
	r().Queries("method", "ResetPassword").Path("/User").HandlerMethod(json, n, "UserPassword")

	return nil
}

// 发送Conversation的更新消息updates
//
// 例子：
//
// > curl 'http://127.0.0.1:23333/conversation' -d '[{"to":{"identity_id":33,"user_id":3,"name":"email1 name","auth_data":"","timezone":"+0800","token":"recipient_email1_token","language":"en_US","provider":"email","external_id":"sender1@gmail.com","external_username":"sender1@gmail.com"},"cross":{"id":123,"by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"sender1@gmail.com","external_username":"sender1@gmail.com"},"title":"Test Cross","description":"test cross description","time":{"begin_at":{"date_word":"","date":"","time_word":"","time":"","timezone":""},"origin":"","output_format":0},"place":{"id":0,"title":"","description":"","lng":"","lat":"","provider":"","external_id":""},"exfee":{"id":0,"name":"","invitations":null}},"post":{"id":1,"by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"sender1@gmail.com","external_username":"sender1@gmail.com"},"content":"email1 post sth","via":"abc","created_at":"2012-10-24 16:31:00"}},{"to":{"identity_id":33,"user_id":3,"name":"email1 name","auth_data":"","timezone":"+0800","token":"recipient_email1_token","language":"en_US","provider":"email","external_id":"sender1@gmail.com","external_username":"sender1@gmail.com"},"cross":{"id":123,"by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"sender1@gmail.com","external_username":"sender1@gmail.com"},"title":"Test Cross","description":"test cross description","time":{"begin_at":{"date_word":"","date":"","time_word":"","time":"","timezone":""},"origin":"","output_format":0},"place":{"id":0,"title":"","description":"","lng":"","lat":"","provider":"","external_id":""},"exfee":{"id":0,"name":"","invitations":null}},"post":{"id":2,"by_identity":{"id":22,"name":"twitter3 name","nickname":"twitter3 nick","bio":"twitter3 bio","timezone":"+0800","connected_user_id":3,"avatar_filename":"http://path/to/twitter3.avatar","provider":"twitter","external_id":"twitter3@domain.com","external_username":"twitter3@domain.com"},"content":"twitter3 post sth","via":"abc","created_at":"2012-10-24 16:40:00"}}]'
//
func (n *Notifier) PostUpdate(params map[string]string, updates model.ConversationUpdates) (int, error) {
	err := n.conversation.Update(updates)
	return len(updates), err
}

// 发送Cross的邀请消息invitations
//
// 例子：
//
// > curl 'http://127.0.0.1:23333/cross/invitation' -d '[{"to":{"identity_id":11,"user_id":1,"name":"email1 name","auth_data":"","timezone":"+0800","token":"recipient_email1_token","language":"en_US","provider":"email","external_id":"email1@domain.com","external_username":"email1@domain.com"},"cross":{"id":123,"by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"email1@domain.com","external_username":"email1@domain.com"},"title":"Test Cross","description":"test cross description","time":{"begin_at":{"date_word":"","date":"2012-10-23","time_word":"","time":"08:45:00","timezone":"+0800"},"origin":"2012-10-23 16:45:00","output_format":0},"place":{"id":0,"title":"","description":"","lng":"","lat":"","provider":"","external_id":""},"exfee":{"id":123,"name":"","invitations":[{"id":11,"host":true,"mates":2,"identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"email1@domain.com","external_username":"email1@domain.com"},"rsvp_status":"NORESPONSE","by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"email1@domain.com","external_username":"email1@domain.com"},"via":""},{"id":22,"host":false,"mates":0,"identity":{"id":12,"name":"email2 name","nickname":"email2 nick","bio":"email2 bio","timezone":"+0800","connected_user_id":2,"avatar_filename":"http://path/to/email2.avatar","provider":"email","external_id":"email2@domain.com","external_username":"email2@domain.com"},"rsvp_status":"NORESPONSE","by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"email1@domain.com","external_username":"email1@domain.com"},"via":""},{"id":33,"host":false,"mates":0,"identity":{"id":22,"name":"twitter3 name","nickname":"twitter3 nick","bio":"twitter3 bio","timezone":"+0800","connected_user_id":3,"avatar_filename":"http://path/to/twitter3.avatar","provider":"twitter","external_id":"twitter3@domain.com","external_username":"twitter3@domain.com"},"rsvp_status":"NORESPONSE","by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"email1@domain.com","external_username":"email1@domain.com"},"via":""},{"id":44,"host":false,"mates":0,"identity":{"id":32,"name":"facebook4 name","nickname":"facebook4 nick","bio":"facebook4 bio","timezone":"+0800","connected_user_id":4,"avatar_filename":"http://path/to/facebook4.avatar","provider":"facebook","external_id":"facebook4@domain.com","external_username":"facebook4@domain.com"},"rsvp_status":"NORESPONSE","by_identity":{"id":22,"name":"twitter3 name","nickname":"twitter3 nick","bio":"twitter3 bio","timezone":"+0800","connected_user_id":3,"avatar_filename":"http://path/to/twitter3.avatar","provider":"twitter","external_id":"twitter3@domain.com","external_username":"twitter3@domain.com"},"via":""}]}}}]'
//
func (n *Notifier) CrossInvite(params map[string]string, invitation model.CrossInvitation) (int, error) {
	return 1, n.cross.Invite(invitation)
}

// 发送Cross的更新汇总消息updates
//
// 例子：
//
// > curl 'http://127.0.0.1:23333/cross/summary' -d '[{"to":{"identity_id":11,"user_id":1,"name":"email1 name","auth_data":"","timezone":"+0800","token":"recipient_email1_token","language":"en_US","provider":"email","external_id":"email1@domain.com","external_username":"email1@domain.com"},"old_cross":{"id":123,"by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"email1@domain.com","external_username":"email1@domain.com"},"title":"Test Cross","description":"test cross description","time":{"begin_at":{"date_word":"","date":"","time_word":"","time":"","timezone":"+0800"},"origin":"","output_format":0},"place":{"id":0,"title":"","description":"","lng":"","lat":"","provider":"","external_id":""},"exfee":{"id":123,"name":"","invitations":[{"id":11,"host":true,"mates":2,"identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"email1@domain.com","external_username":"email1@domain.com"},"rsvp_status":"NORESPONSE","by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"email1@domain.com","external_username":"email1@domain.com"},"via":""},{"id":22,"host":false,"mates":0,"identity":{"id":12,"name":"email2 name","nickname":"email2 nick","bio":"email2 bio","timezone":"+0800","connected_user_id":2,"avatar_filename":"http://path/to/email2.avatar","provider":"email","external_id":"email2@domain.com","external_username":"email2@domain.com"},"rsvp_status":"NORESPONSE","by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"email1@domain.com","external_username":"email1@domain.com"},"via":""},{"id":33,"host":false,"mates":0,"identity":{"id":22,"name":"twitter3 name","nickname":"twitter3 nick","bio":"twitter3 bio","timezone":"+0800","connected_user_id":3,"avatar_filename":"http://path/to/twitter3.avatar","provider":"twitter","external_id":"twitter3@domain.com","external_username":"twitter3@domain.com"},"rsvp_status":"NORESPONSE","by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"email1@domain.com","external_username":"email1@domain.com"},"via":""},{"id":44,"host":false,"mates":0,"identity":{"id":32,"name":"facebook4 name","nickname":"facebook4 nick","bio":"facebook4 bio","timezone":"+0800","connected_user_id":4,"avatar_filename":"http://path/to/facebook4.avatar","provider":"facebook","external_id":"facebook4@domain.com","external_username":"facebook4@domain.com"},"rsvp_status":"ACCEPTED","by_identity":{"id":22,"name":"twitter3 name","nickname":"twitter3 nick","bio":"twitter3 bio","timezone":"+0800","connected_user_id":3,"avatar_filename":"http://path/to/twitter3.avatar","provider":"twitter","external_id":"twitter3@domain.com","external_username":"twitter3@domain.com"},"via":""},{"id":77,"host":false,"mates":0,"identity":{"id":34,"name":"facebook6 name","nickname":"facebook6 nick","bio":"facebook6 bio","timezone":"+0800","connected_user_id":6,"avatar_filename":"http://path/to/facebook6.avatar","provider":"facebook","external_id":"facebook6@domain.com","external_username":"facebook6@domain.com"},"rsvp_status":"NORESPONSE","by_identity":{"id":32,"name":"facebook4 name","nickname":"facebook4 nick","bio":"facebook4 bio","timezone":"+0800","connected_user_id":4,"avatar_filename":"http://path/to/facebook4.avatar","provider":"facebook","external_id":"facebook4@domain.com","external_username":"facebook4@domain.com"},"via":""}]}},"cross":{"id":123,"by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"email1@domain.com","external_username":"email1@domain.com"},"title":"Test Cross","description":"test cross description","time":{"begin_at":{"date_word":"","date":"","time_word":"","time":"","timezone":"+0800"},"origin":"","output_format":0},"place":{"id":0,"title":"","description":"","lng":"","lat":"","provider":"","external_id":""},"exfee":{"id":123,"name":"","invitations":[{"id":22,"host":false,"mates":0,"identity":{"id":12,"name":"email2 name","nickname":"email2 nick","bio":"email2 bio","timezone":"+0800","connected_user_id":2,"avatar_filename":"http://path/to/email2.avatar","provider":"email","external_id":"email2@domain.com","external_username":"email2@domain.com"},"rsvp_status":"ACCEPTED","by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"email1@domain.com","external_username":"email1@domain.com"},"via":""},{"id":33,"host":false,"mates":0,"identity":{"id":22,"name":"twitter3 name","nickname":"twitter3 nick","bio":"twitter3 bio","timezone":"+0800","connected_user_id":3,"avatar_filename":"http://path/to/twitter3.avatar","provider":"twitter","external_id":"twitter3@domain.com","external_username":"twitter3@domain.com"},"rsvp_status":"DECLINED","by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"email1@domain.com","external_username":"email1@domain.com"},"via":""},{"id":44,"host":false,"mates":0,"identity":{"id":32,"name":"facebook4 name","nickname":"facebook4 nick","bio":"facebook4 bio","timezone":"+0800","connected_user_id":4,"avatar_filename":"http://path/to/facebook4.avatar","provider":"facebook","external_id":"facebook4@domain.com","external_username":"facebook4@domain.com"},"rsvp_status":"ACCEPTED","by_identity":{"id":22,"name":"twitter3 name","nickname":"twitter3 nick","bio":"twitter3 bio","timezone":"+0800","connected_user_id":3,"avatar_filename":"http://path/to/twitter3.avatar","provider":"twitter","external_id":"twitter3@domain.com","external_username":"twitter3@domain.com"},"via":""},{"id":55,"host":true,"mates":2,"identity":{"id":21,"name":"twitter1 name","nickname":"twitter1 nick","bio":"twitter1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/twitter1.avatar","provider":"twitter","external_id":"twitter1@domain.com","external_username":"twitter1@domain.com"},"rsvp_status":"NORESPONSE","by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"email1@domain.com","external_username":"email1@domain.com"},"via":""},{"id":66,"host":false,"mates":2,"identity":{"id":33,"name":"facebook5 name","nickname":"facebook5 nick","bio":"facebook5 bio","timezone":"+0800","connected_user_id":5,"avatar_filename":"http://path/to/facebook5.avatar","provider":"facebook","external_id":"facebook5@domain.com","external_username":"facebook5@domain.com"},"rsvp_status":"ACCEPTED","by_identity":{"id":32,"name":"facebook4 name","nickname":"facebook4 nick","bio":"facebook4 bio","timezone":"+0800","connected_user_id":4,"avatar_filename":"http://path/to/facebook4.avatar","provider":"facebook","external_id":"facebook4@domain.com","external_username":"facebook4@domain.com"},"via":""}]}},"by":{"id":32,"name":"facebook4 name","nickname":"facebook4 nick","bio":"facebook4 bio","timezone":"+0800","connected_user_id":4,"avatar_filename":"http://path/to/facebook4.avatar","provider":"facebook","external_id":"facebook4@domain.com","external_username":"facebook4@domain.com"}}]'
//
func (n *Notifier) CrossSummary(params map[string]string, updates model.CrossUpdates) (int, error) {
	return len(updates), n.cross.Summary(updates)
}

// 发送给用户的邀请
//
// 例子：
//
// > curl 'http://127.0.0.1:23333/user/welcome' -d '[{"to":{"identity_id":11,"user_id":1,"name":"email1 name","auth_data":"","timezone":"+0800","token":"recipient_email1_token","language":"en_US","provider":"email","external_id":"email1@domain.com","external_username":"email1@domain.com"},"need_verify":true}]'
//
func (n *Notifier) UserWelcome(params map[string]string, welcome model.UserWelcome) (int, error) {
	return 1, n.user.Welcome(welcome)
}

// 发送给用户的验证请求
//
// 例子：
//
// > curl 'http://127.0.0.1:23333/user/verify' -d '[{"to":{"identity_id":11,"user_id":1,"name":"email1 name","auth_data":"","timezone":"+0800","token":"recipient_email1_token","language":"en_US","provider":"email","external_id":"email1@domain.com","external_username":"email1@domain.com"},"by_name":"by user"}]'
//
func (n *Notifier) UserVerify(params map[string]string, confirmation model.UserVerify) (int, error) {
	return 1, n.user.Verify(confirmation)
}

// 发送给用户的重置密码请求
//
// 例子：
//
// > curl 'http://127.0.0.1:23333/User?method=ResetPassword' -d '[{"to":{"identity_id":11,"user_id":1,"name":"email1 name","auth_data":"","timezone":"+0800","token":"recipient_email1_token","language":"en_US","provider":"email","external_id":"email1@domain.com","external_username":"email1@domain.com"},"by_name":"by user"}]'
//
func (n *Notifier) UserPassword(params map[string]string, verify model.UserVerify) (int, error) {
	return 1, n.user.ResetPassword(verify)
}
