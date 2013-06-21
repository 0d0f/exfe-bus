package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"model"
	"token"
)

const (
	STORE          = "INSERT INTO `tokens` (`key`, `hash`, `user_id`, `scopes`, `client`, `data`, `expires_in`, `created_at`, `touched_at`) VALUES (?, ?, ?, ?, ?, ?, ?, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())"
	FIND_BY_KEY    = "SELECT `key`, hash, user_id, scopes, client, data, touched_at, expires_in FROM `tokens` WHERE expires_in>UNIX_TIMESTAMP() AND `key`=?"
	FIND_BY_HASH   = "SELECT `key`, hash, user_id, scopes, client, data, touched_at, expires_in FROM `tokens` WHERE expires_in>UNIX_TIMESTAMP() AND `hash`=?"
	TOUCH          = "UPDATE `tokens` SET touched_at=UNIX_TIMESTAMP() WHERE expires_in>UNIX_TIMESTAMP()"
	UPDATE_BY_KEY  = "UPDATE `tokens` SET %s WHERE expires_in>UNIX_TIMESTAMP() AND `key`=?"
	UPDATE_BY_HASH = "UPDATE `tokens` SET %s WHERE expires_in>UNIX_TIMESTAMP() AND `hash`=?"
	DELETE         = "DELETE FROM `tokens` WHERE expires_in<=UNIX_TIMESTAMP()"
)

func update(data *string, expiresIn *int64) string {
	ret := ""
	if data != nil {
		ret += fmt.Sprintf("`data`='%s'", *data)
	}
	if expiresIn != nil {
		if ret != "" {
			ret += ", "
		}
		ret += fmt.Sprintf("`expires_in`=%d", *expiresIn)
	}
	return ret
}

func where(key, hash *string) string {
	ret := ""
	if key != nil {
		ret += fmt.Sprintf(" AND `key`=\"%s\"", *key)
	}
	if hash != nil {
		ret += fmt.Sprintf(" AND `hash`=\"%s\"", *hash)
	}
	return ret
}

type TokenRepo struct {
	db *sql.DB
}

func NewTokenRepo(config *model.Config, db *sql.DB) (*TokenRepo, error) {
	ret := &TokenRepo{
		db: db,
	}
	return ret, nil
}

func (r *TokenRepo) Store(token token.Token) error {
	_, err := r.db.Exec(STORE, token.Key, token.Hash, token.UserId, token.Scopes, token.Client, token.Data, token.ExpiresIn)
	if err != nil {
		return err
	}
	r.db.Exec(DELETE)
	return err
}

func (r *TokenRepo) FindByKey(key string) ([]token.Token, error) {
	rows, err := r.db.Query(FIND_BY_KEY, key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []token.Token
	var token token.Token
	for rows.Next() {
		err := rows.Scan(&token.Key, &token.Hash, &token.UserId, &token.Scopes, &token.Client, &token.Data, &token.TouchedAt, &token.ExpiresIn)
		if err != nil {
			return ret, err
		}
		ret = append(ret, token)
	}
	return ret, nil
}

func (r *TokenRepo) FindByHash(hash string) ([]token.Token, error) {
	rows, err := r.db.Query(FIND_BY_HASH, hash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []token.Token
	var token token.Token
	for rows.Next() {
		err := rows.Scan(&token.Key, &token.Hash, &token.UserId, &token.Scopes, &token.Client, &token.Data, &token.TouchedAt, &token.ExpiresIn)
		if err != nil {
			return ret, err
		}
		ret = append(ret, token)
	}
	return ret, nil
}

func (r *TokenRepo) Touch(key, hash *string) error {
	sql := TOUCH + where(key, hash)
	_, err := r.db.Exec(sql)
	return err
}

func (r *TokenRepo) UpdateByKey(key string, data *string, expiresIn *int64) (int64, error) {
	set := update(data, expiresIn)
	if set == "" {
		return 0, nil
	}
	sql := fmt.Sprintf(UPDATE_BY_KEY, set)
	result, err := r.db.Exec(sql, key)
	if err != nil {
		return 0, err
	}
	ret, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func (r *TokenRepo) UpdateByHash(hash string, data *string, expiresIn *int64) (int64, error) {
	set := update(data, expiresIn)
	if set == "" {
		return 0, nil
	}
	sql := fmt.Sprintf(UPDATE_BY_HASH, set)
	result, err := r.db.Exec(sql, hash)
	if err != nil {
		return 0, err
	}
	ret, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return ret, nil
}
