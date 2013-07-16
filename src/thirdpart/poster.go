package thirdpart

import (
	"fmt"
	"github.com/googollee/go-broadcast"
	"github.com/googollee/go-rest"
	"io"
	"io/ioutil"
	"net/http"
)

func init() {
	rest.RegisterMarshaller("plain/text", new(PlainText))
}

type IPoster interface {
	Provider() string

	Post(from, to string, text string) (messageId string, err error)
}

type Response struct {
	Id    string `json:"id"`
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type Poster struct {
	rest.Service `prefix:"/v3/poster" mime:"plain/text"`

	Post  rest.Processor `path:"/:provider/*id" method:"POST"`
	Watch rest.Streaming `path:"/" method:"WATCH"`

	posters   map[string]IPoster
	watchChan *broadcast.Broadcast
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
	if text == "" {
		return ""
	}
	provider := m.Vars()["provider"]
	id := m.Vars()["id"]
	poster, ok := m.posters[provider]
	if !ok {
		m.Error(http.StatusBadRequest, m.DetailError(1, "invalid provider: %s", provider))
		return ""
	}

	ret, err := poster.Post(m.Request().URL.Query().Get("from"), id, text)
	if err != nil {
		m.Error(http.StatusInternalServerError, m.DetailError(2, "%s", err))
	}
	return ret
}

func (m Poster) HandleWatch(stream rest.Stream) {

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

func (p PlainText) Marshal(w io.Writer, name string, v interface{}) error {
	_, err := w.Write([]byte(fmt.Sprintf("%s", v)))
	return err
}

func (p PlainText) Error(code int, message string) error {
	return fmt.Errorf("(%d)%s", code, message)
}
