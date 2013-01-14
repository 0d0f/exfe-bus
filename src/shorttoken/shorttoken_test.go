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

func (r *TestTokenRepo) Find(key string, resource string) ([]Token, error) {
	ret := make([]Token, 0)
	for _, token := range r.store {
		if key != "" && token.Key != key {
			continue
		}
		if resource != "" && token.Resource != resource {
			continue
		}
		if time.Now().After(token.ExpireAt) {
			continue
		}
		ret = append(ret, token)
	}
	if len(ret) == 0 {
		return nil, nil
	}
	return ret, nil
}

func (r *TestTokenRepo) Touch(key, resource string) error {
	for _, token := range r.store {
		if token.Key == key && token.Resource == resource {
			fmt.Println(token.TouchedAt)
			token.TouchedAt = time.Now()
			r.store[token.Key] = token
			fmt.Println(token.TouchedAt)
		}
	}
	return nil
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
		assert.Equal(t, len(token), 1)
		assert.Equal(t, token[0].Key, tk)
		assert.Equal(t, token[0].Data, "data")
	}

	{
		token, err := mgr.Get("", resource)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(token), 1)
		assert.Equal(t, token[0].Key, tk)
		assert.Equal(t, token[0].Data, "data")
	}

	{
		token, err := mgr.Get(tk, resource)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(token), 1)
		assert.Equal(t, token[0].Key, tk)
		assert.Equal(t, token[0].Data, "data")
	}

	{
		err = mgr.UpdateData(tk, "abc")
		assert.Equal(t, err, nil)
	}

	{
		token, err := mgr.Get(tk, resource)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(token), 1)
		assert.Equal(t, token[0].Key, tk)
		assert.Equal(t, token[0].Data, "abc")
	}

	{
		err := mgr.Refresh(tk, resource, 3*time.Second)
		assert.Equal(t, err, nil)

		fmt.Println("touch")
		mgr.Get(tk, resource)

		time.Sleep(2 * time.Second)

		token, err := mgr.Get(tk, resource)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(token), 1)
		touch1 := token[0].TouchedAt

		token, err = mgr.Get(tk, resource)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(token), 1)
		touch2 := token[0].TouchedAt

		fmt.Println(touch1, touch2)
		assert.NotEqual(t, touch1, touch2)
		fmt.Println("touched")
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
		assert.Equal(t, len(token), 1)
		assert.Equal(t, token[0].Key, tk)
		assert.Equal(t, token[0].Data, "data")
	}

	err = mgr.Refresh(tk, "", 0)
	assert.Equal(t, err, nil)

	{
		_, err := mgr.Get(tk, "")
		assert.NotEqual(t, err, nil)
	}
}
