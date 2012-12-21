package email

import (
	"github.com/stretchrcom/testify/assert"
	"net/mail"
	"os"
	"testing"
)

func TestGetCrossID(t *testing.T) {
	bot := NewEmailBot(nil, nil, nil)

	{
		addr := mail.Address{
			Name:    "abc",
			Address: "x+10002@exfe.com",
		}
		id := bot.getCrossId([]*mail.Address{&addr})
		assert.Equal(t, id, "10002")
	}

	{
		addr := mail.Address{
			Name:    "abc",
			Address: "x+c10002@exfe.com",
		}
		id := bot.getCrossId([]*mail.Address{&addr})
		assert.Equal(t, id, "10002")
	}
}

func TestEmail(t *testing.T) {
	f, err := os.Open("test.email")
	if err != nil {
		t.Fatalf("can't open test.email: %s", err)
	}
	message, err := mail.ReadMessage(f)
	if err != nil {
		t.Fatalf("test.email format error: %s", err)
	}
	bot := NewEmailBot(nil, nil, nil)
	id, content, err := bot.GetIDFromInput(message)
	assert.Equal(t, err, nil)

	post := content.(*Email)
	assert.Equal(t, post.CrossID, "100241")
	assert.Equal(t, id, "googollee@gmail.com")
	assert.Equal(t, post.Text, "被发现bug了……\n\n我来测试一下…………\n\nT_T")
}
