package main

import (
	"testing"
	"strings"
	"fmt"
)

func TestHashFromCount(t *testing.T) {
	type TestData struct {
		count int64
		hash string
	}

	datas := []TestData{
		{0, "AA"},
		{25, "AZ"},
		{26, "A0"},
		{35, "A9"},
		{36, "BA"},
		{26*36-1, "Z9"},
	}

	for _, d := range datas {
		h, err := HashFromCount(d.count)
		if err != nil {
			t.Errorf("Error: %s", err)
		}
		if h != d.hash {
			t.Errorf("Count %d should be hash %s, but got %s", d.count, d.hash, h)
		}
	}
}

func TestHashCreate(t *testing.T) {
	handler := NewHashHandler("", 0, "")
	h, _ := handler.Create("123", "http://123/a")
	url, _ := handler.Get("123", h)
	if url != "http://123/a" {
		t.Errorf("User id %s's hash %s expect url %s, but got %s", "123", h, "http://123/a", url)
	}
	h1, _ := handler.FindByUrl("123", url)
	if h1 != h {
		t.Errorf("user id %s's url %s expect hash %s, but got %s", "123", url, h, h1)
	}

	uph := strings.ToUpper(h)
	upurl, _ := handler.Get("123", uph)
	if upurl != url {
		t.Errorf("hash handler should not care about case")
	}
	lowerh := strings.ToLower(h)
	lowerurl, _ := handler.Get("123", lowerh)
	if lowerurl != url {
		t.Errorf("hash handler should not care about case")
	}
}

func TestHashUpdate(t *testing.T) {
	handler := NewHashHandler("", 0, "")
	for _, userid := range []string{"234", "345"} {
		for _, crossid := range []string{"a", "b", "c", "d"} {
			_, _ = handler.Create(userid, fmt.Sprintf("http://%s/%s", userid, crossid))
		}
	}

	hash, _ := handler.FindLatestHash("234")
	url, _ := handler.Get("234", hash)
	if url != "http://234/a" {
		t.Errorf("User id %s last hash should get url %s, but got url %s", "234", "http://234/a", url)
	}
	hash, _ = handler.FindLatestHash("234")
	url, _ = handler.Get("234", hash)
	if url != "http://234/b" {
		t.Errorf("User id %s last hash should get url %s, but got url %s", "234", "http://234/b", url)
	}
	hash, _ = handler.FindLatestHash("234")
	url, _ = handler.Get("234", hash)
	if url != "http://234/c" {
		t.Errorf("User id %s last hash should get url %s, but got url %s", "234", "http://234/c", url)
	}
	hash, _ = handler.FindLatestHash("234")
	url, _ = handler.Get("234", hash)
	if url != "http://234/d" {
		t.Errorf("User id %s last hash should get url %s, but got url %s", "234", "http://234/d", url)
	}
}

