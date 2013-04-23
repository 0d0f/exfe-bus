package broker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func RestHttp(method, url, mime string, arg interface{}, reply interface{}) (int, error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(arg)
	if err != nil {
		return -1, err
	}
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return -1, err
	}
	req.Header.Set("Content-Type", mime)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, fmt.Errorf(resp.Status)
	}
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
		resp.Body.Close()
		return nil, fmt.Errorf(resp.Status)
	}
	return resp, nil
}
