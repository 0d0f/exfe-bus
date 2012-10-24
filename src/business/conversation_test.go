package business

import (
	"github.com/stretchrcom/testify/assert"
	"model"
	"testing"
)

var post1 = model.Post{
	ID:        1,
	By:        email1,
	Content:   "email1 post sth",
	Via:       "abc",
	CreatedAt: "2012-10-24 16:31:00",
}

var post2 = model.Post{
	ID:        2,
	By:        twitter3,
	Content:   "twitter3 post sth",
	Via:       "abc",
	CreatedAt: "2012-10-24 16:40:00",
}

func TestConversationUpdateToSelf(t *testing.T) {
	update1 := model.ConversationUpdate{
		To:    remail1,
		Cross: cross,
		Post:  post1,
	}
	updates := []model.ConversationUpdate{update1}

	c := NewConversation(localTemplate, &config)
	private, public, err := c.getContent(updates)
	assert.Equal(t, err.Error(), "can't parse posts: no need send self")
	assert.Equal(t, private, "")
	assert.Equal(t, public, "")
}

func TestConversationUpdateEmail(t *testing.T) {
	update1 := model.ConversationUpdate{
		To:    remail1,
		Cross: cross,
		Post:  post1,
	}
	update2 := model.ConversationUpdate{
		To:    remail1,
		Cross: cross,
		Post:  post2,
	}
	updates := []model.ConversationUpdate{update1, update2}

	c := NewConversation(localTemplate, &config)
	private, public, err := c.getContent(updates)
	t.Logf("err: %s", err)
	t.Errorf("private:-----start------\n%s\n-------end-------", private)
	t.Errorf("public:-----start------\n%s\n-------end-------", public)
}

func TestConversationUpdateTwitter(t *testing.T) {
	update1 := model.ConversationUpdate{
		To:    rtwitter1,
		Cross: cross,
		Post:  post1,
	}
	update2 := model.ConversationUpdate{
		To:    rtwitter1,
		Cross: cross,
		Post:  post2,
	}
	updates := []model.ConversationUpdate{update1, update2}

	expectPrivate := `email1 name: email1 post sth \((“Test Cross” http://site/url/#!token=recipient_twitter1_token)\)
twitter3 name: twitter3 post sth \((“Test Cross” http://site/url/#!token=recipient_twitter1_token)\)
`
	c := NewConversation(localTemplate, &config)
	private, public, err := c.getContent(updates)
	assert.Equal(t, err, nil)
	assert.Equal(t, private, expectPrivate)
	assert.Equal(t, public, "")
}
