package thirdpart

import (
	"fmt"
	"model"
)

type Thirdpart struct {
	senders  map[string]Sender
	updaters map[string]Updater
	config   *model.Config
}

func New(config *model.Config) *Thirdpart {
	return &Thirdpart{
		senders:  make(map[string]Sender),
		updaters: make(map[string]Updater),
		config:   config,
	}
}

func (t *Thirdpart) AddSender(sender Sender) {
	t.senders[sender.Provider()] = sender
}

func (t *Thirdpart) Send(to *model.Recipient, privateMessage, publicMessage string, data *model.InfoData) (string, error) {
	if to.ExternalID == "" {
		go func() {
			err := t.UpdateIdentity(to)
			if err != nil {
				t.config.Log.Crit("update %s identity error: %s", to, err)
			}
		}()
	}
	sender, ok := t.senders[to.Provider]
	if !ok {
		return "", fmt.Errorf("can't find %s sender", to.Provider)
	}
	return sender.Send(to, privateMessage, publicMessage, data)
}

func (t *Thirdpart) AddUpdater(updater Updater) {
	t.updaters[updater.Provider()] = updater
}

func (t *Thirdpart) UpdateFriends(to *model.Recipient) error {
	updater, ok := t.updaters[to.Provider]
	if !ok {
		return fmt.Errorf("can't find %s updater", to.Provider)
	}
	return updater.UpdateFriends(to)
}

func (t *Thirdpart) UpdateIdentity(to *model.Recipient) error {
	updater, ok := t.updaters[to.Provider]
	if !ok {
		return fmt.Errorf("can't find %s updater", to.Provider)
	}
	return updater.UpdateIdentity(to)
}
