package main

import (
	"email/service"
	"exfe/service"
	"fmt"
	"log"
	"net/smtp"
	"old_gobus"
	"time"
)

func main() {
	log.SetPrefix("exfe.email.notify")
	log.Print("Service start")

	c := exfe_service.InitConfig()

	server := gobus.CreateServer(c.Redis.Netaddr, c.Redis.Db, c.Redis.Password, "email")

	server.Register(
		&email_service.EmailSenderService{
			Server: fmt.Sprintf("%s:%d", c.Email.Host, c.Email.Port),
			Auth:   smtp.PlainAuth("", c.Email.User, c.Email.Password, c.Email.Host),
		})

	server.Serve(c.Email.Time_out * time.Second)
}
