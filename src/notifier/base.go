package notifier

import (
	"bytes"
	"fmt"
	"formatter"
	"model"
)

func GetContent(localTemplate *formatter.LocalTemplate, template string, to model.Recipient, arg interface{}) (string, error) {
	templateName := fmt.Sprintf("%s/%s", to.Provider, template)
	if !localTemplate.IsExist(to.Language, templateName) {
		templateName = fmt.Sprintf("_default/%s", template)
	}

	ret := bytes.NewBuffer(nil)
	err := localTemplate.Execute(ret, to.Language, templateName, arg)
	if err != nil {
		return "", fmt.Errorf("template(%s) failed: %s", templateName, err)
	}

	return ret.String(), nil
}
