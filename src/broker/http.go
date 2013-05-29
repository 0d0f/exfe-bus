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
)

var HttpClient *http.Client

func init() {
	HttpClient = http.DefaultClient
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
	resp, err := Http(method, url, mime, buf.Bytes())
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

func Http(method, url, mime string, body []byte) (io.ReadCloser, error) {
	buf := bytes.NewBuffer(body)
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mime)
	resp, err := HttpClient.Do(req)
	return HttpResponse(resp, err)
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
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, HttpError{resp.StatusCode, string(b)}
	}
	return resp.Body, nil
}
