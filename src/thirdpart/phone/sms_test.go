package sms

import (
	"model"
	"testing"
	"thirdpart"
)

type testSender struct {
	to       string
	contents []string
}

func (s *testSender) Codes() []string {
	return []string{"+1", "+86"}
}

func (s *testSender) Send(phone string, contents []string) (string, error) {
	s.to = phone
	s.contents = contents
	return "1", nil
}

var to = &model.Recipient{
	ExternalID:       "+8613488802890",
	ExternalUsername: "to",
	AuthData:         "",
	Provider:         "sms",
	IdentityID:       321,
	UserID:           2,
}

var data = &model.InfoData{
	CrossID: 12345,
	Type:    model.TypeCrossUpdate,
}

func TestSend(t *testing.T) {
	testSender := new(testSender)
	config := new(model.Config)
	sms := new(Sms)
	sms.config = config
	sms.senders = make(map[string]Sender)
	sms.senders["+1"] = testSender
	sms.senders["+86"] = testSender

	var tester thirdpart.Sender
	tester = sms

	{
		_, err := tester.Send(to, `\(AAAAAAAA name1\), \(AAAAAAAA name2\) and \(AAAAAAAA name3\) are accepted on \("some cross"\), \(IIIII name1\), \(IIIII name2\) and \(IIIII name3\) interested, \(UUUU name1\), \(UUUU name2\) and \(UUUU name3\) are unavailable, \(PPPPPPP name1\), \(PPPPPPP name2\) and \(PPPPPPP name3\) are pending. \(3 of 10 accepted\). https://exfe.com/#!token=932ce5324321433253`)
		if err != nil {
			t.Fatalf("send error: %s", err)
		}
		results := []string{
			`AAAAAAAA name1, AAAAAAAA name2 and AAAAAAAA name3 are accepted on "some cross", IIIII name1, IIIII name2 and IIIII name3 interested, (1/3)`,
			`UUUU name1, UUUU name2 and UUUU name3 are unavailable, PPPPPPP name1, PPPPPPP name2 and PPPPPPP name3 are pending. 3 of 10 accepted. (2/3)`,
			`https://exfe.com/#!token=932ce5324321433253 (3/3)`,
		}
		if got, expect := len(testSender.contents), len(results); got != expect {
			t.Errorf("got: %d, expect: %d", got, expect)
		}
		for i, r := range results {
			if got, expect := testSender.contents[i], r; got != expect {
				t.Errorf("%d got: %s, expect: %s", i, got, expect)
			}
		}
	}

	{
		_, err := tester.Send(to, `Post: abc \(("Title" https://exfe.com/#!token=932ce5324321433253)\)`)
		if err != nil {
			t.Fatalf("send error: %s", err)
		}
		results := []string{
			`Post: abc ("Title" https://exfe.com/#!token=932ce5324321433253)`,
		}
		if got, expect := len(testSender.contents), len(results); got != expect {
			t.Errorf("got: %d, expect: %d", got, expect)
		}
		for i, r := range results {
			if got, expect := testSender.contents[i], r; got != expect {
				t.Errorf("%d got: %s, expect: %s", i, got, expect)
			}
		}
	}

	{
		_, err := tester.Send(to, "post line1\npost line2\n\n")
		if err != nil {
			t.Fatalf("send error: %s", err)
		}
		results := []string{
			`post line1`,
			`post line2`,
		}
		if got, expect := len(testSender.contents), len(results); got != expect {
			t.Errorf("got: %d, expect: %d", got, expect)
		}
		for i, r := range results {
			if got, expect := testSender.contents[i], r; got != expect {
				t.Errorf("%d got: %s, expect: %s", i, got, expect)
			}
		}
	}

	{
		_, err := tester.Send(to, `Googol Lee: 测试时间 (“看电影 007” \(https://exfe.com/#!token=cd48a91ee3c2afb545d32f301b342510\))
aadfdafdas https://exfe.com/fdafa`)
		if err != nil {
			t.Fatalf("send error: %s", err)
		}
		results := []string{
			`Googol Lee: 测试时间 (“看电影 007” (1/2)`,
			`https://exfe.com/#!token=cd48a91ee3c2afb545d32f301b342510) (2/2)`,
			`aadfdafdas https://exfe.com/fdafa`,
		}
		if got, expect := len(testSender.contents), len(results); got != expect {
			t.Errorf("got: %d, expect: %d", got, expect)
		}
		for i, r := range results {
			if got, expect := testSender.contents[i], r; got != expect {
				t.Errorf("%d got: %s, expect %s", i, got, expect)
			}
		}
	}
}
