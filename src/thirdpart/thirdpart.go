package thirdpart

import (
	"fmt"
	"model"
)

var Unreachable = fmt.Errorf("Unreachable")

type Thirdpart struct {
	senders       map[string]Sender
	updaters      map[string]Updater
	photographers map[string]Photographer
	config        *model.Config
}

func New(config *model.Config) *Thirdpart {
	return &Thirdpart{
		senders:       make(map[string]Sender),
		updaters:      make(map[string]Updater),
		photographers: make(map[string]Photographer),
		config:        config,
	}
}

func (t *Thirdpart) AddSender(sender Sender) {
	t.senders[sender.Provider()] = sender
}

func (t *Thirdpart) Send(to *model.Recipient, text string) (string, error) {
	sender, ok := t.senders[to.Provider]
	if !ok {
		return "", fmt.Errorf("can't find %s sender", to.Provider)
	}
	return sender.Send(to, text)
}

func (t *Thirdpart) AddUpdater(updater Updater) {
	t.updaters[updater.Provider()] = updater
}

func (t *Thirdpart) UpdateFriends(to *model.Recipient) error {
	updater, ok := t.updaters[to.Provider]
	if !ok {
		return fmt.Errorf("can't find %s updater", to)
	}
	return updater.UpdateFriends(to)
}

func (t *Thirdpart) UpdateIdentity(to *model.Recipient) error {
	updater, ok := t.updaters[to.Provider]
	if !ok {
		return fmt.Errorf("can't find %s updater", to)
	}
	return updater.UpdateIdentity(to)
}

func (t *Thirdpart) AddPhotographer(photographer Photographer) {
	t.photographers[photographer.Provider()] = photographer
}

func (t *Thirdpart) GrabPhotos(to model.Recipient, albumID string) ([]model.Photo, error) {
	photographer, ok := t.photographers[to.Provider]
	if !ok {
		return nil, fmt.Errorf("can't find %s photographer", to)
	}
	return photographer.Grab(to, albumID)
}

func (t *Thirdpart) GetPhotos(to model.Recipient, pictureIDs []string) ([]string, error) {
	photographer, ok := t.photographers[to.Provider]
	if !ok {
		return nil, fmt.Errorf("can't find %s photographer", to)
	}
	return photographer.Get(to, pictureIDs)
}
