package twitter

import (
	"exfe/service"
	"fmt"
	"gobot"
	"gobus"
	"twitter/service"
)

type Bot struct {
	bus      *gobus.Client
	config   *exfe_service.Config
	helpText string
}

func NewBot(c *exfe_service.Config) *Bot {
	return &Bot{
		bus:      gobus.CreateClient(c.Redis.Netaddr, c.Redis.Db, c.Redis.Password, "twitter"),
		config:   c,
		helpText: fmt.Sprintf("WRONG SYNTAX. Please enclose the 2-character mark in your reply to indicate mentioning 'X', e.g.:\n@%s Sure, be there or be square! #Z4", c.Twitter.Screen_name),
	}
}

func (b *Bot) GenerateContext(id string) bot.Context {
	return NewContext(id, b)
}

func (b *Bot) GetIDFromInput(input interface{}) (id string, content interface{}, e error) {
	tweet, ok := input.(*Tweet)
	if !ok {
		return "", nil, fmt.Errorf("input's type is not *Tweet")
	}
	i := tweet.ToInput()
	id = i.ID
	content = i
	if id == "" {
		return "", nil, fmt.Errorf("input no id")
	}
	return
}

func (b *Bot) SendHelp(screenName string) error {
	f := &twitter_service.FriendshipsExistsArg{
		ClientToken:  b.config.Twitter.Client_token,
		ClientSecret: b.config.Twitter.Client_secret,
		AccessToken:  b.config.Twitter.Access_token,
		AccessSecret: b.config.Twitter.Access_secret,
		UserA:        screenName,
		UserB:        b.config.Twitter.Screen_name,
	}
	var isFriend bool
	err := b.bus.Do("GetFriendship", f, &isFriend, 10)
	if err != nil {
		return fmt.Errorf("Can't require user %s friendship: %s", screenName, err)
	}

	if isFriend {
		dm := &twitter_service.DirectMessagesNewArg{
			ClientToken:  b.config.Twitter.Client_token,
			ClientSecret: b.config.Twitter.Client_secret,
			AccessToken:  b.config.Twitter.Access_token,
			AccessSecret: b.config.Twitter.Access_secret,
			Message:      b.helpText,
			ToUserName:   &screenName,
		}
		b.bus.Send("SendDM", dm, 5)
	} else {
		tweet := &twitter_service.StatusesUpdateArg{
			ClientToken:  b.config.Twitter.Client_token,
			ClientSecret: b.config.Twitter.Client_secret,
			AccessToken:  b.config.Twitter.Access_token,
			AccessSecret: b.config.Twitter.Access_secret,
			Tweet:        fmt.Sprintf("@%s %s", screenName, b.helpText),
		}
		b.bus.Send("SendTweet", tweet, 5)
	}
	return nil
}
