package token

import (
	"github.com/googollee/go-rest"
	"net/http"
	"time"
)

type Repo interface {
	Store(token Token) error
	FindByKey(key string) ([]Token, error)
	FindByHash(hash string) ([]Token, error)
	Touch(key, hash *string) error
	UpdateByKey(key string, data *string, expiresIn *int64) (int64, error)
	UpdateByHash(hash string, data *string, expiresIn *int64) (int64, error)
}

type Manager struct {
	rest.Service `prefix:"/v3/tokens"`

	create          rest.SimpleNode `route:"" method:"POST"`
	resourceGet     rest.SimpleNode `route:"/resources" method:"GET"`
	resourceGet_    rest.SimpleNode `route:"/resources" method:"POST"`
	resourceUpdate  rest.SimpleNode `route:"/resource" method:"PUT"`
	resourceUpdate_ rest.SimpleNode `route:"/resource" method:"POST"`
	keyGet          rest.SimpleNode `route:"/key/:key" method:"GET"`
	keyUpdate       rest.SimpleNode `route:"/key/:key" method:"PUT"`
	keyUpdate_      rest.SimpleNode `route:"/key/:key" method:"POST"`

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
func (t Manager) Create(ctx rest.Context, arg CreateArg) {
	token := arg.Token
	token.Hash = hashResource(arg.Resource)
	token.TouchedAt = time.Now().Unix()
	token.ExpiresAt = time.Now().Add(time.Duration(arg.ExpireAfterSeconds) * time.Second).Unix()
	var gentype string
	ctx.Bind("type", &gentype)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, "%s", err)
		return
	}

	generator, ok := t.generators[gentype]
	if !ok {
		ctx.Return(http.StatusBadRequest, "invalid type %s", gentype)
		return
	}
	for i := 0; i < 3; i++ {
		generator(&token)
		tokens, err := t.repo.FindByKey(token.Key)
		if err != nil {
			ctx.Return(http.StatusInternalServerError, "%s", err)
			return
		}
		if len(tokens) == 0 {
			goto NEXIST
		}
	}
	ctx.Return(http.StatusConflict, "key collided")
	return

NEXIST:

	err := t.repo.Store(token)
	if err != nil {
		ctx.Return(http.StatusInternalServerError, "%s", err)
		return
	}
	token.compatible()
	ctx.Render(token)
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
func (t Manager) KeyGet(ctx rest.Context) {
	var key string
	ctx.Bind("key", &key)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, "%s", err)
		return
	}
	tokens, err := t.repo.FindByKey(key)
	if err != nil {
		ctx.Return(http.StatusInternalServerError, "%s", err)
		return
	}
	if len(tokens) == 0 {
		ctx.Return(http.StatusNotFound, "can't find token with key %s", key)
		return
	}
	err = t.repo.Touch(&key, nil)
	if err != nil {
		ctx.Return(http.StatusInternalServerError, "%s", err)
		return
	}
	for i := range tokens {
		tokens[i].compatible()
	}
	ctx.Render(tokens)
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
func (t Manager) ResourceGet(ctx rest.Context, resource string) {
	hash := hashResource(resource)
	tokens, err := t.repo.FindByHash(hash)
	if err != nil {
		ctx.Return(http.StatusInternalServerError, "%s", err)
		return
	}
	if len(tokens) == 0 {
		ctx.Return(http.StatusNotFound, "can't find token with resource %s", resource)
		return
	}
	err = t.repo.Touch(nil, &hash)
	if err != nil {
		ctx.Return(http.StatusInternalServerError, "%s", err)
		return
	}
	for i := range tokens {
		tokens[i].compatible()
	}
	ctx.Render(tokens)
}

func (t Manager) ResourceGet_(ctx rest.Context, resource string) {
	t.ResourceGet(ctx, resource)
}

type UpdateArg struct {
	Data               *string `json:"data"`
	ExpireAfterSeconds *int    `json:"expire_after_seconds"`
	Resource           string  `json:"resource"`
	ExpiresAt          *int64  `json:"-"`
}

func (a *UpdateArg) convert() {
	if a.ExpireAfterSeconds != nil {
		a.ExpiresAt = new(int64)
		*a.ExpiresAt = time.Now().Add(time.Duration(*a.ExpireAfterSeconds) * time.Second).Unix()
	}
}

// 更新key对应的token的data信息或者expire after seconds
//
// 例子：
//
//     > curl "http://127.0.0.1:23333/v3/tokens/key/0303" -d '{"data":"xyz","expire_after_seconds":13}'
func (t Manager) KeyUpdate(ctx rest.Context, arg UpdateArg) {
	var key string
	ctx.Bind("key", &key)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, "%s", err)
		return
	}
	arg.convert()
	n, err := t.repo.UpdateByKey(key, arg.Data, arg.ExpiresAt)
	if err != nil {
		ctx.Return(http.StatusInternalServerError, "%s", err)
		return
	}
	if n == 0 {
		ctx.Return(http.StatusNotFound, "can't find token with key %s", key)
	}
}

func (t Manager) KeyUpdate_(ctx rest.Context, arg UpdateArg) {
	t.KeyUpdate(ctx, arg)
}

// 更新resource对应的token的expire after seconds
//
// 例子：
//
//     > curl "http://127.0.0.1:23333/v3/tokens/resource" -d '{"resource":"abc", "expire_after_seconds":13}'
func (t Manager) ResourceUpdate(ctx rest.Context, arg UpdateArg) {
	arg.convert()
	hash := hashResource(arg.Resource)
	n, err := t.repo.UpdateByHash(hash, arg.Data, arg.ExpiresAt)
	if err != nil {
		ctx.Return(http.StatusInternalServerError, "%s", err)
		return
	}
	if n == 0 {
		ctx.Return(http.StatusNotFound, "can't find token with resource %s", arg.Resource)
	}
}

func (t Manager) ResourceUpdate_(ctx rest.Context, arg UpdateArg) {
	t.ResourceUpdate(ctx, arg)
}
