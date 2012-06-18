package exfe_service

import (
	"strings"
	"github.com/googollee/godis"
	"time"
	"exfe/model"
	"gomail"
	"fmt"
	"bytes"
	"text/template"
	"log"
	"gobus"
	"os"
)

type CrossEmail struct {
	log *log.Logger
	queue *gobus.TailDelayQueue
	config *Config
	client *gobus.Client
}

func NewCrossEmail(config *Config) *CrossEmail {
	provider := "email"
	log := log.New(os.Stderr, fmt.Sprintf("exfe.cross.%s", provider), log.LstdFlags)

	arg := []OneIdentityUpdateArg{}
	redis := godis.New(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password)
	queue, err := gobus.NewTailDelayQueue(getProviderQueueName(provider), config.Cross.Delay[provider], arg, redis)
	if err != nil {
		panic(err)
	}

	return &CrossEmail{
		log: log,
		queue: queue,
		config: config,
		client: gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, provider),
	}
}

func (e *CrossEmail) Serve() {
	for {
		t, err := e.queue.NextWakeup()
		if err != nil {
			e.log.Printf("next wakeup error: %s", err)
			break
		}
		time.Sleep(t)
		args, err := e.queue.Pop()
		if err != nil {
			e.log.Printf("pop from delay queue failed: %s", err)
			continue
		}
		if args != nil {
			updates := args.([]OneIdentityUpdateArg)

			by_identities := make([]*exfe_model.Identity, 0, 0)
			posts := make([]*exfe_model.Post, 0, 0)
			var old_cross *exfe_model.Cross
			for _, update := range updates {
				by_identities = append(by_identities, &update.By_identity)
				if old_cross == nil && update.Old_cross != nil {
					old_cross = update.Old_cross
				}
				if update.Post != nil {
					posts = append(posts, update.Post)
				}
			}

			arg := &ProviderArg{
				Cross: &updates[len(updates)-1].Cross,
				Old_cross: old_cross,
				To_identity: &updates[0].To_identity,
				By_identities: by_identities,
				Posts: posts,
				Config: e.config,
			}

			if arg.Old_cross == nil {
				e.sendMail(arg, "cross_invitation")
			} else {
				e.sendMail(arg, "cross_update")
			}
		}
	}
}

func (e *CrossEmail) sendMail(arg *ProviderArg, filename string) {
	html := bytes.NewBuffer(nil)
	tmpl := template.Must(template.ParseFiles(fmt.Sprintf("./template/default/%s.html", filename)))
	err := tmpl.Execute(html, arg)
	if err != nil {
		e.log.Printf("template exec error:", err)
		return
	}
	htmls := strings.SplitN(html.String(), "\n\n", 2)

	ics := bytes.NewBuffer(nil)
	tmpl = template.Must(template.ParseFiles("./template/default/cross.ics"))
	err = tmpl.Execute(ics, arg)
	if err != nil {
		e.log.Printf("template exec error:", err)
		return
	}

	mailarg := gomail.Mail{
		To: []gomail.MailUser{gomail.MailUser{arg.To_identity.External_id, arg.To_identity.Name}},
		From: gomail.MailUser{"x@exfe.com", "x@exfe.com"},
		Subject: htmls[0],
		Html: htmls[1],
		FileParts: []gomail.FilePart{
			gomail.FilePart{fmt.Sprintf("x-%d.ics", arg.Cross.Id), ics.Bytes()},
		},
	}

	e.client.Send("EmailSend", &mailarg, 5)
}
