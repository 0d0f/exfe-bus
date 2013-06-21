package token

import (
	"github.com/stretchrcom/testify/assert"
	"testing"
)

func TestShortGenerator(t *testing.T) {
	token := Token{
		Hash: hashResource("resource"),
		Data: []byte("data"),
	}
	GenerateShortToken(&token)
	assert.Equal(t, len(token.Key), 4)
	for _, c := range token.Key {
		if !('0' <= c && c <= '9') {
			t.Errorf("token %s must all number", token.Key)
		}
	}
}

func TestLongGenerator(t *testing.T) {
	token := Token{
		Hash: hashResource("resource"),
		Data: []byte("data"),
	}
	GenerateLongToken(&token)
	assert.Equal(t, len(token.Key), 64)
}
