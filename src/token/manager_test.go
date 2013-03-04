package token

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

func (r *TestTokenRepo) UpdateData(token Token, data string) error {
	tokens, err := r.Find(token)
	if err != nil {
		return fmt.Errorf("can't find")
	}
	for _, t := range tokens {
		t.Data = data
		r.store[t.Key] = t
	}
	return nil
}

func (r *TestTokenRepo) UpdateExpireAt(token Token, expireAt time.Time) error {
	tokens, err := r.Find(token)
	if err != nil {
		return fmt.Errorf("can't find")
	}
	for _, t := range tokens {
		t.ExpireAt = expireAt
		r.store[t.Key] = t
	}
	return nil
}

func (r *TestTokenRepo) Find(token Token) ([]Token, error) {
	ret := make([]Token, 0)
	for _, t := range r.store {
		if token.Key != "" && t.Key != token.Key {
			continue
		}
		if token.Hash != "" && t.Hash != token.Hash {
			continue
		}
		if time.Now().After(t.ExpireAt) {
			continue
		}
		ret = append(ret, t)
	}
	if len(ret) == 0 {
		return nil, nil
	}
	return ret, nil
}

func (r *TestTokenRepo) Touch(token Token) error {
	tokens, err := r.Find(token)
	if err != nil {
		return fmt.Errorf("can't find")
	}
	for _, t := range tokens {
		t.TouchedAt = time.Now()
		r.store[t.Key] = t
	}
	return nil
}

func TestShortToken(t *testing.T) {
	repo := &TestTokenRepo{
		store: make(map[string]Token),
	}
	mgr := New(repo)
	resource := "resource"
	tk := ""

	{
		_, err := mgr.Get(tk, resource)
		assert.NotEqual(t, err, nil)
	}

	token, err := mgr.Create("short", resource, "data", time.Second)
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

	token, err = mgr.Create("long", resource, "data", time.Second)
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
