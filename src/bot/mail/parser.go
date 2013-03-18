package mail

import (
	"bot/ics"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/googollee/go-encoding-ex"
	"io"
	"io/ioutil"
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
		"Sent from",
		`^From:`,
		`^On (.*) wrote:`,
		"发自我的 iPhone",
		`EXFE ·X· <x\+[a-zA-Z0-9]*@exfe.com>`,
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

type Parser struct {
	from         *mail.Address
	addrList     []*mail.Address
	messageID    string
	referenceIDs []string
	subject      string
	config       *model.Config
	idRegexp     *regexp.Regexp
	domain       string
	date         time.Time

	content     string
	contentMime string
	event       *ics.Event
}

func NewParser(msg *mail.Message, config *model.Config) (*Parser, error) {
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
	if s, err := encodingex.DecodeEncodedWord(subject); err == nil {
		subject = s
	}

	date, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", msg.Header.Get("Date"))
	if err != nil {
		date = time.Now()
	}

	ret := &Parser{
		from:         from,
		addrList:     addrList,
		messageID:    msgID,
		referenceIDs: ids,
		subject:      subject,
		date:         date.UTC(),
		config:       config,
		domain:       config.Email.Domain,
		idRegexp:     regexp.MustCompile(config.Email.Prefix + "\\+([0-9a-zA-Z]+)@"),
	}
	err = ret.init(msg.Body, msg.Header)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (h *Parser) init(r io.Reader, header mail.Header) error {
	mime, pairs := parseContentType(header.Get("Content-Type"))

	switch mime {
	case "multipart/mixed":
		fallthrough
	case "multipart/alternative":
		parts := multipart.NewReader(r, pairs["boundary"])
		for part, e := parts.NextPart(); e == nil; part, e = parts.NextPart() {
			err := h.init(part, mail.Header(part.Header))
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	case "text/plain":
		if h.contentMime == "text/plain" {
			return nil
		}
		h.contentMime = "text/plain"
		content, err := getPartBody(r, header.Get("Content-Transfer-Encoding"), pairs["charset"])
		if err != nil {
			return err
		}
		h.content = parsePlain(content)
	case "text/html":
		if h.contentMime == "text/plain" || h.contentMime == "text/html" {
			return nil
		}
		h.contentMime = "text/html"
		content, err := getPartBody(r, header.Get("Content-Transfer-Encoding"), pairs["charset"])
		if err != nil {
			return err
		}
		content = parseHtml(content)
		h.content = parsePlain(content)
	case "text/calendar":
		fallthrough
	case "application/ics":
		if h.event != nil {
			return nil
		}
		body, err := getPartBody(r, header.Get("Content-Transfer-Encoding"), pairs["charset"])
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
	fmt.Printf("%+v\n", h.event)
	fmt.Println(ret)
	return ret
}

func (h *Parser) Date() string {
	return h.date.Format("2006-01-02 15:04:05")
}

func (h *Parser) GetCross() (cross model.Cross) {
	if h.HasICS() {
		cross = h.convertEventToCross(*h.event, h.from)
	} else {
		cross.Title = h.subject
		cross.Description = h.content
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
		if strings.HasSuffix(addr.Address, h.domain) {
			continue
		}
		if _, ok := check[fmt.Sprintf("%s@email", addr.Address)]; ok {
			continue
		}
		cross.Exfee.Invitations = append(cross.Exfee.Invitations, model.Invitation{
			Host: addr.Address == h.from.Address,
			Via:  "email",
			Identity: model.Identity{
				Provider:         "email",
				Name:             addr.Name,
				ExternalID:       addr.Address,
				ExternalUsername: addr.Address,
			},
			By:         cross.By,
			RsvpStatus: model.RsvpNoresponse,
		})
	}
	return
}

func (h *Parser) GetPost() string {
	fmt.Println(h.contentMime, h.content)
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
	if msg.Header.Get(k) != "" {
		ret, _ = msg.Header.AddressList(k)
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
			Date:     event.Start.UTC().Format("2006-01-02"),
			Timezone: "+00:00",
		},
		OutputFormat: model.TimeFormat,
	}
	format := "2006-01-02 15:04:05"
	if event.DateStart {
		format = "2006-01-02"
	} else {
		time.BeginAt.Time = event.Start.UTC().Format("15:04:05")
	}
	time.Origin = fmt.Sprintf("%s - %s", event.Start.UTC().Format(format), event.End.UTC().Format(format))
	invitations := make([]model.Invitation, 0)
	attendees := make(map[string]bool)
	by := model.Identity{
		Name:             from.Name,
		ExternalID:       from.Address,
		ExternalUsername: from.Address,
		Provider:         "email",
	}
	for _, a := range append(event.Attendees, event.Organizer) {
		if strings.HasSuffix(a.Email, h.domain) {
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
		rsvp := model.RsvpNoresponse
		switch a.PartStat {
		case "ACCEPTED":
			rsvp = model.RsvpAccepted
		case "DECLINED":
			rsvp = model.RsvpDeclined
		}
		invitations = append(invitations, model.Invitation{
			Host:       host,
			RsvpStatus: rsvp,
			Identity:   identity,
			By:         by,
		})
	}

	return model.Cross{
		By:          by,
		Title:       event.Summary,
		Place:       place,
		Description: desc,
		Time:        time,
		Exfee: model.Exfee{
			Invitations: invitations,
		},
	}
}

func getPartBody(r io.Reader, encoder string, charset string) (string, error) {
	switch strings.ToLower(encoder) {
	case "base64":
		r = encodingex.NewIgnoreReader(r, []byte(" \r\n"))
		r = base64.NewDecoder(base64.StdEncoding, r)
	case "quoted-printable":
		r = encodingex.NewQEncodingDecoder(r)
	default:
		return "", fmt.Errorf("can't decode %s", encoder)
	}
	if charset = strings.ToLower(charset); charset != "utf-8" {
		if charset == "gb2312" {
			charset = "gbk"
		}
		var err error
		r, err = encodingex.NewIconvReadCloser(r, "utf-8", charset)
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
		p := strings.Split(part, "=")
		pairs[p[0]] = strings.Trim(p[1], "\"' ")
	}
	return mime, pairs
}
