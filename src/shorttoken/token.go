package shorttoken

import (
	"crypto/md5"
	"fmt"
	"io"
	"time"
)

type Token struct {
	Key       string
	Resource  string
	Data      string
	ExpireAt  time.Time
	CreatedAt time.Time
}

func hashResource(resource string) string {
	hash := md5.New()
	io.WriteString(hash, resource)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
