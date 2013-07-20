package iom

import (
	"encoding/base64"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strings"
	"time"
)

var a float64 = 1.0
var b float64 = 0.0
var Inf float64 = a / b

const CharPos = "ABCDEFGHJKLMNPQRSTUVWXYZ0123456789" // A-Z without IO, 0-9
const SizePerPosition = len(CharPos)
const MaxSize = SizePerPosition * SizePerPosition

var ErrorOutOfRange = fmt.Errorf("count out of range")

func HashFromCount(count int64) (string, error) {
	first := count / int64(SizePerPosition)
	if first >= 24 {
		return "", ErrorOutOfRange
	}
	second := count - first*int64(SizePerPosition)
	return fmt.Sprintf("%c%c", CharPos[first], CharPos[second]), nil
}

type Iom struct {
	redis *redis.Pool
}

func NewIom(redis *redis.Pool) *Iom {
	return &Iom{
		redis: redis,
	}
}

func (h *Iom) Get(userid string, hash string) (string, error) {
	hash = strings.ToUpper(hash)
	conn := h.redis.Get()
	defer conn.Close()

	url64, err := redis.Bytes(conn.Do("GET", hashKey(userid, hash)))
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

func (h *Iom) FindByData(userid string, data string) (string, error) {
	conn := h.redis.Get()
	defer conn.Close()

	hash, err := redis.Bytes(conn.Do("GET", dataKey(userid, data)))
	if err == nil {
		err = h.Update(userid, string(hash))
	}
	return string(hash), err
}

func (h *Iom) Update(userid string, hash string) error {
	conn := h.redis.Get()
	defer conn.Close()
	_, err := conn.Do("ZADD", timeKey(userid), time.Now().UnixNano(), hash)
	return err
}

func (h *Iom) FindLatestHash(userid string) (string, error) {
	conn := h.redis.Get()
	defer conn.Close()
	reply, err := redis.Values(conn.Do("ZRANGEBYSCORE", timeKey(userid), "-inf", "+inf", "LIMIT", "0", "1"))
	if err != nil {
		return "", err
	}
	var ret string
	if _, err := redis.Scan(reply, &ret); err != nil {
		return "", err
	}
	return ret, nil
}

func (h *Iom) Create(userid string, data string) (string, error) {
	conn := h.redis.Get()
	defer conn.Close()
	data64 := base64.URLEncoding.EncodeToString([]byte(data))
	count, err := redis.Int64(conn.Do("ZCOUNT", timeKey(userid), -Inf, +Inf))
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
		conn.Do("DEL", dataKey(userid, data))
	} else {
		hash, err = HashFromCount(count)
		if err != nil {
			return "", nil
		}
	}
	if _, err = conn.Do("SET", hashKey(userid, hash), data64); err != nil {
		return "", err
	}
	if _, err = conn.Do("SET", dataKey(userid, data), hash); err != nil {
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
