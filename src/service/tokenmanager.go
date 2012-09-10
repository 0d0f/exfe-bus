package main

import (
	"github.com/googollee/go-logger"
	"gobus"
	"time"
	"tokenmanager"
)

type TokenManager struct {
	tokenRepo *TokenRepository
	manager   *tokenmanager.TokenManager
	log       *logger.SubLogger
	config    *Config
}

func NewTokenManager(config *Config) (*TokenManager, error) {
	repo, err := NewTokenRepository(config)
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
	log := mng.log.SubCode()
	log.Debug("generate with %+v", arg)
	expire := time.Duration(arg.ExpireAfterSeconds) * time.Second
	if arg.ExpireAfterSeconds < 0 {
		expire = tokenmanager.NeverExpire
	}
	token, err := mng.manager.GenerateToken(arg.Resource, arg.Data, expire)
	if err != nil {
		log.Info("generate token with resource(%s), expire(%ds) fail: %s", arg.Resource, arg.ExpireAfterSeconds, err)
		return err
	}
	*reply = token.String()
	log.Debug("return token: %s", *reply)
	return nil
}

// 根据token返回Token对象
//
// 例子：
//
//     > curl http://127.0.0.1:23333/TokenManager?method=Get -d '"ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a"'
//     {"token":"ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a","data":"","is_expire":true}
func (mng *TokenManager) Get(meta *gobus.HTTPMeta, token *string, reply *tokenmanager.Token) error {
	log := mng.log.SubCode()
	log.Debug("get with token: %s", *token)
	tk, err := mng.manager.GetToken(*token)

	if err != nil {
		log.Info("get resource with token(%s) fail: %s", *token, err)
		return err
	} else {
		log.Debug("return token: %+v", tk)
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
	log := mng.log.SubCode()
	log.Debug("find with resource: %s", *resource)
	tks, err := mng.manager.FindTokens(*resource)
	if err != nil {
		log.Info("find tokens with resource(%s) fail: %s", *resource, err)
		return err
	}

	log.Debug("return tokens: %+v", tks)
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
func (mng *TokenManager) Update(meta *gobus.HTTPMeta, arg *TokenUpdateArg, reply *int) (err error) {
	log := mng.log.SubCode()
	log.Debug("update token: %s with data: %s", arg.Token, arg.Data)
	err = mng.manager.UpdateData(arg.Token, arg.Data)
	if err != nil {
		log.Info("update token(%s) with data(%s) fail: %s", arg.Token, arg.Data, err)
	} else {
		log.Debug("success")
	}
	return
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
func (mng *TokenManager) Verify(meta *gobus.HTTPMeta, args *TokenVerifyArg, reply *TokenVerifyReply) (err error) {
	log := mng.log.SubCode()
	log.Debug("verify with token: %s, resource: %s", args.Token, args.Resource)
	reply.Matched, reply.Token, err = mng.manager.VerifyToken(args.Token, args.Resource)
	if err != nil {
		log.Info("verify with token(%s)&resource(%s) fail: %s", args.Token, args.Resource, err)
	} else {
		log.Debug("return: %v", reply)
	}
	if !reply.Matched {
		reply.Token = nil
	}
	return err
}

// 删除token，如果成功返回0
//
// 例子：
//
//     > curl http://127.0.0.1:23333/TokenManager?method=Delete -d '"ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a"'
//     0
func (mng *TokenManager) Delete(meta *gobus.HTTPMeta, token *string, reply *int) (err error) {
	log := mng.log.SubCode()
	log.Debug("delete token: %s", *token)
	err = mng.manager.DeleteToken(*token)
	if err != nil {
		log.Info("delete token(%s) fail: %s", *token, err)
	} else {
		log.Debug("ok")
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
func (mng *TokenManager) Refresh(meta *gobus.HTTPMeta, args *TokenRefreshArg, reply *int) (err error) {
	log := mng.log.SubCode()
	log.Debug("refresh token: %s, expire: %ds", args.Token, args.ExpireAfterSeconds)
	expire := time.Duration(args.ExpireAfterSeconds) * time.Second
	if args.ExpireAfterSeconds < 0 {
		expire = tokenmanager.NeverExpire
	}
	err = mng.manager.RefreshToken(args.Token, expire)
	if err != nil {
		log.Info("refresh token(%s) with expire(%ds) fail: %s", args.Token, args.ExpireAfterSeconds, err)
	} else {
		log.Debug("ok")
	}
	return err
}

// 立刻使token过期。
//
// 例子：
//
//     > curl http://127.0.0.1:23333/TokenManager?method=Expire -d '"ab56b4d92b40713acc5af89985d4b786c027b1ee301059618fb364abafd43f4a"'
//     0
func (mng *TokenManager) Expire(meta *gobus.HTTPMeta, token *string, reply *int) (err error) {
	log := mng.log.SubCode()
	log.Debug("expire token: %s", *token)
	err = mng.manager.ExpireToken(*token)
	if err != nil {
		log.Info("expire token(%s) fail: %s", *token, err)
	} else {
		log.Debug("ok")
	}
	return err
}

// 立刻使key对应的所有token过期。
//
// 例子：
//
//     > curl http://127.0.0.1:23333/TokenManager?method=ExpireAll -d '"ab56b4d92b40713acc5af89985d4b786"'
//     0
func (mng *TokenManager) ExpireAll(meta *gobus.HTTPMeta, key *string, reply *int) (err error) {
	log := mng.log.SubCode()
	log.Debug("expire all tokens: %s", *key)
	err = mng.manager.ExpireTokensByKey(*key)
	if err != nil {
		log.Info("expire all tokens(%s) fail: %s", *key, err)
	} else {
		log.Debug("ok")
	}
	return err
}
