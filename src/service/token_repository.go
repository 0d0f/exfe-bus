package main

import (
	_ "code.google.com/p/go-mysql-driver/mysql"
	"database/sql"
	"fmt"
	"model"
	"time"
	"tokenmanager"
)

type TokenRepository struct {
	db     *sql.DB
	config *model.Config
}

func NewTokenRepository(config *model.Config) (*TokenRepository, error) {
	ret := &TokenRepository{
		config: config,
	}
	err := ret.Connect()
	return ret, err
}

func (r *TokenRepository) Connect() error {
	if r.db != nil {
		r.db.Close()
	}
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&keepalive=1",
		r.config.DB.Username, r.config.DB.Password, r.config.DB.Addr, r.config.DB.Port, r.config.DB.DbName))
	if err != nil {
		return err
	}
	_, err = db.Query("SELECT 1")
	if err != nil {
		db.Close()
		return err
	}
	r.db = db
	return nil
}

// CREATE TABLE `%s` (`id` SERIAL NOT NULL, `key` CHAR(32) NOT NULL, `rand` CHAR(32) NOT NULL, `created_at` DATETIME NOT NULL, `expire_at` DATETIME NOT NULL, `data` TEXT NOT NULL)

func (r *TokenRepository) Create(token *tokenmanager.Token) error {
	const INSERT = "INSERT INTO `%s` VALUES (null, ?, ?, ?, ?, ?)"
	sql := fmt.Sprintf(INSERT, r.config.TokenManager.TableName)
	_, err := r.db.Exec(sql, token.Key, token.Rand, r.timeToString(&token.CreatedAt), r.timeToString(token.ExpireAt), token.Data)
	return err
}

func (r *TokenRepository) Store(token *tokenmanager.Token) error {
	const UPDATE = "UPDATE `%s` SET expire_at=?, data=? WHERE %s.key=? AND %s.rand=?"
	sql := fmt.Sprintf(UPDATE, r.config.TokenManager.TableName, r.config.TokenManager.TableName, r.config.TokenManager.TableName)
	_, err := r.db.Exec(sql, r.timeToString(token.ExpireAt), token.Data, token.Key, token.Rand)
	return err
}

func (r *TokenRepository) FindByKey(key string) ([]*tokenmanager.Token, error) {
	const FIND = "SELECT rand, created_at, expire_at, data FROM `%s` WHERE %s.key=?"
	sql := fmt.Sprintf(FIND, r.config.TokenManager.TableName, r.config.TokenManager.TableName)
	rows, err := r.db.Query(sql, key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ret := make([]*tokenmanager.Token, 0, 0)
	for rows.Next() {
		var rand string
		var createdAtStr string
		var expireAtStr string
		var data string
		err := rows.Scan(&rand, &createdAtStr, &expireAtStr, &data)
		if err != nil {
			return ret, err
		}
		createdAt, _ := time.Parse("2006-01-02 15:04:05", createdAtStr)
		token := tokenmanager.Token{
			Key:       key,
			Rand:      rand,
			CreatedAt: createdAt,
			Data:      data,
		}
		if expireAtStr != "0000-00-00 00:00:00" {
			expireAt, _ := time.Parse("2006-01-02 15:04:05", expireAtStr)
			token.ExpireAt = &expireAt
		}
		ret = append(ret, &token)
	}
	return ret, nil
}

func (r *TokenRepository) FindByToken(key, rand string) (*tokenmanager.Token, error) {
	const FIND = "SELECT created_at, expire_at, data FROM `%s` WHERE %s.key=? AND %s.rand=?"
	sql := fmt.Sprintf(FIND, r.config.TokenManager.TableName, r.config.TokenManager.TableName, r.config.TokenManager.TableName)
	rows, err := r.db.Query(sql, key, rand)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var createdAtStr string
	var expireAtStr string
	var data string
	err = rows.Scan(&createdAtStr, &expireAtStr, &data)
	if err != nil {
		return nil, err
	}
	createdAt, _ := time.Parse("2006-01-02 15:04:05", createdAtStr)
	token := tokenmanager.Token{
		Key:       key,
		Rand:      rand,
		CreatedAt: createdAt,
		Data:      data,
	}
	if expireAtStr != "0000-00-00 00:00:00" {
		expireAt, _ := time.Parse("2006-01-02 15:04:05", expireAtStr)
		token.ExpireAt = &expireAt
	}
	return &token, nil
}

func (r *TokenRepository) UpdateDataByToken(key, rand, data string) error {
	const UPDATE = "UPDATE `%s` SET %s.data=? WHERE %s.key=? AND %s.rand=?"
	sql := fmt.Sprintf(UPDATE, r.config.TokenManager.TableName, r.config.TokenManager.TableName, r.config.TokenManager.TableName, r.config.TokenManager.TableName)
	_, err := r.db.Exec(sql, data, key, rand)
	return err
}

func (r *TokenRepository) UpdateExpireAtByToken(key, rand string, expireAt *time.Time) error {
	const UPDATE = "UPDATE `%s` SET %s.expire_at=? WHERE %s.key=? AND %s.rand=?"
	sql := fmt.Sprintf(UPDATE, r.config.TokenManager.TableName, r.config.TokenManager.TableName, r.config.TokenManager.TableName, r.config.TokenManager.TableName)
	_, err := r.db.Exec(sql, r.timeToString(expireAt), key, rand)
	return err
}

func (r *TokenRepository) UpdateExpireAtByKey(key string, expireAt *time.Time) error {
	const UPDATE = "UPDATE `%s` SET %s.expire_at=? WHERE %s.key=?"
	sql := fmt.Sprintf(UPDATE, r.config.TokenManager.TableName, r.config.TokenManager.TableName, r.config.TokenManager.TableName)
	_, err := r.db.Exec(sql, r.timeToString(expireAt), key)
	return err
}

func (r *TokenRepository) timeToString(t *time.Time) string {
	if t == nil {
		return "0000-00-00 00:00:00"
	}
	return t.UTC().Format(time.RFC3339)
}

func (r *TokenRepository) DeleteByToken(key, rand string) error {
	const DELETE = "DELETE FROM `%s` WHERE %s.key=? AND %s.rand=?"
	sql := fmt.Sprintf(DELETE, r.config.TokenManager.TableName, r.config.TokenManager.TableName, r.config.TokenManager.TableName)
	_, err := r.db.Exec(sql, key, rand)
	return err
}
