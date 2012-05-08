package twitter_job

import (
	"time"
	"gobus"
	"fmt"
	"github.com/simonz05/godis"
)

type Cross_update_merge struct {
	Config *Config
	Client *gobus.Client
	Queue *gobus.LastDelayQueue
}

const CrossUpdateMergeQueue = "gobus:delayqueue:crossupdate"

func NewCrossUpdateMerge(redis *godis.Client) *Cross_update_merge {
	t := new(CrossUpdateArg)
	return &Cross_update_merge{
		Queue: gobus.NewLastDelayQueue(CrossUpdateMergeQueue, 30, t, redis),
	}
}

func (s *Cross_update_merge) Serve() {
	for {
		t, err := s.Queue.NextWakeup()
		if err != nil {
			fmt.Println("Last delay queue failed:", err)
			break
		}
		time.Sleep(t)
		args, err := s.Queue.Pop()
		if len(args) > 0 {
			s.Batch(args)
		}
	}
}

func (s *Cross_update_merge) Batch(args []interface{}) {
}
