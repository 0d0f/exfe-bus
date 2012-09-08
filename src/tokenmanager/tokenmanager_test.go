package tokenmanager

import (
	"fmt"
	"testing"
	"time"
)

type TestTokenRepo struct {
	store map[string]Token
}

func (r *TestTokenRepo) Create(token *Token) error {
	r.store[token.String()] = *token
	return nil
}

func (r *TestTokenRepo) Store(token *Token) error {
	r.store[token.String()] = *token
	return nil
}

func (r *TestTokenRepo) FindByKey(key string) ([]*Token, error) {
	ret := make([]*Token, 0, 0)
	for _, token := range r.store {
		if token.Key != key {
			continue
		}
		ret = append(ret, &token)
	}
	return ret, nil
}

func (r *TestTokenRepo) FindByToken(key, rand string) (*Token, error) {
	k := fmt.Sprintf("%s%s", key, rand)
	token, ok := r.store[k]
	if !ok {
		return nil, nil
	}
	return &token, nil
}

func (r *TestTokenRepo) UpdateDataByToken(key, rand, data string) error {
	token := fmt.Sprintf("%s%s", key, rand)
	v, ok := r.store[token]
	if !ok {
		fmt.Errorf("no token")
	}
	v.Data = data
	r.store[token] = v
	return nil
}

func (r *TestTokenRepo) UpdateExpireAtByToken(key, rand string, expire_at *time.Time) error {
	token := fmt.Sprintf("%s%s", key, rand)
	v, ok := r.store[token]
	if !ok {
		fmt.Errorf("no token")
	}
	v.Key = key
	r.store[token] = v
	return nil
}

func (r *TestTokenRepo) UpdateExpireAtByKey(key string, expire_at *time.Time) error {
	for k, v := range r.store {
		if v.Key == key {
			v.ExpireAt = expire_at
			r.store[k] = v
		}
	}
	return nil
}

func (r *TestTokenRepo) DeleteByToken(key, rand string) error {
	delete(r.store, fmt.Sprintf("%s%s", key, rand))
	return nil
}

func TestTokenManager(t *testing.T) {
	repo := &TestTokenRepo{
		store: make(map[string]Token),
	}
	mgr := New(repo)
	resource := "resource"
	tk := "fjadsklfjkldasfdasiffjuoru21urjew"

	{
		_, err := mgr.GetToken(tk)
		if err == nil {
			t.Fatalf("get resource should failed")
		}

		ok, _, err := mgr.VerifyToken(tk, resource)
		if err == nil || ok {
			t.Errorf("tk(%s) verify with resource(%s) should failed", tk, resource)
		}
	}

	token, err := mgr.GenerateToken(resource, "", time.Second)
	if err != nil {
		t.Fatalf("generate token failed: %s", err)
	}
	tk = token.String()

	{
		token, err := mgr.GetToken(tk)
		if err != nil {
			t.Fatalf("get resource failed: %s", err)
		}
		if got, expect := token.Data, ""; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}

		tks, err := mgr.FindTokens(resource)
		if err != nil {
			t.Errorf("get tokens failed: %s", err)
		}
		if got, expect := len(tks), 1; got != expect {
			t.Errorf("got: %d, expect: %d", got, expect)
		}
		if got, expect := tks[0].String(), tk; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
		if got, expect := tks[0].IsExpired(), false; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}

		ok, token, err := mgr.VerifyToken(tk, resource)
		if err != nil {
			t.Errorf("tk(%s) verify with resource(%s) failed: %s", tk, resource, err)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
		if got, expect := token.Data, ""; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}

		err = mgr.UpdateData(tk, "abc")
		if err != nil {
			t.Errorf("tk(%s) update data failed: %s", tk, err)
		}

		ok, token, err = mgr.VerifyToken(tk, resource)
		if err != nil {
			t.Errorf("tk(%s) verify with resource(%s) failed: %s", tk, resource, err)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
		if got, expect := token.Data, "abc"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}

		err = mgr.UpdateData(tk, "123")
		if err != nil {
			t.Errorf("tk(%s) update data failed: %s", tk, err)
		}

		ok, token, err = mgr.VerifyToken(tk, resource)
		if err != nil {
			t.Errorf("tk(%s) verify with resource(%s) failed: %s", tk, resource, err)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
		if got, expect := token.Data, "123"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}

	time.Sleep(time.Second * 2)

	{
		token, _ := mgr.GetToken(tk)
		if !token.IsExpired() {
			t.Fatalf("get resource should expired")
		}

		ok, token, _ := mgr.VerifyToken(tk, resource)
		if !token.IsExpired() {
			t.Errorf("tk(%s) verify with resource(%s) should expired", tk, resource)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
	}

	mgr.RefreshToken(tk, time.Second)

	{
		_, err := mgr.GetToken(tk)
		if err != nil {
			t.Fatalf("get resource failed: %s", err)
		}

		ok, _, err := mgr.VerifyToken(tk, resource)
		if err != nil {
			t.Errorf("tk(%s) verify with resource(%s) failed: %s", tk, resource, err)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
	}

	mgr.ExpireToken(tk)

	{
		token, _ := mgr.GetToken(tk)
		if !token.IsExpired() {
			t.Fatalf("get resource should expired")
		}

		ok, token, _ := mgr.VerifyToken(tk, resource)
		if !token.IsExpired() {
			t.Errorf("tk(%s) verify with resource(%s) should expired", tk, resource)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
	}

	mgr.RefreshToken(tk, NeverExpire)

	{
		_, err := mgr.GetToken(tk)
		if err != nil {
			t.Fatalf("get resource failed: %s", err)
		}

		ok, _, err := mgr.VerifyToken(tk, resource)
		if err != nil {
			t.Errorf("tk(%s) verify with resource(%s) failed: %s", tk, resource, err)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
	}

	mgr.ExpireTokensByKey(tk[:32])

	{
		token, _ := mgr.GetToken(tk)
		if !token.IsExpired() {
			t.Fatalf("get resource should expired")
		}

		ok, token, _ := mgr.VerifyToken(tk, resource)
		if !token.IsExpired() {
			t.Errorf("tk(%s) verify with resource(%s) should expired", tk, resource)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
	}

	mgr.RefreshToken(tk, NeverExpire)

	{
		_, err := mgr.GetToken(tk)
		if err != nil {
			t.Fatalf("get resource failed: %s", err)
		}

		ok, _, err := mgr.VerifyToken(tk, resource)
		if err != nil {
			t.Errorf("tk(%s) verify with resource(%s) failed: %s", tk, resource, err)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
	}

	time.Sleep(time.Second * 2)

	{
		_, err := mgr.GetToken(tk)
		if err != nil {
			t.Fatalf("get resource failed: %s", err)
		}

		ok, _, err := mgr.VerifyToken(tk, resource)
		if err != nil {
			t.Errorf("tk(%s) verify with resource(%s) failed: %s", tk, resource, err)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
	}

	err = mgr.DeleteToken(tk)
	if err != nil {
		t.Errorf("delete fail: %s", err)
	}
}
