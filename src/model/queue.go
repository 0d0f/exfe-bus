package model

import (
	"fmt"
	"strconv"
	"strings"
)

type QueuePush struct {
	Service  string      `json:"service"`
	Method   string      `json:"method"`
	MergeKey string      `json:"merge_key"`
	Priority string      `json:"priority"`
	Tos      []Recipient `json:"tos"` // it will expand and overwrite "to" field in data
	Data     interface{} `json:"data"`

	Ontime int64 `json:"-"`
}

func (q *QueuePush) String() string {
	return fmt.Sprintf("{service:%s method:%s key:%s priority:%s to:%s}", q.Service, q.Method, q.MergeKey, q.Priority, q.Tos)
}

func (q *QueuePush) Init(priority map[string]uint) error {
	ontime, ok := priority[q.Priority]
	if ok {
		q.Ontime = int64(ontime)
	} else {
		var err error
		q.Ontime, err = strconv.ParseInt(q.Priority, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid priority(%s): %s", q.Priority, err)
		}
	}
	return nil
}

func (q *QueuePush) Key(recipient Recipient) string {
	return fmt.Sprintf("%s,%s,%s,%s", q.Service, q.Method, q.MergeKey, recipient.ID())
}

func QueueParseKey(key string) (service, method, merge_key, recipient string, err error) {
	splits := strings.SplitN(key, ",", 4)
	if len(splits) != 4 {
		err = fmt.Errorf("invalid key(%s)", key)
		return
	}
	service, method, method, recipient = splits[0], splits[1], splits[2], splits[3]
	return
}

type EachArg struct {
	Key  string
	Data interface{}
}

func (q *QueuePush) Each() chan EachArg {
	c := make(chan EachArg)
	mapData, ok := q.Data.(map[string]interface{})

	go func() {
		for _, to := range q.Tos {
			data := q.Data
			if ok {
				mapData["to"] = to
				data = mapData
			}
			c <- EachArg{q.Key(to), data}
		}
		close(c)
	}()

	return c
}
