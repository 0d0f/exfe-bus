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

func (iom *Iom) SetRoute(route gobus.RouteCreater) {
	json := new(gobus.JSON)
	route().Methods("GET").Path("/iom/{user_id}/{hash}").HandlerFunc(gobus.Must(gobus.Method(json, iom, "Get")))
	route().Methods("POST").Path("/iom/user/{user_id}").HandlerFunc(gobus.Must(gobus.Method(json, iom, "Create")))
}

// 获取用户user_id名下的hash对应的资源。
//
// 例子：
//
//     > curl http://127.0.0.1:23333/iom/124/aa -d '{"user_id":"124","hash":"aa"}'
//     "abc"
func (iom *Iom) Get(params map[string]string) (string, error) {
	userID := params["user_id"]
	hash := params["hash"]
	ret, err := iom.handler.Get(userID, hash)
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
//     > curl http://127.0.0.1:23333/iom/user/124 -d '"abc"'
//     "AA"
func (iom *Iom) Create(params map[string]string, data string) (string, error) {
	userID := params["user_id"]
	ret, err := iom.handler.FindByData(userID, data)
	if err != nil {
		ret, err = iom.handler.Create(userID, data)
	}
	return ret, err
}
