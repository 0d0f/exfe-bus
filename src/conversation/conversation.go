package conversation

import (
	"fmt"
	"model"
	convmodel "model/conversation"
	"regexp"
	"strings"
	"time"
)

type Repo interface {
	FindIdentity(identity model.Identity) (model.Identity, error)
	FindCross(id uint64) (model.Cross, error)
	SendUpdate(tos []model.Recipient, cross model.Cross, post model.Post) error

	SavePost(post convmodel.Post) (uint64, error)
	FindPosts(exfeeID uint64, refURI, sinceTime, untilTime string, minID, maxID uint64) ([]convmodel.Post, error)
	DeletePost(refID string, postID uint64) error

	SetUnreadCount(uri string, userID int64, count int) error
	AddUnreadCount(uri string, userID int64, count int) error
	GetUnreadCount(uri string, userID int64) (int, error)
}

type Conversation struct {
	repo       Repo
	mentionRe  *regexp.Regexp
	urlRe      *regexp.Regexp
	relationRe *regexp.Regexp
}

func New(repo Repo) *Conversation {
	return &Conversation{
		repo:       repo,
		mentionRe:  regexp.MustCompile(`@([^@ ]*)(@[a-zA-Z0-9_.]*)?(@[a-zA-Z0-9_]*)?`),
		urlRe:      regexp.MustCompile(`(http|https)://[a-zA-Z0-9%!\.#_/+\-\\]*`),
		relationRe: regexp.MustCompile(`{{[a-zA-Z0-9_]:.*?}}`),
	}
}

func (c *Conversation) NewPost(crossID uint64, post model.Post, via string, createdAt int64) (model.Post, error) {
	cross, err := c.repo.FindCross(crossID)
	if err != nil {
		return post, err
	}
	postIdentity, err := c.repo.FindIdentity(post.By)
	if err != nil {
		return post, err
	}
	by, err := cross.Exfee.FindInvitedUser(postIdentity)
	if err != nil {
		return post, err
	}
	if createdAt <= 0 {
		createdAt = time.Now().Unix()
	}
	t := time.Unix(createdAt, 0)
	content, relationship := c.parseRelationship(post.Content, cross.Exfee)

	p := convmodel.Post{
		Meta: model.Meta{
			CreatedAt:    t,
			By:           by.Identity,
			Relationship: relationship,
		},
		Content: content,
		Via:     via,
		ExfeeID: cross.Exfee.ID,
		RefURI:  fmt.Sprintf("cross://%d", cross.ID),
	}
	id, err := c.repo.SavePost(p)
	if err != nil {
		return post, err
	}
	p.ID = id
	p.URI = fmt.Sprintf("post://%d", id)
	ret := p.ToPost()
	c.repo.AddUnreadCount(p.RefURI, by.Identity.UserID, 1)
	return ret, nil
}

func (c *Conversation) FindPosts(crossID uint64, clearUserID int64, sinceTime, untilTime string, minID, maxID uint64) ([]model.Post, error) {
	refID := fmt.Sprintf("cross://%d", crossID)
	cross, err := c.repo.FindCross(crossID)
	if err != nil {
		return nil, err
	}
	posts, err := c.repo.FindPosts(cross.Exfee.ID, refID, sinceTime, untilTime, minID, maxID)
	if err != nil {
		return nil, err
	}
	ret := make([]model.Post, 0)
	for _, p := range posts {
		ret = append(ret, p.ToPost())
	}
	if clearUserID != 0 {
		c.repo.SetUnreadCount(refID, clearUserID, 0)
	}
	return ret, nil
}

func (c *Conversation) DeletePost(crossID, postID uint64) error {
	refID := fmt.Sprintf("cross://%d", crossID)
	return c.repo.DeletePost(refID, postID)
}

func (c *Conversation) GetUnreadCount(crossID, userID int64) (int, error) {
	refID := fmt.Sprintf("cross://%d", crossID)
	return c.repo.GetUnreadCount(refID, userID)
}

func (c *Conversation) parseRelationship(content string, exfee model.Exfee) (string, []model.Relationship) {
	ret := make([]model.Relationship, 0)
	content = c.mentionRe.ReplaceAllStringFunc(content, func(c string) string {
		id := c[1:]
		atIndex := strings.LastIndex(id, "@")
		var identity model.Identity
		if atIndex < 0 {
			identity.ExternalUsername = id
		} else {
			ptIndex := strings.Index(id[atIndex:], ".")
			if ptIndex > 0 {
				// c is email address
				identity.ExternalUsername = id
				identity.ExternalID = id
				identity.Provider = "email"
			} else {
				identity.ExternalUsername = id[:atIndex]
				identity.Provider = id[atIndex+1:]
			}
		}
		inv, err := exfee.FindInvitedUser(identity)
		if err == nil {
			ret = append(ret, model.Relationship{
				Relation: "mention",
				URI:      fmt.Sprintf("identity://%d", inv.Identity.ID),
			})
		}
		return c
	})
	content = c.urlRe.ReplaceAllStringFunc(content, func(c string) string {
		ret = append(ret, model.Relationship{
			Relation: "url",
			URI:      c,
		})
		return fmt.Sprintf("{{url:%s}}", c)
	})
	content = c.relationRe.ReplaceAllStringFunc(content, func(c string) string {
		relation := strings.SplitN(c[2:len(c)-2], ":", 2)
		ret = append(ret, model.Relationship{
			Relation: relation[0],
			URI:      relation[1],
		})
		return c
	})

	return content, ret
}
