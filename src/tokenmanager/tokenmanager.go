package tokenmanager

import (
	"fmt"
	"time"
)

type TokenRepository interface {
	Create(token *Token) error
	Store(token *Token) error
	FindByKey(key string) ([]*Token, error)
	FindByToken(key, rand string) (*Token, error)
	Delete(token *Token) error
}

const TOKEN_KEY_LENGTH = 32

const (
	CREATE          = "CREATE TABLE `%s` (`id` SERIAL NOT NULL, `token` CHAR(64) NOT NULL, `created_at` DATETIME NOT NULL, `expire_at` DATETIME NOT NULL, `resource` TEXT NOT NULL, `data` TEXT NOT NULL)"
	INSERT          = "INSERT INTO `%s` VALUES (null, '%%s', '%%s', '%%s', '%%s', '%%s')"
	SELECT_TOKEN    = "SELECT expire_at, resource, data FROM `%s` WHERE token='%%s'"
	SELECT_RESOURCE = "SELECT token, expire_at FROM `%s` WHERE resource='%%s'"
	UPDATE_EXPIRE   = "UPDATE `%s` SET expire_at='%%s' WHERE token='%%s'"
	UPDATE_DATA     = "UPDATE `%s` SET data='%%s' WHERE token='%%s'"
	DELETE          = "DELETE FROM `%s` WHERE token='%%s'"
)

var NeverExpire = time.Duration(-1)

type TokenManager struct {
	repo TokenRepository
}

func New(repo TokenRepository) *TokenManager {
	return &TokenManager{
		repo: repo,
	}
}

func (m *TokenManager) GenerateToken(resource, data string, expireAfterSecond time.Duration) (*Token, error) {
	expireAt := time.Now().Add(expireAfterSecond)
	token := NewToken(resource, data, &expireAt)
	if expireAfterSecond < 0 {
		token.ExpireAt = nil
	}

	err := m.repo.Create(token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (m *TokenManager) GetToken(token string) (*Token, error) {
	tk, err := m.repo.FindByToken(token[:TOKEN_KEY_LENGTH], token[TOKEN_KEY_LENGTH:])
	if err != nil {
		return nil, err
	}
	if tk == nil {
		return nil, fmt.Errorf("no token found")
	}

	return tk, nil
}

func (m *TokenManager) FindTokens(resource string) (tokens []*Token, err error) {
	md5 := md5Resource(resource)
	tokens, err = m.repo.FindByKey(md5[:])
	return
}

func (m *TokenManager) UpdateData(token, data string) error {
	t, err := m.GetToken(token)
	if err != nil {
		return err
	}
	t.Data = data
	return m.repo.Store(t)
}

func (m *TokenManager) VerifyToken(token, resource string) (bool, *Token, error) {
	t, err := m.GetToken(token)
	if err != nil {
		return false, nil, err
	}
	key := md5Resource(resource)
	return t.Key == key, t, nil
}

func (m *TokenManager) DeleteToken(token string) error {
	t, err := m.GetToken(token)
	if err != nil {
		return err
	}
	return m.repo.Delete(t)
}

func (m *TokenManager) RefreshToken(token string, duration time.Duration) error {
	t, err := m.GetToken(token)
	if err != nil {
		return err
	}
	if duration < 0 {
		t.ExpireAt = nil
	} else {
		expireAt := time.Now().Add(duration)
		t.ExpireAt = &expireAt
	}
	return m.repo.Store(t)
}

func (m *TokenManager) ExpireToken(token string) error {
	return m.RefreshToken(token, 0)
}
