package mail

import (
	"bot/ics"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/googollee/go-aws/s3"
	"github.com/googollee/go-encoding"
	"io"
	"io/ioutil"
	"logger"
	"mime/multipart"
	"model"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

var htmlRegexp []*regexp.Regexp
var replyRegexp []*regexp.Regexp
var replacer *strings.Replacer
var typeMap = map[byte]string{
	'c': "cross_id",
	'e': "exfee_id",
}

func init() {
	for _, html := range []string{
		`(?iU)<script\b.*>.*</script>`,
		`(?iU)<style\b.*>.*</style>`,
		`(?iU)<div class="gmail_quote".*`,
		`(?U)<.*>`,
	} {
		htmlRegexp = append(htmlRegexp, regexp.MustCompile(html))
	}
	for _, reply := range []string{
		"^--",
		"-*Original Message-*",
		"_____*",
		"Sent from",
		`^From:`,
		`^On (.*) wrote:`,
		`EXFE ·X·.* wrote:`,
		`EXFE ·X·.*写道：`,
		"发自我的 iPhone",
		`EXFE ·X·"? *<x\+?[a-zA-Z0-9]*@exfe.com>`,
		`EXFE ·X·"? *<x\+?[a-zA-Z0-9]*@0d0f.com>`,
		`At \d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d.* wrote:`,
		`在 \d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d.*写道：`,
		`^>+`,
	} {
		replyRegexp = append(replyRegexp, regexp.MustCompile(reply))
	}
	replacer = strings.NewReplacer("\r\n", "\n", "\r", "\n")
}

func parseHtml(content string) string {
	for _, remove := range htmlRegexp {
		content = remove.ReplaceAllString(content, "\n")
	}
	return content
}

func parsePlain(content string) string {
	content = replacer.Replace(content)

	lines := make([]string, 0)
LINE:
	for _, l := range strings.Split(content, "\n") {
		l = strings.Trim(l, " \n\t")
		for _, reply := range replyRegexp {
			if reply.MatchString(l) {
				break LINE
			}
		}
		lines = append(lines, l)
	}
	return strings.Trim(strings.Join(lines, "\n"), "\n ")
}

type attach struct {
	url  string
	name string
}

type Parser struct {
	bucket *s3.Bucket

	from         *mail.Address
	addrList     []*mail.Address
	messageID    string
	referenceIDs []string
	subject      string
	config       *model.Config
	idRegexp     *regexp.Regexp
	domain       string
	date         time.Time
	attachments  []attach

	content     string
	contentMime string
	event       *ics.Event
}

func NewParser(msg *mail.Message, config *model.Config, bucket *s3.Bucket) (*Parser, error) {
	addrList := getMailAddress(msg, "From")
	if len(addrList) == 0 {
		return nil, fmt.Errorf("can't find From address")
	}
	from := addrList[0]
	addrList = append(addrList, getMailAddress(msg, "Cc")...)
	addrList = append(addrList, getMailAddress(msg, "To")...)

	ids := getMailIDs(msg, "Message-ID")
	if len(ids) == 0 {
		return nil, fmt.Errorf("can't find Message-ID")
	}
	msgID := ids[0]
	ids = append(ids, getMailIDs(msg, "References")...)

	subject := msg.Header.Get("Subject")
	if s, err := encoding.DecodeEncodedWord(subject); err == nil {
		subject = s
	}

	date, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", msg.Header.Get("Date"))
	if err != nil {
		date, err = time.Parse("Mon, 2 Jan 2006 15:04:05 -0700 MST", msg.Header.Get("Date"))
	}
	if err != nil {
		date, err = time.Parse("Mon, 2 Jan 2006 15:04:05 -0700 (MST)", msg.Header.Get("Date"))
	}
	if err != nil {
		date = time.Now().UTC()
	}

	ret := &Parser{
		bucket: bucket,

		from:         from,
		addrList:     addrList,
		messageID:    msgID,
		referenceIDs: ids,
		subject:      subject,
		date:         date,
		config:       config,
		domain:       config.Email.Domain,
		idRegexp:     regexp.MustCompile(config.Email.Prefix + "\\+([0-9a-zA-Z]+)@"),
	}
	err = ret.init(msg.Body, msg.Header)
	if err != nil {
		return nil, err
	}

	if len(ret.attachments) > 0 {
		for _, attach := range ret.attachments {
			ret.content += fmt.Sprintf("\n[%s](%s)", attach.name, attach.url)

		}
	}
	return ret, nil
}

func (h *Parser) init(r io.Reader, header mail.Header) error {
	mime, pairs := parseContentType(header.Get("Content-Type"))
	mimeType := mime
	if i := strings.Index(mimeType, "/"); i >= 0 {
		mimeType = mimeType[:i]
	}
	boundary := pairs["boundary"]
	charset := pairs["charset"]
	type_, pairs := parseContentType(header.Get("Content-Disposition"))
	replacer := strings.NewReplacer(" ", "_")
	filename := replacer.Replace(pairs["filename"])
	if s, err := encoding.DecodeEncodedWord(filename); err == nil {
		filename = s
	}

	if mimeType == "multipart" {
		parts := multipart.NewReader(r, boundary)
		for part, e := parts.NextPart(); e == nil; part, e = parts.NextPart() {
			if err := h.init(part, mail.Header(part.Header)); err != nil {
				return err
			}
		}
		return nil
	}

	switch mime {
	case "text/plain":
		if h.contentMime == "text/plain" {
			return nil
		}
		h.contentMime = "text/plain"
		content, err := getPartBody(r, header.Get("Content-Transfer-Encoding"), charset)
		if err != nil {
			return err
		}
		h.content = parsePlain(content)
		return nil
	case "text/html":
		if h.contentMime == "text/plain" || h.contentMime == "text/html" {
			return nil
		}
		h.contentMime = "text/html"
		content, err := getPartBody(r, header.Get("Content-Transfer-Encoding"), charset)
		if err != nil {
			return err
		}
		content = parseHtml(content)
		h.content = parsePlain(content)
		return nil
	case "application/octet-stream":
		if type_ != "attachment" || !strings.HasSuffix(filename, ".ics") {
			break
		}
		fallthrough
	case "text/calendar":
		fallthrough
	case "application/ics":
		if h.event != nil {
			return nil
		}
		body, err := getPartBody(r, header.Get("Content-Transfer-Encoding"), charset)
		if err != nil {
			return err
		}
		buf := bytes.NewBufferString(body)
		calendar, err := ics.ParseCalendar(buf)
		if err != nil {
			return err
		}
		if len(calendar.Event) > 0 {
			h.event = &calendar.Event[0]
		}
		return nil
	}

	if type_ == "attachment" {
		content, err := getPartBody(r, header.Get("Content-Transfer-Encoding"), charset)
		if err != nil {
			return err
		}
		buf := bytes.NewBufferString(content)
		hash := md5.New()
		hash.Write(buf.Bytes())
		path := fmt.Sprintf("%x", hash.Sum(nil))
		if h.bucket == nil {
			h.attachments = append(h.attachments, attach{"http://s3/email-attachment/" + path, filename})
		} else {
			obj, err := h.bucket.CreateObject(path, mime)
			if err == nil {
				err = obj.SaveReader(buf, int64(buf.Len()))
				if err != nil {
					logger.ERROR("can't save attachment %s: %s", filename, err)
				} else {
					h.attachments = append(h.attachments, attach{obj.URL(), filename})
				}
			} else {
				logger.ERROR("can't save attachment %s: %s", filename, err)
			}
		}
	}
	return nil
}

func (h *Parser) HasICS() bool {
	return h.event != nil
}

func (h *Parser) GetIDs() []string {
	ret := h.referenceIDs
	if h.HasICS() {
		ret = append(ret, h.event.ID)
	}
	return ret
}

func (h *Parser) Date() string {
	return h.date.UTC().Format("2006-01-02 15:04:05")
}

func (h *Parser) GetCross() (cross model.Cross) {
	if h.HasICS() {
		cross = h.convertEventToCross(*h.event, h.from)
	} else {
		cross.Description = h.content
	}
	cross.Title = h.subject
	cross.Time = &model.CrossTime{
		BeginAt: model.EFTime{
			Timezone: h.date.Format("-07:00"),
		},
	}

	check := make(map[string]bool)
	for _, i := range cross.Exfee.Invitations {
		check[fmt.Sprintf("%s@%s", i.Identity.ExternalID, i.Identity.Provider)] = true
	}
	cross.By = model.Identity{
		Provider:         "email",
		Name:             h.from.Name,
		ExternalID:       h.from.Address,
		ExternalUsername: h.from.Address,
	}
	for _, addr := range h.addrList {
		if pre := len(h.config.Email.Prefix) + 1; len(addr.Address) > pre &&
			(addr.Address[:pre] == h.config.Email.Prefix+"@" || addr.Address[:pre] == h.config.Email.Prefix+"+") &&
			strings.HasSuffix(addr.Address, h.domain) {
			continue
		}
		if strings.HasSuffix(addr.Address, "googlemail.com") {
			continue
		}
		if _, ok := check[fmt.Sprintf("%s@email", addr.Address)]; ok {
			continue
		}
		invitation := model.Invitation{
			Host: addr.Address == h.from.Address,
			Via:  "email",
			Identity: model.Identity{
				Provider:         "email",
				Name:             addr.Name,
				ExternalID:       addr.Address,
				ExternalUsername: addr.Address,
			},
			By: cross.By,
		}
		if addr.Address == h.from.Address {
			invitation.Response = model.Accepted
		}
		cross.Exfee.Invitations = append(cross.Exfee.Invitations, invitation)
	}

	cross.Widgets = []map[string]interface{}{
		map[string]interface{}{
			"type":  "Background",
			"image": "mail.jpg",
		},
	}
	cross.Attribute.State = "draft"

	return
}

func (h *Parser) GetPost() string {
	return h.content
}

func (h *Parser) GetTypeID() (string, string) {
	for _, addr := range h.addrList {
		if !strings.HasSuffix(addr.Address, h.domain) {
			continue
		}
		match := h.idRegexp.FindAllStringSubmatch(addr.Address, -1)
		if len(match) == 0 {
			continue
		}
		to := "cross_id"
		id := match[0][1]
		if t, ok := typeMap[id[0]]; ok {
			to = t
			id = id[1:]
		}
		return to, id
	}
	return "", ""
}

func getMailAddress(msg *mail.Message, k string) []*mail.Address {
	var ret []*mail.Address
	if lists := msg.Header.Get(k); lists != "" {
		for _, l := range strings.Split(lists, ",") {
			l = strings.Trim(l, " ")
			var addr mail.Address
			switch l[0] {
			case '"':
				last := strings.LastIndex(l, "\"")
				if last <= 0 {
					continue
				}
				addr.Name = strings.Trim(l[1:last], " ")
				addr.Address = strings.Trim(l[last+1:], " <>")
			case '=':
				last := strings.LastIndex(l, "=")
				if last <= 0 {
					continue
				}
				var err error
				addr.Name, err = encoding.DecodeEncodedWord(l[1 : last+1])
				if err != nil {
					continue
				}
				addr.Address = strings.Trim(l[last+1:], " <>")
			case '<':
				addr.Address = strings.Trim(l, " <>")
			default:
				last := strings.LastIndex(l, " ")
				if last <= 0 {
					addr.Address = strings.Trim(l, " <>")
				} else {
					addr.Name = strings.Trim(l[:last], " ")
					addr.Address = strings.Trim(l[last+1:], " <>")
				}
			}
			ret = append(ret, &addr)
		}
	}
	return ret
}

func getMailIDs(m *mail.Message, field string) []string {
	ref := m.Header.Get(field)
	if ref == "" {
		return nil
	}
	ret := strings.Split(ref, " ")
	for i, id := range ret {
		ret[i] = strings.Trim(id, " <>")
	}
	return ret
}

func (h *Parser) convertEventToCross(event ics.Event, from *mail.Address) model.Cross {
	places := strings.SplitN(event.Location, "\n", 2)
	place := model.Place{
		Title: places[0],
	}
	if len(places) > 1 {
		place.Description = places[1]
	}
	desc := event.Description
	if event.URL != "" {
		desc += "\n" + event.URL
	}
	time := model.CrossTime{
		BeginAt: model.EFTime{
			Date:     event.Start.Format("2006-01-02"),
			Timezone: event.Start.Format("-07:00"),
		},
		OutputFormat: model.TimeFormat,
	}
	format := "2006-01-02 15:04:05"
	if event.DateStart {
		format = "2006-01-02"
	} else {
		time.BeginAt.Time = event.Start.UTC().Format("15:04:05")
	}
	time.Origin = fmt.Sprintf("%s", event.Start.UTC().Format(format))
	var invitations []model.Invitation
	attendees := make(map[string]bool)
	by := model.Identity{
		Name:             from.Name,
		ExternalID:       from.Address,
		ExternalUsername: from.Address,
		Provider:         "email",
	}
	for _, a := range append(event.Attendees, event.Organizer) {
		if a.Email == "" || strings.HasSuffix(a.Email, h.domain) {
			continue
		}
		if _, ok := attendees[a.Email]; ok {
			continue
		}
		attendees[a.Email] = true
		host := a.Email == event.Organizer.Email
		identity := model.Identity{
			Name:             a.Name,
			ExternalID:       a.Email,
			ExternalUsername: a.Email,
			Provider:         "email",
		}
		rsvp := model.Noresponse
		switch a.PartStat {
		case "ACCEPTED":
			rsvp = model.Accepted
		case "DECLINED":
			rsvp = model.Declined
		}
		invitations = append(invitations, model.Invitation{
			Host:     host,
			Response: rsvp,
			Identity: identity,
			Via:      "email",
			By:       by,
		})
	}

	return model.Cross{
		By:          by,
		Title:       event.Summary,
		Place:       &place,
		Description: desc,
		Time:        &time,
		Exfee: model.Exfee{
			Invitations: invitations,
		},
	}
}

func getPartBody(r io.Reader, encoder string, charset string) (string, error) {
	switch strings.ToLower(encoder) {
	case "base64":
		r = encoding.NewIgnoreReader(r, []byte(" \r\n"))
		r = base64.NewDecoder(base64.StdEncoding, r)
	case "quoted-printable":
		r = encoding.NewQEncodingDecoder(r)
	}
	if charset = strings.ToLower(charset); charset != "" && charset != "utf-8" {
		if charset == "gb2312" {
			charset = "gbk"
		}
		var err error
		r, err = encoding.NewIconvReadCloser(r, "utf-8", charset)
		if err != nil {
			return "", err
		}
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(b), nil
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
		p := strings.SplitN(part, "=", 2)
		pairs[p[0]] = strings.Trim(p[1], "\"' ")
	}
	return mime, pairs
}
