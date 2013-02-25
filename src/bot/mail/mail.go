package mail

import (
	"bytes"
	"code.google.com/p/go-imap/go1/imap"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/googollee/go-encoding-ex"
	"github.com/googollee/go-logger"
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

	go func() {
		defer tomb.Done()

		for {
			select {
			case <-tomb.Dying():
				log.Debug("quit")
			case <-time.After(timeout):

			}
		}
	}()
	return &tomb, nil
}

func ProcessMail(config model.Config, log *logger.SubLogger) {
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
	// var okIds []uint32
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

		if ok, _ := findAddress(fmt.Sprintf("x+c[0-9]*?@%s", config.Email.Domain), addrList); ok {
		} else if ok, _ := findAddress(fmt.Sprintf("x@%s", config.Email.Domain), addrList); ok {
		} else {
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
	decoder := base64.NewDecoder(base64.StdEncoding, msg.Body)
	if charset := strings.ToLower(pairs["charset"]); charset != "utf-8" {
		var err error
		decoder, err = encodingex.NewIconvReadCloser(decoder, "utf-8", charset)
		if err != nil {
			return "", fmt.Errorf("can't convert from charset %s", charset)
		}
	}
	content, err := ioutil.ReadAll(decoder)
	return string(content), err
}
