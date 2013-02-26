package broker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gobus"
	"model"
	"net/http"
)

type Platform struct {
	dispatcher *gobus.Dispatcher
	config     *model.Config
}

func NewPlatform(config *model.Config) (*Platform, error) {
	table, err := gobus.NewTable(config.Dispatcher)
	if err != nil {
		return nil, err
	}
	dispatcher := gobus.NewDispatcher(table)
	return &Platform{
		dispatcher: dispatcher,
		config:     config,
	}, nil
}

func (p *Platform) Send(to model.Recipient, private, public string, info *model.InfoData) (string, error) {
	arg := model.ThirdpartSend{
		PrivateMessage: private,
		PublicMessage:  public,
		Info:           info,
	}
	arg.To = to

	var ids string
	err := p.dispatcher.DoWithTicket(to.Provider, "bus://exfe_service/thirdpart/message", "POST", &arg, &ids)

	if err != nil {
		return "", err
	}
	return ids, nil
}

func (p *Platform) GetHotRecipient(userID int64) ([]model.Recipient, error) {
	return nil, nil

	resp, err := http.Get(fmt.Sprintf("%s/v2/Gobus/HotRecipient?user=%d", p.config.SiteApi, userID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("find identity failed: %s", resp.Status)
	}
	decoder := json.NewDecoder(resp.Body)
	var ret []model.Recipient
	err = decoder.Decode(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (p *Platform) FindIdentity(identity model.Identity) (model.Identity, error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(identity)
	if err != nil {
		return identity, err
	}
	resp, err := http.Post(fmt.Sprintf("%s/v2/Gobus/RevokeIdentity", p.config.SiteApi), "application/json", buf)
	if err != nil {
		return identity, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return identity, fmt.Errorf("find identity failed: %s", resp.Status)
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&identity)
	if err != nil {
		return identity, err
	}
	return identity, nil
}

func (p *Platform) FindCross(id uint64) (model.Cross, error) {
	var ret model.Cross
	resp, err := http.Get(fmt.Sprintf("%s/v2/Gobus/GetCrossById?id=%d", p.config.SiteApi, id))
	if err != nil {
		return ret, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ret, fmt.Errorf("find cross failed: %s", resp.Status)
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func (p *Platform) UploadPhoto(photoxID string, photos []model.Photo) error {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(photos)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("%s/v2/Gobus/AddPhotos/%s", p.config.SiteApi, photoxID), "application/json", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("upload photo fail: %s", resp.Status)
	}
	return nil
}

func (p *Platform) BotCreateCross(cross model.Cross) error {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(cross)
	if err != nil {
		return err
	}
	p.config.Log.Debug("create cross: %s", buf.String())
	return nil
}

func (p *Platform) BotPostConversation(post, to, id string) error {
	p.config.Log.Debug("post (%s) to %s(%s)", post, to, id)
	return nil
}
