package main

import (
	"github.com/googollee/go-logger"
	"github.com/googollee/godis"
	"gobus"
	"iom"
)

type Iom struct {
	handler *iom.Iom
	log     *logger.SubLogger
}

func NewIom(config *Config, redis *godis.Client) *Iom {
	return &Iom{
		handler: iom.NewIom(redis),
		log:     config.Log.SubPrefix("iom"),
	}
}

type IomGetArg struct {
	UserID string `json:"user_id"`
	Hash   string `json:"hash"`
}

// 获取用户user_id名下的hash对应的资源。
//
// 例子：
//
//     > curl http://127.0.0.1:23333/Iom?method=GET -d '{"user_id":"124","hash":"aa"}'
//     "abc"
func (iom *Iom) GET(meta *gobus.HTTPMeta, arg *IomGetArg, reply *string) (err error) {
	log := iom.log.SubCode()
	log.Debug("get with user id: %s, hash: %s", arg.UserID, arg.Hash)
	*reply, err = iom.handler.Get(arg.UserID, arg.Hash)
	if err != nil {
		log.Info("get with user id: %s, hash: %s, failed: %s", arg.UserID, arg.Hash, err)
	} else {
		log.Debug("return data: %s", *reply)
	}
	return
}

type IomPostArg struct {
	UserID string `json:"user_id"`
	Data   string `json:"data"`
}

// 将资源data存入用户user_id名下，并返回对应的hash。如果资源data已经在user_id名下，则直接返回对应的hash。hash不区分大小写。
//
// 例子：
//
//     > curl http://127.0.0.1:23333/Iom?method=POST -d '{"user_id":"124","data":"abc"}'
//     "AA"
func (iom *Iom) POST(meta *gobus.HTTPMeta, arg *IomPostArg, reply *string) (err error) {
	log := iom.log.SubCode()
	log.Debug("post with user id: %s, data: %s", arg.UserID, arg.Data)
	*reply, err = iom.handler.FindByData(arg.UserID, arg.Data)
	if err != nil {
		*reply, err = iom.handler.Create(arg.UserID, arg.Data)
	}
	if err != nil {
		log.Info("post with user id: %s, data: %s, failed: %s", arg.UserID, arg.Data, err)
	} else {
		log.Debug("return hash: %s", *reply)
	}
	return
}
