package notifier

import (
	"bytes"
	"fmt"
	"formatter"
	"gobus"
	"model"
	"service/args"
	"thirdpart"
)

var SendSelfError = fmt.Errorf("no need send self")

type UpdateArg struct {
	To    model.Recipient
	Cross model.Cross
	Posts []*model.Post

	Config *model.Config
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
		To:     to,
		Cross:  cross,
		Posts:  posts,
		Config: config,
	}

	return ret, nil
}

func (a *UpdateArg) Link() string {
	return fmt.Sprintf("%s/#!token=%s", a.Config.SiteUrl, a.To.Token)
}

func (a *UpdateArg) Timezone() string {
	if a.To.Timezone != "" {
		return a.To.Timezone
	}
	return a.Cross.Time.BeginAt.Timezone
}

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

func (c *Conversation) Update(updates []model.ConversationUpdate) error {
	private, public, err := c.getContent(updates)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s:%d", c.config.ExfeService.Addr, c.config.ExfeService.Port)
	client, err := gobus.NewClient(fmt.Sprintf("%s/%s", url, "Thirdpart"))
	if err != nil {
		return fmt.Errorf("can't create gobus client: %s", err)
	}

	arg := args.SendArg{
		To:             &updates[0].To,
		PrivateMessage: private,
		PublicMessage:  public,
		Info: &thirdpart.InfoData{
			CrossID: updates[0].Cross.ID,
			Type:    thirdpart.Conversation,
		},
	}
	var id string
	err = client.Do("Send", &arg, &id)
	if err != nil {
		return fmt.Errorf("send error: %s", err)
	}
	return nil
}

func (c *Conversation) getContent(updates []model.ConversationUpdate) (string, string, error) {
	arg, err := ArgFromUpdates(updates, c.config)
	if err != nil {
		return "", "", err
	}

	messageType, err := thirdpart.MessageTypeFromProvider(arg.To.Provider)
	if err != nil {
		return "", "", err
	}

	templateName := fmt.Sprintf("conversation.%s", messageType)
	private := bytes.NewBuffer(nil)
	err = c.localTemplate.Execute(private, arg.To.Language, templateName, arg)
	if err != nil {
		return "", "", fmt.Errorf("private template(%s) failed: %s", templateName, err)
	}

	return private.String(), "", nil
}
