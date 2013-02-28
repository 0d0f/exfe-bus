package imsg

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"formatter"
	"github.com/googollee/go-logger"
	"model"
	"net"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

type IMsg struct {
	connChan chan Load
	isRun    bool
	config   *model.Config
}

func New(config *model.Config) (*IMsg, error) {
	connChan := make(chan Load)

	cert, err := tls.LoadX509KeyPair(config.Thirdpart.Apn.Cert, config.Thirdpart.Apn.Key)
	if err != nil {
		return nil, fmt.Errorf("load cert error: %s", err)
	}
	c := tls.Config{Certificates: []tls.Certificate{cert}}
	c.Rand = rand.Reader
	addr := "0.0.0.0:25000"

	ret := &IMsg{
		connChan: connChan,
		isRun:    false,
		config:   config,
	}

	go listen(ret, addr, &c, connChan, config.Log)
	return ret, nil
}

func (i *IMsg) Provider() string {
	return "imessage"
}

func (i *IMsg) Send(to *model.Recipient, text string) (string, error) {
	for _, line := range strings.Split(text, "\n") {
		line = strings.Trim(line, " \r\n\t")
		line = tailUrlRegex.ReplaceAllString(line, "")
		line = tailQuoteUrlRegex.ReplaceAllString(line, `)\)`)
		if line == "" {
			continue
		}

		cutter, err := formatter.CutterParse(line, imsgLen)
		if err != nil {
			return "", fmt.Errorf("parse cutter error: %s", err)
		}

		load := Load{
			Type:    Send,
			To:      to.ExternalID,
			Content: line,
		}

		for _, content := range cutter.Limit(140) {
			load.Content = content
			i.connChan <- load
			l := <-i.connChan
			if l.Content != "" {
				i.config.Log.Err("imessage send error: %s", l.Content)
			}
		}
	}
	return "", nil
}

func listen(imsg *IMsg, addr string, c *tls.Config, connChan chan Load, log *logger.Logger) {
	for {
		log.Info("imsg server: listening")
		listener, err := tls.Listen("tcp", addr, c)
		if err != nil {
			log.Err("listen error: %s", err)
			time.Sleep(time.Second)
			continue
		}
		conn, err := listener.Accept()
		if err != nil {
			log.Err("accept error: %s", err)
			time.Sleep(time.Second)
			continue
		}
		listener.Close()
		imsg.isRun = true
		handleClient(conn, connChan, log.SubCode())
		imsg.isRun = false
	}
}

func handleClient(conn net.Conn, connChan chan Load, log *logger.SubLogger) {
	defer conn.Close()
	log.Info("accepted from: %s", conn.RemoteAddr())

	for {
		select {
		case load := <-connChan:
			load.Type = Send
			l, err := json.Marshal(load)
			if err != nil {
				log.Crit("marshal error: %s", err)
				return
			}
			_, err = conn.Write(l)
			if err != nil {
				log.Err("write error: %s", err)
				return
			}
			log.Info("send to %s", load.To)

			reply := make([]byte, 256)
			conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			n, err := conn.Read(reply)
			if err != nil {
				log.Err("read error: %s", err)
				load.Type = Respond
				load.Content = fmt.Sprintf("read error: %s", err)
				connChan <- load
				return
			}
			err = json.Unmarshal(reply[:n], &load)
			if err != nil {
				log.Err("unmashal error: %s", err)
				load.Content = fmt.Sprintf("unmashal error: %s", err)
				connChan <- load
				return
			}
			if load.Type != Respond {
				log.Info("client no respond")
				load.Content = fmt.Sprintf("client no respond")
				connChan <- load
				return
			}
			connChan <- load
		case <-time.After(10 * time.Second):
			var load Load

			load.Type = Ping
			load.To = ""
			load.Content = ""
			l, err := json.Marshal(load)
			if err != nil {
				log.Crit("marshal error: %s", err)
				return
			}
			_, err = conn.Write(l)
			if err != nil {
				log.Err("write error: %s", err)
				return
			}

			reply := make([]byte, 256)
			conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			n, err := conn.Read(reply)
			if err != nil {
				log.Err("read error: %s", err)
				return
			}
			err = json.Unmarshal(reply[:n], &load)
			if err != nil {
				log.Err("unmashal error: %s", err)
				return
			}
			if load.Type != Pong {
				log.Info("client not respone ping")
				return
			}
		}
	}

	log.Info("closed")
}

func imsgLen(content string) int {
	return utf8.RuneCountInString(content)
}

var tailUrlRegex = regexp.MustCompile(` *(http|https):\/\/exfe.com(\/[\w#!:.?+=&%@!\-\/]*)?$`)
var tailQuoteUrlRegex = regexp.MustCompile(` *(http|https):\/\/exfe.com(\/[\w#!:.?+=&%@!\-\/]*)?\)(\\\))$`)
