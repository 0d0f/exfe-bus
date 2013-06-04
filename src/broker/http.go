package broker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mrjones/oauth"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var HttpClient *http.Client

func init() {
	HttpClient = http.DefaultClient
}

func SetProxy(host string) {
	transport := &http.Transport{
		Proxy: func(r *http.Request) (*url.URL, error) {
			sites := []string{"twitter.com", "facebook.com", "dropbox.com"}
			for _, site := range sites {
				if strings.HasSuffix(r.URL.Host, site) {
					return &url.URL{Host: host}, nil
				}
			}
			return nil, nil
		},
	}
	HttpClient = &http.Client{
		Transport: transport,
	}
}

type HttpError struct {
	Code    int
	Message string
}

func (e HttpError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func RestHttp(method, url, mime string, arg interface{}, reply interface{}) (int, error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(arg)
	if err != nil {
		return -1, err
	}
	resp, err := HttpResponse(Http(method, url, mime, buf.Bytes()))
	if err != nil {
		if e, ok := err.(HttpError); ok {
			return e.Code, err
		}
		return -1, err
	}
	defer resp.Close()

	decoder := json.NewDecoder(resp)
	err = decoder.Decode(reply)
	if err != nil {
		return -2, err
	}
	return http.StatusOK, nil
}

func Http(method, url, mime string, body []byte) (*http.Response, error) {
	buf := bytes.NewBuffer(body)
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}
	if mime != "" {
		req.Header.Set("Content-Type", mime)
	}
	return HttpClient.Do(req)
}

func HttpForm(url string, params url.Values) (io.ReadCloser, error) {
	return HttpResponse(HttpClient.PostForm(url, params))
}

type OAuth struct {
	*oauth.Consumer
}

func NewOAuth(client, secret string, provider oauth.ServiceProvider) OAuth {
	consumer := oauth.NewConsumer(client, secret, provider)
	consumer.HttpClient = HttpClient
	return OAuth{consumer}
}

func HttpResponse(resp *http.Response, err error) (io.ReadCloser, error) {
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, HttpError{resp.StatusCode, string(b)}
	}
	return resp.Body, nil
}
