package email

import (
	"fmt"
	"gobot"
	"github.com/googollee/goimap"
	"github.com/sloonz/go-iconv"
	"strings"
	"net/mail"
	"regexp"
	"exfe/service"
	"gobus"
)

var replyLine = [...]string{
	"^--$",
	"^--&nbsp;$",
	"-----Original Message-----",
	"________________________________",
	"Sent from my iPhone",
	"Sent from my BlackBerry",
	`^From:.*[mailto:.*]`,
	`^On (.*) wrote:`,
	"发自我的 iPhone",
	"^[ \t\n\r]*$",
}

var removeLine = [...]string{
	`(?iU)<script\b.*>.*</script>`,
	`(?iU)<style\b.*>.*</style>`,
	`(?U)<.*>`,
}

type EmailBot struct {
	config      *exfe_service.Config
	bus         *gobus.Client
	crossId     *regexp.Regexp
	retReplacer *strings.Replacer
	remover     []*regexp.Regexp
	replyLine   []*regexp.Regexp
}

func NewEmailBot(config *exfe_service.Config) *EmailBot {
	reply := make([]*regexp.Regexp, len(replyLine), len(replyLine))
	for i, l := range replyLine {
		reply[i] = regexp.MustCompile(l)
	}
	remover := make([]*regexp.Regexp, len(removeLine), len(removeLine))
	for i, l := range removeLine {
		remover[i] = regexp.MustCompile(l)
	}
	return &EmailBot{
		config:      config,
		bus:         gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, "email"),
		crossId:     regexp.MustCompile(`^.*?\+([0-9]*)@.*$`),
		retReplacer: strings.NewReplacer("\r\n", "\n", "\r", "\n"),
		remover:     remover,
		replyLine:   reply,
	}
}

func (b *EmailBot) GenerateContext(id string) bot.Context {
	return NewEmailContext(id, b)
}

func (b *EmailBot) GetIDFromInput(input interface{}) (id string, content interface{}, e error) {
	msg, ok := input.(*mail.Message)
	if !ok {
		e = fmt.Errorf("input is not a mail.Message")
		return
	}
	from, err := imap.ParseAddress(msg.Header.Get("From"))
	if err != nil {
		e = fmt.Errorf("can't parse From: %s", err)
		return
	}
	to, err := imap.ParseAddress(msg.Header.Get("To"))
	if err != nil {
		e = fmt.Errorf("can't parse To: %s", err)
		return
	}
	date, err := msg.Header.Date()
	if err != nil {
		e = fmt.Errorf("Get message(%v) time error: id, err")
		return
	}
	text, mediatype, charset, err := imap.GetBody(msg, "text/plain")
	if err != nil {
		e = fmt.Errorf("Get message(%v) body failed: %s", id, err)
		return
	}
	text, err = iconv.Conv(text, "UTF-8", charset)
	if err != nil {
		e = fmt.Errorf("Convert message(%v) from %s to utf8 error: %s", id, charset, err)
		return
	}
	if mediatype != "text/plain" {
		text = b.stripGmail(text)
		text = b.stripHtml(text)
	}
	text = b.stripReply(text)
	crossId, err := b.getCrossId(to)
	if err != nil {
		e = fmt.Errorf("Find cross id from message(%v) To field(%s) error: %s", id, to, err)
		return
	}
	id = from[0].Address
	content = &Email{
		From:      from[0],
		To:        to,
		Subject:   msg.Header.Get("Subject"),
		CrossId:   crossId,
		Date:      date,
		MessageId: msg.Header.Get("Message-Id"),
		Text:      text,
	}
	return
}

func (b *EmailBot) getCrossId(addrs []*mail.Address) (string, error) {
	for _, addr := range addrs {
		ids := b.crossId.FindStringSubmatch(addr.Address)
		if len(ids) > 1 {
			return ids[1], nil
		}
	}
	return "", fmt.Errorf("No valid mail address")
}

func (b *EmailBot) isReplys(line string) bool {
	for _, r := range b.replyLine {
		if r.MatchString(line) {
			return true
		}
	}
	return false
}

func (b *EmailBot) stripReply(content string) string {
	content = b.retReplacer.Replace(content)
	lines := strings.Split(content, "\n")
	ret := make([]string, len(lines), len(lines))

	for i, line := range lines {
		if b.isReplys(line) {
			ret = ret[:i]
			break
		}
		ret[i] = line
	}
	return strings.Join(ret, "\n")
}

func (b *EmailBot) stripGmail(text string) string {
	pos := strings.Index(text, `<div class="gmail_quote"`)
	if pos >= 0 {
		return strings.Trim(text[:pos], " \t\n\r")
	}
	return text
}

func (b *EmailBot) stripHtml(text string) string {
	for _, r := range b.remover {
		text = r.ReplaceAllString(text, "")
	}
	return strings.Trim(text, " \t\n\r")
}
