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

type TokenGenerateArgs struct {
	Resource           string `json:"resource"`
	Data               string `json:"data"`
	ExpireAfterSeconds int    `json:"expire_after_seconds"`
}

// 根据资源resource，数据data和过期时间expire_after_seconds生成一个token。如果expire_after_seconds是-1，则此token无过期时间
//
// 例子：
//
//     > curl http://127.0.0.1:23333/TokenManager?method=Generate -d '{"resource":"abcde","data":"","expire_after_seconds":12}'
//     "ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a"
func (mng *TokenManager) Generate(meta *gobus.HTTPMeta, arg *TokenGenerateArgs, reply *string) error {
	expire := time.Duration(arg.ExpireAfterSeconds) * time.Second
	if arg.ExpireAfterSeconds < 0 {
		expire = tokenmanager.NeverExpire
	}
	token, err := mng.manager.GenerateToken(arg.Resource, arg.Data, expire)
	if err != nil {
		return err
	}
	*reply = token.String()
	return nil
}

// 根据token返回Token对象
//
// 例子：
//
//     > curl http://127.0.0.1:23333/TokenManager?method=Get -d '"ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a"'
//     {"token":"ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a","data":"","is_expire":true}
func (mng *TokenManager) Get(meta *gobus.HTTPMeta, token *string, reply *tokenmanager.Token) error {
	tk, err := mng.manager.GetToken(*token)
	if err != nil {
		return err
	}
	*reply = *tk
	return nil
}

// 根据resource，查找对应的所有Token对象
//
// 例子：
//
//     > curl http://127.0.0.1:23333/TokenManager?method=Find -d '"abcde"'
//     [{"token":"ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a","data":"","is_expire":true},{"token":"ab56b4d92b40713acc5af89985d4b786aa354d0a4f80c96c85c728e67a73b795","data":"","is_expire":true}]
func (mng *TokenManager) Find(meta *gobus.HTTPMeta, resource *string, reply *[]*tokenmanager.Token) error {
	tks, err := mng.manager.FindTokens(*resource)
	if err != nil {
		return err
	}
	*reply = tks
	return err
}

type TokenUpdateArg struct {
	Token string `json:"token"`
	Data  string `json:"data"`
}

// 更新token对应的数据data
//
// 例子：
//
// > curl http://127.0.0.1:23333/TokenManager?method=Update -d '{"token":"ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a","data":"123"}'
// 0
func (mng *TokenManager) Update(meta *gobus.HTTPMeta, arg *TokenUpdateArg, reply *int) error {
	return mng.manager.UpdateData(arg.Token, arg.Data)
}

type TokenVerifyArg struct {
	Token    string `json:"token"`
	Resource string `json:"resource"`
}

type TokenVerifyReply struct {
	Matched bool                `json:"matched"`
	Token   *tokenmanager.Token `json:"token,omitempty"`
}

// 根据token和资源resource来验证两者是否一致matched，并返回Token对象
//
// 例子：
//
//     > curl http://127.0.0.1:23333/TokenManager?method=Verify -d '{"token":"ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a","resource":"abcde"}'
//     {"matched":true,{"token":"ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a","data":"","is_expire":true}}
func (mng *TokenManager) Verify(meta *gobus.HTTPMeta, args *TokenVerifyArg, reply *TokenVerifyReply) error {
	var err error
	reply.Matched, reply.Token, err = mng.manager.VerifyToken(args.Token, args.Resource)
	if !reply.Matched {
		reply.Token = nil
	}
	return err
}

type TokenRefreshArg struct {
	Token              string `json:"token"`
	ExpireAfterSeconds int    `json:"expire_after_seconds"`
}

// 将token的过期时间设为expire_after_seconds秒之后过期。如果expire_after_seconds为-1，则token永不过期。
//
// 例子：
//
//     > curl http://127.0.0.1:23333/TokenManager?method=Refresh -d '{"token":"ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a","expire_after_seconds":-1}'
//     0
func (mng *TokenManager) Refresh(meta *gobus.HTTPMeta, args *TokenRefreshArg, reply *int) error {
	expire := time.Duration(args.ExpireAfterSeconds) * time.Second
	if args.ExpireAfterSeconds < 0 {
		expire = tokenmanager.NeverExpire
	}
	return mng.manager.RefreshToken(args.Token, expire)
}

// 立刻使token过期。
//
// 例子：
//
//     > curl http://127.0.0.1:23333/TokenManager?method=Expire -d '"ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a"'
//     0
func (mng *TokenManager) Expire(meta *gobus.HTTPMeta, token *string, reply *int) error {
	return mng.manager.RefreshToken(*token, 0)
}

// 立刻使resource对应的所有token过期。
//
// 例子：
//
//     > curl http://127.0.0.1:23333/TokenManager?method=ExpireAll -d '"abc"'
//     0
func (mng *TokenManager) ExpireAll(meta *gobus.HTTPMeta, resource *string, reply *int) error {
	return mng.manager.RefreshTokensByResource(*resource, 0)
}
