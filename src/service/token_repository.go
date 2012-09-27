package main

import (
	"model"
	"strings"
	"time"
	"tokenmanager"
)

const (
	CREATE                   = "INSERT INTO `%%table%%` VALUES (null, ?, ?, ?, ?, ?)"
	STORE                    = "UPDATE `%%table%%` SET expire_at=?, data=? WHERE %%table%%.key=? AND %%table%%.rand=?"
	FIND_BY_KEY              = "SELECT rand, created_at, expire_at, data FROM `%%table%%` WHERE %%table%%.key=?"
	FIND_BY_TOKEN            = "SELECT created_at, expire_at, data FROM `%%table%%` WHERE %%table%%.key=? AND %%table%%.rand=?"
	UPDATE_DATA_BY_TOKEN     = "UPDATE `%%table%%` SET %%table%%.data=? WHERE %%table%%.key=? AND %%table%%.rand=?"
	UPDATE_EXPIREAT_BY_TOKEN = "UPDATE `%%table%%` SET %%table%%.expire_at=? WHERE %%table%%.key=? AND %%table%%.rand=?"
	UPDATE_EXPIREAT_BY_KEY   = "UPDATE `%%table%%` SET %%table%%.expire_at=? WHERE %%table%%.key=?"
	DELETE_BY_TOKEN          = "DELETE FROM `%%table%%` WHERE %%table%%.key=? AND %%table%%.rand=?"
)

type TokenRepository struct {
	DBRepository
	create                string
	store                 string
	findByKey             string
	findByToken           string
	updateDataByToken     string
	updateExpireAtByToken string
	updateExpireAtByKey   string
	deleteByToken         string
}

func NewTokenRepository(config *model.Config) (*TokenRepository, error) {
	ret := &TokenRepository{
		create:                strings.Replace(CREATE, `%%table%%`, config.TokenManager.TableName, -1),
		store:                 strings.Replace(STORE, `%%table%%`, config.TokenManager.TableName, -1),
		findByKey:             strings.Replace(FIND_BY_KEY, `%%table%%`, config.TokenManager.TableName, -1),
		findByToken:           strings.Replace(FIND_BY_TOKEN, `%%table%%`, config.TokenManager.TableName, -1),
		updateDataByToken:     strings.Replace(UPDATE_DATA_BY_TOKEN, `%%table%%`, config.TokenManager.TableName, -1),
		updateExpireAtByToken: strings.Replace(UPDATE_EXPIREAT_BY_TOKEN, `%%table%%`, config.TokenManager.TableName, -1),
		updateExpireAtByKey:   strings.Replace(UPDATE_EXPIREAT_BY_KEY, `%%table%%`, config.TokenManager.TableName, -1),
		deleteByToken:         strings.Replace(DELETE_BY_TOKEN, `%%table%%`, config.TokenManager.TableName, -1),
	}
	ret.Config = config
	err := ret.Connect()
	return ret, err
}

// CREATE TABLE `%%table%%` (`id` SERIAL NOT NULL, `key` CHAR(32) NOT NULL, `rand` CHAR(32) NOT NULL, `created_at` DATETIME NOT NULL, `expire_at` DATETIME NOT NULL, `data` TEXT NOT NULL)

func (r *TokenRepository) Create(token *tokenmanager.Token) error {
	_, err := r.Exec(r.create, token.Key, token.Rand, r.timeToString(&token.CreatedAt), r.timeToString(token.ExpireAt), token.Data)
	return err
}

func (r *TokenRepository) Store(token *tokenmanager.Token) error {
	_, err := r.Exec(r.store, r.timeToString(token.ExpireAt), token.Data, token.Key, token.Rand)
	return err
}

func (r *TokenRepository) FindByKey(key string) ([]*tokenmanager.Token, error) {
	rows, err := r.Query(r.findByKey, key)
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
	rows, err := r.Query(r.findByToken, key, rand)
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
	_, err := r.Exec(r.updateDataByToken, data, key, rand)
	return err
}

func (r *TokenRepository) UpdateExpireAtByToken(key, rand string, expireAt *time.Time) error {
	_, err := r.Exec(r.updateExpireAtByToken, r.timeToString(expireAt), key, rand)
	return err
}

func (r *TokenRepository) UpdateExpireAtByKey(key string, expireAt *time.Time) error {
	_, err := r.Exec(r.updateExpireAtByKey, r.timeToString(expireAt), key)
	return err
}

func (r *TokenRepository) DeleteByToken(key, rand string) error {
	_, err := r.Exec(r.deleteByToken, key, rand)
	return err
}

func (r *TokenRepository) timeToString(t *time.Time) string {
	if t == nil {
		return "0000-00-00 00:00:00"
	}
	return t.UTC().Format(time.RFC3339)
}
