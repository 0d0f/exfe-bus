package gobus

import (
	"fmt"
	"testing"
	"time"
)

/////////////////////////////////////////////////

type EmptyJob struct {
}

func (j *EmptyJob) Do(arg int, reply *int) error {
	*reply = arg * arg
	return nil
}

func TestCreateService(t *testing.T) {
	fmt.Println("Test create service")

	queue := "empty"
	service, err := CreateService("", 0, "", queue, &EmptyJob{})
	if err != nil {
		t.Fatal("Create service failed:", err)
	}
	defer func() { service.Close() }()
	go service.Serve(1e9)

	time.Sleep(0.5e9)
	if !service.IsRunning() {
		t.Fatal("Service doesn't run")
	}

	_ = service.Stop()
	if service.IsRunning() {
		t.Fatal("Service is still running")
	}

	service.Clear()
}

func TestCreateClient(t *testing.T) {
	fmt.Println("Test create client")

	queue := "empty"

	service, err := CreateService("", 0, "", queue, &EmptyJob{})
	if err != nil {
		t.Fatal("Create service failed:", err)
	}
	defer func() {
		service.Close()
		service.Clear()
	}()
	go service.Serve(1e9)

	client := CreateClient("", 0, "", queue)
	defer func() { client.Close() }()

	var reply int
	err = client.Do(3, &reply)
	if err != nil {
		t.Errorf("Return call should no error: %s", err)
	}
	if reply != 9 {
		t.Errorf("Reply should be 9, but got: %d", reply)
	}
	service.Stop()
}

/////////////////////////////////////////////////

type Arg struct {
	A string
}

type PtrJob struct {
}

func (j *PtrJob) Do(arg *Arg, reply *string) error {
	*reply = arg.A
	return nil
}

func TestPtrClient(t *testing.T) {
	fmt.Println("Test pointer client")

	queue := "empty"

	service, err := CreateService("", 0, "", queue, &PtrJob{})
	if err != nil {
		t.Fatal("Create service failed:", err)
	}
	defer func() {
		service.Close()
		service.Clear()
	}()
	go service.Serve(1e9)

	client := CreateClient("", 0, "", queue)
	defer func() { client.Close() }()

	var reply string
	err = client.Do(&Arg{
		A: "abc",
	}, &reply)
	if err != nil {
		t.Errorf("Return call should no error: %s", err)
	}
	if reply != "abc" {
		t.Errorf("Reply should be abc, but got: %d", reply)
	}
	service.Stop()
}

/////////////////////////////////////////////////

type BatchJob struct {
	data []int
}

func (j *BatchJob) Batch(args []int) {
	for _, i := range args {
		j.data = append(j.data, i)
	}
}

func TestBatchService(t *testing.T) {
	fmt.Println("Test batch service")

	queue := "batch"
	job := &BatchJob{}
	service, err := CreateBatchService("", 0, "", queue, job)
	if err != nil {
		t.Fatal("Create service failed:", err)
	}
	service.Clear()
	defer service.Close()
	go service.Serve(1e9)

	client := CreateClient("", 0, "", queue)
	defer client.Close()

	for i := 0; i < 10; i++ {
		client.Send(i)
	}

	time.Sleep(1e9)
	service.Stop()

	for i, d := range job.data {
		if i != d {
			t.Errorf("at index %d expect %d, but got %d", i, i, d)
		}
	}

	service.Clear()
}
