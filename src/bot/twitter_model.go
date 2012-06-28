package main

import (
	"time"
	"fmt"
	"strings"
	"regexp"
)

type User struct {
	Id_str      string
	Screen_name string
}

type DirectMessage struct {
	Sender     User
	Created_at string
	Text       string
}

func (d *DirectMessage) screen_name() string {
	return d.Sender.Screen_name
}

func (d *DirectMessage) text() string {
	return strings.Trim(d.Text, " \t\n\r")
}

func (d *DirectMessage) created_at() string {
	t, err := time.Parse("Mon January 02 15:04:05 -0700 2006", d.Created_at)
	if err != nil {
		return d.Created_at
	}
	return t.Format("2006-01-02 15:04:05 -0700")
}

func (d *DirectMessage) external_id() string {
	return d.Sender.Id_str
}

type Tweet struct {
	Entities struct {
		User_mentions []User
	}
	Created_at                string
	Text                      string
	In_reply_to_status_id_str *string
	User                      *User
	Direct_message            *DirectMessage
}

func (t *Tweet) text(screenName string) string {
	if t.Direct_message != nil {
		return t.Direct_message.text()
	}
	text := strings.Trim(t.Text, " \t\n\r")
	if t.Text[0] == '@' {
		t := strings.SplitN(text, " ", 2)
		if len(t) == 2 && t[0] == screenName {
			return t[1]
		}
	}
	return ""
}

func (t *Tweet) screen_name() string {
	if t.Direct_message != nil {
		return t.Direct_message.screen_name()
	}
	return t.User.Screen_name
}

func (t *Tweet) created_at() string {
	if t.Direct_message != nil {
		return t.Direct_message.created_at()
	}
	return t.Created_at
}

func (t *Tweet) external_id() string {
	if t.Direct_message != nil {
		return t.Direct_message.external_id()
	}
	return t.User.Id_str
}

func (t *Tweet) parse(hashPattern *regexp.Regexp, screenName string) (hash, post string) {
	post = t.text(screenName)
	hashs := hashPattern.FindAllString(post, -1)
	if len(hashs) > 0 {
		hash = strings.Trim(hashs[0], " #")
		post = strings.Trim(strings.Replace(post, fmt.Sprintf("#%s", hash), "", -1), " ")
	}
	return
}
