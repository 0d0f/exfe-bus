package main

import (
	"exfe/service"
	"email/service"
	"os"
	"config"
	"flag"
	"fmt"
	"gobus"
	"log/syslog"
	"net/smtp"
	"time"
)

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
		&email_service.EmailSenderService{
			Server: fmt.Sprintf("%s:%d", c.Email.Host, c.Email.Port),
			Auth:   smtp.PlainAuth("", c.Email.User, c.Email.Password, c.Email.Host),
			Log: log,
		})
	defer func() {
		log.Info("Service stop")
		server.Close()
	}()

	server.Serve(time.Duration(c.Email.Time_out))
}
