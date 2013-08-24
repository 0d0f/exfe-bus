package apn

import (
	"apns"
	"broker"
	"encoding/json"
	"fmt"
	"logger"
	"model"
	"regexp"
	"strings"
	"sync"
	"thirdpart"
	"time"
)

type sendArg struct {
	id      uint
	content string
}

type Apn struct {
	id         uint32
	apps       map[string]*apns.Apn
	defaultApp string
	callback   thirdpart.Callback
	locker     sync.RWMutex
}

func New(config *model.Config) (*Apn, error) {
	ret := &Apn{
		id:         0,
		apps:       make(map[string]*apns.Apn),
		defaultApp: config.Thirdpart.Apn.Default,
	}
	for k, app := range config.Thirdpart.Apn.Apps {
		apn, err := apns.New(app.Server, app.Cert, app.Key, broker.NetworkTimeout)
		if err != nil {
			return nil, fmt.Errorf("apn %s error: %s", k, err)
		}
		ret.apps[k] = apn
		go func(name string) {
			for {
				err := apn.Serve()
				if notificationError, ok := err.(apns.NotificationError); ok && notificationError.Status == apns.ErrorInvalidToken {
					ret.locker.RLock()
					ret.callback(ret.postId(name, notificationError.Identifier), notificationError)
					ret.locker.RUnlock()
				} else {
					logger.ERROR("app %s push error: %s", name, err)
				}
			}
		}(k)
	}

	return ret, nil
}

func (a *Apn) postId(app string, id uint32) string {
	return fmt.Sprintf("%s-%d", app, id)
}

func (a *Apn) Provider() string {
	return "iOS"
}

func (a *Apn) SetPosterCallback(callback thirdpart.Callback) (time.Duration, bool) {
	a.locker.Lock()
	defer a.locker.Unlock()
	a.callback = callback
	return 10 * time.Second, true
}

func (a *Apn) Post(from, id, text string) (string, error) {
	text = strings.Trim(text, " \r\n")
	last := strings.LastIndex(text, "\n")
	if last == -1 {
		return "", fmt.Errorf("no payload")
	}
	dataStr := text[last+1:]
	var data interface{}
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		return "", fmt.Errorf("last line of text(%s) can't unmarshal: %s", dataStr, err)
	}
	text = strings.Trim(text[:last], " \r\n")
	text = tailUrlRegex.ReplaceAllString(text, "")

	a.locker.Lock()
	ret := a.id
	a.id++
	a.locker.Unlock()

	payload := apns.Payload{}
	payload.Aps.Alert.Body = text
	payload.Aps.Badge = 1
	payload.Aps.Sound = "default"
	if data != nil {
		payload.SetCustom("args", data)
	}
	notification := apns.Notification{
		DeviceToken: id,
		Identifier:  ret,
		Payload:     &payload,
	}

	spliter := strings.LastIndex(id, "@")
	appName := a.defaultApp
	if spliter >= 0 {
		appName = id[spliter+1:]
		id = id[:spliter]
	}

	retId := a.postId(appName, notification.Identifier)
	app, ok := a.apps[appName]
	if !ok {
		return retId, fmt.Errorf("invalid app name: %s", appName)
	}
	err = app.Send(&notification)
	return retId, err
}

var tailUrlRegex = regexp.MustCompile(` *(http|https):\/\/exfe.com(\/[\w#!:.?+=&%@!\-\/]*)?$`)
