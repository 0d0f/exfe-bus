package notifier

import (
	"fmt"
	"formatter"
	"gobus"
	"model"
)

var SendSelfError = fmt.Errorf("no need send self")

type Conversation struct {
	localTemplate *formatter.LocalTemplate
	config        *model.Config
}

func NewConversation(localTemplate *formatter.LocalTemplate, config *model.Config) *Conversation {
	return &Conversation{
		localTemplate: localTemplate,
		config:        config,
	}
}

func (c *Conversation) Update(updates model.ConversationUpdates) error {
	private, err := c.getConversationContent(updates)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s:%d", c.config.ExfeService.Addr, c.config.ExfeService.Port)
	client, err := gobus.NewClient(fmt.Sprintf("%s/%s", url, "Thirdpart"))
	if err != nil {
		return fmt.Errorf("can't create gobus client: %s", err)
	}

	arg := model.ThirdpartSend{
		PrivateMessage: private,
		PublicMessage:  "",
		Info: &model.InfoData{
			CrossID: updates[0].Cross.ID,
			Type:    model.TypeConversation,
		},
	}
	arg.To = updates[0].To
	var id string
	err = client.Do("Send", &arg, &id)
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

func ArgFromUpdates(updates []model.ConversationUpdate, config *model.Config) (*UpdateArg, error) {
	if updates == nil && len(updates) == 0 {
		return nil, fmt.Errorf("no update info")
	}

	to := updates[0].To
	cross := updates[0].Cross
	posts := make([]*model.Post, len(updates))

	needSend := false
	for i, update := range updates {
		if !to.SameUser(&update.Post.By) {
			needSend = true
		}
		if !to.Equal(&update.To) {
			return nil, fmt.Errorf("updates not send to same recipient: %s, %s", to, update.To)
		}
		if !cross.Equal(&update.Cross) {
			return nil, fmt.Errorf("updates not send to same exfee: %d, %d", cross.ID, update.Cross.ID)
		}
		posts[i] = &updates[i].Post
	}
	if !needSend {
		return nil, SendSelfError
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
