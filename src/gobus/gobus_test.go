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

type InstanceJob struct {
}

func (j *InstanceJob) Do(arg Arg, reply *string) error {
	*reply = arg.A
	return nil
}

func TestInstanceClient(t *testing.T) {
	fmt.Println("Test instance client")

	queue := "empty"

	service, err := CreateService("", 0, "", queue, &InstanceJob{})
	if err != nil {
		t.Fatal("Create service failed:", err)
	}
	defer func() {
		service.Close()
		service.Clear()
	}()
	go service.Serve(1e9)

	client := CreateClient("", 0, "", queue)

	var reply string
	err = client.Do(Arg{
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

func (j *BatchJob) Batch(args []int) error {
	for _, i := range args {
		j.data = append(j.data, i)
	}
	return nil
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

	for i := 0; i < 10; i++ {
		client.Send(i, 3)
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

/////////////////////////////////////////////////

type ErrorJob struct {
}

func (j *ErrorJob) Do(arg Arg, reply *string) error {
	return fmt.Errorf("Internal Error")
}

func TestErrorClient(t *testing.T) {
	fmt.Println("Test error client")

	queue := "empty"

	service, err := CreateService("", 0, "", queue, &ErrorJob{})
	if err != nil {
		t.Fatal("Create service failed:", err)
	}
	defer func() {
		service.Close()
		service.Clear()
	}()
	go service.Serve(1e9)

	client := CreateClient("", 0, "", queue)

	var reply string
	err = client.Do(Arg{
		A: "abc",
	}, &reply)
	if err == nil {
		t.Errorf("Return should return a error")
	}
	if err.Error() != "Internal Error" {
		t.Errorf("Return error: %s, expect: %s", err, "Internal Error")
	}
	service.Stop()
}

/////////////////////////////////////////////////

type PanicJob struct {
}

func (j *PanicJob) Do(arg Arg, reply *string) error {
	panic("Fatal!")
	return fmt.Errorf("Internal Error")
}

func TestPanicClient(t *testing.T) {
	fmt.Println("Test panic client")

	queue := "empty"

	service, err := CreateService("", 0, "", queue, &PanicJob{})
	if err != nil {
		t.Fatal("Create service failed:", err)
	}
	defer func() {
		service.Close()
		service.Clear()
	}()
	go service.Serve(1e9)

	client := CreateClient("", 0, "", queue)

	defer func() {
		p := recover()
		if p != "Fatal!" {
			t.Errorf("Panic got: %s, expect: Fatal!", p)
		}
		service.Stop()
	}()

	var reply string
	err = client.Do(Arg{
		A: "abc",
	}, &reply)
}

/////////////////////////////////////////////////

type TestRetryJob struct {
	count int
	result string
}

func (j *TestRetryJob) Do(arg Arg, reply *string) error {
	j.count++
	if j.count == 1 {
		return fmt.Errorf("Retry")
	}
	if j.count == 2 {
		return fmt.Errorf("Retry again")
	}
	j.result = arg.A
	return nil
}

func TestRetryClient(t *testing.T) {
	fmt.Println("Test retry client")

	queue := "empty"

	job := TestRetryJob{}
	service, err := CreateService("", 0, "", queue, &job)
	if err != nil {
		t.Fatal("Create service failed:", err)
	}

	retry, err := GetDefaultRetryServer("", 0, "")
	if err != nil {
		t.Fatal("Create retry service failed:", err)
	}
	defer func() {
		service.Close()
		retry.Close()
	}()
	go retry.Serve(1e9)
	go service.Serve(0.1e9)

	client := CreateClient("", 0, "", queue)

	defer func() {
		service.Stop()
	}()

	job.count = 0
	err = client.Send(Arg{
		A: "abc",
	}, 4)

	time.Sleep(3e9)
	if job.result != "abc" {
		t.Fatal("Retry failed")
	}

	job.count = 0
	err = client.Send(Arg{
		A: "123",
	}, 1)

	time.Sleep(3e9)
	if job.result != "abc" {
		t.Fatal("Retry failed")
	}
}

