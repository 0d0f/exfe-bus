package token

import (
	"fmt"
	"github.com/googollee/go-rest"
	"github.com/stretchrcom/testify/assert"
	"net/http"
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

func (r *TestTokenRepo) FindByKey(key string) ([]Token, error) {
	ret, ok := r.store[key]
	if !ok {
		return nil, nil
	}
	if time.Now().Unix() >= ret.ExpiresIn {
		return nil, nil
	}
	return []Token{ret}, nil
}

func (r *TestTokenRepo) FindByHash(hash string) ([]Token, error) {
	fmt.Println(hash, r.store)
	var ret []Token
	now := time.Now().Unix()
	for _, t := range r.store {
		if t.Hash != hash {
			continue
		}
		if now >= t.ExpiresIn {
			continue
		}
		ret = append(ret, t)
	}
	return ret, nil
}

func (r *TestTokenRepo) UpdateByKey(key string, data *string, expiresIn *int64) (int64, error) {
	token, ok := r.store[key]
	if !ok {
		return 0, nil
	}
	if data != nil {
		token.Data = *data
	}
	if expiresIn != nil {
		token.ExpiresIn = *expiresIn
	}
	r.store[key] = token
	return 1, nil
}

func (r *TestTokenRepo) UpdateByHash(hash string, data *string, expiresIn *int64) (int64, error) {
	var i int64 = 0
	for k, token := range r.store {
		if token.Hash != hash {
			continue
		}
		i++
		if data != nil {
			token.Data = *data
		}
		if expiresIn != nil {
			token.ExpiresIn = *expiresIn
		}
		r.store[k] = token
	}
	return i, nil
}

func (r *TestTokenRepo) Touch(key, hash *string) error {
	if key != nil {
		token, ok := r.store[*key]
		if !ok {
			return nil
		}
		token.TouchedAt = time.Now().Unix()
		r.store[*key] = token
	}
	if hash != nil {
		for k, token := range r.store {
			if token.Hash != *hash {
				continue
			}
			token.TouchedAt = time.Now().Unix()
			r.store[k] = token
		}
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

	req, err := http.NewRequest("GET", "http://test/?type=short", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := rest.SetTest(mgr, nil, req)
	if err != nil {
		t.Fatal(err)
	}
	arg := CreateArg{}
	arg.Resource = resource
	arg.Data = "data"
	arg.ExpireAfterSeconds = 1
	token := mgr.HandleCreate(arg)
	assert.Equal(t, resp.Code, http.StatusOK)
	tk = token.Key

	{
		resp, _ := rest.SetTest(mgr, map[string]string{"key": tk}, nil)
		tokens := mgr.HandleKeyGet()
		assert.Equal(t, resp.Code, http.StatusOK)
		assert.Equal(t, len(tokens), 1)
		assert.Equal(t, tokens[0].Key, tk)
		assert.Equal(t, string(tokens[0].Data), "data")
	}

	{
		resp, _ := rest.SetTest(mgr, nil, nil)
		tokens := mgr.HandleResourceGet(resource)
		assert.Equal(t, resp.Code, http.StatusOK)
		assert.Equal(t, len(tokens), 1)
		assert.Equal(t, tokens[0].Key, tk)
		assert.Equal(t, string(tokens[0].Data), "data")
	}

	{
		resp, _ := rest.SetTest(mgr, map[string]string{"key": tk}, nil)
		data := "abc"
		mgr.HandleKeyUpdate(UpdateArg{
			Data: &data,
		})
		assert.Equal(t, resp.Code, http.StatusOK)
	}

	{
		resp, _ := rest.SetTest(mgr, map[string]string{"key": tk}, nil)
		tokens := mgr.HandleKeyGet()
		assert.Equal(t, resp.Code, http.StatusOK)
		assert.Equal(t, len(tokens), 1)
		assert.Equal(t, tokens[0].Key, tk)
		assert.Equal(t, string(tokens[0].Data), "abc")
	}

	{
		resp, _ := rest.SetTest(mgr, map[string]string{"key": tk}, nil)
		after := 3
		mgr.HandleKeyUpdate(UpdateArg{
			ExpireAfterSeconds: &after,
		})
		assert.Equal(t, resp.Code, http.StatusOK)

		fmt.Println("touch")
		mgr.HandleKeyGet()

		time.Sleep(2 * time.Second)

		token := mgr.HandleKeyGet()
		assert.Equal(t, resp.Code, http.StatusOK)
		assert.Equal(t, len(token), 1)
		touch1 := token[0].TouchedAt

		token = mgr.HandleKeyGet()
		assert.Equal(t, resp.Code, http.StatusOK)
		assert.Equal(t, len(token), 1)
		touch2 := token[0].TouchedAt

		fmt.Println(touch1, touch2)
		assert.NotEqual(t, touch1, touch2)
		fmt.Println("touched")
	}

	time.Sleep(time.Second * 2)

	{
		resp, _ := rest.SetTest(mgr, map[string]string{"key": tk}, nil)
		_ = mgr.HandleKeyGet()
		assert.NotEqual(t, resp.Code, http.StatusOK)
	}

	req, err = http.NewRequest("GET", "http://test/?type=long", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = rest.SetTest(mgr, nil, req)
	if err != nil {
		t.Fatal(err)
	}
	arg = CreateArg{}
	arg.Resource = resource
	arg.Data = "data"
	arg.ExpireAfterSeconds = 1
	token = mgr.HandleCreate(arg)
	assert.Equal(t, resp.Code, http.StatusOK)
	tk = token.Key

	after := 1
	resp, _ = rest.SetTest(mgr, map[string]string{"key": tk}, nil)
	mgr.HandleKeyUpdate(UpdateArg{
		ExpireAfterSeconds: &after,
	})
	assert.Equal(t, resp.Code, http.StatusOK)

	{
		tokens := mgr.HandleKeyGet()
		assert.Equal(t, resp.Code, http.StatusOK)
		assert.Equal(t, len(tokens), 1)
		assert.Equal(t, tokens[0].Key, tk)
		assert.Equal(t, string(tokens[0].Data), "data")
	}

	after = 0
	mgr.HandleKeyUpdate(UpdateArg{
		ExpireAfterSeconds: &after,
	})
	assert.Equal(t, resp.Code, http.StatusOK)

	{
		mgr.HandleKeyGet()
		assert.Equal(t, resp.Code, http.StatusNotFound)
	}
}
