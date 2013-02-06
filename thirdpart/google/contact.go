package google

import (
	"exfe/model"
	"strings"
)

type Link struct {
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
	Href string `xml:"href,attr"`
}

type Email struct {
	Address string `xml:"address,attr"`
	Primary bool   `xml:"primary,attr"`
}

type Contact struct {
	Name    string    `xml:"title"`
	Updated time.Time `xml:"updated"`
	Emails  []Email   `xml:"email"`
}

type Feed struct {
	Links    []Link    `xml:"link"`
	Contacts []Contact `xml:"entry"`
}

func (c Contact) Identities() []Identity {
	length := len(c.Emails)
	ret := make([]Identity, length)
	for i, e := range c.Emails {
		if c.Name != "" {
			ret[i].Name = c.Name
		} else {
			split := strings.Split(e.Address, "@")
			ret[i].Name = split[0]
		}
		ret[i].External_id = e.Address
		ret[i].External_username = e.Address
	}
	return ret
}

func (c *Client) GetFriends() {
	resp, err := c.http.Get("https://www.google.com/m8/feeds/contacts/default/full?max-results=100")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
}
