package main

import (
	"broker"
	"database/sql"
	"fmt"
	"github.com/googollee/go-multiplexer"
	"model"
	"shorttoken"
	"time"
)

const (
	SHORTTOKEN_STORE           = "INSERT INTO `shorttokens` (`key`, `resource`, `data`, `expire_at`, `created_at`) VALUES (?, ?, ?, ?, ?)"
	SHORTTOKEN_FIND            = "SELECT key, resource, data, expire_at FROM `shorttokens` WHERE 1==1"
	SHORTTOKEN_UPDATE_DATA     = "UPDATE `shorttokens` SET data=? WHERE 1==1"
	SHORTTOKEN_UPDATE_EXPIREAT = "UPDATE `shorttokens` SET tokens.expire_at=? WHERE 1==1"
)

type ShortTokenRepository struct {
	db     *broker.DBMultiplexer
	config *model.Config
}

func NewShortTokenRepository(config *model.Config, db *broker.DBMultiplexer) (*ShortTokenRepository, error) {
	ret := &ShortTokenRepository{
		db:     db,
		config: config,
	}
	return ret, nil
}

// CREATE TABLE `tokens` (`id` SERIAL NOT NULL, `key` CHAR(32) NOT NULL, `rand` CHAR(32) NOT NULL, `created_at` DATETIME NOT NULL, `expire_at` DATETIME NOT NULL, `data` TEXT NOT NULL)

func (r *ShortTokenRepository) Store(token shorttoken.Token) error {
	var err error
	r.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)
		_, err = db.Exec(SHORTTOKEN_STORE, token.Key, token.Resource, token.Data, r.timeToString(token.ExpireAt), r.timeToString(token.CreatedAt))
	})
	return err
}

func (r *ShortTokenRepository) Find(key, resource string) (shorttoken.Token, bool, error) {
	query := SHORTTOKEN_FIND
	if key != "" {
		query = fmt.Sprintf("%s AND key='%s'", query, key)
	}
	if key != "" {
		query = fmt.Sprintf("%s AND resource='%s'", query, resource)
	}
	var err error
	isExist := true
	token := shorttoken.Token{}
	r.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)
		var rows *sql.Rows
		rows, err = db.Query(query, key)
		if err != nil {
			return
		}
		if !rows.Next() {
			isExist = false
			return
		}
		var expireAt string
		err := rows.Scan(&token.Key, &token.Resource, &token.Data, &expireAt)
		if err != nil {
			return
		}
		token.ExpireAt, _ = time.Parse("2006-01-02 15:04:05", expireAt)
	})
	return token, isExist, err
}

func (r *ShortTokenRepository) UpdateData(key, resource, data string) error {
	sql := SHORTTOKEN_UPDATE_DATA
	if key != "" {
		sql = fmt.Sprintf("%s AND key='%s'", sql, key)
	}
	if key != "" {
		sql = fmt.Sprintf("%s AND resource='%s'", sql, resource)
	}
	var err error
	r.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)
		_, err = db.Exec(sql, data)
	})
	return err
}

func (r *ShortTokenRepository) UpdateExpireAt(key, resource string, expireAt time.Time) error {
	sql := SHORTTOKEN_UPDATE_EXPIREAT
	if key != "" {
		sql = fmt.Sprintf("%s AND key='%s'", sql, key)
	}
	if key != "" {
		sql = fmt.Sprintf("%s AND resource='%s'", sql, resource)
	}
	var err error
	r.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)
		_, err = db.Exec(sql, r.timeToString(expireAt))
	})
	return err
}

func (r *ShortTokenRepository) timeToString(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
