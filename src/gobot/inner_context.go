package bot

import (
	"fmt"
	"reflect"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

type BotInterruptError struct {
	err error
}

func (e *BotInterruptError) String() string {
	return fmt.Sprintf("Bot interrupted by: %s", e.err)
}

func (e *BotInterruptError) Error() string {
	return e.String()
}

var typeOfInterrupted = reflect.TypeOf((*BotInterruptError)(nil)).Elem()

type innerContext struct {
	bot       *Bot
	ctx       Context
	pretreats []reflect.Method
	methods   []reflect.Method
}

var BotNoProcessed = fmt.Errorf("no method processed input")

var BotNotMatched = fmt.Errorf("input not match method")

func newInnerContext(bot *Bot, ctx Context, inputType reflect.Type) (*innerContext, error) {
	ret := &innerContext{
		bot: bot,
		ctx: ctx,
	}
	var err error
	ret.pretreats, err = ret.grabMethods(bot.pretreats, inputType)
	if err != nil {
		return nil, err
	}
	ret.methods, err = ret.grabMethods(bot.methods, inputType)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *innerContext) grabMethods(methods []string, inputType reflect.Type) ([]reflect.Method, error) {
	t := reflect.TypeOf(c.ctx)
	ret := make([]reflect.Method, len(methods), len(methods))
	for i, name := range methods {
		m, ok := t.MethodByName(name)
		if !ok {
			return nil, fmt.Errorf("Context doesn't have method %s", name)
		}
		if ins := m.Type.NumIn(); ins != 2 {
			return nil, fmt.Errorf("Context method %s has wrong numer of ins: %d", name, ins)
		}
		if inType := m.Type.In(0); inputType.AssignableTo(inType) {
			return nil, fmt.Errorf("Context method %s input %s can't be assigned by %s", name, inType, inputType)
		}
		if outs := m.Type.NumOut(); outs != 1 {
			return nil, fmt.Errorf("Context method %s has wrong numer of outs: %d", name, outs)
		}
		if returnType := m.Type.Out(0); returnType != typeOfError {
			return nil, fmt.Errorf("Context method %s returns %s not error", name, returnType)
		}
		ret[i] = m
	}
	return ret, nil
}

func (c *innerContext) feed(input interface{}) error {
	for _, m := range c.pretreats {
		err := m.Func.Call([]reflect.Value{reflect.ValueOf(c.ctx), reflect.ValueOf(input)})[0]
		if err.Type() == typeOfInterrupted {
			return err.Interface().(error)
		}
	}
	defer c.ctx.SetLast()
	for _, m := range c.methods {
		err := m.Func.Call([]reflect.Value{reflect.ValueOf(c.ctx), reflect.ValueOf(input)})[0]
		if err.Interface() == nil {
			return nil
		}
		if e := err.Interface().(error); e != BotNotMatched {
			return e
		}
	}
	return BotNoProcessed
}
