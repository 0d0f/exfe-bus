package valve

import (
	"testing"
	"time"
)

type testWorker struct {
	i int
}

func (w testWorker) Do() (interface{}, error) {
	if w.i == 2 {
		time.Sleep(time.Second)
	}
	return w.i, nil
}

func TestValve(t *testing.T) {
	valve := New(3, time.Second)
	go valve.Serve()
	defer valve.Quit()

	var save []interface{}

	c := make(chan int)
	for i := 0; i < 10; i++ {
		go func(i int) {
			ret, err := valve.Do(testWorker{i})
			if err == nil {
				save = append(save, ret)
			} else {
				if err != QueueFull {
					t.Errorf("not expect error: %s", err)
				}
			}
			c <- 1
		}(i)
	}

	time.Sleep(time.Second / 2)
	if len(save) != 1 {
		t.Errorf("should return 1, but %d", len(save))
	}

	time.Sleep(time.Second)
	if len(save) != 2 {
		t.Errorf("should return 2, but %d", len(save))
	}

	time.Sleep(time.Second)
	if len(save) != 2 {
		t.Errorf("should return 2, but %d", len(save))
	}

	time.Sleep(time.Second)
	if len(save) != 3 {
		t.Errorf("should return 3, but %d", len(save))
	}

	for i := 0; i < 10; i++ {
		<-c
	}
}
