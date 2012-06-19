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
	"reflect"
)

type CrossEmail struct {
	log *log.Logger
	queue *gobus.TailDelayQueue
	config *Config
	client *gobus.Client
	tmpl *template.Template
}

var helper = template.FuncMap{
	"last": func(x int, a interface{}) bool {
		return x == reflect.ValueOf(a).Len() - 1
	},
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

	t := template.Must(template.New("invitation").Funcs(helper).ParseFiles("./template/default/cross_invitation.html", "./template/default/cross_update.html", "./template/default/cross.ics"))

	return &CrossEmail{
		log: log,
		queue: queue,
		config: config,
		client: gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, provider),
		tmpl: t,
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

			e.sendMail(arg)
		}
	}
}

func (e *CrossEmail) GetBody(arg *ProviderArg, filename string) (string, string, error) {
	html := bytes.NewBuffer(nil)
	err := e.tmpl.ExecuteTemplate(html, filename, arg)
	if err != nil {
		return "", "", err
	}

	ics := bytes.NewBuffer(nil)
	err = e.tmpl.ExecuteTemplate(ics, "cross.ics", arg)
	if err != nil {
		return "", "", err
	}

	return html.String(), ics.String(), nil
}

func (e *CrossEmail) sendMail(arg *ProviderArg) {
	filename := "cross_invitation.html"
	if arg.Old_cross != nil {
		arg.Diff(e.log)
		filename = "cross_update.html"
	}

	html, ics, err := e.GetBody(arg, filename)
	if err != nil {
		e.log.Printf("template exec error:", err)
		return
	}
	htmls := strings.SplitN(html, "\n\n", 2)

	mailarg := gomail.Mail{
		To: []gomail.MailUser{gomail.MailUser{arg.To_identity.External_id, arg.To_identity.Name}},
		From: gomail.MailUser{"x@exfe.com", "x@exfe.com"},
		Subject: htmls[0],
		Html: htmls[1],
		FileParts: []gomail.FilePart{
			gomail.FilePart{fmt.Sprintf("x-%d.ics", arg.Cross.Id), []byte(ics)},
		},
	}

	e.client.Send("EmailSend", &mailarg, 5)
}
