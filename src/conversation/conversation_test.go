package conversation

import (
	"fmt"
	"github.com/stretchrcom/testify/assert"
	"model"
	convmodel "model/conversation"
	"testing"
	"time"
)

type FakeRepo struct {
	identities []model.Identity
	posts      []convmodel.Post
	cross      model.Cross
	tos        []model.Recipient
	cross_     model.Cross
	post       model.Post
	unread     map[int64]map[string]int
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
	cross := model.Cross{
		ID:    1,
		By:    identities[0],
		Title: "cross title",
		Exfee: model.Exfee{
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
		},
	}
	return &FakeRepo{
		identities: identities,
		cross:      cross,
		unread:     make(map[int64]map[string]int),
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

func (r *FakeRepo) FindCross(id uint64) (model.Cross, error) {
	if id == r.cross.ID {
		return r.cross, nil
	}
	return model.Cross{}, fmt.Errorf("can't find exfee: %d", id)
}

func (r *FakeRepo) SendUpdate(tos []model.Recipient, cross model.Cross, post model.Post) error {
	r.tos = tos
	r.cross_ = cross
	r.post = post
	return nil
}

func (r *FakeRepo) SavePost(post convmodel.Post) (uint64, error) {
	post.ID = uint64(len(r.posts) + 1)
	r.posts = append(r.posts, post)
	return post.ID, nil
}

func (r *FakeRepo) FindPosts(exfeeID uint64, refID, sinceTime, untilTime string, minID, maxID uint64) ([]convmodel.Post, error) {
	since, err := time.Parse("2006-01-02 15:04:05", sinceTime)
	checkSince := true
	if err != nil {
		checkSince = false
	}
	until, err := time.Parse("2006-01-02 15:04:05", untilTime)
	checkUntil := true
	if err != nil {
		checkUntil = false
	}
	checkMax := true
	if maxID == 0 {
		checkMax = false
	}
	checkMin := true
	if minID == 0 {
		checkMin = false
	}
	ret := make([]convmodel.Post, 0)
	for _, p := range r.posts {
		if p.ExfeeID != exfeeID {
			continue
		}
		if p.RefURI != refID {
			continue
		}
		if checkMin && p.ID < minID {
			continue
		}
		if checkMax && p.ID > maxID {
			continue
		}
		if checkSince && p.CreatedAt.Before(since) {
			continue
		}
		if checkUntil && p.CreatedAt.After(until) {
			continue
		}
		ret = append(ret, p)
	}
	return ret, nil
}

func (r *FakeRepo) DeletePost(refID string, postID uint64) error {
	index := -1
	for i, p := range r.posts {
		if p.RefURI != refID {
			continue
		}
		if p.ID == postID {
			index = i
			break
		}
	}
	if index >= 0 {
		r.posts = append(r.posts[:index], r.posts[index+1:]...)
	}
	return nil
}

func (r *FakeRepo) SetUnreadCount(refID string, userID int64, count int) error {
	if _, ok := r.unread[userID]; !ok {
		r.unread[userID] = make(map[string]int)
	}
	r.unread[userID][refID] = count
	return nil
}

func (r *FakeRepo) AddUnreadCount(refID string, userID int64, count int) error {
	if _, ok := r.unread[userID]; !ok {
		r.unread[userID] = make(map[string]int)
		r.unread[userID][refID] = 0
	}
	r.unread[userID][refID] += count
	return nil
}

func (r *FakeRepo) GetUnreadCount(refID string, userID int64) (int, error) {
	if _, ok := r.unread[userID]; !ok {
		return 0, nil
	}
	return r.unread[userID][refID], nil
}

func TestConversation(t *testing.T) {
	repo := newFakeRepo()
	conv := New(repo)
	var minID, maxID uint64
	{
		post, err := conv.NewPost(1, model.Post{
			By: model.Identity{
				ID: 1,
			},
			Content: "@0d0f@twitter look at this image http://instagr.am/xxxx\n cool!",
		}, "web", 0)
		assert.Equal(t, err, nil)
		assert.Equal(t, post.ID, uint64(1))
		assert.Equal(t, post.By.ID, 1)
		assert.Equal(t, post.By.Name, "exfe")
		assert.Equal(t, post.By.Provider, "twitter")
		assert.Equal(t, post.By.ExternalUsername, "exfe")
		assert.Equal(t, post.Content, "@0d0f@twitter look at this image {{url:http://instagr.am/xxxx}}\n cool!")
		assert.Equal(t, post.Via, "web")
		assert.Equal(t, post.CreatedAt, time.Now().UTC().Format("2006-01-02 15:04:05 -0700"))
		assert.Contains(t, fmt.Sprintf("%v", post.Relationship), "mention:identity://3")
		assert.Contains(t, fmt.Sprintf("%v", post.Relationship), "url:http://instagr.am/xxxx")
		assert.Equal(t, post.ExfeeID, uint64(1))
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
		assert.Equal(t, post.ID, uint64(2))
		assert.Equal(t, post.By.ID, 1)
		assert.Equal(t, post.By.Name, "exfe")
		assert.Equal(t, post.By.Provider, "twitter")
		assert.Equal(t, post.By.ExternalUsername, "exfe")
		assert.Equal(t, post.Content, "@0d0f@twitter look at this image {{url:http://instagr.am/xxxx}}\n cool!")
		assert.Equal(t, post.Via, "email")
		assert.Equal(t, post.CreatedAt, "2012-12-20 10:08:49 +0000")
		assert.Contains(t, fmt.Sprintf("%v", post.Relationship), "mention:identity://3")
		assert.Contains(t, fmt.Sprintf("%v", post.Relationship), "url:http://instagr.am/xxxx")
		assert.Equal(t, post.ExfeeID, uint(1))
		maxID = post.ID
	}

	{
		posts, err := conv.FindPosts(1, 0, "2012-12-20 10:08:49", "", 0, 0)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(posts), 2)
		assert.Equal(t, posts[0].ID, minID)
		assert.Equal(t, posts[1].ID, maxID)
	}

	{
		posts, err := conv.FindPosts(1, 0, "", "2012-12-20 10:08:49", 0, 0)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(posts), 1)
		assert.Equal(t, posts[0].ID, maxID)
	}

	{
		posts, err := conv.FindPosts(1, 0, "2012-12-20 10:08:49", time.Now().Format("2006-01-02 15:04:05 +0700"), 0, 0)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(posts), 2)
		assert.Equal(t, posts[0].ID, minID)
		assert.Equal(t, posts[1].ID, maxID)
	}

	{
		posts, err := conv.FindPosts(1, 0, "", "", maxID, 0)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(posts), 1)
		assert.Equal(t, posts[0].ID, maxID)
	}

	{
		posts, err := conv.FindPosts(1, 0, "", "", 0, minID)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(posts), 1)
		assert.Equal(t, posts[0].ID, minID)
	}

	{
		posts, err := conv.FindPosts(1, 0, "", "", minID, maxID)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(posts), 2)
		assert.Equal(t, posts[0].ID, minID)
		assert.Equal(t, posts[1].ID, maxID)
	}

	{
		unreadcount, err := conv.GetUnreadCount(1, 10)
		assert.Equal(t, err, nil)
		assert.Equal(t, unreadcount, 2)

		_, err = conv.FindPosts(1, 10, "", "", 0, minID)
		assert.Equal(t, err, nil)

		unreadcount, err = conv.GetUnreadCount(1, 10)
		assert.Equal(t, err, nil)
		assert.Equal(t, unreadcount, 0)
	}

	{
		posts, err := conv.FindPosts(1, 0, "", "", minID, maxID)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(posts), 2)

		err = conv.DeletePost(1, minID)
		assert.Equal(t, err, nil)

		posts, err = conv.FindPosts(1, 0, "", "", minID, maxID)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(posts), 1)
		assert.Equal(t, posts[0].ID, maxID)
	}
}

func TestParse(t *testing.T) {
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
					Provider:         "email",
					ExternalID:       "abc@domain.com",
					ExternalUsername: "abc@domain.com",
				},
			},
		},
	}
	repo := newFakeRepo()
	conv := New(repo)

	{
		content, relationships := conv.parseRelationship("@exfe@twitter blablabla", exfee)
		assert.Equal(t, content, "@exfe@twitter blablabla")
		assert.Contains(t, fmt.Sprintf("%v", relationships), "mention:identity://1")
	}

	{
		content, relationships := conv.parseRelationship("@abc@domain.com@email blablabla", exfee)
		assert.Equal(t, content, "@abc@domain.com@email blablabla")
		assert.Contains(t, fmt.Sprintf("%v", relationships), "mention:identity://3")
	}

	{
		content, relationships := conv.parseRelationship("@abc@domain.com blablabla", exfee)
		assert.Equal(t, content, "@abc@domain.com blablabla")
		assert.Contains(t, fmt.Sprintf("%v", relationships), "mention:identity://3")
	}

	{
		content, relationships := conv.parseRelationship("@abc@domain.com blablabla http://instagr.am/xxxxx", exfee)
		assert.Equal(t, content, "@abc@domain.com blablabla {{url:http://instagr.am/xxxxx}}")
		assert.Contains(t, fmt.Sprintf("%v", relationships), "mention:identity://3")
		assert.Contains(t, fmt.Sprintf("%v", relationships), "url:http://instagr.am/xxxxx")
	}

	{
		content, relationships := conv.parseRelationship("@exfe blablabla", exfee)
		assert.Equal(t, content, "@exfe blablabla")
		assert.Contains(t, fmt.Sprintf("%v", relationships), "mention:identity://1")
	}
}
