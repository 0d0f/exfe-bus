package main

import (
	"broker"
	"github.com/googollee/go-logger"
	"gobus"
	"model"
	"time"
	"tokenmanager"
)

type TokenManager struct {
	tokenRepo *TokenRepository
	manager   *tokenmanager.TokenManager
	log       *logger.SubLogger
	config    *model.Config
}

func NewTokenManager(config *model.Config, db *broker.DBMultiplexer) (*TokenManager, error) {
	repo, err := NewTokenRepository(config, db)
	if err != nil {
		return nil, err
	}
	return &TokenManager{
		tokenRepo: repo,
		manager:   tokenmanager.New(repo),
		log:       config.Log.SubPrefix("token manager"),
		config:    config,
	}, nil
}

func (mng *TokenManager) SetRoute(r gobus.RouteCreater) error {
	json := new(gobus.JSON)
	r().Methods("POST").Path("/tokenmanager").HandlerMethod(json, mng, "Generate")
	r().Methods("GET").Path("/tokenmanager/token/{token}").HandlerMethod(json, mng, "Get")
	r().Methods("POST").Path("/tokenmanager/token/{token}").HandlerMethod(json, mng, "Update")
	r().Methods("POST").Path("/tokenmanager/token/{token}/verify").HandlerMethod(json, mng, "Verify")
	r().Methods("POST").Path("/tokenmanager/token/{token}/refresh").HandlerMethod(json, mng, "Refresh")
	r().Methods("POST").Path("/tokenmanager/token/{token}/expire").HandlerMethod(json, mng, "Expire")
	r().Methods("POST").Path("/tokenmanager/resource").HandlerMethod(json, mng, "ResourceFind")
	r().Methods("POST").Path("/tokenmanager/resource/expire").HandlerMethod(json, mng, "ExpireAll")

	r().Queries("method", "Generate").Path("/TokenManager").HandlerMethod(json, mng, "Generate")
	r().Queries("method", "Get").Path("/TokenManager").HandlerMethod(json, mng, "Get_")
	r().Queries("method", "Find").Path("/TokenManager").HandlerMethod(json, mng, "ResourceFind")
	r().Queries("method", "Refresh").Path("/TokenManager").HandlerMethod(json, mng, "Refresh_")
	r().Queries("method", "Update").Path("/TokenManager").HandlerMethod(json, mng, "Update_")
	r().Queries("method", "Verify").Path("/TokenManager").HandlerMethod(json, mng, "Verify_")
	r().Queries("method", "Expire").Path("/TokenManager").HandlerMethod(json, mng, "Expire_")
	r().Queries("method", "ExpireAll").Path("/TokenManager").HandlerMethod(json, mng, "ExpireAll")

	return nil
}

type TokenGenerateArgs struct {
	Resource           string `json:"resource"`
	Data               string `json:"data"`
	ExpireAfterSeconds int    `json:"expire_after_seconds"`
}

// 根据资源resource，数据data和过期时间expire_after_seconds生成一个token。如果expire_after_seconds是-1，则此token无过期时间
//
// 例子：
//
//     > curl http://127.0.0.1:23333/tokenmanager -d '{"resource":"abcde","data":"","expire_after_seconds":12}'
//     "ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a"
func (mng *TokenManager) Generate(params map[string]string, arg TokenGenerateArgs) (string, error) {
	expire := time.Duration(arg.ExpireAfterSeconds) * time.Second
	if arg.ExpireAfterSeconds < 0 {
		expire = tokenmanager.NeverExpire
	}
	token, err := mng.manager.GenerateToken(arg.Resource, arg.Data, expire)
	if err != nil {
		return "", err
	}
	return token.String(), nil
}

// 根据token返回Token对象
//
// 例子：
//
//     > curl http://127.0.0.1:23333/tokenmanager/token/ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a
//     {"token":"ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a","data":"","is_expire":true}
func (mng *TokenManager) Get(params map[string]string) (*tokenmanager.Token, error) {
	token := params["token"]
	return mng.manager.GetToken(token)
}

func (mng *TokenManager) Get_(params map[string]string, token string) (*tokenmanager.Token, error) {
	return mng.manager.GetToken(token)
}

// 根据resource，查找对应的所有Token对象
//
// 例子：
//
//     > curl http://127.0.0.1:23333/tokenmanager/resource -d '"abcde"'
//     [{"token":"ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a","data":"","is_expire":true},{"token":"ab56b4d92b40713acc5af89985d4b786aa354d0a4f80c96c85c728e67a73b795","data":"","is_expire":true}]
func (mng *TokenManager) ResourceFind(params map[string]string, resource string) ([]*tokenmanager.Token, error) {
	return mng.manager.FindTokens(resource)
}

// 更新token对应的数据data
//
// 例子：
//
// > curl http://127.0.0.1:23333/tokenmanager/token/ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a -d '"123"'
// 0
func (mng *TokenManager) Update(params map[string]string, data string) (int, error) {
	token := params["token"]
	return 0, mng.manager.UpdateData(token, data)
}

type TokenUpdateArg struct {
	Token string `json:"token"`
	Data  string `json:"data"`
}

func (mng *TokenManager) Update_(params map[string]string, arg TokenUpdateArg) (int, error) {
	return 0, mng.manager.UpdateData(arg.Token, arg.Data)
}

type TokenVerifyReply struct {
	Matched bool                `json:"matched"`
	Token   *tokenmanager.Token `json:"token,omitempty"`
}

// 根据token和资源resource来验证两者是否一致matched，并返回Token对象
//
// 例子：
//
//     > curl http://127.0.0.1:23333/tokenmanager/token/ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a/verify -d '"abcde"'
//     {"matched":true,{"token":"ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a","data":"","is_expire":true}}
func (mng *TokenManager) Verify(params map[string]string, resource string) (TokenVerifyReply, error) {
	token := params["token"]
	matched, tk, err := mng.manager.VerifyToken(token, resource)
	if !matched {
		return TokenVerifyReply{matched, nil}, err
	}
	return TokenVerifyReply{matched, tk}, err
}

type TokenVerifyArg struct {
	Token    string `json:"token"`
	Resource string `json:"resource"`
}

func (mng *TokenManager) Verify_(params map[string]string, arg TokenVerifyArg) (TokenVerifyReply, error) {
	matched, tk, err := mng.manager.VerifyToken(arg.Token, arg.Resource)
	if !matched {
		return TokenVerifyReply{matched, nil}, err
	}
	return TokenVerifyReply{matched, tk}, err
}

// 将token的过期时间设为expire_after_seconds秒之后过期。如果expire_after_seconds为-1，则token永不过期。
//
// 例子：
//
//     > curl http://127.0.0.1:23333/tokenmanager/token/ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a/refresh -d '-1'
//     0
func (mng *TokenManager) Refresh(params map[string]string, expireAfterSeconds int) (int, error) {
	token := params["token"]
	expire := time.Duration(expireAfterSeconds) * time.Second
	if expireAfterSeconds < 0 {
		expire = tokenmanager.NeverExpire
	}
	return 0, mng.manager.RefreshToken(token, expire)
}

type TokenRefreshArg struct {
	Token              string `json:"token"`
	ExpireAfterSeconds int    `json:"expire_after_seconds"`
}

func (mng *TokenManager) Refresh_(params map[string]string, arg TokenRefreshArg) (int, error) {
	expire := time.Duration(arg.ExpireAfterSeconds) * time.Second
	if arg.ExpireAfterSeconds < 0 {
		expire = tokenmanager.NeverExpire
	}
	return 0, mng.manager.RefreshToken(arg.Token, expire)
}

// 立刻使token过期。
//
// 例子：
//
//     > curl http://127.0.0.1:23333/tokenmanager/token/ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a/expire
//     0
func (mng *TokenManager) Expire(params map[string]string) (int, error) {
	token := params["token"]
	return 0, mng.manager.RefreshToken(token, 0)
}

func (mng *TokenManager) Expire_(params map[string]string, token string) (int, error) {
	return 0, mng.manager.RefreshToken(token, 0)
}

// 立刻使resource对应的所有token过期。
//
// 例子：
//
//     > curl http://127.0.0.1:23333/tokenmanager/resource/expire -d '"abc"'
//     0
func (mng *TokenManager) ExpireAll(params map[string]string, resource string) (int, error) {
	return 0, mng.manager.RefreshTokensByResource(resource, 0)
}
