package main

import (
	"flag"
	"config"
	"time"
	"gobus"
	"net/smtp"
	"log"
	"fmt"
	"./pkg/mail"
)

type MailSenderService struct {
	server string
	auth smtp.Auth
}

func (m *MailSenderService) Do(arg mail.Mail, reply *string) error {
	log.Printf("Try to send subject(%s) mail from (%s) to (%s)...", arg.Subject, arg.From.ToString(), arg.ToLine())
	err := smtp.SendMail(m.server, m.auth, arg.From.Mail, arg.ToMail(), arg.Body())

	if err != nil {
		log.Printf("Mail send failed: %s", err)
		return err
	}

	*reply = "ok"
	return nil
}

const (
	queue = "mail:sender"
)

func main() {
	log.SetPrefix("[Sendmail]")
	log.Printf("Service start, queue: %s", queue)

	var configFile string
	flag.StringVar(&configFile, "config", "mail.yaml", "Specify the configuration file")
	flag.Parse()

	config := config.LoadFile(configFile)

	service := gobus.CreateService(
		config.String("redis.netaddr"),
		config.Int("redis.db"),
		config.String("redis.password"),
		queue,
		&MailSenderService{
			server: fmt.Sprintf("%s:%d", config.String("mail.host"), config.Int("mail.port")),
			auth: smtp.PlainAuth("", config.String("mail.user"), config.String("mail.password"), config.String("mail.host")),
		})
	defer func() {
		log.Printf("Service stop, queue: %s", queue)
		service.Close()
		service.Clear()
	}()

	service.Serve(time.Duration(config.Int("service.time_out")))
}
