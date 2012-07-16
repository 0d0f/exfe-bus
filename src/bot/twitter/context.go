package twitter

import (
	"fmt"
	"gobot"
	"io/ioutil"
	"net/http"
	"time"
)

type Context struct {
	*bot.BaseContext
	b       *Bot
	lastIom string
}

func NewContext(id string, b *Bot) *Context {
	return &Context{bot.NewBaseContext(id), b, ""}
}

func (c *Context) Idle(input *Input) error {
	fmt.Println("idle")
	if c.DurationFromLast() > c.b.config.Bot.Iom_timeout*time.Second {
		c.lastIom = ""
	}
	return nil
}

func (c *Context) TweetWithIOM(input *Input) error {
	if input.Iom == "" {
		if c.lastIom == "" {
			return bot.BotNotMatched
		}
		input.Iom = c.lastIom
	}
	fmt.Println("iom: ", input.Iom)
	c.lastIom = input.Iom
	resp, err := http.PostForm(fmt.Sprintf("%s/v2/gobus/PostConversation", c.b.config.Site_api), input.ToUrl())
	if err != nil {
		return fmt.Errorf("message(%s) send to server error: %s", input.ScreenName, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("message(%s) get response body error: %s", input.ScreenName, err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("message(%s) send error(%s): %s", input.ScreenName, resp.Status, string(body))
	}
	return nil
}

func (c *Context) Default(input *Input) error {
	fmt.Println("default")
	return nil
	return c.b.SendHelp(input.ScreenName)
}
