package gobus

import (
	"fmt"
)

type Table map[string]map[string]string

func NewTable(route map[string]map[string]string) Table {
	return Table(route)
}

func (d Table) Find(url, identity string) (string, error) {
	urls, ok := d[url]
	if !ok {
		return "", fmt.Errorf("can't find %s in table", url)
	}
	ret, ok := urls[identity]
	if !ok {
		ret, ok = urls["_default"]
	}
	if !ok {
		return "", fmt.Errorf("can't find identity(%s) or default", identity)
	}

	return ret, nil
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

func (d *Dispatcher) DoWithIdentity(identity, addr, method string, arg, reply interface{}) error {
	url, err := d.table.Find(addr, identity)
	if err != nil {
		return err
	}
	return d.client.Do(url, method, arg, reply)
}

func (d *Dispatcher) Do(addr, method string, arg, reply interface{}) error {
	return d.DoWithIdentity("_default", addr, method, arg, reply)
}
