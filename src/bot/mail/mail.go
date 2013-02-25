package mail

import (
	"broker"
	"bytes"
	"code.google.com/p/go-imap/go1/imap"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"formatter"
	"github.com/googollee/go-encoding-ex"
	"github.com/googollee/go-logger"
	"gobus"
	"io/ioutil"
	"launchpad.net/tomb"
	"mime/multipart"
	"model"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

func Daemon(config model.Config) (*tomb.Tomb, error) {
	var tomb tomb.Tomb

	log := config.Log.SubPrefix("mail")
	timeout := time.Duration(config.Bot.Email.TimeoutInSecond) * time.Second
	platform, err := broker.NewPlatform(&config)
	if err != nil {
		return nil, err
	}
	templ, err := formatter.NewLocalTemplate(config.TemplatePath, config.DefaultLang)
	if err != nil {
		return nil, err
	}
	table, err := gobus.NewTable(config.Dispatcher)
	if err != nil {
		return nil, err
	}
	sender, err := broker.NewSender(&config, gobus.NewDispatcher(table))
	if err != nil {
		return nil, err
	}

	go func() {
		defer tomb.Done()

		for {
			select {
			case <-tomb.Dying():
				log.Debug("quit")
			case <-time.After(timeout):
				ProcessMail(sender, templ, config, platform, log)
			}
		}
	}()
	return &tomb, nil
}

func ProcessMail(sender *broker.Sender, templ *formatter.LocalTemplate, config model.Config, platform *broker.Platform, log *logger.SubLogger) {
	conn, err := imap.DialTLS(config.Bot.Email.IMAPHost, new(tls.Config))
	if err != nil {
		log.Err("can't connect to %s: %s", config.Bot.Email.IMAPHost, err)
		return
	}
	defer conn.Logout(10 * time.Second)

	conn.Data = nil
	if conn.Caps["STARTTLS"] {
		conn.StartTLS(nil)
	}

	_, err = conn.Select("INBOX", false)
	if err != nil {
		log.Err("can't select INBOX: %s", err)
		return
	}

	cmd, err := imap.Wait(conn.Search("UNSEEN"))
	if err != nil {
		log.Err("can't seach UNSEEN: %s", err)
	}
	var ids []uint32
	for _, resp := range cmd.Data {
		ids = append(ids, resp.SearchResults()...)
	}

	var errorIds []uint32
	var okIds []uint32
	for _, id := range ids {
		buf := bytes.NewBuffer(nil)
		set := new(imap.SeqSet)
		set.AddNum(id)
		cmd, err := imap.Wait(conn.Fetch(set, "RFC822"))
		if err != nil {
			log.Err("can't fetch mail %d: %s", id, err)
			errorIds = append(errorIds, id)
			continue
		}
		buf.Write(imap.AsBytes(cmd.Data[0].MessageInfo().Attrs["RFC822"]))

		msg, err := mail.ReadMessage(buf)
		if err != nil {
			log.Err("can't parse mail %d: %s", id, err)
			errorIds = append(errorIds, id)
			continue
		}

		toList, err := msg.Header.AddressList("To")
		if err != nil {
			log.Err("can't get address To of mail %d: %s", id, err)
			errorIds = append(errorIds, id)
		}
		ccList, err := msg.Header.AddressList("Cc")
		if err != nil {
			log.Err("can't get address To of mail %d: %s", id, err)
			errorIds = append(errorIds, id)
		}
		fromList, err := msg.Header.AddressList("From")
		if err != nil {
			log.Err("can't get address To of mail %d: %s", id, err)
			errorIds = append(errorIds, id)
		}
		addrList := append(toList, ccList...)
		addrList = append(addrList, fromList...)

		title := parseTitle(msg.Header.Get("Subject"))
		froms, err := msg.Header.AddressList("From")
		if err != nil || len(froms) == 0 {
			log.Err("invalid from address(%s): %s", msg.Header.Get("From"), err)
			errorIds = append(errorIds, id)
			continue
		}
		from := froms[0]
		content, err := getContent(msg)
		if err != nil {
			errorIds = append(errorIds, id)
			continue
		}
		if ok, args := findAddress(fmt.Sprintf("x+([0-9a-zA-Z]*)@%s", config.Email.Domain), addrList); ok {
			post := content
			addr := args[0]
			to := "cross"
			toId := addr
			typeId := map[uint8]string{
				'c': "cross",
				'e': "exfee",
			}
			if t := addr[0]; !('0' <= t && t <= '9') {
				if len(toId) < 2 {
					log.Err("invalid address %s", t)
					errorIds = append(errorIds, id)
					continue
				}
				var ok bool
				to, ok = typeId[t]
				if !ok {
					log.Err("invalid address %s", t)
					errorIds = append(errorIds, id)
					continue
				}
				toId = toId[1:]
			}

			err := platform.BotPostConversation(post, to, toId)
			if err != nil {
				log.Err("platform BotPostConversation call failed: %s", err)
				errorIds = append(errorIds, id)
				continue
			}
			okIds = append(okIds, id)
		} else if ok, _ := findAddress(fmt.Sprintf("x@%s", config.Email.Domain), addrList); ok {
			cross := model.Cross{
				Title:       title,
				Description: content,
				By: model.Identity{
					Provider:   "email",
					Name:       from.Name,
					ExternalID: from.Address,
				},
			}
			invite := make([]model.Invitation, len(addrList))
			for i, addr := range addrList {
				invite[i] = model.Invitation{
					Host: addr.Address == from.Address,
					Via:  "email",
					Identity: model.Identity{
						Name:       addr.Name,
						ExternalID: addr.Address,
						Provider:   "email",
					},
				}
			}
			cross.Exfee.Invitations = invite
			err = platform.BotCreateCross(cross)
			if err != nil {
				log.Err("platform BotCreateCross call failed: %s", err)
				errorIds = append(errorIds, id)
				continue
			}
			okIds = append(okIds, id)
		} else {
			errorIds = append(errorIds, id)
			buf := bytes.NewBuffer(nil)
			type Email struct {
				From      *mail.Address
				Subject   string
				CrossID   string
				Date      time.Time
				MessageID string
				Text      string

				Config *model.Config
			}
			email := Email{
				From:      from,
				Subject:   title,
				MessageID: msg.Header.Get("Message-ID"),
				Text:      content,
			}
			err := templ.Execute(buf, "en_US", "conversation_reply.email", email)
			if err != nil {
				log.Crit("template(conversation_reply.email) failed: %s", err)
				errorIds = append(errorIds, id)
				continue
			}

			info := &model.InfoData{
				CrossID: 0,
				Type:    model.TypeCrossInvitation,
			}
			to := model.Recipient{
				Provider:         "email",
				ExternalID:       from.Address,
				ExternalUsername: from.Name,
			}
			_, err = sender.Send(to, buf.String(), "", info)
			if err != nil {
				log.Crit("send error: %s", err)
				errorIds = append(errorIds, id)
				continue
			}
		}
	}
	{
		set := new(imap.SeqSet)
		set.AddNum(errorIds...)
		_, err := imap.Wait(conn.Copy(set, "error"))
		if err != nil {
			log.Err("can't copy to error: %s", err)
		}
	}
	{
		set := new(imap.SeqSet)
		set.AddNum(okIds...)
		_, err := imap.Wait(conn.Copy(set, "error"))
		if err != nil {
			log.Err("can't copy to error: %s", err)
		}
	}
	{
		set := new(imap.SeqSet)
		set.AddNum(ids...)
		_, err := imap.Wait(conn.Store(set, "FLAGS", "\\Deleted"))
		if err != nil {
			log.Err("can't remove from inbox: %s", err)
		}
	}
}

func findAddress(pattern string, list []*mail.Address) (bool, []string) {
	r := regexp.MustCompile(pattern)
	for _, addr := range list {
		match := r.FindAllStringSubmatch(addr.Address, -1)
		if len(match) == 0 {
			continue
		}
		return true, match[0][1:]
	}
	return false, nil
}

func parseContentType(contentType string) (string, map[string]string) {
	parts := strings.Split(contentType, ";")
	if len(parts) == 0 {
		return "", nil
	}
	mime := ""
	if strings.Index(parts[0], "=") == -1 {
		mime = parts[0]
	}
	if len(parts) == 1 {
		return mime, nil
	}
	pairs := make(map[string]string)
	for _, part := range parts[1:] {
		part = strings.Trim(part, " \n\t")
		p := strings.Split(part, "=")
		pairs[p[0]] = p[1]
	}
	return mime, pairs
}

func getContent(msg *mail.Message) (string, error) {
	mime, pairs := parseContentType(msg.Header.Get("Content-Type"))
	if mime == "multipart/alternative" {
		parts := multipart.NewReader(msg.Body, pairs["boundary"])
		var err error
		var part *multipart.Part
		for part, err = parts.NextPart(); err == nil; part, err = parts.NextPart() {
			m, p := parseContentType(part.Header.Get("Content-Type"))
			if m == "text/plain" || (m == "text/html" && mime != "text/plain") {
				mime, pairs = m, p
				msg.Body = part
				for k, v := range part.Header {
					msg.Header[k] = v
				}
				break
			}
		}
	}
	if mime != "text/plain" && mime != "text/html" {
		return "", fmt.Errorf("can't find plain or html, mime %s can't process", mime)
	}
	if encoder := msg.Header.Get("Content-Transfer-Encoding"); encoder != "base64" {
		return "", fmt.Errorf("can't decode %s", encoder)
	}
	b, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		return "", err
	}
	b, _ = base64.StdEncoding.DecodeString(string(b))
	if charset := strings.ToLower(pairs["charset"]); charset != "utf-8" {
		buf := bytes.NewBuffer(b)
		reader, err := encodingex.NewIconvReadCloser(buf, "utf-8", charset)
		if err != nil {
			return "", err
		}
		b, err = ioutil.ReadAll(reader)
		if err != nil {
			return "", err
		}
	}

	replacer := strings.NewReplacer("\r\n", "\n", "\r", "\n")
	content := replacer.Replace(string(b))

	if mime == "text/html" {
		content = parseHtml(content)
	}
	content = parsePlain(content)

	return content, nil
}

func parseHtml(content string) string {
	var removeLine = [...]string{
		`(?iU)<script\b.*>.*</script>`,
		`(?iU)<style\b.*>.*</style>`,
		`(?U)<.*>`,
		`(?iU)<div class="gmail_quote".*`,
	}
	for _, remove := range removeLine {
		re := regexp.MustCompile(remove)
		content = re.ReplaceAllString(content, "\n")
	}
	return content
}

func parsePlain(content string) string {
	var replyLine = [...]string{
		"^--",
		"-*Original Message-*",
		"_____*",
		"Sent from",
		"Sent from",
		`^From:`,
		`^On (.*) wrote:`,
		"发自我的 iPhone",
		`EXFE ·X· <x\+[a-zA-Z0-9]*@exfe.com>`,
		`^>+`,
	}

	lines := make([]string, 0)
LINE:
	for _, l := range strings.Split(content, "\n") {
		l = strings.Trim(l, " \n\t")
		for _, reply := range replyLine {
			if ok, err := regexp.MatchString(reply, l); ok || err != nil {
				if err != nil {
					panic(err)
				}
				break LINE
			}
		}
		lines = append(lines, l)
	}
	return strings.Trim(strings.Join(lines, "\n"), "\n ")
}

func parseTitle(title string) string {
	got, charset, err := encodingex.DecodeEncodedWord(title)
	if err != nil {
		return title
	}
	got, err = encodingex.Conv(got, "UTF-8", charset)
	if err != nil {
		return title
	}
	return got
}
