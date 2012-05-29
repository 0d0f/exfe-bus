package main

import (
	"testing"
	"fmt"
)

func TestHashFromCount(t *testing.T) {
	type TestData struct {
		count int64
		hash string
	}

	datas := []TestData{
		{0, "00"},
		{9, "09"},
		{10, "0a"},
		{35, "0z"},
		{36, "10"},
		{36*36-1, "zz"},
	}

	for _, d := range datas {
		h := HashFromCount(d.count)
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

