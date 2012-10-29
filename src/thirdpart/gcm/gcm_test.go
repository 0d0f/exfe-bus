package gcm

import (
	"encoding/json"
	"fmt"
	"github.com/googollee/go-gcm"
	"model"
	"testing"
	"thirdpart"
)

type FakeBroker struct {
	messages []string
}

func (b *FakeBroker) Reset() {
	b.messages = make([]string, 0)
}

func (b *FakeBroker) Send(n *gcm.Message) (*gcm.Response, error) {
	b.messages = append(b.messages, n.Data["text"])

	result := `{"results":[{"message_id":"1:08","registration_id":"%s"}]}`
	result = fmt.Sprintf(result, n.RegistrationIDs[0])
	ret := new(gcm.Response)
	json.Unmarshal([]byte(result), ret)
	return ret, nil
}

var to = &model.Recipient{
	ExternalID:       "54321",
	ExternalUsername: "to",
	AuthData:         "",
	Provider:         "iOS",
	IdentityID:       321,
	UserID:           2,
}

var data = &thirdpart.InfoData{
	CrossID: 12345,
	Type:    thirdpart.CrossUpdate,
}

func TestSend(t *testing.T) {
	broker := new(FakeBroker)
	gcm := New(broker)
	var tester thirdpart.Sender
	tester = gcm

	{
		broker.Reset()
		_, err := tester.Send(to, `\(AAAAAAAA name1\), \(AAAAAAAA name2\) and \(AAAAAAAA name3\) are accepted on \(“some cross”\), \(IIIII name1\), \(IIIII name2\) and \(IIIII name3\) interested, \(UUUU name1\), \(UUUU name2\) and \(UUUU name3\) are unavailable, \(PPPPPPP name1\), \(PPPPPPP name2\) and \(PPPPPPP name3\) are pending. \(3 of 10 accepted\). https://exfe.com/#!token=932ce5324321433253`, "", data)
		if err != nil {
			t.Fatalf("send error: %s", err)
		}
		results := []string{
			`AAAAAAAA name1, AAAAAAAA name2 and AAAAAAAA name3 are accepted on “some cross”, IIIII name1, IIIII name2 and IIIII name3 interested,…(1/2)`,
			`UUUU name1, UUUU name2 and UUUU name3 are unavailable, PPPPPPP name1, PPPPPPP name2 and PPPPPPP name3 are pending. 3 of 10 accepted. (2/2)`,
		}
		if got, expect := len(broker.messages), len(results); got != expect {
			t.Errorf("got: %d, expect: %d", got, expect)
		}
		for i, r := range results {
			if got, expect := broker.messages[i], r; got != expect {
				t.Errorf("%d got: %s, expect %s", i, got, expect)
			}
		}
	}

	{
		broker.Reset()
		_, err := tester.Send(to, `Post: abc ("Title" https://exfe.com/#!token=932ce5324321433253)`, "", data)
		if err != nil {
			t.Fatalf("send error: %s", err)
		}
		results := []string{
			`Post: abc ("Title")`,
		}
		if got, expect := len(broker.messages), len(results); got != expect {
			t.Errorf("got: %d, expect: %d", got, expect)
		}
		for i, r := range results {
			if got, expect := broker.messages[i], r; got != expect {
				t.Errorf("%d got: %s, expect %s", i, got, expect)
			}
		}
	}

	{
		broker.Reset()
		_, err := tester.Send(to, "post line1\npost line2\n\n", "", data)
		if err != nil {
			t.Fatalf("send error: %s", err)
		}
		results := []string{
			`post line1`,
			`post line2`,
		}
		if got, expect := len(broker.messages), len(results); got != expect {
			t.Errorf("got: %d, expect: %d", got, expect)
		}
		for i, r := range results {
			if got, expect := broker.messages[i], r; got != expect {
				t.Errorf("%d got: %s, expect %s", i, got, expect)
			}
		}
	}
}
