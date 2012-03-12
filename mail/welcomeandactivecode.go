package main

import (
	"gobus"
	"config"
	"fmt"
	"bytes"
	"gosque"
	"log"
	"text/template"
)

type WelcomeAndActiveData struct {
	Identityid int64
	External_identity string
	Name string
	Avatar_file_name string
	Activecode string
	Token string

	SiteUrl string
	config *config.Configure
	sendmail *gobus.Client
}

type MailUser struct {
	Mail string
	Name string
}

type SendMail struct {
	To []MailUser
	From MailUser
	Subject string
	Text string
	Html string
}

func LoadTemplate(name string) *template.Template {
	return template.Must(template.ParseFiles(fmt.Sprintf("./template/default/%s", name)))
}

func MailResponseGenerator() interface{} {
	var ret string
	return &ret
}

func WelcomeAndActiveDataGenerator() interface{} {
	return &WelcomeAndActiveData{}
}

func (d *WelcomeAndActiveData) Do() {
	log.Printf("Try to send welcome and active email to %s <%s>...", d.Name, d.External_identity)

	tmpl := LoadTemplate("welcome_and_active.tmpl")
	buf := bytes.NewBuffer(nil)
	tmpl.Execute(buf, d)
	message := buf.String()

	d.sendmail.Do(&SendMail{
		To: []MailUser{
			{
				Mail: d.External_identity,
				Name: d.Name,
			},
		},
		From: MailUser{
			Mail: "x@exfe.com",
			Name: "exfe",
		},
		Subject: "Welcome to EXFE!",
		Html: message,
	})
}

func main() {
	log.SetPrefix("[Welcomeandactvie]")
	log.Printf("Service start")
	config := config.LoadFile("mail.yaml")

	client := gosque.CreateQueue(
		config.String("redis.netaddr"),
		config.Int("redis.db"),
		config.String("redis.password"),
		"resque:queue:email")

	sendmail := gobus.CreateClient(
		config.String("redis.netaddr"),
		config.Int("redis.db"),
		config.String("redis.password"),
		"gobus:queue:mail:sender",
		MailResponseGenerator)

	recv := client.IncomingJob("welcomeandactivecode_job", WelcomeAndActiveDataGenerator, 5e9)
	for {
		select {
		case job := <-recv:
			data := job.(*WelcomeAndActiveData)
			data.config = config
			data.sendmail = sendmail
			data.SiteUrl = config.String("site_url")
			go func() {
				data.Do()
			}()
		}
	}
	log.Printf("Service stop")
}
