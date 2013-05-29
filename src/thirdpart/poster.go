package thirdpart

import (
	"fmt"
	"github.com/googollee/go-rest"
	"io"
	"io/ioutil"
	// "model"
	"net/http"
)

func init() {
	rest.RegisterMarshaller("plain/text", new(PlainText))
}

type IPoster interface {
	Provider() string

	Post(userId string, text string) (messageId string, err error)
	// Send(to *model.Recipient, text string) (messageId string, err error)
}

type Poster struct {
	rest.Service `prefix:"/v3/poster" mime:"plain/text"`

	Post rest.Processor `path:"/:provider/*id" method:"POST"`

	posters map[string]IPoster
}

func NewPoster() (*Poster, error) {
	ret := &Poster{
		posters: make(map[string]IPoster),
	}
	return ret, nil
}

func (m *Poster) Add(poster IPoster) {
	m.posters[poster.Provider()] = poster
}

func (m Poster) HandlePost(text string) string {
	provider := m.Vars()["provider"]
	id := m.Vars()["id"]
	poster, ok := m.posters[provider]
	if !ok {
		m.Error(http.StatusBadRequest, m.DetailError(1, "invalid provider: %s", provider))
		return ""
	}

	// to := &model.Recipient{
	// 	ExternalUsername: id,
	// 	ExternalID:       id,
	// 	Provider:         provider,
	// }
	// ret, err := poster.Send(to, text)
	ret, err := poster.Post(id, text)
	if err != nil {
		m.Error(http.StatusInternalServerError, m.DetailError(2, "%s", err))
	}
	return ret
}

type PlainText struct{}

func (p PlainText) Unmarshal(r io.Reader, v interface{}) error {
	ps, ok := v.(*string)
	if !ok {
		return fmt.Errorf("plain text only can save in string")
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	*ps = string(b)
	return nil
}

func (p PlainText) Marshal(w io.Writer, v interface{}) error {
	_, err := w.Write([]byte(fmt.Sprintf("%s", v)))
	return err
}

func (p PlainText) Error(code int, message string) error {
	return fmt.Errorf("(%d)%s", code, message)
}
