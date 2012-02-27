package gobus

import (
	"fmt"
	"testing"
	"time"
)

func ValueGenerator() interface{} {
	var ret int
	return &ret
}

/////////////////////////////////////////////////

type EmptyService struct {
}

func (s *EmptyService) Do(jobs []interface{}) []interface{} {
	fmt.Println("in empty service do")
	fmt.Println(jobs)
	return append(make([]interface{}, 0), 1)
}

func (s *EmptyService) MaxJobsCount() int {
	return 1
}

func (s *EmptyService) JobGenerator() interface{} {
	var ret int
	return &ret
}

func TestCreateService(t *testing.T) {
	fmt.Println("Test create service")

	queue := "gobus:queue:empty"
	service := CreateService("", 0, "", queue, &EmptyService{})
	defer func() { service.Close() }()
	_ = service.Run(1e9)
	if !service.IsRunning() {
		t.Fatal("Service doesn't run")
	}
	_ = service.Stop()
	time.Sleep(0.5e9)
	if service.IsRunning() {
		t.Fatal("Service is still running")
	}

	service.Empty()
}

func TestCreateClient(t *testing.T) {
	fmt.Println("Test create client")

	queue := "gobus:queue:empty"
	client := CreateClient("", 0, "", queue, ValueGenerator)
	defer func() { client.Close() }()
}

////////////////////////////////////////////////////

type PerJobService struct {
}

func (s *PerJobService) Do(jobs []interface{}) []interface{} {
	if len(jobs) > 1 {
		fmt.Println("PerJobService get jobs count > 1")
	}
	data := *jobs[0].(*int)
	var ret int
	ret = data * data
	return append(make([]interface{}, 0), &ret)
}

func (s *PerJobService) MaxJobsCount() int {
	return 1
}

func (s *PerJobService) JobGenerator() interface{} {
	var ret int
	return &ret
}

func TestPerJobWork(t *testing.T) {
	fmt.Println("Test per job work")

	queue := "gobus:queue:perjobwork"
	service := CreateService("", 0, "", queue, &PerJobService{})
	defer func() {
		service.Close()
		service.Empty()
	}()
	_ = service.Run(1e9)

	{
		client := CreateClient("", 0, "", queue, ValueGenerator)
		defer func() { client.Close() }()


		d := 2
		ret_, err := client.Do(d)
		ret := *ret_.(*int)
		if err != nil {
			t.Fatal(err)
		}
		expect := 4
		if ret != expect {
			t.Error("Expect ", expect, ", got ", ret)
		}
	}

	client := CreateClient("", 0, "", queue, ValueGenerator)
	defer func() { client.Close() }()

	{
		d := 3
		ret_, _ := client.Do(d)
		ret := *ret_.(*int)
		expect := 9
		if ret != expect {
			t.Error("Expect ", expect, ", got ", ret)
		}
	}

	{
		d := -1
		ret_, _ := client.Do(d)
		ret := *ret_.(*int)
		expect := 1
		if ret != expect {
			t.Error("Expect ", expect, ", got ", ret)
		}
	}
}

func TestPerJobTime(t *testing.T) {
	fmt.Println("Test per job time")

	queue := "gobus:queue:perjobtime"
	service := CreateService("", 0, "", queue, &PerJobService{})
	defer func() {
		service.Close()
		service.Empty()
	}()

	client := CreateClient("", 0, "", queue, ValueGenerator)
	defer func() { client.Close() }()

	_ = service.Run(1e9)

	retChan := make(chan int)
	start := time.Now()
	for i:=0; i< 10; i++ {
		go func() {
			ret, _ := client.Do(2)
			retChan<-(*ret.(*int))
		}()
	}

	for i:=0; i< 10; i++ {
		ret := <-retChan
		if ret != 4 {
			t.Error("Got %d, expect: 4", ret)
		}
	}

	stop := time.Now()

	duration := stop.Sub(start)
	if duration > 2 * 1e9 {
		t.Fatal("Run time too long")
	}
}

////////////////////////////////////////////////////

type Per4JobsService struct {
}

func (s *Per4JobsService) Do(jobs []interface{}) []interface{} {
	if len(jobs) > 4 {
		fmt.Println("Per4JobsService get jobs count > 4")
	}

	rets := make([]interface{}, 0)
	for _, job := range jobs {
		data := *job.(*int)
		result := data * data
		rets = append(rets, &result)
	}
	return rets
}

func (s *Per4JobsService) MaxJobsCount() int {
	return 4
}

func (s *Per4JobsService) JobGenerator() interface{} {
	var ret int
	return &ret
}

func TestPer4JobsWork(t *testing.T) {
	fmt.Println("Test per 4 jobs work")

	queue := "gobus:queue:per4jobswork"
	service := CreateService("", 0, "", queue, &Per4JobsService{})
	defer func() {
		service.Close()
		service.Empty()
	}()

	client := CreateClient("", 0, "", queue, ValueGenerator)
	defer func() { client.Close() }()

	_ = service.Run(1e9)

	{
		d := 2
		ret_, _ := client.Do(d)
		ret := *ret_.(*int)
		expect := 4
		if ret != expect {
			t.Error("Expect ", expect, ", got ", ret)
		}
	}

	{
		d := 3
		ret_, _ := client.Do(d)
		ret := *ret_.(*int)
		expect := 9
		if ret != expect {
			t.Error("Expect ", expect, ", got ", ret)
		}
	}

	{
		d := -1
		ret_, _ := client.Do(d)
		ret := *ret_.(*int)
		expect := 1
		if ret != expect {
			t.Error("Expect ", expect, ", got ", ret)
		}
	}
}

func TestPer4JobsTime(t *testing.T) {
	fmt.Println("Test per 4 jobs time")

	queue := "gobus:queue:per4jobstime"
	service := CreateService("", 0, "", queue, &Per4JobsService{})
	defer func() {
		service.Close()
		service.Empty()
	}()

	client := CreateClient("", 0, "", queue, ValueGenerator)
	defer func() { client.Close() }()

	_ = service.Run(1e9)

	retChan := make(chan int)
	{
		start := time.Now()
		for i:=0; i< 4; i++ {
			go func() {
				ret, _ := client.Do(2)
				retChan<-(*ret.(*int))
			}()
		}

		for i:=0; i< 4; i++ {
			<-retChan
		}

		stop := time.Now()

		duration := stop.Sub(start)
		if duration > 2 * 1e9 {
			t.Fatal("Run time too long")
		}
	}

	{
		start := time.Now()
		for i:=0; i< 10; i++ {
			go func() {
				ret, _ := client.Do(2)
				retChan<-(*ret.(*int))
			}()
		}

		for i:=0; i< 10; i++ {
			<-retChan
		}

		stop := time.Now()

		duration := stop.Sub(start)
		if duration > 4 * 1e9 {
			t.Fatal("Run time too long")
		}
	}
}
