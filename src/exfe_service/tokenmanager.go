package main

import (
	"github.com/googollee/go-log"
	"github.com/googollee/go-mysql"
	"gobus"
	"time"
	"tokenmanager"
)

type TokenManager struct {
	manager *tokenmanager.TokenManager
	log     *log.Logger
	config  *Config
}

func NewTokenManager(config *Config, db *mysql.Client) (*TokenManager, error) {
	l, err := log.New(config.loggerOutput, "service bus.token manager", config.loggerFlags)
	if err != nil {
		return nil, err
	}
	return &TokenManager{
		manager: tokenmanager.New(db, config.TokenManager.TableName),
		log:     l,
		config:  config,
	}, nil
}

type TokenGenerateArgs struct {
	Resource           string `json:"resource"`
	ExpireAfterSeconds int    `json:"expire_after_seconds"`
}

func (mng *TokenManager) Generate(meta *gobus.HTTPMeta, arg *TokenGenerateArgs, reply *string) (err error) {
	serial := mng.log.SerialCode()
	mng.log.Debug("(%d)generate with resource: %s, expire: %ds", serial, arg.Resource, arg.ExpireAfterSeconds)
	expire := time.Duration(arg.ExpireAfterSeconds) * time.Second
	if arg.ExpireAfterSeconds < 0 {
		expire = tokenmanager.NeverExpire
	}
	*reply, err = mng.manager.GenerateToken(arg.Resource, expire)
	if err != nil {
		mng.log.Warning("(%d)generate token fail: %s", serial, err)
	} else {
		mng.log.Debug("(%d)return token: %s", serial, *reply)
	}
	return err
}

type TokenGetReply struct {
	Resource  string `json:"resource"`
	IsExpired bool   `json:"is_expired"`
}

func (mng *TokenManager) Get(meta *gobus.HTTPMeta, token *string, reply *TokenGetReply) (err error) {
	serial := mng.log.SerialCode()
	mng.log.Debug("(%d)get with token: %s", serial, *token)
	reply.Resource, err = mng.manager.GetResource(*token)
	reply.IsExpired = err == tokenmanager.ExpiredError
	if err == tokenmanager.ExpiredError {
		err = nil
	}
	if err != nil {
		mng.log.Warning("(%d)get resource fail: %s", serial, err)
	} else {
		mng.log.Debug("(%d)return resource: %s, is expired: %v", serial, reply.Resource, reply.IsExpired)
	}
	return err
}

type TokenVerifyArg struct {
	Token    string `json:"token"`
	Resource string `json:"resource"`
}

type TokenVerifyReply struct {
	Matched   bool `json:"matched"`
	IsExpired bool `json:"is_expired"`
}

func (mng *TokenManager) Verify(meta *gobus.HTTPMeta, args *TokenVerifyArg, reply *TokenVerifyReply) (err error) {
	serial := mng.log.SerialCode()
	mng.log.Debug("(%d)verify with token: %s, resource: %s", serial, args.Token, args.Resource)
	reply.Matched, err = mng.manager.VerifyToken(args.Token, args.Resource)
	reply.IsExpired = err == tokenmanager.ExpiredError
	if err == tokenmanager.ExpiredError {
		err = nil
	}
	if err != nil {
		mng.log.Warning("(%d)verify fail: %s", serial, err)
	} else {
		mng.log.Debug("(%d)return verify matched: %v, is expired: %v", serial, reply.Matched, reply.IsExpired)
	}
	return err
}

func (mng *TokenManager) Delete(meta *gobus.HTTPMeta, token *string, reply *int) (err error) {
	serial := mng.log.SerialCode()
	mng.log.Debug("(%d)delete token: %s", serial, *token)
	err = mng.manager.DeleteToken(*token)
	if err != nil {
		mng.log.Warning("(%d)delete fail: %s", serial, err)
	} else {
		mng.log.Debug("(%d)ok", serial)
	}
	return err
}

type TokenRefreshArg struct {
	Token              string `json:"token"`
	ExpireAfterSeconds int    `json:"expire_after_seconds"`
}

func (mng *TokenManager) Refresh(meta *gobus.HTTPMeta, args *TokenRefreshArg, reply *int) (err error) {
	serial := mng.log.SerialCode()
	mng.log.Debug("(%d)refresh token: %s, expire: %ds", serial, args.Token, args.ExpireAfterSeconds)
	expire := time.Duration(args.ExpireAfterSeconds) * time.Second
	if args.ExpireAfterSeconds < 0 {
		expire = tokenmanager.NeverExpire
	}
	err = mng.manager.RefreshToken(args.Token, expire)
	if err != nil {
		mng.log.Warning("(%d)refresh fail: %s", serial, err)
	} else {
		mng.log.Debug("(%d)ok", serial)
	}
	return err
}

func (mng *TokenManager) Expire(meta *gobus.HTTPMeta, token *string, reply *int) (err error) {
	serial := mng.log.SerialCode()
	mng.log.Debug("(%d)expire token: %s", serial, *token)
	err = mng.manager.ExpireToken(*token)
	if err != nil {
		mng.log.Warning("(%d)expire fail: %s", serial, err)
	} else {
		mng.log.Debug("(%d)ok", serial)
	}
	return err
}
