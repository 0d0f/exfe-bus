package token

import (
	"crypto/md5"
	"fmt"
	"io"
)

type Token struct {
	Key       string `json:"key"`
	Hash      string `json:"hash"`
	UserId    string `json:"user_id"`
	Scopes    string `json:"scopes"`
	Client    string `json:"client"`
	CreatedAt int64  `json:"-"`
	ExpiresAt int64  `json:"expires_at"`
	TouchedAt int64  `json:"touched_at"`
	Data      string `json:"data"`

	ExpireAt int64 `json:"expire_at"`
}

func (t *Token) compatible() {
	t.ExpireAt = t.ExpiresAt
}

func hashResource(resource string) string {
	hash := md5.New()
	io.WriteString(hash, resource)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
