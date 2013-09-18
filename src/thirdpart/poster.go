package thirdpart

import (
	"fmt"
	"github.com/googollee/go-broadcast"
	"github.com/googollee/go-rest"
	"io"
	"io/ioutil"
	"logger"
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

	post     rest.SimpleNode `route:"/message/:provider/*id" method:"POST"`
	response rest.SimpleNode `route:"/response/:provider/*id" method:"POST"`
	watch    rest.Streaming  `route:"" method:"WATCH"`

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
	provider := poster.Provider()
	waiting, defaultOK := poster.SetPosterCallback(func(id string, err error) {
		resp := PostResponse{
			Id: fmt.Sprintf("%s-%s", provider, id),
			Ok: err == nil,
		}
		if !resp.Ok {
			resp.Error = err.Error()
		}
		logger.INFO("poster", provider, "response", id, resp.Ok, resp.Error)
		m.watchChan.Send(resp)
	})
	m.posters[provider] = posterHandler{
		poster:    poster,
		waiting:   waiting,
		defaultOK: defaultOK,
	}
}

func (m Poster) Post(ctx rest.Context, text string) {
	if text == "" {
		ctx.Return(http.StatusBadRequest, "invalid text")
		return
	}
	var provider, id string
	ctx.Bind("provider", &provider)
	ctx.Bind("id", &id)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}
	handler, ok := m.posters[provider]
	if !ok {
		ctx.Return(http.StatusBadRequest, "invalid provider: %s", provider)
		return
	}

	ret, err := handler.poster.Post(ctx.Request().URL.Query().Get("from"), id, text)
	if err != nil {
		logger.INFO("poster", provider, "fail", id, err.Error())
		ctx.Return(http.StatusInternalServerError, err)
		return
	}
	ontime := time.Now().Add(handler.waiting).Unix()
	if handler.waiting > 0 {
		logger.INFO("poster", provider, "waiting", ret, id, "ontime", ontime, fmt.Sprintf("default %v", handler.defaultOK))
		ctx.Response().Header().Set("Ontime", fmt.Sprintf("%d", ontime))
		ctx.Response().Header().Set("Default", fmt.Sprintf("%v", handler.defaultOK))
		ctx.Return(http.StatusAccepted)
	} else {
		logger.INFO("poster", provider, "ok", ret, id)
		ctx.Return(http.StatusOK)
	}
	ctx.Render(fmt.Sprintf("%s-%s", provider, ret))
}

func (m Poster) Response(ctx rest.Context, resp PostResponse) {
	var provider, id string
	ctx.Bind("provider", &provider)
	ctx.Bind("id", &id)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}
	logger.INFO("poster", provider, "response", id, resp.Ok, resp.Error)
	resp.Id = fmt.Sprintf("%s-%s", provider, id)
	m.watchChan.Send(resp)
}

func (m Poster) Watch(ctx rest.StreamContext) {
	c := make(chan interface{})
	err := m.watchChan.Register(c)
	if err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}
	defer m.watchChan.Unregister(c)
	ctx.Return(http.StatusOK)

	for ctx.Ping() == nil {
		select {
		case i := <-c:
			ctx.SetWriteDeadline(time.Now().Add(time.Second))
			if err := ctx.Render(i); err != nil {
				return
			}
		case <-time.After(time.Second):
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
