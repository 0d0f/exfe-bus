package twitter

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"model"
	"net/url"
	"testing"
	"thirdpart"
)

type FakeBroker struct {
	cmds   []string
	paths  []string
	params []url.Values
	id     int
}

func (b *FakeBroker) Reset() {
	b.cmds = make([]string, 0, 0)
	b.paths = make([]string, 0, 0)
	b.params = make([]url.Values, 0, 0)
	b.id = 0
}

func copyValues(value url.Values) url.Values {
	ret := make(url.Values)
	for k, v := range value {
		ret[k] = v
	}
	return ret
}

func (b *FakeBroker) Do(cmd, path string, params url.Values) (io.ReadCloser, error) {
	b.cmds = append(b.cmds, cmd)
	b.paths = append(b.paths, path)
	b.params = append(b.params, copyValues(params))
	b.id++
	ret := bytes.NewBufferString(fmt.Sprintf(`{"id_str":"%d"}`, b.id))
	return ioutil.NopCloser(ret), nil
}

var toPublic = &model.Recipient{
	ExternalID:       "12345",
	ExternalUsername: "publicer",
	AuthData:         "",
	Provider:         "twitter",
	IdentityID:       123,
	UserID:           1,
}

var toPrivate = &model.Recipient{
	ExternalID:       "54321",
	ExternalUsername: "privater",
	AuthData:         "",
	Provider:         "twitter",
	IdentityID:       321,
	UserID:           2,
}

func TestSend(t *testing.T) {
	helper := new(thirdpart.FakeHelper)
	broker := new(FakeBroker)
	broker.Reset()
	twitter := New(nil, nil, broker, helper)

	{
		broker.Reset()
		id, err := twitter.Send(toPublic, "private message", "public message", nil)
		if err != nil {
			t.Fatalf("send fail: %s", err)
		}

		if got, expect := id, "1"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
		if got, expect := broker.cmds[0], "POST"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
		if got, expect := broker.paths[0], "direct_messages/new.json"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
		if got, expect := broker.params[0].Get("text"), "private message"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
		if got, expect := broker.params[0].Get("user_id"), "12345"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}

	{
		broker.Reset()
		id, err := twitter.Send(toPublic, `\(accepted name1\), \(accepted name2\) and \(accepted name3\) are accepted on \(“some cross”\), \(inter name1\), \(inter name2\) and \(inter name3\) interested, \(un name1\), \(un name2\) and \(un name3\) is unavailable, \(pending name1\), \(pending name2\) and \(pending name3\) is pending. \(3 of 10 accepted\). https://exfe.com/#!token=932ce5324321433253`, "public message", nil)
		if err != nil {
			t.Fatalf("send fail: %s", err)
		}

		if got, expect := id, "3"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
		for i := 0; i < 3; i++ {
			if got, expect := broker.cmds[i], "POST"; got != expect {
				t.Errorf("got: %s, expect: %s", got, expect)
			}
			if got, expect := broker.paths[i], "direct_messages/new.json"; got != expect {
				t.Errorf("got: %s, expect: %s", got, expect)
			}
		}
		results := []string{
			"(1/3)accepted name1, accepted name2 and accepted name3 are accepted on “some cross”, inter name1, inter name2 and inter name3",
			"(2/3)interested, un name1, un name2 and un name3 is unavailable, pending name1, pending name2 and pending name3 is pending.",
			"(3/3)3 of 10 accepted. https://exfe.com/#!token=932ce5324321433253",
		}
		for i := 0; i < 3; i++ {
			if got, expect := broker.params[i].Get("user_id"), "12345"; got != expect {
				t.Errorf("got: %s, expect: %s", got, expect)
			}
			if got, expect := broker.params[i].Get("text"), results[i]; got != expect {
				t.Errorf("got: %s, expect: %s", got, expect)
			}
		}
	}
}
