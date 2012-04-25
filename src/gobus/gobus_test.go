package gobus

import (
	"fmt"
	"testing"
	"time"
)

const serverName = "test"

func closeAndClearServer(q *Server) {
	q.Stop()
	q.Close()
	q.ClearQueue()
}

/////////////////////////////////////////////////

type EmptyJob struct {
}

func (j *EmptyJob) Test(arg int, reply *int) error {
	*reply = arg * arg
	return nil
}

func (j *EmptyJob) Batch(args []int) error {
	return nil
}

func TestCreateService(t *testing.T) {
	fmt.Println("Test create service")

	server := CreateServer("", 0, "", serverName)
	defer closeAndClearServer(server)

	server.Register(&EmptyJob{})
	go server.Serve(1e9)

	time.Sleep(0.5e9)
	if !server.IsRunning() {
		t.Fatal("server doesn't run")
	}

	_ = server.Stop()
	if server.IsRunning() {
		t.Fatal("server is still running")
	}
}

func TestCreateClient(t *testing.T) {
	fmt.Println("Test create client")

	server := CreateServer("", 0, "", serverName)
	defer closeAndClearServer(server)

	server.Register(&EmptyJob{})
	go server.Serve(1e9)

	client := CreateClient("", 0, "", serverName)

	var reply int
	err := client.Do("Test", 3, &reply, 3)
	if err != nil {
		t.Errorf("Return call should no error: %s", err)
	}
	if reply != 9 {
		t.Errorf("Reply should be 9, but got: %d", reply)
	}

	err = client.Send("Batch", 3, 5)
	if err != nil {
		t.Errorf("Return call should no error: %s", err)
	}

	err = client.Do("Batch", 3, &reply, 3)
	if err.Error() != "Can't find service: Batch(arg, reply)" {
		t.Errorf("Error should: Can't find service, but got: %s", err)
	}
}

/////////////////////////////////////////////////

type Arg struct {
	A string
}

type Job struct {
}

func (j *Job) PtrTest(arg *Arg, reply *string) error {
	*reply = arg.A
	return nil
}

func (j *Job) InstanceTest(arg Arg, reply *string) error {
	*reply = arg.A
	return nil
}

func TestPtrClient(t *testing.T) {
	fmt.Println("Test pointer client")

	server := CreateServer("", 0, "", serverName)
	defer closeAndClearServer(server)

	server.Register(&Job{})
	go server.Serve(1e9)

	client := CreateClient("", 0, "", serverName)

	var reply string

	reply = ""
	err := client.Do("PtrTest", &Arg{"abc"}, &reply, 3)
	if err != nil {
		t.Errorf("Return call should no error: %s", err)
	}
	if reply != "abc" {
		t.Errorf("Reply should be abc, but got: %d", reply)
	}

	reply = ""
	err = client.Do("InstanceTest", Arg{"abc"}, &reply, 3)
	if err != nil {
		t.Errorf("Return call should no error: %s", err)
	}
	if reply != "abc" {
		t.Errorf("Reply should be abc, but got: %d", reply)
	}
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

	server := CreateServer("", 0, "", serverName)
	defer closeAndClearServer(server)

	job := BatchJob{}
	server.Register(&job)
	go server.Serve(1e9)

	client := CreateClient("", 0, "", serverName)

	for i := 0; i < 10; i++ {
		client.Send("Batch", i, 3)
	}

	time.Sleep(1e9)

	for i, d := range job.data {
		if i != d {
			t.Errorf("at index %d expect %d, but got %d", i, i, d)
		}
	}
}

/////////////////////////////////////////////////

type ErrorJob struct {
}

func (j *ErrorJob) Error(arg Arg, reply *string) error {
	return fmt.Errorf("Internal Error")
}

func (j *ErrorJob) Panic(arg Arg, reply *string) error {
	panic("Fatal!")
	return fmt.Errorf("Internal Error")
}

func TestErrorClient(t *testing.T) {
	fmt.Println("Test error client")

	server := CreateServer("", 0, "", serverName)
	defer closeAndClearServer(server)
	server.Register(&ErrorJob{})
	go server.Serve(1e9)

	client := CreateClient("", 0, "", serverName)

	var reply string
	err := client.Do("Error", Arg{"abc"}, &reply, 3)
	if err == nil {
		t.Errorf("Return should return a error")
	}
	if err.Error() != "Internal Error" {
		t.Errorf("Return error: %s, expect: %s", err, "Internal Error")
	}

	defer func() {
		p := recover()
		if p != "Fatal!" {
			t.Errorf("Panic got: %s, expect: Fatal!", p)
		}
	}()

	client.Do("Panic", Arg{"abc"}, &reply, 3)
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

func (j *TestRetryJob) Batch(args []Arg) error {
	j.count++
	if j.count == 1 {
		j.result += args[0].A
		return fmt.Errorf("Retry")
	}
	if j.count == 2 {
		j.result += args[0].A
		return fmt.Errorf("Retry again")
	}
	for _, a := range args {
		j.result += a.A
	}
	return nil
}

func TestRetryClient(t *testing.T) {
	fmt.Println("Test retry client")

	job := TestRetryJob{}
	server := CreateServer("", 0, "", serverName)
	retry := DefaultRetryServer("", 0, "")
	defer closeAndClearServer(server)
	defer closeAndClearServer(retry)
	server.Register(&job)
	go retry.Serve(1e9)
	go server.Serve(0.1e9)

	client := CreateClient("", 0, "", serverName)

	{
		job.count = 0
		job.result = ""
		client.Send("Do", Arg{"abc"}, 4)

		time.Sleep(3e9)
		if job.result != "abc" {
			t.Fatal("Retry failed")
		}
	}

	{
		job.count = 0
		job.result = ""
		client.Send("Do", Arg{"123"}, 1)

		time.Sleep(3e9)
		if job.result != "" {
			t.Fatal("Retry failed")
		}
	}

	{
		job.count = 0
		job.result = ""
		client.Send("Batch", Arg{"x"}, 1)
		for i:=0; i<10; i++ {
			client.Send("Batch", Arg{fmt.Sprintf("%d", i)}, 4)
		}

		time.Sleep(3e9)
		fmt.Println(job.result)
		if job.result != "x00123456789" {
			t.Fatal("Retry failed")
		}
	}
}

/////////////////////////////////////////////////

type PhpJob struct {
	data []int
}

func (j *PhpJob) Batch(args []int) error {
	for _, i := range args {
		j.data = append(j.data, i)
	}
	return nil
}

func (j *PhpJob) Do(args int, reply *int) error {
	j.data = append(j.data, args)
	return nil
}

func TestPhpService(t *testing.T) {
	fmt.Println("Test batch service")

	server := CreateServer("", 0, "", "php")
	defer closeAndClearServer(server)

	job := PhpJob{}
	server.Register(&job)
	go server.Serve(0.1e9)

	time.Sleep(1e9)

	if len(job.data) != 2 {
		t.Errorf("need run php_test.php first")
		return
	}
	for i := range job.data {
		if job.data[i] != 3 {
			t.Errorf("data[%d] got: %d, expect: 3", i, job.data[i])
		}
	}
}
