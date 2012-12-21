package conversation

import (
	"fmt"
	"github.com/stretchrcom/testify/assert"
	"model"
	"testing"
	"time"
)

type FakeRepo struct {
	identities []model.Identity
	exfee      model.Exfee
}

func newFakeRepo() *FakeRepo {
	identities := []model.Identity{
		model.Identity{
			ID:               1,
			Name:             "exfe",
			Nickname:         "exfe",
			Timezone:         "+0800",
			UserID:           10,
			Provider:         "twitter",
			ExternalID:       "234",
			ExternalUsername: "exfe",
		},
		model.Identity{
			ID:               2,
			Name:             "steve",
			Timezone:         "+0800",
			UserID:           10,
			Provider:         "email",
			ExternalID:       "steve@domain.com",
			ExternalUsername: "steve@domain.com",
		},
		model.Identity{
			ID:               3,
			Name:             "0d0f",
			Timezone:         "+0800",
			UserID:           11,
			Provider:         "twitter",
			ExternalID:       "123",
			ExternalUsername: "0d0f",
		},
	}
	exfee := model.Exfee{
		ID:   1,
		Name: "abc",
		Invitations: []model.Invitation{
			model.Invitation{
				ID:   1,
				Host: true,
				Identity: model.Identity{
					ID:               1,
					Name:             "exfe",
					Nickname:         "exfe",
					Timezone:         "+0800",
					UserID:           10,
					Provider:         "twitter",
					ExternalID:       "234",
					ExternalUsername: "exfe",
				},
			},
			model.Invitation{
				ID:   2,
				Host: false,
				Identity: model.Identity{
					ID:               3,
					Name:             "0d0f",
					Timezone:         "+0800",
					UserID:           11,
					Provider:         "twitter",
					ExternalID:       "123",
					ExternalUsername: "0d0f",
				},
			},
		},
	}
}

func (r *FakeRepo) FindIdentity(identity model.Identity) (model.Identity, error) {
	for _, i := range r.identities {
		if identity.ID == i.ID {
			return i, nil
		}
		if identity.Provider == i.Provider {
			if identity.ExternalID == i.ExternalID {
				return i, nil
			}
			if identity.ExternalUsername == i.ExternalUsername {
				return i, nil
			}
		}
	}
	return model.Identity{}, fmt.Errorf("can't find identity: %s", identity)
}

func (r *FakeRepo) FindExfee(id uint64) (model.Exfee, error) {
	if id == r.exfee.ID {
		return r.exfee, nil
	}
	return model.Exfee{}, fmt.Errorf("can't find exfee: %d", id)
}

func (r *FakeRepo) SendUpdate(tos []model.Recipient, cross model.Cross, post model.Post) error {
	return nil
}

func TestConversation(t *testing.T) {
	repo := newFakeRepo()
	conv := New(repo)
	var minID, maxID uint
	{
		post, err := conv.NewPost(1, model.Post{
			By: model.Identity{
				ID: 1,
			},
			Content: "@0d0f@twitter look at this image http://instagr.am/xxxx\n cool!",
		}, "web", 0)
		assert.Equal(t, err, nil)
		assert.Equal(t, post.By.ID, 1)
		assert.Equal(t, post.By.Name, "exfe")
		assert.Equal(t, post.By.Provider, "twitter")
		assert.Equal(t, post.By.ExternalUsername, "234")
		assert.Equal(t, post.Content, "@0d0f@twitter look at this image {{url:http://instagr.am/xxxx}}\n cool!")
		assert.Equal(t, post.Via, "web")
		assert.Equal(t, post.CreatedAt, time.Now().UTC().Format("2006-01-02 15:04:05 +0700"))
		assert.Contains(t, fmt.Sprintf("%v", post.Relationship), "mention:identity://3")
		assert.Contains(t, fmt.Sprintf("%v", post.Relationship), "url:http://instagr.am/xxxx")
		assert.Equal(t, post.ExfeeID, 1)
		minID = post.ID
	}

	{
		post, err := conv.NewPost(1, model.Post{
			By: model.Identity{
				ExternalUsername: "steve@domain.com",
				Provider:         "email",
			},
			Content: "@0d0f@twitter look at this image http://instagr.am/xxxx\n cool!",
		}, "email", 1355998129)
		assert.Equal(t, err, nil)
		assert.Equal(t, post.By.ID, 2)
		assert.Equal(t, post.By.Name, "exfe")
		assert.Equal(t, post.By.Provider, "twitter")
		assert.Equal(t, post.By.ExternalUsername, "234")
		assert.Equal(t, post.Content, "@0d0f@twitter look at this image {{url:http://instagr.am/xxxx}}\n cool!")
		assert.Equal(t, post.Via, "web")
		assert.Equal(t, post.CreatedAt, "2012-12-20 10:08:49 +0000")
		assert.Contains(t, fmt.Sprintf("%v", post.Relationship), "mention:identity://3")
		assert.Contains(t, fmt.Sprintf("%v", post.Relationship), "url:http://instagr.am/xxxx")
		assert.Equal(t, post.ExfeeID, 1)
		maxID = post.ID
	}

	{
		posts, err := conv.GetPost(1, 0, "2012-12-20 10:08:49", "", 0, 0)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(posts), 2)
		assert.Equal(t, posts[0].ID, maxID)
		assert.Equal(t, posts[1].ID, minID)
	}

	{
		posts, err := conv.GetPost(1, 0, "", "2012-12-20 10:08:49", 0, 0)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(posts), 1)
		assert.Equal(t, posts[0].ID, maxID)
	}

	{
		posts, err := conv.GetPost(1, 0, "2012-12-20 10:08:49", time.Now().Format("2006-01-02 15:04:05 +0700"), 0, 0)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(posts), 2)
		assert.Equal(t, posts[0].ID, maxID)
		assert.Equal(t, posts[1].ID, minID)
	}

	{
		posts, err := conv.GetPost(1, 0, "", "", maxID, 0)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(posts), 1)
		assert.Equal(t, posts[0].ID, maxID)
	}

	{
		posts, err := conv.GetPost(1, 0, "", "", 0, minID)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(posts), 1)
		assert.Equal(t, posts[0].ID, minID)
	}

	{
		posts, err := conv.GetPost(1, 0, "", "", minID, maxID)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(posts), 2)
		assert.Equal(t, posts[0].ID, maxID)
		assert.Equal(t, posts[1].ID, minID)
	}

	{
		unreadcount, err := conv.getunreadcount(1, 10)
		assert.Equal(t, err, nil)
		assert.Equal(t, unreadcount, 2)

		_, err = conv.GetPost(1, 10, "", "", 0, minID)
		assert.Equal(t, err, nil)

		unreadcount, err = conv.getunreadcount(1, 10)
		assert.Equal(t, err, nil)
		assert.Equal(t, unreadcount, 0)
	}

	{
		posts, err := conv.GetPost(1, 0, "", "", minID, maxID)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(posts), 2)

		err = conv.DeletePost(1, minID)
		assert.Equal(t, err, nil)

		posts, err = conv.GetPost(1, 0, "", "", minID, maxID)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(posts), 1)
		assert.Equal(t, posts[0].ID, maxID)
	}
}
