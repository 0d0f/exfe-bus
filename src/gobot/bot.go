package bot

import (
	"reflect"
)

type Botter interface {
	GetIDFromInput(input interface{}) (id string, content interface{}, err error)
	GenerateContext(id string) Context
}

type Bot struct {
	botter    Botter
	pretreats []string
	methods   []string
	contexts  map[string]*innerContext
}

func NewBot(botter Botter) *Bot {
	return &Bot{
		botter:   botter,
		methods:  make([]string, 0, 0),
		contexts: make(map[string]*innerContext),
	}
}

func (b *Bot) RegisterPretreat(method string) {
	b.pretreats = append(b.pretreats, method)
}

func (b *Bot) Register(method string) {
	b.methods = append(b.methods, method)
}

func (b *Bot) Feed(input interface{}) error {
	id, content, err := b.botter.GetIDFromInput(input)
	if err != nil {
		return err
	}
	ctx, ok := b.contexts[id]
	if !ok {
		ctx, err = newInnerContext(b, b.botter.GenerateContext(id), reflect.TypeOf(input))
		if err != nil {
			return err
		}
		b.contexts[id] = ctx
	}

	return ctx.feed(content)
}

func (b *Bot) IDs() []string {
	ret := make([]string, len(b.contexts), len(b.contexts))
	i := 0
	for id, _ := range b.contexts {
		ret[i] = id
	}
	return ret
}

func (b *Bot) GetContext(id string) interface{} {
	return b.contexts[id].ctx
}
