package photostreaming

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"model"
	"net/http"
)

type FileMeta struct {
	Height   string `json:"height"`
	Width    string `json:"width"`
	Checksum string `json:"checksum"`
}

type PhotoMeta struct {
	PhotoGuid   string     `json:"photoGuid"`
	DateCreated string     `json:"dateCreated"`
	Caption     string     `json:"caption"`
	Derivatives []FileMeta `json:"derivatives"`
}

type StreamingList struct {
	Photos []PhotoMeta `json:"photos"`
}

type UrlRequest struct {
	PhotoGuids []strings `json:"photoGuids"`
}

type UrlMeta struct {
	UrlExpiry   string `json:"url_expiry"`
	UrlPath     string `json:"url_path"`
	UrlLocation string `json:"url_location"`
}

type LocationMeta struct {
	Scheme string   `json:"scheme"`
	Hosts  []string `json:"hosts"`
}

type UrlList struct {
	Locations map[string]LocationMeta `json:"locations"`
	Items     map[string]UrlMeta      `json:"items"`
}

type Photostreaming struct {
	domain string
}

func New(config *model.Config) (*Photostreaming, error) {
	return &Photostreaming{
		domain: config.Thirdpart.Photostreaming.Domain,
	}, nil
}

func (p *Photostreaming) Provider() string {
	return "photostreaming"
}

func (p *Photostreaming) Grab(to model.Recipient, albumID string) ([]model.Photo, error) {
	list, err := p.getList(albumID)
	if err != nil {
		return nil, fmt.Errorf("get streaming failed: %s", err)
	}
	guids := make([]string, len(list.Photos))
	for i, photo := range list.Photos {
		guids[i] = photo.PhotoGuid
	}
	urls, err := p.getUrls(albumID, guids)
	if err != nil {
		return nil, fmt.Errorf("get urls failed: %s", err)
	}

}

func (p *Photostreaming) getList(albumID string) (list StreamingList, err error) {
	url := fmt.Sprintf("https://%s/%s/sharedstreams/webstream", p.domain, albumID)
	buf := bytes.NewBufferString(`{"streamCtag":null}`)
	var body io.ReadCloser
	body, err = p.request(url, buf)
	if err != nil {
		return
	}
	defer body.Close()

	decoder := json.NewDecoder(body)
	err = decoder.Decode(&list)
	return
}

func (p *Photostreaming) getUrls(albumID string, guids []string) (list UrlList, err error) {
	req := UrlRequest{
		Photostreaming: guids,
	}
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err = encoder.Encode(req)
	url := fmt.Sprintf("https://%s/%s/sharedstreams/webasseturls", p.domain, albumID)
	var body io.ReadCloser
	body, err = p.request(url, buf)
	if err != nil {
		return
	}
	defer body.Close()

	decoder := json.NewDecoder(body)
	err = decoder.Decode(&list)
	return
}

func (p *Photostreaming) request(url string, reader io.Reader) (io.ReadCloser, error) {
	resp, err := http.Post(url, "application/json", reader)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		content, err := ioutil.ReadAll()
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s: %s", resp.Status, content)
	}
	return resp.Body, nil
}
