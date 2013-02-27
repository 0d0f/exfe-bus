package thirdpart

import (
	"model"
	"testing"
)

type Faker struct {
	provider        string
	tos             []*model.Recipient
	privateMessages []string
	publicMessages  []string
	datas           []*model.InfoData
}

func (f *Faker) Reset() {
	f.tos = make([]*model.Recipient, 0)
	f.privateMessages = make([]string, 0)
	f.publicMessages = make([]string, 0)
	f.datas = make([]*model.InfoData, 0)
}

func (f *Faker) Provider() string {
	return f.provider
}

func (f *Faker) Send(to *model.Recipient, privateMessage string, publicMessage string, data *model.InfoData) (id string, err error) {
	f.tos = append(f.tos, to)
	f.privateMessages = append(f.privateMessages, privateMessage)
	f.publicMessages = append(f.publicMessages, publicMessage)
	f.datas = append(f.datas, data)

	return "1", nil
}

func (f *Faker) UpdateFriends(to *model.Recipient) error {
	f.tos = append(f.tos, to)
	return nil
}

func (f *Faker) UpdateIdentity(to *model.Recipient) error {
	f.tos = append(f.tos, to)
	return nil
}

var to1 = &model.Recipient{
	ExternalID:       "12345",
	ExternalUsername: "to1",
	AuthData:         "",
	Provider:         "faker1",
	IdentityID:       123,
	UserID:           1,
}

var to2 = &model.Recipient{
	ExternalID:       "54321",
	ExternalUsername: "to2",
	AuthData:         "",
	Provider:         "faker2",
	IdentityID:       321,
	UserID:           2,
}

func TestThirdpartSender(t *testing.T) {
	faker1 := &Faker{
		provider: "faker1",
	}
	faker2 := &Faker{
		provider: "faker2",
	}
	config := new(model.Config)
	third := New(config)
	third.AddSender(faker1)
	third.AddSender(faker2)

	{
		faker1.Reset()
		faker2.Reset()

		third.Send(to1, "private to 1", "public to 1", nil)
		if len(faker1.tos) != 1 {
			t.Fatalf("faker1 should received 1 message")
		}
		if got, expect := faker1.tos[0].ExternalUsername, "to1"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
		if got, expect := faker1.privateMessages[0], "private to 1"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
		if got, expect := faker1.publicMessages[0], "public to 1"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}

		third.Send(to2, "private to 2", "public to 2", nil)
		if len(faker2.tos) != 1 {
			t.Fatalf("faker2 should received 1 message")
		}
		if got, expect := faker2.tos[0].ExternalUsername, "to2"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
		if got, expect := faker2.privateMessages[0], "private to 2"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
		if got, expect := faker2.publicMessages[0], "public to 2"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}
}

func TestThirdpartUpdate(t *testing.T) {
	faker1 := &Faker{
		provider: "faker1",
	}
	faker2 := &Faker{
		provider: "faker2",
	}
	config := new(model.Config)
	third := New(config)
	third.AddUpdater(faker1)
	third.AddUpdater(faker2)

	{
		faker1.Reset()
		faker2.Reset()

		third.UpdateFriends(to1)
		if len(faker1.tos) != 1 {
			t.Fatalf("faker1 should received 1 message")
		}
		if got, expect := faker1.tos[0].ExternalUsername, "to1"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}

		third.UpdateFriends(to2)
		if len(faker2.tos) != 1 {
			t.Fatalf("faker2 should received 1 message")
		}
		if got, expect := faker2.tos[0].ExternalUsername, "to2"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}

	{
		faker1.Reset()
		faker2.Reset()

		third.UpdateIdentity(to1)
		if len(faker1.tos) != 1 {
			t.Fatalf("faker1 should received 1 message")
		}
		if got, expect := faker1.tos[0].ExternalUsername, "to1"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}

		third.UpdateIdentity(to2)
		if len(faker2.tos) != 1 {
			t.Fatalf("faker2 should received 1 message")
		}
		if got, expect := faker2.tos[0].ExternalUsername, "to2"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}
}
