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
	service := CreateService("", 0, "", queue, &EmptyService{}, -1)
	defer func() { service.Close() }()
	go func() {
		service.Run(1e9)
	}()
	time.Sleep(0.5e9)
	if !service.IsRunning() {
		t.Fatal("Service doesn't run")
	}
	_ = service.Stop()
	time.Sleep(0.5e9)
	if service.IsRunning() {
		t.Fatal("Service is still running")
	}

	service.Clear()
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
	service := CreateService("", 0, "", queue, &PerJobService{}, -1)
	defer func() {
		service.Close()
		service.Clear()
	}()

	go func() {
		service.Run(1e9)
	}()

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
	service := CreateService("", 0, "", queue, &PerJobService{}, -1)
	defer func() {
		service.Close()
		service.Clear()
	}()

	client := CreateClient("", 0, "", queue, ValueGenerator)
	defer func() { client.Close() }()

	go func() {
		service.Run(1e9)
	}()

	retChan := make(chan int)
	start := time.Now()
	for i := 0; i < 10; i++ {
		go func() {
			ret, _ := client.Do(2)
			retChan <- (*ret.(*int))
		}()
	}

	for i := 0; i < 10; i++ {
		ret := <-retChan
		if ret != 4 {
			t.Error("Got %d, expect: 4", ret)
		}
	}

	stop := time.Now()

	duration := stop.Sub(start)
	if duration > 2*1e9 {
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
	service := CreateService("", 0, "", queue, &Per4JobsService{}, -1)
	defer func() {
		service.Close()
		service.Clear()
	}()

	client := CreateClient("", 0, "", queue, ValueGenerator)
	defer func() { client.Close() }()

	go func() {
		service.Run(1e9)
	}()

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
	service := CreateService("", 0, "", queue, &Per4JobsService{}, -1)
	defer func() {
		service.Close()
		service.Clear()
	}()

	client := CreateClient("", 0, "", queue, ValueGenerator)
	defer func() { client.Close() }()

	go func() {
		service.Run(1e9)
	}()

	retChan := make(chan int)
	{
		start := time.Now()
		for i := 0; i < 4; i++ {
			go func() {
				ret, _ := client.Do(2)
				retChan <- (*ret.(*int))
			}()
		}

		for i := 0; i < 4; i++ {
			<-retChan
		}

		stop := time.Now()

		duration := stop.Sub(start)
		if duration > 2*1e9 {
			t.Fatal("Run time too long")
		}
	}

	{
		start := time.Now()
		for i := 0; i < 10; i++ {
			go func() {
				ret, _ := client.Do(2)
				retChan <- (*ret.(*int))
			}()
		}

		for i := 0; i < 10; i++ {
			<-retChan
		}

		stop := time.Now()

		duration := stop.Sub(start)
		if duration > 4*1e9 {
			t.Fatal("Run time too long")
		}
	}
}

//////////////////////////////////////////

type SendService struct {
	out chan int
}

func (s *SendService) Do(jobs []interface{}) []interface{} {
	if len(jobs) > 1 {
		fmt.Println("SendService get jobs count > 1")
	}

	data := *jobs[0].(*int)
	s.out <- (data * data)
	return nil
}

func (s *SendService) MaxJobsCount() int {
	return 1
}

func (s *SendService) JobGenerator() interface{} {
	var ret int
	return &ret
}

func TestSendWork(t *testing.T) {
	fmt.Println("Test send work")

	queue := "gobus:queue:sendwork"
	out := make(chan int)
	service := CreateService("", 0, "", queue, &SendService{out: out}, -1)
	defer func() {
		service.Close()
		service.Clear()
	}()

	client := CreateClient("", 0, "", queue, ValueGenerator)
	defer func() { client.Close() }()

	go func() {
		service.Run(1e9)
	}()

	{
		d := 2
		err := client.Send(d)
		if err != nil {
			t.Fatal("Send error:", err)
		}
		i := <-out
		if i != 4 {
			t.Error("expect: 4, got:", i)
		}
	}
}

func TestSendTime(t *testing.T) {
	fmt.Println("Test send time")

	queue := "gobus:queue:sendtime"
	out := make(chan int)
	service := CreateService("", 0, "", queue, &SendService{out: out}, -1)
	defer func() {
		service.Close()
		service.Clear()
	}()

	client := CreateClient("", 0, "", queue, ValueGenerator)
	defer func() { client.Close() }()

	go func() {
		service.Run(1e9)
	}()

	{
		start := time.Now()
		for i := 0; i < 10; i++ {
			go func() {
				client.Send(2)
			}()
		}

		stop := time.Now()

		duration := stop.Sub(start)
		if duration > 1e9 {
			t.Fatal("Run time too long")
		}

		for i := 0; i < 10; i++ {
			<-out
		}
	}
}

//////////////////////////////////////////////

type LimitService struct {
}

func (s *LimitService) Do(jobs []interface{}) []interface{} {
	time.Sleep(0.5 * 1e9)
	i := 2
	return []interface{}{&i}
}

func (s *LimitService) MaxJobsCount() int {
	return 1
}

func (s *LimitService) JobGenerator() interface{} {
	var i int
	return &i
}

func TestLimitWork(t *testing.T) {
	fmt.Println("Test limit work")

	queue := "gobus:queue:limitwork"
	service := CreateService("", 0, "", queue, &LimitService{}, 4)
	defer func() {
		service.Close()
		service.Clear()
	}()

	client := CreateClient("", 0, "", queue, ValueGenerator)
	defer func() { client.Close() }()

	go func() {
		service.Run(1e9)
	}()

	c := make(chan int)
	for i := 0; i < 10; i++ {
		go func() {
			client.Do(1)
			c <- 1
		}()
	}

	{
		start := time.Now()
		for i := 0; i < 10; i++ {
			<-c
		}
		stop := time.Now()
		duration := stop.Sub(start)
		if duration > 5*1e9 {
			t.Fatal("Run time too long")
		}
		if duration < 4*1e9 {
			t.Fatal("Run time too fast")
		}
	}
}
