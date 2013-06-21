package token

import (
	"crypto/md5"
	"fmt"
	"github.com/eaigner/hood"
	"io"
)

type Token struct {
	Id        hood.Id
	Key       string       `json:"key" sql:"pk,size(64)" validate:"presence"`
	Hash      string       `json:"hash" sql:"size(32)"`
	UserId    string       `json:"user_id" sql:"size(128)"`
	Scopes    []byte       `json:"scopes"`
	Client    []byte       `json:"client"`
	CreatedAt hood.Created `json:"-"`
	ExpiresIn int64        `json:"expires_in"`
	TouchedAt int64        `json:"touched_at"`
	Data      []byte       `json:"data"`

	ExpireAt int64 `json:"expire_at" sql:"-"`
}

func (t *Token) compatible() {
	t.ExpireAt = t.ExpiresIn
}

func hashResource(resource string) string {
	hash := md5.New()
	io.WriteString(hash, resource)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
