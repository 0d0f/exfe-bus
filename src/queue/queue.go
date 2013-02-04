package main

import (
	"broker"
	"delayrepo"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
	"gobus"
	"launchpad.net/tomb"
	"model"
	"strconv"
	"strings"
)

type Queue struct {
	instantCallback func(string, [][]byte)
	heads           map[uint]*delayrepo.Head
	tails           map[uint]*delayrepo.Tail
	config          *model.Config
	priority        map[string]uint
	dispatcher      *gobus.Dispatcher
	log             *logger.SubLogger
	tombs           []*tomb.Tomb
}

func NewQueue(config *model.Config, redis broker.Redis) (*Queue, error) {
	table, err := gobus.NewTable(config.Dispatcher)
	if err != nil {
		return nil, err
	}
	ret := &Queue{
		heads:      make(map[uint]*delayrepo.Head),
		tails:      make(map[uint]*delayrepo.Tail),
		config:     config,
		priority:   config.ExfeQueue.Priority,
		dispatcher: gobus.NewDispatcher(table),
		log:        config.Log.SubPrefix("queue"),
		tombs:      make([]*tomb.Tomb, 0),
	}
	ret.instantCallback = ret.callback("instant")
	for _, delay := range ret.priority {
		if delay == 0 {
			continue
		}

		{
			name := fmt.Sprintf("delayrepo:head_%ds", delay)
			repo := delayrepo.NewHead(name, delay, redis)
			log := config.Log.SubPrefix(name)
			tomb := delayrepo.ServRepository(log, repo, ret.callback(name))
			ret.tombs = append(ret.tombs, tomb)
			ret.heads[delay] = repo
		}

		{
			name := fmt.Sprintf("delayrepo:tail_%ds", delay)
			repo := delayrepo.NewTail(name, delay, redis)
			log := config.Log.SubPrefix(name)
			tomb := delayrepo.ServRepository(log, repo, ret.callback(name))
			ret.tombs = append(ret.tombs, tomb)
			ret.tails[delay] = repo
		}
	}
	return ret, nil
}

func (q *Queue) SetRoute(r gobus.RouteCreater) error {
	json := new(gobus.JSON)
	return r().Methods("POST").Path("/").HandlerMethod(json, q, "Push")
}

type Push struct {
	Service    string            `json:"service"`
	Priority   string            `json:"priority"`
	Delay      string            `json:"delay"`
	GroupKey   string            `json:"group_key"`
	Recipients []model.Recipient `json:"recipients"`
	Data       interface{}       `json:"data"`
}

func (a Push) String() string {
	return fmt.Sprintf("{service:%s priority:%s type:%s key:%s recipients:%s}", a.Service, a.Priority, a.Delay, a.GroupKey, a.Recipients)
}

// 将data以delay type的合并方式，放入priority队列，之后发送给service服务。合并关键字group key，接收者recipients
//
// priority取值：
//
// - instant
// - urgent
// - normal
//
// priority也可以使用数字，表示经过多少秒后合并发送。
//
// 例子：
//
// > curl 'http://127.0.0.1:23334/' -d '{"service":"bus://exfe_service/notifier/conversation",
//       "priority": "urgent",
//       "delay_type": "head",
//       "group_key":"email_cross123",
//       "recipients":[{"identity_id":12,"user_id":3,"name":"email1 name","auth_data":"","timezone":"+0800","token":"recipient_email1_token","language":"en_US","provider":"email","external_id":"sender1@gmail.com","external_username":"sender1@gmail.com"},{"identity_id":12,"user_id":3,"name":"email1 name","auth_data":"","timezone":"+0800","token":"recipient_email1_token","language":"en_US","provider":"email","external_id":"sender2@hotmail.com","external_username":"sender2@hotmail.com"}],
//       "data":{"cross":{"id":123,"by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"sender2@hotmail.com","external_username":"sender2@hotmail.com"},"title":"Test Cross","description":"test cross description","time":{"begin_at":{"date_word":"","date":"","time_word":"","time":"","timezone":""},"origin":"","output_format":0},"place":{"id":0,"title":"","description":"","lng":"","lat":"","provider":"","external_id":""},"exfee":{"id":0,"name":"","invitations":null}},"post":{"id":1,"by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"sender2@hotmail.com","external_username":"sender2@hotmail.com"},"content":"email1 post sth","via":"abc","created_at":"2012-10-24 16:31:00"}}}'
func (q *Queue) Push(param map[string]string, arg Push) (int, error) {
	if len(arg.Recipients) == 0 {
		return -1, fmt.Errorf("recipient is empty")
	}
	delay, ok := q.priority[arg.Priority]
	if !ok {
		d, err := strconv.ParseUint(arg.Priority, 10, 32)
		if err != nil {
			return -1, err
		}
		delay = uint(d)
	}

	if delay == 0 {
		return q.instant(arg.Service, arg.GroupKey, arg.Recipients, arg.Data)
	}
	var repo delayrepo.Repository = nil
	switch arg.Delay {
	case "head":
		repo = q.heads[delay]
	case "tail":
		repo = q.tails[delay]
	}
	if repo == nil {
		return -1, fmt.Errorf("invalid delay type: %s", arg.Delay)
	}
	return q.delay(repo, arg.Service, arg.GroupKey, arg.Recipients, arg.Data)
}

func (q *Queue) instant(service, groupKey string, recipients []model.Recipient, data interface{}) (int, error) {
	ret := len(recipients)
	keys := make([]string, ret)
	for index, to := range recipients {
		keys[index] = fmt.Sprintf("%s,%s,%s(%s)@%s", service, groupKey, to.ExternalID, to.ExternalUsername, to.Provider)
	}

	datas := make([][]byte, ret)
	d, ok := data.(map[string]interface{})
	ret = 0
	for i, _ := range keys {
		if ok {
			d["to"] = recipients[i]
		}
		var err error
		datas[i], err = json.Marshal(d)
		if err != nil {
			return ret, err
		}
		ret++
	}

	go func() {
		for index, key := range keys {
			q.instantCallback(key, [][]byte{datas[index]})
		}
	}()

	return ret, nil
}

func (q *Queue) delay(repo delayrepo.Repository, service, groupKey string, recipients []model.Recipient, data interface{}) (int, error) {
	ret := 0
	d, ok := data.(map[string]interface{})
	for _, to := range recipients {
		if ok {
			d["to"] = to
		}
		b, err := json.Marshal(d)
		if err != nil {
			return ret, fmt.Errorf("can't marshal input data: %s", err)
		}
		key := fmt.Sprintf("%s,%s,%s(%s)@%s", service, groupKey, to.ExternalID, to.ExternalUsername, to.Provider)
		err = repo.Push(key, b)
		if err != nil {
			return ret, fmt.Errorf("push to repo failed: %s", err)
		}
		ret++
	}

	return ret, nil
}

func (q *Queue) callback(name string) func(string, [][]byte) {
	log := q.log.SubPrefix(name)
	return func(key string, datas [][]byte) {
		names := strings.SplitN(key, ",", 3)
		if len(names) != 3 {
			log.Crit("can't split service and method from key: %s", key)
			return
		}
		service, key := names[0], names[1]

		arg := make([]interface{}, 0)
		for _, data := range datas {
			var d interface{}
			err := json.Unmarshal(data, &d)
			if err != nil {
				log.Err("can't unmarshal %s(%+v)", err, data)
				continue
			}
			if key != "" {
				arg = append(arg, d)
			} else {
				var i interface{}
				err := q.dispatcher.Do(service, "POST", arg, &i)
				if err != nil {
					j, _ := json.Marshal(arg)
					log.Err("call %s failed(%s) with %s", service, err, string(j))
				}
			}
		}
		if key != "" {
			var i interface{}
			err := q.dispatcher.Do(service, "POST", arg, &i)
			if err != nil {
				j, _ := json.Marshal(arg)
				log.Err("call %s failed(%s) with %s", service, err, string(j))
			}
		}
	}
}
