package gobus

import (
	"github.com/googollee/godis"
	"time"
)

const interval = 2

var redis = godis.New("", 0, "")
var duration = interval * time.Second

type testData struct {
	Id   string
	Data int
}

func (d testData) KeyForQueue() string {
	return d.Id
}
