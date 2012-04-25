package twitter_job

import (
	"exfe"
	"gobus"
)

type CrossUpdateArg struct {
	Cross exfe.Cross
	Old_cross exfe.Cross
}

func (a CrossUpdateArg) UpdateMessage() {
}

type Cross_update struct {
	Config *Config
	Client *gobus.Client
}

