package gobus

import (
	"fmt"
)

type Dispatcher map[string]map[string]string

func NewDispatcher(route map[string]map[string]string) Dispatcher {
	ret := Dispatcher(route)
	updateKeys := make([]string, 0)
	for k, _ := range ret {
		if l := len(k); l > 0 && k[l-1] != '/' {
			updateKeys = append(updateKeys, k)
		}
	}
	for _, k := range updateKeys {
		ret[k+"/"] = ret[k]
		delete(ret, k)
	}
	return ret
}

func (d Dispatcher) Find(url, identity string) (string, error) {
	if l := len(url); l > 0 && url[l-1] != '/' {
		url = url + "/"
	}
	urls, ok := d[url]
	if !ok {
		return "", fmt.Errorf("can't find dispatcher")
	}
	ret, ok := urls[identity]
	if !ok {
		ret, ok = urls["_default"]
	}
	if !ok {
		return "", fmt.Errorf("can't find identity")
	}
	if l := len(ret); l > 0 && ret[l-1] != '/' {
		ret = ret + "/"
	}

	return ret, nil
}
