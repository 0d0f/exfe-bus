package main

import (
	"broker"
	"database/sql"
	"fmt"
	"github.com/googollee/go-multiplexer"
	"model"
	"time"
	"token"
)

const (
	TOKEN_STORE           = "INSERT INTO `tokens_` (`key`, `hash`, `data`, `expire_at`, `created_at`) VALUES (?, ?, ?, ?, ?)"
	TOKEN_FIND            = "SELECT `key`, hash, data, touched_at, expire_at FROM `tokens_` WHERE expire_at>UTC_TIMESTAMP()"
	TOKEN_UPDATE_DATA     = "UPDATE `tokens_` SET data=? WHERE expire_at>UTC_TIMESTAMP()"
	TOKEN_UPDATE_EXPIREAT = "UPDATE `tokens_` SET expire_at=? WHERE expire_at>UTC_TIMESTAMP()"
	TOKEN_TOUCH           = "UPDATE `tokens_` SET touched_at=NOW() WHERE expire_at>UTC_TIMESTAMP()"
)

func where(token token.Token) string {
	ret := ""
	if token.Key != "" {
		ret += fmt.Sprintf(" AND `key`=\"%s\"", token.Key)
	}
	if token.Hash != "" {
		ret += fmt.Sprintf(" AND `hash`=\"%s\"", token.Hash)
	}
	return ret
}

type TokenRepo struct {
	db     *broker.DBMultiplexer
	config *model.Config
}

func NewTokenRepo(config *model.Config, db *broker.DBMultiplexer) (*TokenRepo, error) {
	ret := &TokenRepo{
		db:     db,
		config: config,
	}
	return ret, nil
}

func (r *TokenRepo) Store(token token.Token) (err error) {
	e := r.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)
		_, err = db.Exec(TOKEN_STORE, token.Key, token.Hash, token.Data, r.timeToString(token.ExpireAt), r.timeToString(token.CreatedAt))
	})
	if e != nil {
		r.config.Log.Crit("sql error: %s", e)
		err = e
	}
	return
}

func (r *TokenRepo) Touch(token token.Token) (err error) {
	query := TOKEN_TOUCH + where(token)
	e := r.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)
		_, err = db.Exec(query)
	})
	if e != nil {
		r.config.Log.Crit("sql error: %s", e)
		err = e
	}
	return
}

func (r *TokenRepo) Find(token token.Token) (ret []token.Token, err error) {
	query := TOKEN_FIND + where(token)
	e := r.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)
		var rows *sql.Rows
		rows, err = db.Query(query)
		if err != nil {
			return
		}
		defer rows.Close()

		for rows.Next() {
			var touchedAt string
			var expireAt string
			err := rows.Scan(&token.Key, &token.Hash, &token.Data, &touchedAt, &expireAt)
			if err != nil {
				return
			}
			token.TouchedAt, _ = time.Parse("2006-01-02 15:04:05", touchedAt)
			token.ExpireAt, _ = time.Parse("2006-01-02 15:04:05", expireAt)
			ret = append(ret, token)
		}
	})
	if e != nil {
		r.config.Log.Crit("sql error: %s", e)
		err = e
	}
	return
}

func (r *TokenRepo) UpdateData(token token.Token, data string) (err error) {
	query := TOKEN_UPDATE_DATA + where(token)
	e := r.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)
		_, err = db.Exec(query, data)
	})
	if e != nil {
		r.config.Log.Crit("sql error: %s", e)
		err = e
	}
	return
}

func (r *TokenRepo) UpdateExpireAt(token token.Token, expireAt time.Time) (err error) {
	query := TOKEN_UPDATE_EXPIREAT + where(token)
	e := r.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)
		_, err = db.Exec(query, r.timeToString(expireAt))
	})
	if e != nil {
		r.config.Log.Crit("sql error: %s", e)
		err = e
	}
	return err
}

func (r *TokenRepo) timeToString(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
