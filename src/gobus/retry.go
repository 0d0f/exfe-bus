package gobus

import (
	"strings"
	"fmt"
)

type RetryJob struct {
	netaddr string
	db int
	password string
}

func (j *RetryJob) Resend(args []failedType) error {
	for _, arg := range args {
		sp := strings.Split(arg.Meta.Id, ":")
		queue := sp[len(sp) - 2]
		client := CreateClient(j.netaddr, j.db, j.password, queue)

		client.connRedis()
		defer client.closeRedis()

		err := client.send(&arg.Meta)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func DefaultRetryServer(netaddr string, db int, password string) *Server {
	job := RetryJob{netaddr, db, password}
	server := CreateServer(netaddr, db, password, "failed")
	server.queueName = "gobus:failed"
	server.Register(&job)
	return server
}
