package main

import (
	"github.com/googollee/go-logger"
	"github.com/googollee/go-mysql"
	"gobus"
	"time"
	"tokenmanager"
)

type TokenManager struct {
	manager *tokenmanager.TokenManager
	log     *logger.SubLogger
	config  *Config
}

func NewTokenManager(config *Config, db *mysql.Client) (*TokenManager, error) {
	return &TokenManager{
		manager: tokenmanager.New(db, config.TokenManager.TableName),
		log:     config.Log.SubPrefix("token manager"),
		config:  config,
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
//     "deae3cee0be68e2ae2c590f0a1b5bb032168477d2d2c2a515b652042331b0220"
func (mng *TokenManager) Generate(meta *gobus.HTTPMeta, arg *TokenGenerateArgs, reply *string) (err error) {
	log := mng.log.SubCode()
	log.Debug("generate with %+v", arg)
	expire := time.Duration(arg.ExpireAfterSeconds) * time.Second
	if arg.ExpireAfterSeconds < 0 {
		expire = tokenmanager.NeverExpire
	}
	*reply, err = mng.manager.GenerateToken(arg.Resource, arg.Data, expire)
	if err != nil {
		log.Info("generate token with resource(%s), expire(%ds) fail: %s", arg.Resource, arg.ExpireAfterSeconds, err)
	} else {
		log.Debug("return token: %s", *reply)
	}
	return err
}

type TokenGetReply struct {
	Resource  string `json:"resource"`
	Data      string `json:"data"`
	IsExpired bool   `json:"is_expired"`
}

// 根据token返回资源resource，数据data和是否过期is_expired
//
// 例子：
//
//     > curl http://127.0.0.1:23333/TokenManager?method=Get -d '"deae3cee0be68e2ae2c590f0a1b5bb032168477d2d2c2a515b652042331b0220"'
//     {"resource":"abcde","data":"","is_expired":true}
func (mng *TokenManager) Get(meta *gobus.HTTPMeta, token *string, reply *TokenGetReply) (err error) {
	log := mng.log.SubCode()
	log.Debug("get with token: %s", *token)
	reply.Resource, reply.Data, err = mng.manager.GetResource(*token)
	reply.IsExpired = err == tokenmanager.ExpiredError
	if err == tokenmanager.ExpiredError {
		err = nil
	}
	if err != nil {
		log.Info("get resource with token(%s) fail: %s", *token, err)
	} else {
		log.Debug("return resource: %s, is expired: %v", reply.Resource, reply.IsExpired)
	}
	return err
}

// 更新token对应的数据data
//
// 例子：
//
// > curl http://127.0.0.1:23333/TokenManager?method=Update -d '{"token":"deae3cee0be68e2ae2c590f0a1b5bb032168477d2d2c2a515b652042331b0220","data":"123"}'
// 0

type TokenUpdateArg struct {
	Token string `json:"token"`
	Data  string `json:"data"`
}

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
	Matched   bool   `json:"matched"`
	IsExpired bool   `json:"is_expired"`
	Data      string `json:"data"`
}

// 根据token和资源resource来验证两者是否一致matched，并返回token是否过期is_expired和token对应的数据data
//
// 例子：
//
//     > curl http://127.0.0.1:23333/TokenManager?method=Verify -d '{"token":"deae3cee0be68e2ae2c590f0a1b5bb032168477d2d2c2a515b652042331b0220","resource":"abcde"}'
//     {"matched":true,"is_expired":false,"data":""}
func (mng *TokenManager) Verify(meta *gobus.HTTPMeta, args *TokenVerifyArg, reply *TokenVerifyReply) (err error) {
	log := mng.log.SubCode()
	log.Debug("verify with token: %s, resource: %s", args.Token, args.Resource)
	reply.Matched, reply.Data, err = mng.manager.VerifyToken(args.Token, args.Resource)
	reply.IsExpired = err == tokenmanager.ExpiredError
	if err == tokenmanager.ExpiredError {
		err = nil
	}
	if err != nil {
		log.Info("verify with token(%s)&resource(%s) fail: %s", args.Token, args.Resource, err)
	} else {
		log.Debug("return verify matched: %v, is expired: %v", reply.Matched, reply.IsExpired)
	}
	return err
}

// 删除token，如果成功返回0
//
// 例子：
//
//     > curl http://127.0.0.1:23333/TokenManager?method=Delete -d '"deae3cee0be68e2ae2c590f0a1b5bb032168477d2d2c2a515b652042331b0220"'
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
//     > curl http://127.0.0.1:23333/TokenManager?method=Refresh -d '{"token":"deae3cee0be68e2ae2c590f0a1b5bb032168477d2d2c2a515b652042331b0220","expire_after_seconds":-1}'
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
//     > curl http://127.0.0.1:23333/TokenManager?method=Expire -d '"deae3cee0be68e2ae2c590f0a1b5bb032168477d2d2c2a515b652042331b0220"'
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
