package imsg

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
	"gobus"
	"model"
	"net"
	"time"
)

type IMsg struct {
	connChan chan Load
	isRun    bool
}

func NewiMsg(config *model.Config) (*IMsg, error) {
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
	}

	go listen(ret, addr, &c, connChan, config.Log)
	return ret, nil
}

func (i *IMsg) Send(meta *gobus.HTTPMeta, load *Load, r *int) error {
	if !i.isRun {
		return fmt.Errorf("%s", "can't connect to client sender")
	}
	i.connChan <- *load
	l := <-i.connChan
	if l.Content != "" {
		return fmt.Errorf("%s", l.Content)
	}
	return nil
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
