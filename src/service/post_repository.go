package main

import (
	"broker"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-multiplexer"
	"gobus"
	"model"
	convmodel "model/conversation"
	"net/http"
	"time"
)

const (
	POST_CREATE = "INSERT INTO `posts` (by_id, created_at, relationship, content, via, exfee_id, ref_uri) VALUES (?, ?, ?, ?, ?, ?, ?)"
	POST_FIND   = "SELECT id, by_id, created_at, relationship, content, via, exfee_id, ref_uri FROM `posts` WHERE del=0 AND exfee_id=? AND ref_uri=?"
	POST_DELETE = "UPDATE `posts` SET posts.del=1 WHERE id=? AND ref_uri=?"
)

type PostRepository struct {
	db         *broker.DBMultiplexer
	redis      *broker.RedisMultiplexer
	dispatcher *gobus.Dispatcher
	config     *model.Config
}

func NewPostRepository(config *model.Config, db *broker.DBMultiplexer, redis *broker.RedisMultiplexer, dispatcher *gobus.Dispatcher) (*PostRepository, error) {
	ret := &PostRepository{
		db:         db,
		redis:      redis,
		dispatcher: dispatcher,
		config:     config,
	}
	return ret, nil
}

func (r *PostRepository) SavePost(post convmodel.Post) (uint64, error) {
	createdAt := post.CreatedAt.UTC().Format("2006-01-02 15:04:05")
	relationship, err := json.Marshal(post.Relationship)
	if err != nil {
		return 0, fmt.Errorf("can't marshal relationship: %s", err)
	}

	var id int64
	r.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)
		var ret sql.Result
		ret, err = db.Exec(POST_CREATE, post.By.ID, createdAt, relationship, post.Content, post.Via, post.ExfeeID, post.RefURI)
		if err != nil {
			return
		}
		id, err = ret.LastInsertId()
	})
	return uint64(id), err
}

func (r *PostRepository) FindPosts(exfeeID uint64, refURI, sinceTime, untilTime string, minID, maxID uint64) ([]convmodel.Post, error) {
	query := POST_FIND
	if sinceTime != "" {
		query = fmt.Sprintf("%s AND created_at>='%s'", query, sinceTime)
	}
	if untilTime != "" {
		query = fmt.Sprintf("%s AND created_at<='%s'", query, untilTime)
	}
	if minID != 0 {
		query = fmt.Sprintf("%s AND id>=%d", query, minID)
	}
	if maxID != 0 {
		query = fmt.Sprintf("%s AND id<=%d", query, maxID)
	}
	query = fmt.Sprintf("%s ORDER BY created_at ASC", query)

	var err error
	ret := make([]convmodel.Post, 0)
	r.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)
		var rows *sql.Rows
		rows, err = db.Query(query, exfeeID, refURI)
		if err != nil {
			return
		}
		defer rows.Close()

		for rows.Next() {
			post := convmodel.Post{}
			var createdAt string
			var relationship []byte
			err = rows.Scan(&post.ID, &post.By.ID, &createdAt, &relationship, &post.Content, &post.Via, &post.ExfeeID, &post.RefURI)
			if err != nil {
				return
			}
			post.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
			err = json.Unmarshal(relationship, &post.Relationship)
			if err != nil {
				return
			}
			post.By, err = r.FindIdentity(post.By)
			if err != nil {
				return
			}
			ret = append(ret, post)
		}
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *PostRepository) DeletePost(refID string, postID uint64) error {
	var err error
	r.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)
		_, err = db.Exec(POST_DELETE, postID, refID)
	})
	return err
}

func (r *PostRepository) SetUnreadCount(uri string, userID int64, count int) error {
	key := fmt.Sprintf("unreadpost:u%d:%s", userID, uri)
	return r.redis.Set(key, count)
}

func (r *PostRepository) AddUnreadCount(uri string, userID int64, count int) error {
	key := fmt.Sprintf("unreadpost:u%d:%s", userID, uri)
	_, err := r.redis.Incrby(key, int64(count))
	return err
}

func (r *PostRepository) GetUnreadCount(uri string, userID int64) (int, error) {
	key := fmt.Sprintf("unreadpost:u%d:%s", userID, uri)
	ret, err := r.redis.Get(key)
	if err != nil {
		return 0, err
	}
	return int(ret.Int64()), err
}

func (r *PostRepository) FindIdentity(identity model.Identity) (model.Identity, error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(identity)
	if err != nil {
		return identity, err
	}
	resp, err := http.Post(fmt.Sprintf("%s/v2/Gobus/RevokeIdentity", r.config.SiteApi), "application/json", buf)
	if err != nil {
		return identity, err
	}
	if resp.StatusCode != 200 {
		return identity, fmt.Errorf("find identity failed: %s", resp.Status)
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&identity)
	if err != nil {
		return identity, err
	}
	return identity, nil
}

func (r *PostRepository) FindCross(id uint64) (model.Cross, error) {
	var ret model.Cross
	resp, err := http.Get(fmt.Sprintf("%s/v2/Gobus/GetCrossById?id=%d", r.config.SiteApi, id))
	if err != nil {
		return ret, err
	}
	if resp.StatusCode != 200 {
		return ret, fmt.Errorf("find cross failed: %s", resp.Status)
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func (r *PostRepository) SendUpdate(tos []model.Recipient, cross model.Cross, post model.Post) error {
	arg := make(map[string]interface{})
	arg["service"] = "Conversation"
	arg["method"] = "Update"
	arg["merge_key"] = fmt.Sprintf("c%d", cross.ID)
	arg["tos"] = tos
	arg["data"] = model.ConversationUpdate{
		Cross: cross,
		Post:  post,
	}

	var i int
	err := r.dispatcher.Do("bus://exfe_queue/Instance", "Push", arg, &i)
	if err != nil {
		return err
	}
	return nil
}
