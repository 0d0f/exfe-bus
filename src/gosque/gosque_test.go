package gosque

import (
	"fmt"
	"time"
	"testing"
)

type Test struct {
	t *testing.T
	count int
}

func (t *Test) Perform(arg string) {
	expect := fmt.Sprintf("Job-%d", t.count)
	got := arg
	if got != expect {
		t.t.Errorf("Failed arg: expect: %s, got: %s", expect, got)
	}
	t.count++
}

func GenerateJobs(gosque *Client, t *testing.T) {
	for i := 0; i < 10; i++ {
		job := fmt.Sprintf("Job-%d", i)
		err := gosque.Enqueue("test_queue", "test", job)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func TestGetJob(t *testing.T) {
	gosque := CreateQueue("", 0, "", "test_queue")

	test := &Test{
		t: t,
	}
	fmt.Println("here")
	gosque.Register(test)
	defer func() { gosque.Close() }()
	GenerateJobs(gosque, t)

	go gosque.Serve(1e9)

	time.Sleep(5e9)

	if test.count != 10 {
		t.Errorf("Job count should be 10, got: %d", test.count)
	}
}
