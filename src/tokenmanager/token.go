package tokenmanager

import (
	"crypto/md5"
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
	token := &t
	json := fmt.Sprintf("{\"token\":\"%s\",\"data\":\"%s\",\"is_expire\":%v}", token.String(), token.Data, token.IsExpired())
	return []byte(json), nil
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
