package main

import (
	"broker"
	"database/sql"
	"github.com/googollee/go-multiplexer"
	"model"
	"time"
	"tokenmanager"
)

const (
	CREATE                   = "INSERT INTO `tokens_` VALUES (null, ?, ?, ?, ?, ?, ?)"
	STORE                    = "UPDATE `tokens_` SET expire_at=?, data=? WHERE tokens_.key=? AND tokens_.rand=?"
	FIND_BY_KEY              = "SELECT rand, touched_at, created_at, expire_at, data FROM `tokens_` WHERE tokens_.key=?"
	FIND_BY_TOKEN            = "SELECT touched_at, created_at, expire_at, data FROM `tokens_` WHERE tokens_.key=? AND tokens_.rand=?"
	UPDATE_DATA_BY_TOKEN     = "UPDATE `tokens_` SET tokens_.data=? WHERE tokens_.key=? AND tokens_.rand=?"
	UPDATE_EXPIREAT_BY_TOKEN = "UPDATE `tokens_` SET tokens_.expire_at=? WHERE tokens_.key=? AND tokens_.rand=?"
	UPDATE_EXPIREAT_BY_KEY   = "UPDATE `tokens_` SET tokens_.expire_at=? WHERE tokens_.key=?"
	DELETE_BY_TOKEN          = "DELETE FROM `tokens_` WHERE tokens_.key=? AND tokens_.rand=?"
	TOUCH                    = "UPDATE `tokens_` SET touched_at=NOW() WHERE tokens_.key=? AND tokens_.rand=?"
)

type TokenRepository struct {
	db     *broker.DBMultiplexer
	config *model.Config
}

func NewTokenRepository(config *model.Config, db *broker.DBMultiplexer) (*TokenRepository, error) {
	ret := &TokenRepository{
		db:     db,
		config: config,
	}
	return ret, nil
}

// CREATE TABLE `tokens_` (`id` SERIAL NOT NULL, `key` CHAR(32) NOT NULL, `rand` CHAR(32) NOT NULL, `created_at` DATETIME NOT NULL, `expire_at` DATETIME NOT NULL, `data` TEXT NOT NULL)

func (r *TokenRepository) Create(token *tokenmanager.Token) error {
	var err error
	r.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)
		_, err = db.Exec(CREATE, token.Key, token.Rand, r.timeToString(&token.TouchedAt), r.timeToString(&token.CreatedAt), r.timeToString(token.ExpireAt), token.Data)
	})
	return err
}

func (r *TokenRepository) Store(token *tokenmanager.Token) error {
	var err error
	r.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)
		_, err = db.Exec(STORE, r.timeToString(token.ExpireAt), token.Data, token.Key, token.Rand)
	})
	return err
}

func (r *TokenRepository) FindByKey(key string) ([]*tokenmanager.Token, error) {
	var err error
	ret := make([]*tokenmanager.Token, 0, 0)
	r.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)
		var rows *sql.Rows
		rows, err = db.Query(FIND_BY_KEY, key)
		if err != nil {
			return
		}
		defer rows.Close()

		for rows.Next() {
			var touchedAtStr string
			var createdAtStr string
			var expireAtStr string
			token := tokenmanager.Token{
				Key: key,
			}
			err := rows.Scan(&token.Rand, &touchedAtStr, &createdAtStr, &expireAtStr, &token.Data)
			if err != nil {
				return
			}
			token.TouchedAt, _ = time.Parse("2006-01-02 15:04:05", touchedAtStr)
			token.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
			if expireAtStr != "0000-00-00 00:00:00" {
				time, _ := time.Parse("2006-01-02 15:04:05", expireAtStr)
				token.ExpireAt = &time
			}
			ret = append(ret, &token)
		}
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *TokenRepository) FindByToken(key, rand string) (*tokenmanager.Token, error) {
	var err error
	ret := &tokenmanager.Token{
		Key:  key,
		Rand: rand,
	}
	r.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)
		var rows *sql.Rows
		rows, err = db.Query(FIND_BY_TOKEN, key, rand)
		if err != nil {
			return
		}
		defer rows.Close()

		if !rows.Next() {
			ret = nil
			return
		}

		var touchedAtStr string
		var createdAtStr string
		var expireAtStr string
		err = rows.Scan(&touchedAtStr, &createdAtStr, &expireAtStr, &ret.Data)
		if err != nil {
			return
		}
		ret.TouchedAt, _ = time.Parse("2006-01-02 15:04:05", touchedAtStr)
		ret.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if expireAtStr != "0000-00-00 00:00:00" {
			expireAt, _ := time.Parse("2006-01-02 15:04:05", expireAtStr)
			ret.ExpireAt = &expireAt
		}

		db.Exec(TOUCH, key, rand)
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *TokenRepository) UpdateDataByToken(key, rand, data string) error {
	var err error
	r.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)
		_, err = db.Exec(UPDATE_DATA_BY_TOKEN, data, key, rand)
	})
	return err
}

func (r *TokenRepository) UpdateExpireAtByToken(key, rand string, expireAt *time.Time) error {
	var err error
	r.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)
		_, err = db.Exec(UPDATE_EXPIREAT_BY_TOKEN, r.timeToString(expireAt), key, rand)
	})
	return err
}

func (r *TokenRepository) UpdateExpireAtByKey(key string, expireAt *time.Time) error {
	var err error
	r.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)
		_, err = db.Exec(UPDATE_EXPIREAT_BY_KEY, r.timeToString(expireAt), key)
	})
	return err
}

func (r *TokenRepository) DeleteByToken(key, rand string) error {
	var err error
	r.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)
		_, err = db.Exec(DELETE_BY_TOKEN, key, rand)
	})
	return err
}

func (r *TokenRepository) timeToString(t *time.Time) string {
	if t == nil {
		return "0000-00-00 00:00:00"
	}
	return t.UTC().Format(time.RFC3339)
}
