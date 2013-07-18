package notifier

import (
	"broker"
	"bytes"
	"errors"
	"fmt"
	"formatter"
	"model"
)

var noneedSend = errors.New("no need send")

func GenerateContent(localTemplate *formatter.LocalTemplate, template string, poster, lang string, arg interface{}) (string, error) {
	switch poster {
	case "facebook":
		fallthrough
	case "google":
		poster = "email"
	}

	templateName := fmt.Sprintf("%s/%s", poster, template)
	if !localTemplate.IsExist(lang, templateName) {
		templateName = fmt.Sprintf("_default/%s", template)
	}

	ret := bytes.NewBuffer(nil)
	err := localTemplate.Execute(ret, lang, templateName, arg)
	if err != nil {
		return "", fmt.Errorf("template(%s/%s) failed: %s", lang, templateName, err)
	}

	if ret.Len() == 0 {
		return "", noneedSend
	}

	return ret.String(), nil
}

func SendAndSave(localTemplate *formatter.LocalTemplate, platform *broker.Platform, to model.Recipient, arg interface{}, template, failUrl string) error {
	text, err := GenerateContent(localTemplate, template, to.Provider, to.Language, arg)
	if err != nil {
		return err
	}
	id, ontime, defaultOk, err := platform.Send(to, text)
	if err != nil {
		return err
	}
	WaitResponse(id, ontime, defaultOk, to, failUrl, arg)
	return nil
}
