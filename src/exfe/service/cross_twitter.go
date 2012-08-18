package exfe_service

import (
	"exfe/model"
	"fmt"
	"twitter/service"
)

type CrossTwitter struct {
	CrossProviderBase
}

func NewCrossTwitter(config *Config) (ret *CrossTwitter) {
	ret = &CrossTwitter{
		CrossProviderBase: NewCrossProviderBase("twitter", config),
	}
	ret.handler = ret
	return
}

func (s *CrossTwitter) Handle(arg *ProviderArg) {
	s.sendNewCross(arg)
	s.sendCrossChange(arg)
	s.sendExfeeChange(arg)
}

func (s *CrossTwitter) getIdentityInfo(id *exfe_model.Identity) {
	var twitterId *string
	if id.External_id != "" {
		twitterId = &id.External_id
	}
	// get to_identity info
	s.client.Send("GetInfo", &twitter_service.UsersShowArg{
		ClientToken:  s.config.Twitter.Client_token,
		ClientSecret: s.config.Twitter.Client_secret,
		AccessToken:  s.config.Twitter.Access_token,
		AccessSecret: s.config.Twitter.Access_secret,
		ScreenName:   &id.External_username,
		UserId:       twitterId,
		IdentityId:   &id.Id,
	}, 5)
}

func (s *CrossTwitter) checkFriend(to *exfe_model.Identity) (isFriend bool) {
	f := &twitter_service.FriendshipsExistsArg{
		ClientToken:  s.config.Twitter.Client_token,
		ClientSecret: s.config.Twitter.Client_secret,
		AccessToken:  s.config.Twitter.Access_token,
		AccessSecret: s.config.Twitter.Access_secret,
		UserA:        to.External_username,
		UserB:        s.config.Twitter.Screen_name,
	}
	err := s.client.Do("GetFriendship", f, &isFriend, 10)
	if err != nil {
		s.log.Printf("Can't require identity %d friendship: %s", to.Id, err)
		isFriend = false
	}
	return
}

func (s *CrossTwitter) sendTweet(arg *ProviderArg, message, url string) {
	urls := []string{}
	if url != "" {
		urls = append(urls, url)
	}
	tweet := &twitter_service.StatusesUpdateArg{
		ClientToken:  s.config.Twitter.Client_token,
		ClientSecret: s.config.Twitter.Client_secret,
		AccessToken:  s.config.Twitter.Access_token,
		AccessSecret: s.config.Twitter.Access_secret,
		Tweet:        fmt.Sprintf("@%s %s", arg.To_identity.External_username, message),
		Urls:         urls,
	}
	s.client.Send("SendTweet", tweet, 5)
}

func (s *CrossTwitter) sendDM(arg *ProviderArg, t, url string) {
	urls := []string{}
	if url != "" {
		urls = append(urls, url)
	}
	dm := &twitter_service.DirectMessagesNewArg{
		ClientToken:  s.config.Twitter.Client_token,
		ClientSecret: s.config.Twitter.Client_secret,
		AccessToken:  s.config.Twitter.Access_token,
		AccessSecret: s.config.Twitter.Access_secret,
		Message:      t,
		Urls:         urls,
		ToUserName:   &arg.To_identity.External_username,
		IdentityId:   &arg.To_identity.Id,
	}
	s.client.Send("SendDM", dm, 5)
}

func (s *CrossTwitter) sendNewCross(arg *ProviderArg) {
	if arg.Old_cross != nil {
		return
	}

	s.sendInvitation(arg)
}

func (s *CrossTwitter) sendInvitation(arg *ProviderArg) {
	s.getIdentityInfo(arg.To_identity)
	isFriend := s.checkFriend(arg.To_identity)

	if isFriend {
		msg, err := arg.TextPrivateInvitation()
		if err != nil {
			s.log.Printf("template error: %s", err)
		}
		s.sendDM(arg, msg, arg.Cross.LinkTo(s.config.Site_url, arg.Token()))
	} else {
		msg, err := arg.TextPublicInvitation()
		if err != nil {
			s.log.Printf("template error: %s", err)
		}
		s.sendTweet(arg, msg, fmt.Sprintf("%s/%s", arg.Cross.Link(s.config.Site_url), arg.Token()[1:4]))
	}
}

func (s *CrossTwitter) sendCrossChange(arg *ProviderArg) {
	if arg.Old_cross == nil {
		return
	}

	if arg.IsTitleChanged() {
		msg, err := arg.TextTitleChange()
		if err != nil {
			s.log.Printf("template error: %s", err)
		}
		s.send(arg, msg)
	}
	if arg.IsTimeChanged() || arg.IsPlaceChanged() {
		msg, err := arg.TextCrossChange()
		if err != nil {
			s.log.Printf("template error: %s", err)
		}
		s.send(arg, msg)
	}
}

func (s *CrossTwitter) sendExfeeChange(arg *ProviderArg) {
	if arg.Old_cross == nil {
		return
	}
	accepted, declined, newlyInvited, removed := arg.Diff(s.log)

	var msg string
	var err error
	needSend := false
	if len(accepted) > 0 {
		needSend = true
		msg, err = arg.TextAccepted()
		if err != nil {
			s.log.Printf("template error: %s", err)
		}
	}
	if len(declined) > 0 {
		needSend = true
		msg, err = arg.TextDeclined()
		if err != nil {
			s.log.Printf("template error: %s", err)
		}
	}
	if len(newlyInvited) > 0 {
		needSend = true
		if _, ok := newlyInvited[arg.To_identity.DiffId()]; ok {
			s.sendInvitation(arg)
			return
		} else {
			msg, _ = arg.TextNewlyInvited()
		}
	}
	if len(removed) > 0 {
		needSend = true
		if _, ok := removed[arg.To_identity.DiffId()]; ok {
			msg, err = arg.TextQuit()
			if err != nil {
				s.log.Printf("template error: %s", err)
			}
		} else {
			msg, err = arg.TextRemoved()
			if err != nil {
				s.log.Printf("template error: %s", err)
			}
		}
	}
	if needSend {
		s.send(arg, msg)
	}
}

func (s *CrossTwitter) send(arg *ProviderArg, msg string) {
	s.getIdentityInfo(arg.To_identity)
	isFriend := s.checkFriend(arg.To_identity)

	if isFriend {
		s.sendDM(arg, msg, arg.Cross.Link(s.config.Site_url))
	} else {
		s.sendTweet(arg, msg, arg.Cross.Link(s.config.Site_url))
	}
}
