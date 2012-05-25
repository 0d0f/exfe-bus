package exfe_service

import (
	"github.com/simonz05/godis"
	"time"
	"exfe/model"
	"gomail"
	"fmt"
	"bytes"
	"text/template"
	"log/syslog"
	"gobus"
)

type CrossEmail struct {
	log *syslog.Writer
	queue *gobus.TailDelayQueue
	config *Config
	client *gobus.Client
}

func NewCrossEmail(config *Config) (ret *CrossEmail) {
	provider := "email"
	var err error
	ret.log, err = syslog.New(syslog.LOG_DEBUG, fmt.Sprintf("exfe.cross.%s", provider))
	if err != nil {
		panic(err)
	}

	arg := []OneIdentityUpdateArg{}
	redis := godis.New(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password)
	ret.queue, err = gobus.NewTailDelayQueue(getProviderQueueName(provider), config.Cross.Delay[provider], arg, redis)
	if err != nil {
		panic(err)
	}

	ret.config = config
	ret.client = gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, provider)
	return
}

func (e *CrossEmail) Serve() {
	for {
		t, err := e.queue.NextWakeup()
		if err != nil {
			e.log.Crit(fmt.Sprintf("next wakeup error: %s", err))
			break
		}
		time.Sleep(t)
		args, err := e.queue.Pop()
		if err != nil {
			e.log.Err(fmt.Sprintf("pop from delay queue failed: %s", err))
			continue
		}
		if args != nil {
			updates := args.([]OneIdentityUpdateArg)

			cross := &updates[len(updates)-1].Cross
			to := findInvitation(&updates[0].To_identity, &updates[0].Cross)
			posts := make([]*exfe_model.Post, 0, 0)

			var old_cross *exfe_model.Cross
			for _, update := range updates {
				if old_cross == nil && update.Old_cross != nil {
					old_cross = update.Old_cross
				}
				if update.Post != nil {
					posts = append(posts, update.Post)
				}
			}

			e.sendMail(to, cross, old_cross, posts)
		}
	}
}

type CrossTemplateData struct {
	To *exfe_model.Invitation
	Cross *exfe_model.Cross

	Old_cross *exfe_model.Cross
	Posts []*exfe_model.Post

	Site_url string
	App_url string
}

func findInvitation(to *exfe_model.Identity, cross *exfe_model.Cross) *exfe_model.Invitation {
	for i := range cross.Exfee.Invitations {
		if to.Id == cross.Exfee.Invitations[i].Identity.Id {
			return &cross.Exfee.Invitations[i]
		}
	}
	return nil
}

func (e *CrossEmail) sendMail(to *exfe_model.Invitation, cross, old_cross *exfe_model.Cross, posts []*exfe_model.Post) {
	data := CrossTemplateData{to, cross, old_cross, posts, e.config.Site_url, "appurl"}

	buf := bytes.NewBuffer(nil)
	tmpl := template.Must(template.ParseFiles("./template/default/cross_email.html"))
	err := tmpl.Execute(buf, data)
	if err != nil {
		e.log.Err(fmt.Sprintf("template exec error:", err))
		return
	}

	arg := gomail.Mail{
		To: []gomail.MailUser{gomail.MailUser{to.Identity.External_id, to.Identity.Name}},
		From: gomail.MailUser{"x@exfe.com", "x@exfe.com"},
		Subject: cross.Title,
		Html: buf.String(),
	}

	fmt.Println(buf.String())
	e.client.Send("EmailSend", &arg, 5)
}
