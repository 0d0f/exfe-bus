package bot

import (
	"testing"
	"strings"
	"fmt"
	"regexp"
	"time"
)

type TesterContext struct {
	*BaseContext
	bot   *BotTester
	mark  string
	input string
}

func NewTesterContext(bot *BotTester, id string) *TesterContext {
	ret := &TesterContext{NewBaseContext(id), bot, "", ""}
	return ret
}

func (ctx *TesterContext) CheckTimeout(input string) error {
	if ctx.DurationFromLast() > time.Second {
		ctx.mark = ""
	}
	return nil
}

func (ctx *TesterContext) ProcessWithMark(input string) error {
	markers := ctx.bot.markerPattern.FindStringSubmatch(string(input))
	if len(markers) == 0 {
		return BotNotMatched
	}
	ctx.mark = markers[1]
	ctx.input = string(input)
	return nil
}

func (ctx *TesterContext) Default(input string) error {
	ctx.input = string(input)
	return nil
}

type BotTester struct {
	markerPattern *regexp.Regexp
}

func NewBotTester() *BotTester {
	return &BotTester{
		markerPattern: regexp.MustCompile("#([^ ]*)"),
	}
}

func (b *BotTester) GetIDFromInput(input interface{}) (id string, content interface{}, err error) {
	inputText, ok := input.(string)
	if !ok {
		err = fmt.Errorf("Input not a string")
		return
	}
	splits := strings.SplitN(inputText, " ", 2)
	id = splits[0]
	if len(splits) == 2 {
		content = splits[1]
	} else {
		content = ""
	}
	return
}

func (b *BotTester) GenerateContext(id string) Context {
	return NewTesterContext(b, id)
}

func TestBot(t *testing.T) {
	bot := NewBot(NewBotTester())
	bot.RegisterPretreat("CheckTimeout")
	bot.Register("ProcessWithMark")
	bot.Register("Default")

	bot.Feed("@abc default")
	ids := bot.IDs()
	if ids[0] != "@abc" {
		t.Fatalf("bot should only have one id(@abc), but got: %s", ids)
	}
	context1 := bot.GetContext("@abc").(*TesterContext)
	expect := "@abc"
	if context1.ID() != expect {
		t.Errorf("context screen name should be: %s, got: %s", expect, context1.ID())
	}
	expect = "default"
	if context1.input != expect {
		t.Errorf("context input should be: %s, got: %s", expect, context1.input)
	}
	expect = ""
	if context1.mark != expect {
		t.Errorf("context last iom should be: %s, got: %s", expect, context1.mark)
	}

	bot.Feed("@123 #xyz iom")
	context2 := bot.GetContext("@123").(*TesterContext)
	expect = "@123"
	if context2.ID() != expect {
		t.Errorf("context screen name should be: %s, got: %s", expect, context2.ID())
	}
	expect = "#xyz iom"
	if context2.input != expect {
		t.Errorf("context input should be: %s, got: %s", expect, context2.input)
	}
	expect = "xyz"
	if context2.mark != expect {
		t.Errorf("context last iom should be: %s, got: %s", expect, context1.mark)
	}

	bot.Feed("@abc #foo 123")
	expect = "#foo 123"
	if context1.input != expect {
		t.Errorf("context input should be: %s, got: %s", expect, context1.input)
	}
	expect = "foo"
	if context1.mark != expect {
		t.Errorf("context last iom should be: %s, got: %s", expect, context1.mark)
	}

	time.Sleep(time.Second / 2)
	bot.Feed("@abc 456")
	expect = "456"
	if context1.input != expect {
		t.Errorf("context input should be: %s, got: %s", expect, context1.input)
	}
	expect = "foo"
	if context1.mark != expect {
		t.Errorf("context last iom should be: %s, got: %s", expect, context1.mark)
	}
	expect = "xyz"
	if context2.mark != expect {
		t.Errorf("context last iom should be: %s, got: %s", expect, context2.mark)
	}

	time.Sleep(time.Second / 2)
	bot.Feed("@123")
	expect = "foo"
	if context1.mark != expect {
		t.Errorf("context last iom should be: %s, got: %s", expect, context1.mark)
	}
	expect = ""
	if context2.mark != expect {
		t.Errorf("context last iom should be: %s, got: %s", expect, context2.mark)
	}

	time.Sleep(time.Second)
	bot.Feed("@abc")
	expect = ""
	if context1.mark != expect {
		t.Errorf("context last iom should be: %s, got: %s", expect, context1.mark)
	}

	bot.Feed("@abc 890")
	expect = "890"
	if context1.input != expect {
		t.Errorf("context input should be: %s, got: %s", expect, context1.input)
	}
	expect = ""
	if context1.mark != expect {
		t.Errorf("context last iom should be: %s, got: %s", expect, context1.mark)
	}
}
