package apn

import (
	"fmt"
	"github.com/virushuo/Go-Apns"
	"model"
	"testing"
	"thirdpart"
)

type FakeBroker struct {
	errChan       chan apns.NotificationError
	notifications []*apns.Notification
}

func (b *FakeBroker) Reset() {
	b.notifications = make([]*apns.Notification, 0)
	b.errChan = make(chan apns.NotificationError)
}

func (b *FakeBroker) GetErrorChan() <-chan apns.NotificationError {
	return b.errChan
}

func (b *FakeBroker) Send(n *apns.Notification) error {
	b.notifications = append(b.notifications, n)
	return nil
}

func errHandler(err apns.NotificationError) {
	fmt.Println(err)
}

var to = &model.Recipient{
	ExternalID:       "54321",
	ExternalUsername: "to",
	AuthData:         "",
	Provider:         "iOS",
	IdentityID:       321,
	UserID:           2,
}

var data = &model.InfoData{
	CrossID: 12345,
	Type:    model.TypeCrossUpdate,
}

func TestSend(t *testing.T) {
	broker := new(FakeBroker)
	apn := New(broker, errHandler)
	var tester thirdpart.Sender
	tester = apn

	{
		broker.Reset()
		_, err := tester.Send(to, `\(AAAAAAAA name1\), \(AAAAAAAA name2\) and \(AAAAAAAA name3\) are accepted on \(“some cross”\), \(IIIII name1\), \(IIIII name2\) and \(IIIII name3\) interested, \(UUUU name1\), \(UUUU name2\) and \(UUUU name3\) are unavailable, \(PPPPPPP name1\), \(PPPPPPP name2\) and \(PPPPPPP name3\) are pending. \(3 of 10 accepted\). https://exfe.com/#!token=932ce5324321433253`, "", data)
		if err != nil {
			t.Fatalf("send error: %s", err)
		}
		results := []string{
			`AAAAAAAA name1, AAAAAAAA name2 and AAAAAAAA name3 are accepted on “some cross”, IIIII name1, IIIII name2 and IIIII name3…(1/3)`,
			`interested, UUUU name1, UUUU name2 and UUUU name3 are unavailable, PPPPPPP name1, PPPPPPP name2 and PPPPPPP name3 are pending.…(2/3)`,
			`3 of 10 accepted. (3/3)`,
		}
		if got, expect := len(broker.notifications), len(results); got != expect {
			t.Errorf("got: %d, expect: %d", got, expect)
		}
		for i, r := range results {
			if got, expect := broker.notifications[i].Payload.Aps.Alert, r; got != expect {
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
		if got, expect := len(broker.notifications), len(results); got != expect {
			t.Errorf("got: %d, expect: %d", got, expect)
		}
		for i, r := range results {
			if got, expect := broker.notifications[i].Payload.Aps.Alert, r; got != expect {
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
		if got, expect := len(broker.notifications), len(results); got != expect {
			t.Errorf("got: %d, expect: %d", got, expect)
		}
		for i, r := range results {
			if got, expect := broker.notifications[i].Payload.Aps.Alert, r; got != expect {
				t.Errorf("%d got: %s, expect %s", i, got, expect)
			}
		}
	}
}
