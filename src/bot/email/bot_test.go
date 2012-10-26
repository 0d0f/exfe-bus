package email

import (
	"exfe/service"
	"github.com/stretchrcom/testify/assert"
	"net/mail"
	"os"
	"testing"
)

func TestEmail(t *testing.T) {
	var config exfe_service.Config
	f, err := os.Open("test.email")
	if err != nil {
		t.Fatalf("can't open test.email: %s", err)
	}
	message, err := mail.ReadMessage(f)
	if err != nil {
		t.Fatalf("test.email format error: %s", err)
	}
	bot := NewEmailBot(&config)
	id, content, err := bot.GetIDFromInput(message)
	assert.Equal(t, err, nil)

	post := content.(*Email)
	assert.Equal(t, id, "googollee@gmail.com")
	assert.Equal(t, post.Text, "被发现bug了……\n\n我来测试一下…………\n\nT_T")
}
