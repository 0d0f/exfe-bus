package gosque

import (
	"testing"
	"fmt"
)

type Job struct {
	Name string
}

func GenerateJobs(gosque *Client, max int) {
	for i:=0; i<max; i++ {
		job := Job{
			fmt.Sprintf("Job-%d", i),
		}
		err := gosque.PutJob(job)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func GenerateJobTemplate() interface{} {
	return &Job{"123"}
}

func TestGetJob(t *testing.T) {
	const queue = "resque:test:job"
	const maxJobs = 10
	gosque := CreateQueue("", 0, "", queue)
	defer func() { gosque.Close() }()
	GenerateJobs(gosque, maxJobs)

	jobRecv := gosque.IncomingJob(GenerateJobTemplate, 5e9)
	for i:=0; i<10; i++ {
		job := (<-jobRecv).(*Job)
		expectName := fmt.Sprintf("Job-%d", i)
		if job.Name != expectName {
			t.Errorf("Failed Job.Name, expect: %s, got: %s", expectName, job.Name)
		}
	}
}
