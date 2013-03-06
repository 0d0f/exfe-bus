package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Printf("Usage: %s [pid file] [addr: exfe.com:25000] [cert] [key]\n", os.Args[0])
		return
	}
	pid := os.Args[1]
	addr := os.Args[2]
	certFile := os.Args[3]
	keyFile := os.Args[4]

	p, err := os.Create(pid)
	if err != nil {
		fmt.Printf("can't create pid %s: %s\n", pid, err)
		os.Exit(-1)
	}
	p.Write([]byte(fmt.Sprintf("%d", os.Getpid())))
	p.Close()

	log, err := logger.New(logger.Stderr, "imsg sender")
	if err != nil {
		panic(err)
		return
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Crit("loadkeys error: %s", err)
		return
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

	conn, err := tls.Dial("tcp", addr, &config)
	if err != nil {
		log.Crit("dial error: %s", err)
		return
	}
	defer conn.Close()
	log.Info("connected to: ", conn.RemoteAddr())

	buf := make([]byte, 512)

	for {
		conn.SetReadDeadline(time.Now().Add(time.Minute))
		n, err := conn.Read(buf)
		if err != nil {
			log.Crit("conn read: %s", err)
			return
		}
		var load Load
		err = json.Unmarshal(buf[:n], &load)
		if err != nil {
			log.Err("unmashal: %s", err)
			return
		}

		switch load.Type {
		case Ping:
			load.Type = Pong
			l, err := json.Marshal(load)
			if err != nil {
				log.Crit("marshal: %s", err)
				return
			}
			_, err = conn.Write(l)
			if err != nil {
				log.Crit("write: %s", err)
				return
			}
		case Send:
			log.Info("received send to %s", load.To)

			err := SendiMsg(load.To, load.Content)
			load.Type = Respond
			load.Content = ""
			if err != nil {
				load.Content = err.Error()
			}
			l, err := json.Marshal(load)
			if err != nil {
				log.Crit("marshal: %s", err)
				return
			}
			_, err = conn.Write(l)
			log.Info("respond: %s", load.Content)
			if err != nil {
				log.Crit("write: %s", err)
				return
			}
		}
	}

	log.Info("exiting")
}

func SendiMsg(to, content string) error {
	cmd := exec.Command("osascript", "/usr/local/bin/exfe_imsg.applescript", to, content)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	e, err := ioutil.ReadAll(stderr)
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		str := string(e)
		i := strings.LastIndex(str, "error: ")
		if i == -1 {
			return fmt.Errorf("%s", str)
		}
		return fmt.Errorf("%s", str[i+len("error: "):])
	}
	return nil
}
