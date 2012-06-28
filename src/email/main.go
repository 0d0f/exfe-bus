package main

import (
	"exfe/service"
	"email/service"
	"fmt"
	"gobus"
	"log"
	"net/smtp"
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
	defer func() {
		log.Print("Service stop")
		server.Close()
	}()

	server.Serve(c.Email.Time_out * time.Second)
}
