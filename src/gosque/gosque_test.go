package gosque

import (
	"fmt"
	"time"
	"testing"
)

type Test1 struct {
	t *testing.T
	count int
}

type Test1Arg struct {
	Arg string
}

func (t *Test1Arg) Queue() string {
	return "test_queue"
}

func (t *Test1) Perform(arg *Test1Arg) {
	expect := fmt.Sprintf("Job1-%d", t.count)
	got := arg.Arg
	if got != expect {
		t.t.Errorf("Failed arg: expect: %s, got: %s", expect, got)
	}
	t.count++
}

type Test2 struct {
	t *testing.T
	count int
}

type Test2Arg struct {
	Arg string
}

func (t *Test2Arg) Queue() string {
	return "test_queue"
}

func (t *Test2) Perform(arg Test2Arg) {
	expect := fmt.Sprintf("Job2-%d", t.count)
	got := arg.Arg
	if got != expect {
		t.t.Errorf("Failed arg: expect: %s, got: %s", expect, got)
	}
	t.count++
}

func GenerateJobs(gosque *Client, t *testing.T) {
	for i := 0; i < 10; i++ {
		job := fmt.Sprintf("Job1-%d", i)
		err := gosque.Enqueue("test1", &Test1Arg{job})
		if err != nil {
			fmt.Println(err)
		}

		job = fmt.Sprintf("Job2-%d", i)
		err = gosque.Enqueue("test2", &Test2Arg{job})
		if err != nil {
			fmt.Println(err)
		}
	}
}

func TestGetJob(t *testing.T) {
	gosque := CreateQueue("", 0, "", "test_queue")

	test1 := &Test1{
		t: t,
	}
	test2 := &Test2{
		t: t,
	}

	err := gosque.Register(test1)
	t.Log(err)
	gosque.Register(test2)
	defer func() { gosque.Close() }()

	GenerateJobs(gosque, t)

	go gosque.Serve(1e9)

	time.Sleep(5e9)

	if test1.count != 10 {
		t.Errorf("Job count should be 10, got: %d", test1.count)
	}
	if test2.count != 10 {
		t.Errorf("Job count should be 10, got: %d", test2.count)
	}
}
