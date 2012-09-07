package tokenmanager

import (
	"encoding/json"
	"testing"
	"time"
)

func TestToken(t *testing.T) {
	token := NewToken("resource", "data", nil)

	if got, expect := token.IsExpired(), false; got != expect {
		t.Errorf("got: %v, expect: %v", got, expect)
	}

	ti := time.Now()
	token.ExpireAt = &ti
	if got, expect := token.IsExpired(), true; got != expect {
		t.Errorf("got: %v, expect: %v", got, expect)
	}

	token_str := token.String()
	t.Log(token_str)
	if got, expect := len(token_str), 64; got != expect {
		t.Errorf("got: %v, expect: %v", got, expect)
	}

	ti = ti.Add(time.Minute)
	token.ExpireAt = &ti
	if got, expect := token.IsExpired(), false; got != expect {
		t.Errorf("got: %v, expect: %v", got, expect)
	}

	{
		token1 := NewToken("resource", "data123", nil)
		if got, expect := token.Key == token1.Key, true; got != expect {
			t.Errorf("got: %v, expect: %v", got, expect)
		}
		if got, expect := token.Data == token1.Data, false; got != expect {
			t.Errorf("got: %v, expect: %v", got, expect)
		}
	}

	j, err := json.Marshal(token)
	if err != nil {
		t.Errorf("marshal json fail: %s", err)
	}

	type TokenJson struct {
		Token    string `json:"token"`
		Data     string `json:"data"`
		IsExpire bool   `json:"is_expire"`
	}
	var tkFromJson TokenJson
	err = json.Unmarshal(j, &tkFromJson)
	if err != nil {
		t.Errorf("unmashal json fail: %s", err)
	}
	if got, expect := tkFromJson.Token, token.String(); got != expect {
		t.Errorf("got: %s, expect: %s", got, expect)
	}
	if got, expect := tkFromJson.Data, token.Data; got != expect {
		t.Errorf("got: %s, expect: %s", got, expect)
	}
	if got, expect := tkFromJson.IsExpire, token.IsExpired(); got != expect {
		t.Errorf("got: %v, expect: %v", got, expect)
	}
}
