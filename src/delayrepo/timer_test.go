package delayrepo

import (
	"broker"
	"fmt"
	"github.com/stretchrcom/testify/assert"
	"sort"
	"testing"
	"time"
)

type timeData struct {
	timer int64
	key   string
}

type timeDatas []timeData

func (s timeDatas) Len() int {
	return len(s)
}

func (s timeDatas) Less(i, j int) bool {
	return s[i].timer < s[j].timer
}

func (s timeDatas) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type FakeStorage struct {
	timer timeDatas
	array map[string][][]byte
}

func newFakeStorage() *FakeStorage {
	return &FakeStorage{
		timer: make(timeDatas, 0),
		array: make(map[string][][]byte),
	}
}

func (s *FakeStorage) Save(iupdate broker.UpdateType, ontime int64, key string, data []byte) error {
	exist := false
	for i := range s.timer {
		if s.timer[i].key == key {
			s.timer[i].timer = ontime
			exist = true
			break
		}
	}
	if !exist {
		s.timer = append(s.timer, timeData{
			timer: ontime,
			key:   key,
		})
	}
	sort.Sort(s.timer)
	s.array[key] = append(s.array[key], data)
	return nil
}

func (s *FakeStorage) Load(key string) ([][]byte, error) {
	return s.array[key], nil
}

func (s *FakeStorage) Ontime(key string) (int64, error) {
	for i := range s.timer {
		if s.timer[i].key == key {
			return s.timer[i].timer, nil
		}
	}
	return 0, nil
}

func (s *FakeStorage) Next() (string, error) {
	if len(s.timer) == 0 {
		return "", nil
	}
	return s.timer[0].key, nil
}

func TestTimer(t *testing.T) {
	s := newFakeStorage()
	timer, err := NewTimer(s)
	if err != nil {
		t.Fatal(err)
	}
	ontime := time.Now().Add(time.Second).Unix()
	err = timer.push(broker.Always, ontime, "123", []byte("a"))
	if err != nil {
		t.Fatal(err)
	}
	err = timer.push(broker.Always, ontime, "123", []byte("b"))
	if err != nil {
		t.Fatal(err)
	}
	err = timer.push(broker.Always, ontime, "123", []byte("c"))
	if err != nil {
		t.Fatal(err)
	}

	wait, err := timer.NextWakeup()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("waiting:", wait)

	time.Sleep(wait)

	key, data, err := timer.pop()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, key, "123")
	assert.Equal(t, fmt.Sprintf("%v", data), "[[97] [98] [99]]")
}

func TestEmptyTimer(t *testing.T) {
	s := newFakeStorage()
	timer, err := NewTimer(s)
	if err != nil {
		t.Fatal(err)
	}
	wait, err := timer.NextWakeup()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, wait < 0, true)
}

func TestTimerUpdate(t *testing.T) {
	s := newFakeStorage()
	timer, err := NewTimer(s)
	if err != nil {
		t.Fatal(err)
	}
	ontime := time.Now().Add(time.Second * 10).Unix()
	err = timer.push(broker.Always, ontime, "123", []byte("a"))
	if err != nil {
		t.Fatal(err)
	}
	wait, err := timer.NextWakeup()
	if err != nil {
		t.Fatal(err)
	}
	if wait < time.Second {
		t.Fatalf("wait too short: %s", wait)
	}

	ontime = time.Now().Unix()
	err = timer.push(broker.Always, ontime, "123", []byte("b"))
	if err != nil {
		t.Fatal(err)
	}
	wait, err = timer.NextWakeup()
	if err != nil {
		t.Fatal(err)
	}
	if wait > time.Second {
		t.Fatalf("wait too long: %s", wait)
	}

	time.Sleep(wait)

	key, data, err := timer.pop()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, key, "123")
	assert.Equal(t, fmt.Sprintf("%v", data), "[[97] [98]]")
}
