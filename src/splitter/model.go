package splitter

import (
	"fmt"
	"model"
)

type Pack struct {
	MergeKey string
	Method   string
	Service  string
	Data     map[string]interface{}
}

type BigPack struct {
	Recipients []model.Recipient
	MergeKey   string
	Method     string
	Service    string
	Data       map[string]interface{}
}

func (b BigPack) Each() chan Pack {
	ret := make(chan Pack)
	pack := Pack{
		Method:  b.Method,
		Service: b.Service,
	}

	go func() {
		for _, to := range b.Recipients {
			pack.MergeKey = fmt.Sprintf("%s_%d", b.MergeKey, to.IdentityID)
			pack.Data = b.Data
			pack.Data["to"] = to
			ret <- pack
		}
		close(ret)
	}()

	return ret
}
