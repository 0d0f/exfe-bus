package thirdpart

import (
	"fmt"
	"github.com/googollee/go-aws/smtp"
	"github.com/googollee/go-logger"
	"net"
)

type SmtpInstance struct {
	conn *smtp.Client
	log  *logger.SubLogger
}

func NewSmtpSenderInstance(log *logger.Logger, host string, a smtp.Auth) (*SmtpInstance, error) {
	s, err := smtp.Dial(fmt.Sprintf("%s:25", host))
	if err != nil {
		return nil, err
	}
	if ok, _ := s.Extension("STARTTLS"); ok {
		if err = s.StartTLS(nil); err != nil {
			return nil, err
		}
	}
	err = s.Auth(a)
	if err != nil {
		return nil, err
	}
	sublog := log.SubPrefix(host)
	sublog.Debug("sender connect to %s", host)
	return &SmtpInstance{
		conn: s,
		log:  sublog,
	}, nil
}

func NewSmtpCheckerInstance(host string, log *logger.Logger) (*SmtpInstance, error) {
	mx, err := net.LookupMX(host)
	if err != nil {
		return nil, fmt.Errorf("lookup mail exchange fail: %s", err)
	}
	if len(mx) == 0 {
		return nil, fmt.Errorf("can't find mail exchange of %s", host)
	}
	s, err := smtp.Dial(fmt.Sprintf("%s:25", mx[0].Host))
	if err != nil {
		return nil, fmt.Errorf("dial to mail exchange %s fail: %s", mx[0].Host, err)
	}
	sublog := log.SubPrefix(host)
	sublog.Debug("checker connect to %s", host)
	return &SmtpInstance{
		conn: s,
		log:  sublog,
	}, nil
}

func (i *SmtpInstance) Ping() error {
	i.log.Debug("ping")
	return i.conn.Reset()
}

func (i *SmtpInstance) Close() error {
	i.log.Debug("close")
	return i.conn.Quit()
}

func (i *SmtpInstance) Error(err error) {
	i.log.Err("%s", err)
	return
}
