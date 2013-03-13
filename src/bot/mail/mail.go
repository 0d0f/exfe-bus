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
	"io/ioutil"
	"launchpad.net/tomb"
	"mime/multipart"
	"model"
	"net"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

func GetID(m *mail.Message, field string) string {
	ret := m.Header.Get(field)
	return strings.Trim(ret, " <>")
}

func GetReferencesID(m *mail.Message) []string {
	ref := m.Header.Get("References")
	if ref == "" {
		return nil
	}
	ret := strings.Split(ref, " ")
	for i, id := range ret {
		ret[i] = strings.Trim(id, " <>")
	}
	return ret
}

type MessageIDSaver interface {
	Save(id []string, crossID string) error
	Check(id []string) (string, bool, error)
}

var typeId = map[uint8]string{
	'c': "cross_id",
	'e': "exfee_id",
}

type Worker struct {
	Tomb tomb.Tomb

	config   *model.Config
	log      *logger.SubLogger
	templ    *formatter.LocalTemplate
	platform *broker.Platform
	saver    MessageIDSaver

	htmlRegexp  []*regexp.Regexp
	replyRegexp []*regexp.Regexp
}

func New(config *model.Config, templ *formatter.LocalTemplate, platform *broker.Platform, saver MessageIDSaver) (*Worker, error) {
	var htmlRegexp []*regexp.Regexp
	for _, html := range []string{
		`(?iU)<script\b.*>.*</script>`,
		`(?iU)<style\b.*>.*</style>`,
		`(?iU)<div class="gmail_quote".*`,
		`(?U)<.*>`,
	} {
		re, err := regexp.Compile(html)
		if err != nil {
			return nil, fmt.Errorf("can't compile %s html: %s", html, err)
		}
		htmlRegexp = append(htmlRegexp, re)
	}
	var replyRegexp []*regexp.Regexp
	for _, reply := range []string{
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
	} {
		re, err := regexp.Compile(reply)
		if err != nil {
			return nil, fmt.Errorf("can't compile %s reply: %s", reply, err)
		}
		replyRegexp = append(replyRegexp, re)
	}

	return &Worker{
		config:   config,
		log:      config.Log.SubPrefix("mail"),
		templ:    templ,
		platform: platform,
		saver:    saver,

		htmlRegexp:  htmlRegexp,
		replyRegexp: replyRegexp,
	}, nil
}

func (w *Worker) Daemon() {
	defer w.Tomb.Done()

	timeout := time.Duration(w.config.Bot.Email.TimeoutInSecond) * time.Second

	for true {
		select {
		case <-w.Tomb.Dying():
			w.log.Notice("quitted")
			return
		case <-time.After(timeout):
			w.process()
		}
	}
}

func (w *Worker) process() {
	w.log.Debug("process...")
	w.log.Notice("process...")

	conn, imapConn, err := w.login()
	if err != nil {
		w.log.Err("can't connect to %s: %s", w.config.Bot.Email.IMAPHost, err)
		return
	}
	w.log.Notice("logined")
	defer imapConn.Logout(broker.NetworkTimeout)
	conn.SetDeadline(time.Now().Add(broker.ProcessTimeout))

	_, err = imapConn.Select("INBOX", false)
	if err != nil {
		w.log.Err("can't select INBOX: %s", err)
		return
	}
	w.log.Notice("selected inbox")

	cmd, err := imap.Wait(imapConn.Search("UNSEEN"))
	if err != nil {
		w.log.Err("can't seach UNSEEN: %s", err)
		return
	}
	w.log.Notice("done search unseen")
	var ids []uint32
	for _, resp := range cmd.Data {
		ids = append(ids, resp.SearchResults()...)
	}

	var errorIds []uint32
	var okIds []uint32
	for _, id := range ids {
		conn.SetDeadline(time.Now().Add(broker.ProcessTimeout))

		msg, err := w.getMail(imapConn, id)
		if err != nil {
			w.log.Err("can't get mail %d: %s", id, err)
			errorIds = append(errorIds, id)
			continue
		}
		w.log.Notice("got %d", id)
		err = w.parseMail(msg)
		if err != nil {
			w.log.Err("parse mail %d failed: %s", id, err)
			errorIds = append(errorIds, id)
			continue
		}
		w.log.Notice("parsed %d", id)
		okIds = append(okIds, id)
		w.log.Notice("handled %d", id)
	}
	w.log.Debug("id:%v, ok:%v, err:%v", ids, okIds, errorIds)
	if err := w.copy(imapConn, okIds, "posted"); err != nil {
		w.log.Err("can't copy %v to posted: %s", errorIds, err)
	}
	w.log.Notice("copied %v to posted", okIds)
	if err := w.copy(imapConn, errorIds, "error"); err != nil {
		w.log.Err("can't copy %v to error: %s", errorIds, err)
	}
	w.log.Notice("copied %v to error", errorIds)
	if err := w.delete(imapConn, ids); err != nil {
		w.log.Err("can't remove %v from inbox: %s", ids, err)
	}
	w.log.Notice("removed %v", ids)
}

func (w *Worker) copy(conn *imap.Client, ids []uint32, folder string) error {
	if len(ids) == 0 {
		return nil
	}
	set := new(imap.SeqSet)
	set.AddNum(ids...)
	_, err := imap.Wait(conn.Copy(set, folder))
	return err
}

func (w *Worker) delete(conn *imap.Client, ids []uint32) error {
	if len(ids) == 0 {
		return nil
	}
	set := new(imap.SeqSet)
	set.AddNum(ids...)
	_, err := imap.Wait(conn.Store(set, "FLAGS", "\\Deleted"))
	return err
}

func (w *Worker) login() (net.Conn, *imap.Client, error) {
	c, err := net.DialTimeout("tcp", w.config.Bot.Email.IMAPHost, broker.NetworkTimeout)
	if err != nil {
		return nil, nil, err
	}
	tlsConn := tls.Client(c, nil)

	conn, err := imap.NewClient(tlsConn, strings.Split(w.config.Bot.Email.IMAPHost, ":")[0], broker.NetworkTimeout)
	if err != nil {
		return nil, nil, err
	}
	c.SetDeadline(time.Now().Add(broker.NetworkTimeout))

	conn.Data = nil
	if conn.Caps["STARTTLS"] {
		conn.StartTLS(nil)
	}

	if conn.State() == imap.Login {
		_, err = conn.Login(w.config.Bot.Email.IMAPUser, w.config.Bot.Email.IMAPPassword)
		if err != nil {
			return nil, nil, err
		}
	}

	return c, conn, nil
}

func (w *Worker) getMail(conn *imap.Client, id uint32) (*mail.Message, error) {
	buf := bytes.NewBuffer(nil)
	set := new(imap.SeqSet)
	set.AddNum(id)

	cmd, err := imap.Wait(conn.Fetch(set, "RFC822"))
	if err != nil {
		return nil, err
	}
	buf.Write(imap.AsBytes(cmd.Data[0].MessageInfo().Attrs["RFC822"]))

	msg, err := mail.ReadMessage(buf)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (w *Worker) getMailAddress(msg *mail.Message, k string) []*mail.Address {
	if msg.Header.Get(k) != "" {
		list, err := msg.Header.AddressList(k)
		if err == nil {
			return list
		}
	}
	return nil
}

func (w *Worker) parseMail(msg *mail.Message) error {
	var addrList []*mail.Address
	if list := w.getMailAddress(msg, "To"); len(list) != 0 {
		addrList = append(addrList, list...)
	}
	if list := w.getMailAddress(msg, "Cc"); len(list) != 0 {
		addrList = append(addrList, list...)
	}
	var from *mail.Address
	if list := w.getMailAddress(msg, "From"); len(list) != 0 {
		addrList = append(addrList, list...)
		from = list[0]
	} else {
		return fmt.Errorf("no from field")
	}

	refIDs := GetReferencesID(msg)
	msgID := GetID(msg, "Message-ID")
	if msgID != "" {
		refIDs = append(refIDs, msgID)
	}
	replyID := GetID(msg, "In-Reply-To")
	if replyID != "" {
		refIDs = append(refIDs, replyID)
	}
	crossID, crossExist, err := w.saver.Check(refIDs)
	if err != nil {
		return err
	}

	subject := msg.Header.Get("Subject")
	if s, err := encodingex.DecodeEncodedWord(subject); err == nil {
		subject = s
	}
	content, err := w.getContent(msg)
	if err != nil {
		return err
	}

	code := 500
	if ok, args := findAddress(fmt.Sprintf("x\\+([0-9a-zA-Z]+)@%s", w.config.Email.Domain), addrList); ok {
		code, err = w.sendPost(args[0], from, addrList, content)
		w.saver.Save(refIDs, args[0])
	} else if ok, _ := findAddress(fmt.Sprintf("x@%s", w.config.Email.Domain), addrList); ok {
		if !crossExist {
			var id uint64
			id, code, err = w.createCross(from, addrList, subject, content)
			if err == nil {
				crossID = fmt.Sprintf("%d", id)
				w.saver.Save(refIDs, crossID)
			}
		} else {
			code, err = w.sendPost(crossID, from, addrList, content)
			w.saver.Save(refIDs, crossID)
		}
	} else {
		code = http.StatusBadRequest
		err = fmt.Errorf("can't parse mail list: %v", addrList)
	}
	if err != nil {
		w.sendHelp(code, err, msgID, from, subject, content)
		return err
	}
	return nil
}

func (w *Worker) sendPost(arg string, from *mail.Address, addrs []*mail.Address, post string) (int, error) {
	w.log.Debug("send post(%s) from(%s) to x+%s", post, from.Address, arg)
	to := "cross"
	toId := arg

	if t, ok := typeId[arg[0]]; ok {
		to = t
		toId = arg[1:]
	}

	code, err := w.platform.BotPostConversation(from.Address, post, to, toId)
	if err != nil {
		return code, err
	}

	var identities []model.Identity
	self := fmt.Sprintf("x+%s@%s", arg, w.config.Email.Domain)
	for _, addr := range addrs {
		if addr.Address == self {
			continue
		}
		identities = append(identities, model.Identity{
			Provider:         "email",
			Name:             addr.Name,
			ExternalID:       addr.Address,
			ExternalUsername: addr.Address,
		})
	}
	if len(identities) == 0 {
		return code, err
	}
	code, err = w.platform.BotCrossInvite(to, toId, identities)
	return code, err
}

func (w *Worker) createCross(from *mail.Address, list []*mail.Address, title, desc string) (uint64, int, error) {
	w.log.Debug("create x(%s) for %v", title, list)
	cross := model.Cross{
		Title:       title,
		Description: desc,
		By: model.Identity{
			Provider:         "email",
			Name:             from.Name,
			ExternalID:       from.Address,
			ExternalUsername: from.Address,
		},
	}
	invite := make([]model.Invitation, len(list))
	for i, addr := range list {
		invite[i] = model.Invitation{
			Host: addr.Address == from.Address,
			Via:  "email",
			Identity: model.Identity{
				Provider:         "email",
				Name:             addr.Name,
				ExternalID:       addr.Address,
				ExternalUsername: addr.Address,
			},
		}
	}
	cross.Exfee.Invitations = invite
	id, status, err := w.platform.BotCrossGather(cross)
	if err != nil {
		return 0, status, err
	}
	return id, status, err
}

func (w *Worker) sendHelp(code int, err error, msgID string, from *mail.Address, subject, content string) error {
	w.log.Debug("send help to %v", from)
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
		Subject:   subject,
		MessageID: msgID,
		Text:      content,
	}
	err = w.templ.Execute(buf, "en_US", "email/conversation_reply", email)
	if err != nil {
		w.log.Crit("template(conversation_reply.email) failed: %s", err)
	}

	to := model.Recipient{
		Provider:         "email",
		Name:             from.Name,
		ExternalID:       from.Address,
		ExternalUsername: from.Address,
	}
	_, err = w.platform.Send(to, buf.String())
	return err
}

func (w *Worker) getContent(msg *mail.Message) (string, error) {
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
		content = w.parseHtml(content)
	}
	content = w.parsePlain(content)

	return content, nil
}

func (w *Worker) parseHtml(content string) string {
	for _, remove := range w.htmlRegexp {
		content = remove.ReplaceAllString(content, "\n")
	}
	return content
}

func (w *Worker) parsePlain(content string) string {
	lines := make([]string, 0)
LINE:
	for _, l := range strings.Split(content, "\n") {
		l = strings.Trim(l, " \n\t")
		for _, reply := range w.replyRegexp {
			if reply.MatchString(l) {
				break LINE
			}
		}
		lines = append(lines, l)
	}
	return strings.Trim(strings.Join(lines, "\n"), "\n ")
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
