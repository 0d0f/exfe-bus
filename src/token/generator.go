package token

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func GenerateShortToken(token *Token) {
	key := fmt.Sprintf("%04d", rand.Int31n(10000))
	token.Key = key
}

func GenerateLongToken(token *Token) {
	rand := hashResource(fmt.Sprintf("%s%x", time.Now().String(), randBytes()))
	token.Key = fmt.Sprintf("%s%s", rand, token.Hash)
}

func randBytes() (ret [32]byte) {
	for i := range ret {
		ret[i] = byte(rand.Int31n(math.MaxInt8))
	}
	return ret
}
