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
	UpdateDataByToken(key, rand, data string) error
	UpdateExpireAtByToken(key, rand string, expireAt *time.Time) error
	UpdateExpireAtByKey(key string, expireAt *time.Time) error
	DeleteByToken(key, rand string) error
}

const TOKEN_KEY_LENGTH = 32

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
	tk, err := m.repo.FindByToken(m.splitToken(token))
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
	key, rand := m.splitToken(token)
	return m.repo.UpdateDataByToken(key, rand, data)
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
	return m.repo.DeleteByToken(m.splitToken(token))
}

func (m *TokenManager) RefreshToken(token string, duration time.Duration) error {
	var expireAt *time.Time
	if duration < 0 {
		expireAt = nil
	} else {
		t := time.Now().Add(duration)
		expireAt = &t
	}
	key, rand := m.splitToken(token)
	return m.repo.UpdateExpireAtByToken(key, rand, expireAt)
}

func (m *TokenManager) ExpireToken(token string) error {
	return m.RefreshToken(token, 0)
}

func (m *TokenManager) ExpireTokensByKey(key string) error {
	t := time.Now()
	return m.repo.UpdateExpireAtByKey(key, &t)
}

func (m *TokenManager) splitToken(token string) (key, rand string) {
	return token[:TOKEN_KEY_LENGTH], token[TOKEN_KEY_LENGTH:]
}
