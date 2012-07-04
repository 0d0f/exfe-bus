package main

import (
	"strings"
	"os"
	"exfe/service"
	"time"
	"regexp"
	"github.com/googollee/goimap"
	"github.com/sloonz/go-iconv"
	"log"
	"net/mail"
	"net/url"
	"net/http"
	"io/ioutil"
	"fmt"
	"gobus"
	"gomail"
)

var mailConfig *exfe_service.Config
var mailPattern *regexp.Regexp
var emailLog *log.Logger
var mailBus *gobus.Client

func InitEmail(c *exfe_service.Config) {
	mailConfig = c
	mailPattern = regexp.MustCompile(`^.*?\+([0-9]*)@.*$`)
	emailLog = log.New(os.Stderr, "bot.mail", log.LstdFlags)
	mailBus = gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, "email")
}

func getCrossId(addrs []*mail.Address) (string, error) {
	for _, addr := range addrs {
		ids := mailPattern.FindStringSubmatch(addr.Address)
		if len(ids) > 1 {
			return ids[1], nil
		}
	}
	return "", fmt.Errorf("No valid mail address")
}

func processEmail(quit chan int) {
	conn, _ := imap.NewClient(mailConfig.Bot.Imap_host)
	defer func() {
		conn.Logout()
		conn.Close()
	}()

	_ = conn.Login(mailConfig.Bot.Imap_user, mailConfig.Bot.Imap_password)
	conn.Select(imap.Inbox)
	for {
		ids, _ := conn.Search("unseen")
		for _, id := range ids {
			fmt.Println("Process message", id)
			conn.StoreFlag(id, imap.Seen)

			msg, err := conn.GetMessage(id)
			if err != nil {
				emailLog.Printf("Get message(%v) error: %s", id, err)
				return
			}
			tos, err := imap.ParseAddress(msg.Header.Get("To"))
			if err != nil {
				emailLog.Printf("Parse message(%v) To field(%s) error: %s", id, msg.Header.Get("To"), err)
				continue
			}
			froms, err := imap.ParseAddress(msg.Header.Get("From"))
			if err != nil {
				emailLog.Printf("Parse message(%v) From field(%s) error: %s", id, msg.Header.Get("From"), err)
				continue
			}
			content, mediatype, charset, err := imap.GetBody(msg, "text/plain")
			if err != nil {
				emailLog.Printf("Get message(%v) body failed: %s", id, err)
				continue
			}
			content, err = iconv.Conv(content, "UTF-8", charset)
			if err != nil {
				emailLog.Printf("Convert message(%v) from %s to utf8 error: %s", id, charset, err)
				continue
			}
			date, err := msg.Header.Date()
			if err != nil {
				emailLog.Printf("Get message(%v) time error: id, err")
				continue
			}
			if mediatype != "text/plain" {
				content = stripGmail(content)
				content = stripHtml(content)
			}
			content = stripReply(content)

			crossId, err := getCrossId(tos)
			if err != nil {
				emailLog.Printf("Find cross id from message(%v) To field(%s) error: %s", id, tos, err)
				sendErrorMail(froms[0], msg.Header.Get("Subject"), content)
				continue
			}

			params := make(url.Values)
			params.Add("cross_id", crossId)
			params.Add("content", content)
			params.Add("external_id", froms[0].Address)
			params.Add("provider", "email")
			params.Add("time", date.Format("2006-01-02 15:04:05 -0700"))
			resp, err := http.PostForm(fmt.Sprintf("%s/v2/gobus/PostConversation", mailConfig.Site_api), params)
			if err != nil {
				emailLog.Printf("Send post to server error: %s", err)
				continue
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				emailLog.Printf("Get response body error: %s", err)
				continue
			}
			if resp.StatusCode == 500 {
				emailLog.Printf("Server inner error: %s", string(body))
				continue
			}
			if resp.StatusCode == 400 {
				emailLog.Printf("User status error: %s", string(body))
				continue
			}
		}
		time.Sleep(mailConfig.Bot.Imap_time_out * time.Second)
	}
	quit <- 1
}

func isReplys(line string) bool {
	if line == "--" || line == "--&nbsp;" {
		return true
	}
	if strings.Index(line, "-----Original Message-----") >= 0 {
		return true
	}
	if strings.Index(line, "________________________________") >= 0 {
		return true
	}
	if strings.Index(line, "Sent from my iPhone") >= 0 {
		return true
	}
	if strings.Index(line, "Sent from my BlackBerry") >= 0 {
		return true
	}
	if matched, err := regexp.MatchString(`^From:.*[mailto:.*]`, line); matched && err == nil {
		return true
	}
	if matched, err := regexp.MatchString(`^On (.*) wrote:`, line); matched && err == nil {
		return true
	}
	if strings.Index(line, "发自我的 iPhone") >= 0 {
		return true
	}
	if strings.Trim(line, " \t\n\r") == "" {
		return true
	}
	return false
}

func stripReply(content string) string {
	content = strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(content)
	lines := strings.Split(content, "\n")
	ret := make([]string, len(lines), len(lines))

	for i, line := range lines {
		if isReplys(line) {
			ret = ret[:i]
			break
		}
		ret[i] = line
	}
	return strings.Join(ret, "\n")
}

func stripGmail(text string) string {
	pos := strings.Index(text, `<div class="gmail_quote"`)
	if pos >= 0 {
		return strings.Trim(text[:pos], " \t\n\r")
	}
	return text
}

func stripHtml(text string) string {
	reg := regexp.MustCompile(`(?iU)<script\b.*>.*</script>`)
	text = reg.ReplaceAllString(text, "")
	reg = regexp.MustCompile(`(?iU)<style\b.*>.*</style>`)
	text = reg.ReplaceAllString(text, "")
	reg = regexp.MustCompile(`(?U)<.*>`)
	text = reg.ReplaceAllString(text, "")
	return strings.Trim(text, " \t\n\r")
}

func sendErrorMail(to *mail.Address, subject, content string) {
	body := fmt.Sprintf("Sorry for the inconvenience, but email you just sent to EXFE was not sent from an attendee identity to the X (cross). Please try again from the correct email address.\n -- $s",
		content)
	mailarg := gomail.Mail{
		To:      []gomail.MailUser{gomail.MailUser{to.Address, to.Name}},
		From:    gomail.MailUser{"x@exfe.com", "x@exfe.com"},
		Subject: fmt.Sprintf("Re: %s", subject),
		Text:    body,
	}
	mailBus.Send("EmailSend", &mailarg, 5)
}
