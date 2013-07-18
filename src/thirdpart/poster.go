package thirdpart

import (
	"fmt"
	"github.com/googollee/go-broadcast"
	"github.com/googollee/go-rest"
	"io"
	"io/ioutil"
	"model"
	"net/http"
	"time"
)

func init() {
	rest.RegisterMarshaller("plain/text", new(PlainText))
}

type Callback func(id string, err error)

type posterHandler struct {
	poster    IPoster
	waiting   time.Duration
	defaultOK bool
}

type IPoster interface {
	Provider() string
	SetPosterCallback(f Callback) (waiting time.Duration, defaultOK bool)
	Post(from, to string, text string) (messageId string, err error)
}

type PostResponse struct {
	Id    string `json:"id"`
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type Poster struct {
	rest.Service `prefix:"/v3/poster" mime:"plain/text"`

	Post     rest.Processor `path:"/message/:provider/*id" method:"POST"`
	Response rest.Processor `path:"/response/:provider/*id" method:"POST"`
	Watch    rest.Streaming `path:"" method:"WATCH"`

	config    *model.Config
	posters   map[string]posterHandler
	watchChan *broadcast.Broadcast
}

func NewPoster() (*Poster, error) {
	ret := &Poster{
		posters:   make(map[string]posterHandler),
		watchChan: broadcast.NewBroadcast(10),
	}
	return ret, nil
}

func (m *Poster) Add(poster IPoster) {
	waiting, defaultOK := poster.SetPosterCallback(func(id string, err error) {
		resp := PostResponse{
			Id:    fmt.Sprintf("%s-%s", poster.Provider(), id),
			Ok:    err == nil,
			Error: err.Error(),
		}
		m.HandleResponse(resp)
	})
	m.posters[poster.Provider()] = posterHandler{
		poster:    poster,
		waiting:   waiting,
		defaultOK: defaultOK,
	}
}

func (m Poster) HandlePost(text string) string {
	if text == "" {
		return ""
	}
	provider := m.Vars()["provider"]
	id := m.Vars()["id"]
	handler, ok := m.posters[provider]
	if !ok {
		m.Error(http.StatusBadRequest, m.DetailError(1, "invalid provider: %s", provider))
		return ""
	}

	ret, err := handler.poster.Post(m.Request().URL.Query().Get("from"), id, text)
	if err != nil {
		m.Error(http.StatusInternalServerError, m.DetailError(2, "%s", err))
		return ""
	}
	if handler.waiting > 0 {
		m.Header().Set("Ontime", fmt.Sprintf("%d", time.Now().Add(handler.waiting).Unix()))
		m.Header().Set("Default", fmt.Sprintf("%v", handler.defaultOK))
		m.WriteHeader(http.StatusAccepted)
	} else {
		m.WriteHeader(http.StatusOK)
	}
	return fmt.Sprintf("%s-%s", provider, ret)
}

func (m Poster) HandleResponse(resp PostResponse) {
	resp.Id = fmt.Sprintf("%s-%s", m.Vars()["provider"], m.Vars()["id"])
	m.watchChan.Send(resp)
}

func (m Poster) HandleWatch(stream rest.Stream) {
	c := make(chan interface{})
	err := m.watchChan.Register(c)
	if err != nil {
		m.Error(http.StatusBadRequest, err)
		return
	}
	defer m.watchChan.Unregister(c)
	m.WriteHeader(http.StatusOK)

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
