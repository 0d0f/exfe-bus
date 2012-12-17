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
		return "", fmt.Errorf("can't find table")
	}
	ret, ok := urls[identity]
	if !ok {
		ret, ok = urls["_default"]
	}
	if !ok {
		return "", fmt.Errorf("can't find identity or default")
	}

	return ret, nil
}

type Dispatcher struct {
	table Table
}

func NewDispatcher(table Table) *Dispatcher {
	return &Dispatcher{
		table: table,
	}
}

func (d *Dispatcher) DoWithIdentity(identity, addr, method string, arg, reply interface{}) error {
	url, err := d.table.Find(addr, identity)
	fmt.Println(url)
	if err != nil {
		return err
	}
	client, err := NewClient(url)
	if err != nil {
		return err
	}
	return client.Do(method, arg, reply)
}

func (d *Dispatcher) Do(addr, method string, arg, reply interface{}) error {
	return d.DoWithIdentity("_default", addr, method, arg, reply)
}
