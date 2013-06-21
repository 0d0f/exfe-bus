package token

import (
	"fmt"
	"github.com/googollee/go-rest"
	"net/http"
	"time"
)

type Repo interface {
	Store(token Token) error
	FindByKey(key string) ([]Token, error)
	FindByHash(hash string) ([]Token, error)
	Touch(key, hash *string) error
	UpdateByKey(key string, data []byte, expiresIn *int64) (int, error)
	UpdateByHash(hash string, data []byte, expiresIn *int64) (int, error)
}

type Manager struct {
	rest.Service `prefix:"/v3/tokens"`

	Create          rest.Processor `method:"POST" path:""`
	KeyGet          rest.Processor `method:"GET" path:"/key/:key"`
	ResourceGet     rest.Processor `method:"GET" path:"/resources"`
	ResourceGet_    rest.Processor `method:"POST" path:"/resources" func:"HandleResourceGet"`
	KeyUpdate       rest.Processor `method:"PUT" path:"/key/:key"`
	KeyUpdate_      rest.Processor `method:"POST" path:"/key/:key" func:"HandleKeyUpdate"`
	ResourceUpdate  rest.Processor `method:"PUT" path:"/resource"`
	ResourceUpdate_ rest.Processor `method:"POST" path:"/resource" func:"HandleResourceUpdate"`

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

type CreateArg struct {
	Token
	Resource           string `json:"resource"`
	ExpireAfterSeconds int    `json:"expire_after_seconds"`
}

// 根据resource，data，user id，scopes，client和expire after seconds创建一个token
//
// 例子：
//
//     > curl "http://127.0.0.1:23333/v3/tokens?type=long" -d '{"data":"abc","resource":"123","user_id":"123","scopes":"exfe://user/verification","client":"c","expire_after_seconds":300}'
//
// 返回：
//
//     {"key":"0303","data":"abc","touched_at":21341234,"expire_at":66354}
func (t Manager) HandleCreate(arg CreateArg) Token {
	token := arg.Token
	token.Hash = hashResource(arg.Resource)
	token.TouchedAt = time.Now().Unix()
	token.ExpiresIn = time.Now().Add(time.Duration(arg.ExpireAfterSeconds) * time.Second).Unix()
	gentype := t.Request().URL.Query().Get("type")

	generator, ok := t.generators[gentype]
	if !ok {
		t.Error(http.StatusBadRequest, fmt.Errorf("invalid type %s", gentype))
		return token
	}
	for i := 0; i < 3; i++ {
		generator(&token)
		tokens, err := t.repo.FindByKey(token.Key)
		if err != nil {
			t.Error(http.StatusInternalServerError, err)
			return token
		}
		if len(tokens) == 0 {
			goto NEXIST
		}
	}
	t.Error(http.StatusConflict, fmt.Errorf("key collided"))
	return token

NEXIST:

	err := t.repo.Store(token)
	if err != nil {
		t.Error(http.StatusInternalServerError, err)
		return token
	}
	token.compatible()
	return token
}

// 根据key获得一个token，如果token不存在，返回错误
//
// 例子：
//
//     > curl "http://127.0.0.1:23333/v3/tokens/key/0303"
//
// 返回：
//
//     [{"key":"0303","data":"abc","touched_at":21341234,"expire_at":66354}]
func (t Manager) HandleKeyGet() []Token {
	key := t.Vars()["key"]
	tokens, err := t.repo.FindByKey(key)
	if err != nil {
		t.Error(http.StatusInternalServerError, err)
		return nil
	}
	if len(tokens) == 0 {
		t.Error(http.StatusNotFound, fmt.Errorf("can't find token with key %s", key))
		return nil
	}
	err = t.repo.Touch(&key, nil)
	if err != nil {
		t.Error(http.StatusInternalServerError, err)
		return nil
	}
	for i := range tokens {
		tokens[i].compatible()
	}
	return tokens
}

// 根据resource获得一个token，如果token不存在，返回错误
//
// 例子：
//
//     > curl "http://127.0.0.1:23333/v3/tokens/resources" -d '"abc"'
//
// 返回：
//
//     [{"key":"0303","data":"abc","touched_at":21341234,"expire_at":66354}]
func (t Manager) HandleResourceGet(resource string) []Token {
	hash := hashResource(resource)
	tokens, err := t.repo.FindByHash(hash)
	if err != nil {
		t.Error(http.StatusInternalServerError, err)
		return nil
	}
	if len(tokens) == 0 {
		t.Error(http.StatusNotFound, fmt.Errorf("can't find token with resource %s", resource))
		return nil
	}
	err = t.repo.Touch(nil, &hash)
	if err != nil {
		t.Error(http.StatusInternalServerError, err)
		return nil
	}
	for i := range tokens {
		tokens[i].compatible()
	}
	return tokens
}

type UpdateArg struct {
	Data               *string `json:"data"`
	ExpireAfterSeconds *int    `json:"expire_after_seconds"`
	Resource           string  `json:"resource"`
	ExpiresIn          *int64  `json:"-"`
}

func (a *UpdateArg) convert() {
	if a.ExpireAfterSeconds != nil {
		a.ExpiresIn = new(int64)
		*a.ExpiresIn = time.Now().Add(time.Duration(*a.ExpireAfterSeconds) * time.Second).Unix()
	}
}

// 更新key对应的token的data信息或者expire after seconds
//
// 例子：
//
//     > curl "http://127.0.0.1:23333/v3/tokens/key/0303" -d '{"data":"xyz","expire_after_seconds":13}'
func (t Manager) HandleKeyUpdate(arg UpdateArg) {
	arg.convert()
	key := t.Vars()["key"]
	var data []byte
	if arg.Data != nil {
		data = []byte(*arg.Data)
	}
	n, err := t.repo.UpdateByKey(key, data, arg.ExpiresIn)
	if err != nil {
		t.Error(http.StatusInternalServerError, err)
		return
	}
	if n == 0 {
		t.Error(http.StatusNotFound, fmt.Errorf("can't find token with key %s", key))
	}
}

// 更新resource对应的token的expire after seconds
//
// 例子：
//
//     > curl "http://127.0.0.1:23333/v3/tokens/resource" -d '{"resource":"abc", "expire_after_seconds":13}'
func (t Manager) HandleResourceUpdate(arg UpdateArg) {
	arg.convert()
	hash := hashResource(arg.Resource)
	var data []byte
	if arg.Data != nil {
		data = []byte(*arg.Data)
	}
	n, err := t.repo.UpdateByHash(hash, data, arg.ExpiresIn)
	if err != nil {
		t.Error(http.StatusInternalServerError, err)
		return
	}
	if n == 0 {
		t.Error(http.StatusNotFound, fmt.Errorf("can't find token with resource %s", arg.Resource))
	}
}
