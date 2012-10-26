package exfe_service

import (
	"bytes"
	"email/service"
	"encoding/base64"
	"exfe/model"
	"fmt"
	"github.com/googollee/godis"
	"log"
	"net/mail"
	"old_gobus"
	"os"
	"reflect"
	"strings"
	"text/template"
	"time"
)

type CrossEmail struct {
	log    *log.Logger
	queue  *gobus.TailDelayQueue
	config *Config
	client *gobus.Client
	tmpl   *template.Template
}

var helper = map[string]interface{}{
	"last": func(x int, a interface{}) bool {
		return x == reflect.ValueOf(a).Len()-1
	},
	"limit": func(s string, max int) string {
		if max > len(s) {
			max = len(s)
		}
		return s[0:max]
	},
	"replace": func(s, old, new string) string {
		return strings.Replace(s, old, new, -1)
	},
	"treplace": func(old, new, s string) string {
		return strings.Replace(s, old, new, -1)
	},
	"base64": func(s string) string {
		return base64.StdEncoding.EncodeToString([]byte(s))
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
		log:    log,
		queue:  queue,
		config: config,
		client: gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, provider),
		tmpl:   t,
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
			firstUpdate := true
			for _, update := range updates {
				by_identities = append(by_identities, &update.By_identity)
				if update.Post != nil {
					posts = append(posts, update.Post)
					continue
				}
				if firstUpdate {
					old_cross = update.Old_cross
					firstUpdate = false
				}
			}

			e.log.Printf("Got %d updates to %s", len(updates), updates[0].To_identity.ExternalId())
			arg := &ProviderArg{
				Cross:         &updates[len(updates)-1].Cross,
				Old_cross:     old_cross,
				To_identity:   &updates[0].To_identity,
				By_identities: by_identities,
				Posts:         posts,
				Config:        e.config,
			}

			e.sendMail(arg, firstUpdate == false)
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
	if arg.Cross.Time.Begin_at.Date != "" {
		err = e.tmpl.ExecuteTemplate(ics, "cross.ics", arg)
		if err != nil {
			return "", "", err
		}
	}

	return html.String(), ics.String(), nil
}

func (e *CrossEmail) sendMail(arg *ProviderArg, hasUpdate bool) {
	_, _, newlyInvited, _ := arg.Diff(e.log)
	_, ok := newlyInvited[arg.To_identity.DiffId()]

	if hasUpdate && (arg.Old_cross == nil || ok) {
		filename := "cross_invitation.html"
		html, ics, err := e.GetBody(arg, filename)
		if err != nil {
			e.log.Printf("template exec error: %s", err)
			return
		}
		htmls := strings.SplitN(html, "//////////////////////////////////\n\n", 3)

		mail_addr := fmt.Sprintf("x+%d@%s", arg.Cross.Id, e.config.EmailDomain)
		mailarg := &email_service.MailArg{
			To:         []*mail.Address{&mail.Address{arg.To_identity.Name, arg.To_identity.External_id}},
			From:       &mail.Address{e.config.EmailName, mail_addr},
			Subject:    strings.Trim(htmls[0], " \n\r\t"),
			Text:       strings.Trim(htmls[1], " \n\r\t"),
			Html:       strings.Trim(htmls[2], " \n\r\t"),
			References: []string{fmt.Sprintf("<%s>", mail_addr)},
		}
		if ics != "" {
			mailarg.FileParts = []email_service.FilePart{
				email_service.FilePart{fmt.Sprintf("%s.ics", arg.Cross.Title), []byte(ics)},
			}
		}

		e.client.Send("EmailSend", &mailarg, 5)
	}
	if ok {
		return
	}

	filename := "cross_update.html"
	if !arg.IsCrossChanged() && len(arg.Posts) == 0 {
		return
	}
	if len(arg.Posts) > 0 && !arg.IsCrossChanged() {
		selfPost := true
		for _, p := range arg.Posts {
			if p.By_identity.DiffId() != arg.To_identity.DiffId() {
				selfPost = false
				break
			}
		}
		if selfPost {
			return
		}
	}

	html, ics, err := e.GetBody(arg, filename)
	if err != nil {
		e.log.Printf("template exec error: %s", err)
		return
	}
	htmls := strings.SplitN(html, "//////////////////////////////////\n\n", 3)

	mail_addr := fmt.Sprintf("x+%d@%s", arg.Cross.Id, e.config.EmailDomain)
	mailarg := &email_service.MailArg{
		To:         []*mail.Address{&mail.Address{arg.To_identity.Name, arg.To_identity.External_id}},
		From:       &mail.Address{e.config.EmailName, mail_addr},
		Subject:    strings.Trim(htmls[0], " \n\r\t"),
		Text:       strings.Trim(htmls[1], " \n\r\t"),
		Html:       strings.Trim(htmls[2], " \n\r\t"),
		References: []string{fmt.Sprintf("<%s>", mail_addr)},
	}
	if ics != "" {
		mailarg.FileParts = []email_service.FilePart{
			email_service.FilePart{fmt.Sprintf("%s.ics", arg.Cross.Title), []byte(ics)},
		}
	}

	e.client.Send("EmailSend", &mailarg, 5)
}
