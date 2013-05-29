package model

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"logger"
	"net/http"
)

type UpdateInfo struct {
	UpdatedAt string `json:"updated_at"`
}

type Cross struct {
	ID          int64                    `json:"id,omitempty"`
	By          Identity                 `json:"by_identity,omitempty"`
	Title       string                   `json:"title,omitempty"`
	Description string                   `json:"description,omitempty"`
	Time        *CrossTime               `json:"time,omitempty"`
	Place       *Place                   `json:"place,omitempty"`
	Exfee       Exfee                    `json:"exfee,omitempty"`
	Updated     map[string]UpdateInfo    `json:"updated,omitempty"`
	Widgets     []map[string]interface{} `json:"widget"`
}

func (c Cross) Equal(other *Cross) bool {
	return c.ID == other.ID
}

func (c Cross) String() string {
	return fmt.Sprintf("Cross:%d", c.ID)
}

func (c Cross) Ics(config *Config, to Recipient) string {
	url := fmt.Sprintf("%s/v2/ics/crosses?token=%s", config.SiteApi, to.Token)
	logger.DEBUG("ics: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
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

type CrossInvitation struct {
	To      Recipient `json:"to"`
	CrossId int64     `json:"cross_id"`
	Cross   Cross     `json:"cross"`
	By      Identity  `json:"by"`

	Config *Config `json:"-"`
}

func (a CrossInvitation) String() string {
	return fmt.Sprintf("{to:%s cross:%d}", a.To, a.Cross.ID)
}

func (a *CrossInvitation) Parse(config *Config) (err error) {
	a.Config = config
	return nil
}

func (a CrossInvitation) ToIn(invitations []Invitation) bool {
	for _, i := range invitations {
		if a.To.SameUser(&i.Identity) {
			return true
		}
	}
	return false
}

func (a CrossInvitation) Link() string {
	return fmt.Sprintf("%s/#!token=%s", a.Config.SiteUrl, a.To.Token)
}

func (a CrossInvitation) PublicLink() string {
	return fmt.Sprintf("%s/#!%d/%s", a.Config.SiteUrl, a.Cross.ID, a.To.Token[1:5])
}

func (a CrossInvitation) Timezone() string {
	if a.To.Timezone != "" {
		return a.To.Timezone
	}
	return a.Cross.Time.BeginAt.Timezone
}

func (a CrossInvitation) IsCreator() bool {
	return a.To.SameUser(&a.Cross.By)
}

func (a CrossInvitation) LongDescription() bool {
	if len(a.Cross.Description) > 200 {
		return true
	}
	return false
}

func (a CrossInvitation) ListInvitations() string {
	l := len(a.Cross.Exfee.Invitations)
	max := 3
	ret := ""
	for i := 0; i < 3 && i < l; i++ {
		if i > 0 {
			ret += ", "
		}
		ret += a.Cross.Exfee.Invitations[i].Identity.Name
	}
	if l > max {
		ret += "..."
	}
	return ret
}
