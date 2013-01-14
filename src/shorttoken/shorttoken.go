package shorttoken

import (
	"fmt"
	"math"
	"math/rand"
	"model"
	"time"
)

type Repo interface {
	Store(token Token) error
	UpdateData(key, resource, data string) error
	UpdateExpireAt(key, resource string, expireAt time.Time) error
	Find(key string, resource string) ([]Token, error)
	Touch(key, resource string) error
}

type ShortToken struct {
	repo   Repo
	max    int32
	fmt    string
	random *rand.Rand
}

func New(repo Repo, length int) *ShortToken {
	return &ShortToken{
		repo:   repo,
		max:    int32(math.Pow10(length)),
		fmt:    fmt.Sprintf("%%0%dd", length),
		random: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (t *ShortToken) Create(resource, data string, after time.Duration) (model.Token, error) {
	key := ""
	for i := 0; i < 3; i++ {
		key = fmt.Sprintf(t.fmt, t.random.Int31n(t.max))
		tokens, err := t.repo.Find(key, "")
		if err != nil {
			return model.Token{}, err
		}
		if tokens == nil {
			goto NEXIST
		}
	}
	return model.Token{}, fmt.Errorf("key collided")
NEXIST:
	token := Token{
		Key:       key,
		Resource:  hashResource(resource),
		Data:      data,
		ExpireAt:  time.Now().Add(after),
		CreatedAt: time.Now(),
	}
	t.repo.Store(token)
	return model.Token{
		Key:  key,
		Data: data,
	}, nil
}

func (t *ShortToken) Get(key, resource string) ([]model.Token, error) {
	if key == "" && resource == "" {
		return nil, fmt.Errorf("key and resource should not both empty")
	}
	md5 := hashResource(resource)
	if resource == "" {
		md5 = ""
	}
	tokens, err := t.repo.Find(key, md5)
	if err != nil {
		return nil, err
	}
	if tokens == nil || len(tokens) == 0 {
		return nil, fmt.Errorf("can't find token with key(%s) or resource(%s)", key, resource)
	}
	if key != "" && md5 != "" {
		t.repo.Touch(key, md5)
	}
	ret := make([]model.Token, len(tokens))
	for i, token := range tokens {
		ret[i].Key = token.Key
		ret[i].Data = token.Data
		ret[i].TouchedAt = token.TouchedAt.UTC().Format("2006-01-02 15:04:05")
	}
	return ret, nil
}

func (t *ShortToken) UpdateData(key, data string) error {
	return t.repo.UpdateData(key, "", data)
}

func (t *ShortToken) Refresh(key, resource string, after time.Duration) error {
	if resource != "" {
		resource = hashResource(resource)
	}
	return t.repo.UpdateExpireAt(key, resource, time.Now().Add(after))
}
