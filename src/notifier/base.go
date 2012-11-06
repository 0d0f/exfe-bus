package notifier

import (
	"bytes"
	"fmt"
	"formatter"
	"model"
	"thirdpart"
)

func GetContent(localTemplate *formatter.LocalTemplate, template string, arg model.ServiceToInterface) (string, error) {
	to := arg.ToRecipient()
	t, err := thirdpart.MessageTypeFromProvider(to.Provider)
	if err != nil {
		return "", err
	}
	templateName := fmt.Sprintf("%s.%s", template, t)

	ret := bytes.NewBuffer(nil)
	err = localTemplate.Execute(ret, to.Language, templateName, arg)
	if err != nil {
		return "", fmt.Errorf("template(%s) failed: %s", templateName, err)
	}

	return ret.String(), nil
}
