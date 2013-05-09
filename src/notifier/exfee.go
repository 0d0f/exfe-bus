package notifier

import (
	"broker"
	"fmt"
	"formatter"
	"model"
)

type Exfee struct {
	localTemplate *formatter.LocalTemplate
	config        *model.Config
	platform      *broker.Platform
}

func NewExfee(localTemplate *formatter.LocalTemplate, config *model.Config, platform *broker.Platform) *Exfee {
	return &Exfee{
		localTemplate: localTemplate,
		config:        config,
		platform:      platform,
	}
}

func (e Exfee) V3Conversation(updates []model.ConversationUpdate) error {
	arg, err := ArgFromConversations(updates, e.config)
	if err != nil {
		return err
	}

	to := arg.To
	text, err := GenerateContent(e.localTemplate, "exfee_conversation", to.Provider, to.Language, arg)
	if err != nil {
		return err
	}
	_, err = e.platform.Send(to, text)
	if err != nil {
		return err
	}
	return nil
}

type ConversationArg struct {
	model.ThirdpartTo
	Posts []*model.Post
}

func ArgFromConversations(updates []model.ConversationUpdate, config *model.Config) (*ConversationArg, error) {
	if updates == nil && len(updates) == 0 {
		return nil, fmt.Errorf("no update info")
	}

	to := updates[0].To
	posts := make([]*model.Post, len(updates))

	for i, update := range updates {
		if !to.Equal(&update.To) {
			return nil, fmt.Errorf("updates not send to same recipient: %s, %s", to, update.To)
		}
		posts[i] = &updates[i].Post
	}

	ret := &ConversationArg{
		Posts: posts,
	}
	ret.To = to
	err := ret.Parse(config)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (a ConversationArg) Link() string {
	return fmt.Sprintf("%s/#!token=%s", a.Config.SiteUrl, a.To.Token)
}

func (a ConversationArg) Timezone() string {
	if a.To.Timezone != "" {
		return a.To.Timezone
	}
	return "+00:00"
}

func (a ConversationArg) Bys() []*model.Identity {
	var ret []*model.Identity
	for _, post := range a.Posts {
		fmt.Println(post)
		isSame := false
		for _, i := range ret {
			if i.SameUser(post.By) {
				isSame = true
				break
			}
		}
		if !isSame {
			ret = append(ret, &post.By)
		}
	}
	fmt.Println(ret)
	return ret
}
