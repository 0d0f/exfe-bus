package main

import (
	"gomail"
	"exfe/service"
	"os"
	"config"
	"flag"
	"fmt"
	"gobus"
	"log/syslog"
	"net/smtp"
	"time"
)

type MailSenderService struct {
	server string
	auth   smtp.Auth
	log    *syslog.Writer
}

func (m *MailSenderService) SendMail(arg *gomail.Mail, reply *string) error {
	m.log.Info(fmt.Sprintf("Send mail: subject(%s) from (%s) to (%s)", arg.Subject, arg.From.ToString(), arg.ToLine()))
	to := arg.ToMail()
	body := arg.Body()
	err := smtp.SendMail(m.server, m.auth, arg.From.Mail, to, body)

	if err != nil {
		m.log.Err(fmt.Sprintf("Mail send failed: %s", err))
		return err
	}
	m.log.Info("send ok")

	*reply = "ok"
	return nil
}

func main() {
	log, err := syslog.New(syslog.LOG_INFO, "exfe.email.notify")
	if err != nil {
		panic(err)
	}
	log.Info("Service start")

	var pidfile string
	var configFile string

	flag.StringVar(&pidfile, "pid", "", "Specify the pid file")
	flag.StringVar(&configFile, "config", "exfe.json", "Specify the configuration file")
	flag.Parse()

	var c exfe_service.Config
	config.LoadFile(configFile, &c)

	if pidfile != "" {
		pid, err := os.Create(pidfile)
		if err != nil {
			log.Crit(fmt.Sprintf("Can't create pid(%s): %s", pidfile, err))
			return
		}
		pid.WriteString(fmt.Sprintf("%d", os.Getpid()))
	}

	server := gobus.CreateServer(c.Redis.Netaddr, c.Redis.Db, c.Redis.Password, "email")

	server.Register(
		&MailSenderService{
			server: fmt.Sprintf("%s:%d", c.Email.Host, c.Email.Port),
			auth:   smtp.PlainAuth("", c.Email.User, c.Email.Password, c.Email.Host),
			log: log,
		})
	defer func() {
		log.Info("Service stop")
		server.Close()
	}()

	server.Serve(time.Duration(c.Email.Time_out))
}
