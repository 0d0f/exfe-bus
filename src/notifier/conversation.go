package notifier

import (
	"broker"
	"fmt"
	"formatter"
	"model"
)

type Conversation struct {
	localTemplate *formatter.LocalTemplate
	config        *model.Config
	sender        *broker.Sender
}

func NewConversation(localTemplate *formatter.LocalTemplate, config *model.Config, sender *broker.Sender) *Conversation {
	return &Conversation{
		localTemplate: localTemplate,
		config:        config,
		sender:        sender,
	}
}

func (c *Conversation) Update(updates model.ConversationUpdates) error {
	to := updates[0].To
	if to.Provider == "twitter" {
		c.config.Log.Debug("not send to twitter: %s", to)
		return nil
	}

	private, err := c.getConversationContent(updates)
	if err != nil {
		return err
	}

	_, err = c.sender.Send(to, private, "", &model.InfoData{
		CrossID: updates[0].Cross.ID,
		Type:    model.TypeConversation,
	})

	if err != nil {
		return fmt.Errorf("send error: %s", err)
	}
	return nil
}

func (c *Conversation) getConversationContent(updates []model.ConversationUpdate) (string, error) {
	arg, err := ArgFromUpdates(updates, c.config)
	if err != nil {
		return "", err
	}

	content, err := GetContent(c.localTemplate, "conversation", arg.To, arg)
	if err != nil {
		return "", err
	}

	return content, nil
}

type UpdateArg struct {
	model.ThirdpartTo
	Cross model.Cross
	Posts []*model.Post
}

func (a UpdateArg) Link() string {
	return fmt.Sprintf("%s/#!token=%s", a.Config.SiteUrl, a.To.Token)
}

func ArgFromUpdates(updates []model.ConversationUpdate, config *model.Config) (*UpdateArg, error) {
	if updates == nil && len(updates) == 0 {
		return nil, fmt.Errorf("no update info")
	}

	to := updates[0].To
	cross := updates[0].Cross
	posts := make([]*model.Post, len(updates))
	selfUpdates := true

	for i, update := range updates {
		if !to.Equal(&update.To) {
			return nil, fmt.Errorf("updates not send to same recipient: %s, %s", to, update.To)
		}
		if !cross.Equal(&update.Cross) {
			return nil, fmt.Errorf("updates not send to same exfee: %d, %d", cross.ID, update.Cross.ID)
		}
		if !to.SameUser(&update.Post.By) {
			selfUpdates = false
		}
		posts[i] = &updates[i].Post
	}

	if selfUpdates {
		return nil, fmt.Errorf("not send with all self updates")
	}

	ret := &UpdateArg{
		Cross: cross,
		Posts: posts,
	}
	ret.To = to
	err := ret.Parse(config)
	if err != nil {
		return nil, nil
	}

	return ret, nil
}

func (a UpdateArg) Timezone() string {
	if a.To.Timezone != "" {
		return a.To.Timezone
	}
	return a.Cross.Time.BeginAt.Timezone
}
