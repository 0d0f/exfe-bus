package twitter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var iomGrabber = regexp.MustCompile("(?U)^(.* )?#([a-zA-Z][0-9a-zA-Z])( .*)?$")

type StreamingReader struct {
	buf      [2048]byte
	bufBegin int
	bufEnd   int
	reader   io.Reader
}

func NewStreamingReader(r io.Reader) *StreamingReader {
	return &StreamingReader{
		bufBegin: 0,
		bufEnd:   0,
		reader:   r,
	}
}

func (r *StreamingReader) ReadTweet() (tweet *Tweet, err error) {
	cache := make([]byte, 0, 0)
	for {
		end := r.find(r.buf[r.bufBegin:r.bufEnd], '\r')
		if end < 0 {
			cache = append(cache, r.buf[r.bufBegin:r.bufEnd]...)
			r.bufBegin = 0
			r.bufEnd, err = r.reader.Read(r.buf[:])
			if err != nil {
				return
			}
		} else {
			cache = append(cache, r.buf[r.bufBegin:r.bufBegin+end]...)
			r.bufBegin = r.bufBegin + end + 1
			if len(cache) > 1 {
				break
			}
			cache = cache[0:0]
		}
	}

	buf := bytes.NewBuffer(cache)
	decoder := json.NewDecoder(buf)
	tweet = new(Tweet)
	err = decoder.Decode(tweet)
	return
}

func (s *StreamingReader) find(data []byte, c rune) int {
	for i, d := range data {
		if rune(d) == c {
			return i
		}
	}
	return -1
}

type User struct {
	ID_        string `json:"id_str"`
	ScreenName string `json:"screen_name"`
}

type DirectMessage struct {
	Sender    User   `json:"sender"`
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
}

type Tweet struct {
	Entities struct {
		UserMentions []User `json:"user_mentions"`
	} `json:"entities"`
	CreatedAt         string         `json:"created_at"`
	Text              string         `json:"text"`
	InReplyToStatusId *string        `json:"in_reply_to_user_id_str"`
	User              User           `json:"user"`
	DirectMessage     *DirectMessage `json:"direct_message"`
}

func (t *Tweet) ToInput() *Input {
	ret := new(Input)
	if t.DirectMessage != nil {
		ret.ID = t.DirectMessage.Sender.ID_
		ret.ScreenName = t.DirectMessage.Sender.ScreenName
		ret.Text = t.DirectMessage.Text
		ret.CreatedAt = t.DirectMessage.CreatedAt
	} else {
		ret.ID = t.User.ID_
		ret.ScreenName = t.User.ScreenName
		ret.Text = t.Text
		ret.CreatedAt = t.CreatedAt
	}
	t_, err := time.Parse("Mon Jan 02 15:04:05 -0700 2006", ret.CreatedAt)
	if err != nil {
		t_ = time.Now()
	}
	ret.CreatedAt = t_.Format("2006-01-02 15:04:05 -0700")
	iom := iomGrabber.FindStringSubmatch(ret.Text)
	if len(iom) > 0 {
		ret.Iom = iom[2]
		ret.Text = strings.Replace(ret.Text, fmt.Sprintf("#%s", ret.Iom), "", 1)
	}
	return ret
}

type Input struct {
	ID         string
	ScreenName string
	Text       string
	Iom        string
	CreatedAt  string
}

func (i *Input) ToUrl() (ret url.Values) {
	ret = make(url.Values)
	ret.Add("iom", i.Iom)
	ret.Add("content", i.Text)
	ret.Add("external_id", i.ID)
	ret.Add("provider", "twitter")
	ret.Add("time", i.CreatedAt)
	return
}
