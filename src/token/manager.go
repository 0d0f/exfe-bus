package token

import (
	"fmt"
	"model"
	"time"
)

type Repo interface {
	Store(token Token) error
	UpdateData(token Token, data string) error
	UpdateExpireAt(token Token, expireAt time.Time) error
	Find(token Token) ([]Token, error)
	Touch(token Token) error
}

type Manager struct {
	repo       Repo
	generators map[string]func(*Token)
}

func New(repo Repo) *Manager {
	generators := map[string]func(*Token){
		"short": GenerateShortToken,
		"long":  GenerateLongToken,
	}
	return &Manager{
		repo:       repo,
		generators: generators,
	}
}

func (t *Manager) Create(gentype, resource, data string, after time.Duration) (model.Token, error) {
	token := Token{
		Hash:      hashResource(resource),
		Data:      data,
		ExpireAt:  time.Now().Add(after),
		CreatedAt: time.Now(),
	}
	generator, ok := t.generators[gentype]
	if !ok {
		return model.Token{}, fmt.Errorf("invalid type %s", gentype)
	}
	for i := 0; i < 3; i++ {
		generator(&token)
		tokens, err := t.repo.Find(token)
		if err != nil {
			return model.Token{}, err
		}
		if tokens == nil {
			goto NEXIST
		}
	}
	return model.Token{}, fmt.Errorf("key collided")

NEXIST:
	err := t.repo.Store(token)
	if err != nil {
		return model.Token{}, err
	}
	return token.Token(), nil
}

func (t *Manager) Get(key, resource string) ([]model.Token, error) {
	if key == "" && resource == "" {
		return nil, fmt.Errorf("key and resource should not both empty")
	}
	token := Token{}
	token.Key = key
	if resource != "" {
		token.Hash = hashResource(resource)
	}
	tokens, err := t.repo.Find(token)
	if err != nil {
		return nil, err
	}
	if len(tokens) == 0 {
		return nil, fmt.Errorf("can't find token with key(%s) or resource(%s)", key, resource)
	}
	err = t.repo.Touch(token)
	if err != nil {
		return nil, err
	}
	ret := make([]model.Token, len(tokens))
	for i, token := range tokens {
		ret[i] = token.Token()
	}
	return ret, nil
}

func (t *Manager) UpdateData(key, data string) error {
	token := Token{
		Key: key,
	}
	return t.repo.UpdateData(token, data)
}

func (t *Manager) Refresh(key, resource string, after time.Duration) error {
	hash := ""
	if resource != "" {
		hash = hashResource(resource)
	}
	token := Token{
		Key:  key,
		Hash: hash,
	}
	return t.repo.UpdateExpireAt(token, time.Now().Add(after))
}
