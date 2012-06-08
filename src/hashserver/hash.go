package main

import (
	"time"
	"strings"
	"fmt"
	"github.com/googollee/godis"
	"encoding/base64"
)

var a float64 = 1.0
var b float64 = 0.0
var Inf float64 = a / b
const CharPos = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // A-Z, 0-9
const SizePerPosition = len(CharPos)
const MaxSize = SizePerPosition * SizePerPosition
var ErrorOutOfRange = fmt.Errorf("count out of range")

func HashFromCount(count int64) (string, error) {
	first := count / int64(SizePerPosition)
	if first >= 26 {
		return "", ErrorOutOfRange
	}
	second := count - first * int64(SizePerPosition)
	return fmt.Sprintf("%c%c", CharPos[first], CharPos[second]), nil
}

type HashHandler struct {
	redis *godis.Client
}

func NewHashHandler(netaddr string, db int, password string) *HashHandler {
	return &HashHandler{
		redis: godis.New(netaddr, db, password),
	}
}

func (h *HashHandler) Get(userid string, hash string) (string, error) {
	hash = strings.ToUpper(hash)
	url64, err := h.redis.Get(hashKey(userid, hash))
	if err != nil {
		return "", err
	}
	err = h.Update(userid, hash)
	if err != nil {
		return "", err
	}
	url, err := base64.URLEncoding.DecodeString(string(url64))
	return string(url), err
}

func (h *HashHandler) FindByData(userid string, data string) (string, error) {
	hash, err := h.redis.Get(dataKey(userid, data))
	if err == nil {
		err = h.Update(userid, string(hash))
	}
	return string(hash), err
}

func (h *HashHandler) Update(userid string, hash string) error {
	_, err := h.redis.Zadd(timeKey(userid), time.Now().UnixNano(), hash)
	return err
}

func (h *HashHandler) FindLatestHash(userid string) (string, error) {
	reply, err := h.redis.Zrangebyscore(timeKey(userid), "-Inf", "+Inf", "LIMIT", "offset", "1")
	return string(reply.Elems[0].Elem), err
}

func (h *HashHandler) Create(userid string, data string) (string, error) {
	data64 := base64.URLEncoding.EncodeToString([]byte(data))
	count, err := h.redis.Zcount(timeKey(userid), -Inf, +Inf)
	if err != nil {
		return "", err
	}

	var hash string
	if count > int64(MaxSize) {
		hash, err = h.FindLatestHash(userid)
		if err != nil {
			return "", err
		}
		data, err = h.Get(userid, hash)
		if err != nil {
			return "", err
		}
		h.redis.Del(dataKey(userid, data))
	} else {
		hash, err = HashFromCount(count)
		if err != nil {
			return "", nil
		}
	}
	err = h.redis.Set(hashKey(userid, hash), data64)
	if err != nil {
		return "", err
	}
	err = h.redis.Set(dataKey(userid, data), hash)
	if err != nil {
		return "", err
	}
	err = h.Update(userid, hash)
	return hash, err
}

func hashKey(userid, hash string) string {
	return fmt.Sprintf("hash:%s:%s", userid, hash)
}

func dataKey(userid, data string) string {
	url64 := base64.URLEncoding.EncodeToString([]byte(data))
	return fmt.Sprintf("hash_url:%s:%s", userid, url64)
}

func timeKey(userid string) string {
	return fmt.Sprintf("hash_time:%s", userid)
}
