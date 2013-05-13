package broker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func RestHttp(method, url, mime string, arg interface{}, reply interface{}) (int, error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(arg)
	if err != nil {
		return -1, err
	}
	resp, err := Http(method, url, mime, buf.Bytes())
	if err != nil {
		if resp != nil {
			return resp.StatusCode, err
		} else {
			return -1, err
		}
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(reply)
	if err != nil {
		return resp.StatusCode, err
	}
	return resp.StatusCode, nil
}

func Http(method, url, mime string, body []byte) (*http.Response, error) {
	buf := bytes.NewBuffer(body)
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mime)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return resp, fmt.Errorf("(%s)%s", resp.Status, string(b))
	}
	return resp, nil
}
