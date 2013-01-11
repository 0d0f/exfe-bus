package message

import (
	"fmt"
	"github.com/googollee/go-logger"
	"github.com/stretchrcom/testify/assert"
	"model"
	"testing"
)

type FakeDispatcher struct {
	tickets []string
	urls    []string
	methods []string
	args    []interface{}
	replys  []interface{}
}

func NewFakeDispatcher() *FakeDispatcher {
	return &FakeDispatcher{
		tickets: make([]string, 0),
		urls:    make([]string, 0),
		methods: make([]string, 0),
		args:    make([]interface{}, 0),
		replys:  make([]interface{}, 0),
	}
}

func (f *FakeDispatcher) Reset() {
	f.tickets = make([]string, 0)
	f.urls = make([]string, 0)
	f.methods = make([]string, 0)
	f.args = make([]interface{}, 0)
	f.replys = make([]interface{}, 0)
}

func (f *FakeDispatcher) DoWithTicket(ticket, url, method string, arg, reply interface{}) error {
	f.tickets = append(f.tickets, ticket)
	f.urls = append(f.urls, url)
	f.methods = append(f.methods, method)
	f.args = append(f.args, arg)
	f.replys = append(f.replys, reply)
	return nil
}

type FakePlatform struct{}

func (p *FakePlatform) GetHotRecipient(userID int64) ([]model.Recipient, error) {
	return nil, nil
}

func TestMessage(t *testing.T) {
	config := new(model.Config)
	log, err := logger.New(logger.Stderr, "message test")
	if err != nil {
		t.Fatal(err)
	}
	config.Log = log
	dispatcher := NewFakeDispatcher()
	platform := new(FakePlatform)
	msg, err := New(config, dispatcher, platform)
	if err != nil {
		t.Fatal(err)
	}

	var recipients [4]model.Recipient
	recipients[0].Provider = "twitter"
	recipients[1].Provider = "iOS"
	recipients[2].Provider = "email"
	recipients[3].Provider = "Android"

	err = msg.Send("bus://exfe_service/notifier/conversation", "cross123", recipients[:], 1)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(dispatcher.tickets), 2)
	assert.Equal(t, len(dispatcher.urls), 2)
	assert.Equal(t, len(dispatcher.methods), 2)
	assert.Equal(t, len(dispatcher.args), 2)
	assert.Equal(t, len(dispatcher.replys), 2)
	assert.Contains(t, fmt.Sprintf("%v", dispatcher.urls), "bus://exfe_queue/instant")
	assert.Contains(t, fmt.Sprintf("%v", dispatcher.urls), "bus://exfe_queue/head10")
}
