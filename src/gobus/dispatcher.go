package gobus

import (
	"fmt"
	"regexp"
	"strings"
)

type urls struct {
	matches  map[*regexp.Regexp]string
	default_ string
}

type Table map[string]urls

func NewTable(route map[string]map[string]string) (Table, error) {
	ret := make(Table)
	for bus, v := range route {
		u := urls{
			matches:  make(map[*regexp.Regexp]string),
			default_: "",
		}
		for pattern, dst := range v {
			if pattern == "_default" {
				u.default_ = dst
				continue
			}
			reg, err := regexp.Compile(pattern)
			if err != nil {
				return nil, fmt.Errorf("%s is not valid regexp: %s", pattern, err)
			}
			u.matches[reg] = dst
		}
		ret[bus] = u
	}
	return ret, nil
}

func (d Table) Find(url, ticket string) (string, error) {
	urls, ok := d[url]
	path := ""
	if !ok {
		for prefix, u := range d {
			if lp, lu := len(prefix), len(url); lp <= lu && url[:lp] == prefix {
				ok = true
				urls = u
				path = url[len(prefix):]
				break
			}
		}
		if !ok {
			return "", fmt.Errorf("can't find %s in table", url)
		}
	}
	matched := ""
	for re, dst := range urls.matches {
		if re.MatchString(ticket) {
			matched = dst
		}
	}
	if matched == "" {
		matched = urls.default_
	}
	if matched == "" {
		return "", fmt.Errorf("can't match ticket(%s) or default", ticket)
	}

	return matched + path, nil
}

type Dispatcher struct {
	table  Table
	client *Client
}

func NewDispatcher(table Table) *Dispatcher {
	return &Dispatcher{
		table:  table,
		client: NewClient(new(JSON)),
	}
}

func (d *Dispatcher) DoWithTicket(ticket, addr, method string, arg, reply interface{}) error {
	if !strings.HasPrefix(addr, "http") {
		url, err := d.table.Find(addr, ticket)
		if err != nil {
			return err
		}
		addr = url
	}
	return d.client.Do(addr, method, arg, reply)
}

func (d *Dispatcher) Do(addr, method string, arg, reply interface{}) error {
	return d.DoWithTicket("_default", addr, method, arg, reply)
}
