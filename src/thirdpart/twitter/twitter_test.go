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

func (b *FakeBroker) Do(accessToken *thirdpart.Token, cmd, path string, params url.Values) (io.ReadCloser, error) {
	if path == "direct_messages/new.json" && params.Get("user_id") == "12345" {
		return nil, fmt.Errorf(`{"code":150}`)
	}
	b.cmds = append(b.cmds, cmd)
	b.paths = append(b.paths, path)
	b.params = append(b.params, copyValues(params))
	b.id++
	ret := bytes.NewBufferString(fmt.Sprintf(`{"id_str":"%d"}`, b.id))
	if path == "direct_messages/new.json" {
		ret = bytes.NewBufferString(fmt.Sprintf(`{"id_str":"%d","recipient":{"id_str":"7890","screen_name":"twitter_id","profile_image_url":"http://twitter/avatar","description":"desc","name":"twitterName"}}`, b.id))
	}
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
	config := new(model.Config)
	twitter := New(config, broker, helper)
	var tester thirdpart.Sender
	tester = twitter

	{
		broker.Reset()
		id, err := tester.Send(toPublic, "private message", "public message", nil)
		if err != nil {
			t.Fatalf("send fail: %s", err)
		}

		if got, expect := id, "1"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
		if got, expect := broker.cmds[0], "POST"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
		if got, expect := broker.paths[0], "statuses/update.json"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
		if got, expect := broker.params[0].Get("status"), "@publicer public message"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}
	{
		broker.Reset()
		id, err := tester.Send(toPrivate, "private message", "public message", nil)
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
		if got, expect := broker.params[0].Get("user_id"), "54321"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}

	{
		broker.Reset()
		id, err := tester.Send(toPrivate, `\(AAAAAAAA name1\), \(AAAAAAAA name2\) and \(AAAAAAAA name3\) are accepted on \(“some cross”\), \(IIIII name1\), \(IIIII name2\) and \(IIIII name3\) interested, \(UUUU name1\), \(UUUU name2\) and \(UUUU name3\) are unavailable, \(PPPPPPP name1\), \(PPPPPPP name2\) and \(PPPPPPP name3\) are pending. \(3 of 10 accepted\). https://exfe.com/#!token=932ce5324321433253`, "public message", nil)
		if err != nil {
			t.Fatalf("send fail: %s", err)
		}

		if got, expect := id, "1,2,3"; got != expect {
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
			`AAAAAAAA name1, AAAAAAAA name2 and AAAAAAAA name3 are accepted on “some cross”, IIIII name1, IIIII name2 and IIIII name3 interested, (1/3)`,
			`UUUU name1, UUUU name2 and UUUU name3 are unavailable, PPPPPPP name1, PPPPPPP name2 and PPPPPPP name3 are pending. 3 of 10 accepted. (2/3)`,
			`https://exfe.com/#!token=932ce5324321433253 (3/3)`,
		}
		for i := 0; i < 3; i++ {
			if got, expect := broker.params[i].Get("user_id"), "54321"; got != expect {
				t.Errorf("got: %s, expect: %s", got, expect)
			}
			if got, expect := broker.params[i].Get("text"), results[i]; got != expect {
				t.Errorf("got: %s, expect: %s", got, expect)
			}
		}
	}

	{
		broker.Reset()
		id, err := tester.Send(toPublic, "private", `\(AAAAAAAA name1\), \(AAAAAAAA name2\) and \(AAAAAAAA name3\) are accepted on \(“some cross”\), \(IIIII name1\), \(IIIII name2\) and \(IIIII name3\) interested, \(UUUU name1\), \(UUUU name2\) and \(UUUU name3\) are unavailable, \(PPPPPPP name1\), \(PPPPPPP name2\) and \(PPPPPPP name3\) are pending. \(3 of 10 accepted\). https://exfe.com/#!token=932ce5324321433253`, nil)
		if err != nil {
			t.Fatalf("send fail: %s", err)
		}

		if got, expect := id, "1,2,3"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
		for i := 0; i < 3; i++ {
			if got, expect := broker.cmds[i], "POST"; got != expect {
				t.Errorf("got: %s, expect: %s", got, expect)
			}
			if got, expect := broker.paths[i], "statuses/update.json"; got != expect {
				t.Errorf("got: %s, expect: %s", got, expect)
			}
		}
		results := []string{
			`@publicer AAAAAAAA name1, AAAAAAAA name2 and AAAAAAAA name3 are accepted on “some cross”, IIIII name1, IIIII name2 and IIIII name3 (1/3)`,
			`@publicer interested, UUUU name1, UUUU name2 and UUUU name3 are unavailable, PPPPPPP name1, PPPPPPP name2 and PPPPPPP name3 are (2/3)`,
			`@publicer pending. 3 of 10 accepted. https://exfe.com/#!token=932ce5324321433253 (3/3)`,
		}
		for i := 0; i < 3; i++ {
			if got, expect := broker.params[i].Get("status"), results[i]; got != expect {
				t.Errorf("got: %s, expect: %s", got, expect)
			}
		}
	}

	{
		broker.Reset()
		id, err := tester.Send(toPrivate, "post line1\npost line2\n\n\n", "public message", nil)
		if err != nil {
			t.Fatalf("send fail: %s", err)
		}

		if got, expect := id, "1,2"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
		for i := 0; i < 2; i++ {
			if got, expect := broker.cmds[i], "POST"; got != expect {
				t.Errorf("got: %s, expect: %s", got, expect)
			}
			if got, expect := broker.paths[i], "direct_messages/new.json"; got != expect {
				t.Errorf("got: %s, expect: %s", got, expect)
			}
		}
		results := []string{
			`post line1`,
			`post line2`,
		}
		for i := 0; i < 2; i++ {
			if got, expect := broker.params[i].Get("user_id"), "54321"; got != expect {
				t.Errorf("got: %s, expect: %s", got, expect)
			}
			if got, expect := broker.params[i].Get("text"), results[i]; got != expect {
				t.Errorf("got: %s, expect: %s", got, expect)
			}
		}
	}

	{
		broker.Reset()
		id, err := tester.Send(toPublic, "private", "post line1\npost line2\n\n\n", nil)
		if err != nil {
			t.Fatalf("send fail: %s", err)
		}

		if got, expect := id, "1,2"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
		for i := 0; i < 2; i++ {
			if got, expect := broker.cmds[i], "POST"; got != expect {
				t.Errorf("got: %s, expect: %s", got, expect)
			}
			if got, expect := broker.paths[i], "statuses/update.json"; got != expect {
				t.Errorf("got: %s, expect: %s", got, expect)
			}
		}
		results := []string{
			`@publicer post line1`,
			`@publicer post line2`,
		}
		for i := 0; i < 2; i++ {
			if got, expect := broker.params[i].Get("status"), results[i]; got != expect {
				t.Errorf("got: %s, expect: %s", got, expect)
			}
		}
	}
}
