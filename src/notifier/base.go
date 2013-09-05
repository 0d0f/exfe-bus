package notifier

import (
	"broker"
	"bytes"
	"errors"
	"fmt"
	"formatter"
	"logger"
	"model"
)

var noneedSend = errors.New("no need send")

func GenerateContent(localTemplate *formatter.LocalTemplate, template string, poster, lang string, arg interface{}) (string, error) {
	templateName := fmt.Sprintf("%s/%s", poster, template)
	if !localTemplate.IsExist(lang, templateName) {
		return "", fmt.Errorf("template(%s) not found in %s/en_us", template, lang)
	}

	ret := bytes.NewBuffer(nil)
	err := localTemplate.Execute(ret, lang, templateName, arg)
	if err != nil {
		return "", fmt.Errorf("template(%s/%s) failed: %s", lang, templateName, err)
	}

	if ret.Len() == 0 {
		return "", fmt.Errorf("template(%s/%s) no need send", lang, templateName)
	}

	return ret.String(), nil
}

// TODO: to是个指针有些过于精妙了。这个指针指向failArg里的to字段，在每次pop时会自动改变failArg里To字段的值，保证wait response的正确。
//       由于interface{}也有可能传值，所以调用者传递failArg时，应该显式使用&保证传址而不是传值。
func SendAndSave(localTemplate *formatter.LocalTemplate, platform *broker.Platform, to *model.Recipient, arg interface{}, template, failUrl string, failArg interface{}) {
	var id string
	var ontime int64
	var defaultOk bool
	needResponse := false
	for run := true; run; run = len(to.Fallbacks) > 0 {
		fallback := to.PopRecipient()
		text, err := GenerateContent(localTemplate, template, fallback.Provider, fallback.Language, arg)
		if err != nil {
			logger.ERROR("generate content failed: %s with %#v", err, arg)
			continue
		}
		id, ontime, defaultOk, err = platform.Send(fallback, text)
		if err != nil {
			logger.INFO("notifier", id, template, to, "error", err)
			if len(to.Fallbacks) == 0 {
				logger.DEBUG("notifier %s send to %s failed: %s with %s", template, fallback, err, text)
			}
			continue
		}
		needResponse = true
		break
	}
	if !needResponse {
		return
	}
	logger.INFO("notifier", id, template, to)
	WaitResponse(id, ontime, defaultOk, *to, failUrl, failArg)
}
