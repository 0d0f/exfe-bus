package thirdpart

import (
	"fmt"
	"model"
)

type Thirdpart struct {
	senders  map[string]Sender
	updaters map[string]Updater
}

func New() *Thirdpart {
	return &Thirdpart{
		senders:  make(map[string]Sender),
		updaters: make(map[string]Updater),
	}
}

func (t *Thirdpart) AddSender(sender Sender) {
	t.senders[sender.Provider()] = sender
}

func (t *Thirdpart) Send(to *model.Recipient, privateMessage, publicMessage string, data *InfoData) (string, error) {
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
