package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"model"
	"net/http"
)

type Platform struct {
	config *model.Config
}

func NewPlatform(config *model.Config) (*Platform, error) {
	return &Platform{
		config: config,
	}, nil
}

func (p *Platform) GetHotRecipient(userID int64) ([]model.Recipient, error) {
	resp, err := http.Get(fmt.Sprintf("%s/v2/Gobus/HotRecipient?user=%d", p.config.SiteApi, userID))
	if err != nil {
		return nil, err
	}
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
