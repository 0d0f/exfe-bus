package thirdpart

import (
	"fmt"
	"github.com/googollee/go-broadcast"
	"github.com/googollee/go-rest"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func init() {
	rest.RegisterMarshaller("plain/text", new(PlainText))
}

type Callback func(id string, err error)

type IPoster interface {
	Provider() string
	SetCallback(f Callback)
	Post(from, to string, text string) (messageId string, waiting bool, err error)
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
		posters:   make(map[string]IPoster),
		watchChan: broadcast.NewBroadcast(-1),
	}
	return ret, nil
}

func (m *Poster) Add(poster IPoster) {
	poster.SetCallback(func(id string, err error) {
		m.callback(fmt.Sprintf("%s_%s", poster.Provider(), id), err)
	})
	m.posters[poster.Provider()] = poster
}

func (m *Poster) callback(id string, err error) {
	response := Response{
		Id:    id,
		Ok:    err == nil,
		Error: err.Error(),
	}
	m.watchChan.Send(response)
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

	ret, waiting, err := poster.Post(m.Request().URL.Query().Get("from"), id, text)
	if err != nil {
		m.Error(http.StatusInternalServerError, m.DetailError(2, "%s", err))
	}
	if waiting {
		m.WriteHeader(http.StatusAccepted)
	} else {
		m.WriteHeader(http.StatusOK)
	}
	return ret
}

func (m Poster) HandleWatch(stream rest.Stream) {
	c := make(chan interface{})
	err := m.watchChan.Register(c)
	if err != nil {
		m.Error(http.StatusBadRequest, err)
		return
	}
	defer m.watchChan.Unregister(c)

	for {
		select {
		case i := <-c:
			stream.SetWriteDeadline(time.Now().Add(time.Second))
			err = stream.Write(i)
		case <-time.After(time.Second):
			err = stream.Ping()
		}
		if err != nil {
			return
		}
	}
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
