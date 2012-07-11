package bot

import (
	"time"
)

type Context interface {
	ID() string
	SetLast()
	DurationFromLast() time.Duration
}

type BaseContext struct {
	id   string
	last time.Time
}

func NewBaseContext(id string) *BaseContext {
	return &BaseContext{
		id: id,
	}
}

func (c *BaseContext) ID() string {
	return c.id
}

func (c *BaseContext) DurationFromLast() time.Duration {
	return time.Now().Sub(c.last)
}

func (c *BaseContext) SetLast() {
	c.last = time.Now()
}
