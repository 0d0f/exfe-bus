package token

import (
	"crypto/md5"
	"fmt"
	"io"
	"model"
	"time"
)

type Token struct {
	Key       string
	Hash      string
	Data      string
	TouchedAt time.Time
	ExpireAt  time.Time
	CreatedAt time.Time
}

func (t Token) Token() model.Token {
	return model.Token{
		Key:       t.Key,
		Data:      t.Data,
		TouchedAt: t.TouchedAt.UTC().Format("2006-01-02 15:04:05"),
	}
}

func hashResource(resource string) string {
	hash := md5.New()
	io.WriteString(hash, resource)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
