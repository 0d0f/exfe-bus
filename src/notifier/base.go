package notifier

import (
	"bytes"
	"fmt"
	"formatter"
)

func GenerateContent(localTemplate *formatter.LocalTemplate, template string, poster, lang string, arg interface{}) (string, error) {
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
		return "", fmt.Errorf("no need send")
	}

	return ret.String(), nil
}
