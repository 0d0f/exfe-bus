package shorttoken

import (
	"fmt"
	"github.com/stretchrcom/testify/assert"
	"testing"
	"time"
)

type TestTokenRepo struct {
	store map[string]Token
}

func (r *TestTokenRepo) Store(token Token) error {
	r.store[token.Key] = token
	return nil
}

func (r *TestTokenRepo) UpdateData(key, resource, data string) error {
	for k, token := range r.store {
		if token.Key == key {
			token.Data = data
			r.store[k] = token
			return nil
		}
		if token.Resource == resource {
			token.Data = data
			r.store[k] = token
			return nil
		}
	}
	return fmt.Errorf("can't find")
}

func (r *TestTokenRepo) UpdateExpireAt(key, resource string, expireAt time.Time) error {
	for k, token := range r.store {
		if token.Key == key {
			token.ExpireAt = expireAt
			r.store[k] = token
			return nil
		}
		if token.Resource == resource {
			token.ExpireAt = expireAt
			r.store[k] = token
			return nil
		}
	}
	return fmt.Errorf("can't find")
}

func (r *TestTokenRepo) Find(key string, resource string) (Token, bool, error) {
	for _, token := range r.store {
		find := true
		if key != "" && token.Key != key {
			find = false
		}
		if resource != "" && token.Resource != resource {
			find = false
		}
		if find {
			if time.Now().After(token.ExpireAt) {
				return Token{}, false, nil
			}
			return token, true, nil
		}
	}
	return Token{}, false, nil
}

func TestShortToken(t *testing.T) {
	repo := &TestTokenRepo{
		store: make(map[string]Token),
	}
	mgr := New(repo, 4)
	resource := "resource"
	tk := "0432"

	{
		_, err := mgr.Get(tk, resource)
		assert.NotEqual(t, err, nil)
	}

	token, err := mgr.Create(resource, "data", time.Second)
	assert.Equal(t, err, nil)
	tk = token.Key

	{
		token, err := mgr.Get(tk, "")
		assert.Equal(t, err, nil)
		assert.Equal(t, token.Key, tk)
		assert.Equal(t, token.Data, "data")
	}

	{
		token, err := mgr.Get("", resource)
		assert.Equal(t, err, nil)
		assert.Equal(t, token.Key, tk)
		assert.Equal(t, token.Data, "data")
	}

	{
		token, err := mgr.Get(tk, resource)
		assert.Equal(t, err, nil)
		assert.Equal(t, token.Key, tk)
		assert.Equal(t, token.Data, "data")
	}

	{
		err = mgr.UpdateData(tk, "abc")
		assert.Equal(t, err, nil)
	}

	{
		token, err := mgr.Get(tk, resource)
		assert.Equal(t, err, nil)
		assert.Equal(t, token.Key, tk)
		assert.Equal(t, token.Data, "abc")
	}

	time.Sleep(time.Second * 2)

	{
		_, err := mgr.Get(tk, "")
		assert.NotEqual(t, err, nil)
	}

	token, err = mgr.Create(resource, "data", time.Second)
	assert.Equal(t, err, nil)
	tk = token.Key

	err = mgr.Refresh(tk, "", time.Second)
	assert.Equal(t, err, nil)

	{
		token, err := mgr.Get(tk, "")
		assert.Equal(t, err, nil)
		assert.Equal(t, token.Key, tk)
		assert.Equal(t, token.Data, "data")
	}

	err = mgr.Refresh(tk, "", 0)
	assert.Equal(t, err, nil)

	{
		_, err := mgr.Get(tk, "")
		assert.NotEqual(t, err, nil)
	}
}
