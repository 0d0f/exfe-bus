package notifier

import (
	"bytes"
	"fmt"
	"formatter"
	"model"
	"thirdpart"
)

type ArgBase struct {
	To model.Recipient `json:"to"`

	Config       *model.Config         `json:"-"`
	templateType thirdpart.MessageType `json:"-"`
}

func (a *ArgBase) Parse(config *model.Config) (err error) {
	a.Config = config
	a.templateType, err = thirdpart.MessageTypeFromProvider(a.To.Provider)
	return
}

func (a ArgBase) ToIn(invitations []model.Invitation) bool {
	for _, i := range invitations {
		if a.To.SameUser(&i.Identity) {
			return true
		}
	}
	return false
}

func (a ArgBase) Link() string {
	return fmt.Sprintf("%s/#!token=%s", a.Config.SiteUrl, a.To.Token)
}

func (a ArgBase) ToRecipient() model.Recipient {
	return a.To
}

func (a ArgBase) Type() thirdpart.MessageType {
	return a.templateType
}

type ArgBaseInterface interface {
	ToRecipient() model.Recipient
	Type() thirdpart.MessageType
}

func GetContent(localTemplate *formatter.LocalTemplate, template string, arg ArgBaseInterface) (string, error) {
	templateName := fmt.Sprintf("%s.%s", template, arg.Type())

	ret := bytes.NewBuffer(nil)
	err := localTemplate.Execute(ret, arg.ToRecipient().Language, templateName, arg)
	if err != nil {
		return "", fmt.Errorf("template(%s) failed: %s", templateName, err)
	}

	return ret.String(), nil
}
