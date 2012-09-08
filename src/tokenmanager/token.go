package tokenmanager

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"time"
)

type Token struct {
	Key       string
	Rand      string
	Data      string
	ExpireAt  *time.Time
	CreatedAt time.Time
}

type tokenJson struct {
	Token    string `json:"token"`
	Data     string `json:"data"`
	IsExpire bool   `json:"is_expire"`
}

// expireAt == nil for never expire.
func NewToken(resource, data string, expireAt *time.Time) *Token {
	hash := md5.New()
	io.WriteString(hash, fmt.Sprintf("%s%x", time.Now().String(), randBytes()))
	return &Token{
		Key:       md5Resource(resource),
		Rand:      fmt.Sprintf("%x", hash.Sum(nil)),
		Data:      data,
		ExpireAt:  expireAt,
		CreatedAt: time.Now(),
	}
}

func (t *Token) IsExpired() bool {
	if t.ExpireAt == nil {
		return false
	}
	return t.ExpireAt.Sub(time.Now()) <= 0
}

func (t *Token) String() string {
	return fmt.Sprintf("%s%s", t.Key, t.Rand)
}

func (t Token) MarshalJSON() ([]byte, error) {
	j := tokenJson{
		Token:    (&t).String(),
		Data:     t.Data,
		IsExpire: (&t).IsExpired(),
	}
	return json.Marshal(j)
}

var rander = rand.New(rand.NewSource(time.Now().UnixNano()))

func randBytes() (ret [32]byte) {
	for i := range ret {
		ret[i] = byte(rander.Int31n(math.MaxInt8))
	}
	return ret
}

func md5Resource(resource string) string {
	hash := md5.New()
	io.WriteString(hash, resource)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
