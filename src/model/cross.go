package model

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"logger"
	"net/http"
)

type UpdateInfo struct {
	UpdatedAt string   `json:"updated_at"`
	By        Identity `json:"by_identity"`
}

type Cross struct {
	ID          uint64                   `json:"id,omitempty"`
	By          Identity                 `json:"by_identity,omitempty"`
	Title       string                   `json:"title,omitempty"`
	Description string                   `json:"description,omitempty"`
	Time        *CrossTime               `json:"time,omitempty"`
	Place       *Place                   `json:"place,omitempty"`
	Exfee       Exfee                    `json:"exfee,omitempty"`
	Updated     map[string]UpdateInfo    `json:"updated,omitempty"`
	Widgets     []map[string]interface{} `json:"widget,omitempty"`
	CreatedAt   string                   `json:"created_at,omitempty"`
	UpdatedAt   string                   `json:"updated_at,omitempty"`
	Attribute   struct {
		State string `json:"state"`
	} `json:"attribute"`
}

func (c Cross) Equal(other *Cross) bool {
	return c.ID == other.ID
}

func (c Cross) String() string {
	return fmt.Sprintf("Cross:%d", c.ID)
}

func (c Cross) Ics(config *Config, to Recipient) string {
	url := fmt.Sprintf("%s/v2/ics/crosses?token=%s", config.SiteApi, to.Token)
	resp, err := http.Get(url)
	if err != nil {
		logger.ERROR("get %s error: %s", url, err)
		return ""
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.ERROR("get %s error: %s", url, err)
		return ""
	}
	if resp.StatusCode != http.StatusOK {
		logger.ERROR("get %s error: (%s)%s", url, resp.Status, string(b))
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

func (c Cross) Background(config *Config) string {
	for _, w := range c.Widgets {
		if t, ok := w["type"].(string); !ok || t != "Background" {
			continue
		}
		if img, ok := w["image"]; ok && img != "" {
			return fmt.Sprintf("%s/static/img/xbg/%s", config.SiteUrl, img)
		}
	}
	return fmt.Sprintf("%s/static/img/xbg/default.jpg", config.SiteUrl)
}

func (c Cross) Timezone(to Recipient) string {
	if to.Timezone != "" {
		return to.Timezone
	}
	return c.Time.BeginAt.Timezone
}

func (c Cross) Link(to Recipient, config *Config) string {
	return fmt.Sprintf("%s/#!token=%s", config.SiteUrl, to.Token)
}

func (c Cross) PublicLink(to Recipient, config *Config) string {
	token := to.Token
	if len(token) > 5 {
		token = token[1:5]
	}
	return fmt.Sprintf("%s/#!%d/%s", config.SiteUrl, c.ID, token)
}

type CrossDigestRequest struct {
	To        Recipient `json:"to"`
	CrossId   int64     `json:"cross_id"`
	UpdatedAt string    `json:"updated_at"`
}

type CrossUpdate struct {
	To       Recipient `json:"to"`
	OldCross Cross     `json:"old_cross"`
	Cross    Cross     `json:"cross"`
	By       Identity  `json:"by"`
}

type CrossUpdates []CrossUpdate

func (u CrossUpdates) String() string {
	if len(u) == 0 {
		return "{updates:0}"
	}
	return fmt.Sprintf("{to:%s with:%s updates:%d}", u[0].To, u[0].Cross, len(u))
}
