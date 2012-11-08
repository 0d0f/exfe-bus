package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"thirdpart/imsg"
	"time"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Printf("Usage: %s [addr: exfe.com:25000] [cert] [key] [time out in second]\n", os.Args[0])
		return
	}
	addr := os.Args[1]
	certFile := os.Args[2]
	keyFile := os.Args[3]
	t, err := strconv.Atoi(os.Args[4])
	if err != nil {
		fmt.Printf("Usage: %s [addr: exfe.com:25000] [cert] [key] [time out in second]\n", os.Args[0])
		return
	}
	timeout := time.Duration(t) * time.Second
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
		var load imsg.Load
		err = json.Unmarshal(buf[:n], &load)
		if err != nil {
			log.Err("unmashal: %s", err)
			return
		}

		switch load.Type {
		case imsg.Ping:
			log.Info("received ping")

			load.Type = imsg.Pong
			l, err := json.Marshal(load)
			if err != nil {
				log.Crit("marshal: %s", err)
				return
			}
			_, err = conn.Write(l)
			log.Info("pong")
			if err != nil {
				log.Crit("write: %s", err)
				return
			}
		case imsg.Send:
			log.Info("received send to %s", load.To)

			err := SendiMsg(load.To, load.Content)
			load.Type = imsg.Respond
			load.Content = ""
			if err != nil {
				load.Content = err.Error()
			}
			l, err := json.Marshal(load)
			if err != nil {
				log.Crit("marshal: %s", err)
				return
			}
			time.Sleep(timeout)
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
