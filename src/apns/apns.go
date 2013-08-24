package apns

import (
	"bytes"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

var ErrTimeout = errors.New("timeout")

type Notification struct {
	DeviceToken        string
	Identifier         uint32
	ExpireAfterSeconds int

	Payload *Payload
}

type sendArg struct {
	notification *Notification
	ret          chan error
}

type errorType struct {
	body []byte
	err  NotificationError
}

// An Apn contain a ErrorChan channle when connected to apple server. When a notification sent wrong, you can get the error infomation from this channel.
type Apn struct {
	server   string
	conf     *tls.Config
	conn     *tls.Conn
	timeout  time.Duration
	sendChan chan sendArg
	quitChan chan int
	locker   sync.RWMutex
}

// New Apn with cert_filename and key_filename.
func New(server, certFile, keyFile string, timeout time.Duration) (*Apn, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	certificate := []tls.Certificate{cert}
	conf := &tls.Config{
		Certificates: certificate,
	}

	ret := &Apn{
		server:   server,
		conf:     conf,
		timeout:  timeout,
		sendChan: make(chan sendArg),
	}
	return ret, nil
}

func (a *Apn) Serve() error {
	if err := a.connect(); err != nil {
		return err
	}
	if !a.canServe() {
		return errors.New("has served")
	}
	defer a.Close()

	errChan := make(chan NotificationError)
	defer func() { close(errChan) }()

	go a.waitError(errChan)

	for {
		select {
		case <-a.quitChan:
			break
		case arg := <-a.sendChan:
			arg.ret <- a.send(arg.notification)
		case err := <-errChan:
			return err
		}
	}
	return nil
}

func (a *Apn) Send(notification *Notification) error {
	arg := sendArg{
		notification: notification,
		ret:          make(chan error),
	}
	a.sendChan <- arg
	select {
	case ret := <-arg.ret:
		return ret
	case <-time.After(a.timeout):
	}
	return ErrTimeout
}

func (a *Apn) Close() error {
	a.locker.Lock()
	defer a.locker.Unlock()

	if a.quitChan != nil {
		close(a.quitChan)
		a.quitChan = nil
	}
	return a.conn.Close()
}

func (a *Apn) canServe() bool {
	a.locker.Lock()
	defer a.locker.Unlock()

	if a.quitChan != nil {
		return false
	}
	a.quitChan = make(chan int)
	return true
}

func (a *Apn) waitError(errChan chan NotificationError) {
	p := make([]byte, 6, 6)
	n, err := a.conn.Read(p)
	select {
	case errChan <- NewNotificationError(p[:n], err):
	default:
	}
	return
}

func (a *Apn) connect() error {
	conn, err := net.Dial("tcp", a.server)
	if err != nil {
		return err
	}

	tlsConn := tls.Client(conn, a.conf)
	err = tlsConn.Handshake()
	if err != nil {
		return err
	}

	a.conn = tlsConn

	return nil
}

func (a *Apn) send(notification *Notification) error {
	if a.conn == nil {
		return fmt.Errorf("not connected")
	}

	tokenbin, err := hex.DecodeString(notification.DeviceToken)
	if err != nil {
		return err
	}

	payloadbyte, err := json.Marshal(notification.Payload)
	if err != nil {
		return err
	}
	expiry := time.Now().Add(time.Duration(notification.ExpireAfterSeconds) * time.Second).Unix()

	buffer := bytes.NewBuffer([]byte{})
	writer := newBigWriter(buffer)
	writer.Write(uint8(1))
	writer.Write(uint32(notification.Identifier))
	writer.Write(uint32(expiry))
	writer.Write(uint16(len(tokenbin)))
	writer.Write(tokenbin)
	writer.Write(uint16(len(payloadbyte)))
	writer.Write(payloadbyte)
	if writer.Error != nil {
		return writer.Error
	}
	pushPackage := buffer.Bytes()

	a.conn.SetWriteDeadline(time.Now().Add(a.timeout))
	_, err = a.conn.Write(pushPackage)
	if err != nil {
		return err
	}
	return nil
}
