package main

import (
	"broker"
	"github.com/googollee/go-logger"
	"gobus"
	"iom"
	"model"
)

type Iom struct {
	handler *iom.Iom
	log     *logger.SubLogger
}

func NewIom(config *model.Config, redis broker.Redis) *Iom {
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
	*reply, err = iom.handler.Get(arg.UserID, arg.Hash)
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
	*reply, err = iom.handler.FindByData(arg.UserID, arg.Data)
	if err != nil {
		*reply, err = iom.handler.Create(arg.UserID, arg.Data)
	}
	return
}
