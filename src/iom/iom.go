package iom

import (
	"encoding/base64"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/googollee/go-rest"
	"net/http"
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
	rest.Service `prefix:"/iom"`

	get    rest.SimpleNode `route:"/:user_id/:hash" method:"GET"`
	create rest.SimpleNode `route:"/user/:user_id" method:"POST"`

	redis  *redis.Pool
	prefix string
}

func NewIom(redis *redis.Pool, prefix string) *Iom {
	return &Iom{
		redis:  redis,
		prefix: prefix,
	}
}

func (h *Iom) Get(ctx rest.Context) {
	var userID string
	var hash string
	ctx.Bind("user_id", &userID)
	ctx.Bind("hash", &hash)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, "%s", err)
		return
	}
	ret, err := h.grab(userID, hash)
	if err != nil {
		ctx.Return(http.StatusInternalServerError, "%s", err)
		return
	}
	ctx.Render(ret)
}

func (h *Iom) Create(ctx rest.Context, data string) {
	var userID string
	ctx.Bind("user_id", &userID)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, "%s", err)
		return
	}
	ret, err := h.findByData(userID, data)
	if err != nil {
		ret, err = h.make(userID, data)
	}
	if err != nil {
		ctx.Return(http.StatusInternalServerError, "%s", err)
		return
	}
	ctx.Render(ret)
}

func (h *Iom) grab(userID string, hash string) (string, error) {
	hash = strings.ToUpper(hash)
	conn := h.redis.Get()
	defer conn.Close()

	url64, err := redis.Bytes(conn.Do("GET", h.hashKey(userID, hash)))
	if err != nil {
		return "", err
	}
	err = h.update(userID, hash)
	if err != nil {
		return "", err
	}
	url, err := base64.URLEncoding.DecodeString(string(url64))
	return string(url), err
}

func (h *Iom) findByData(userID string, data string) (string, error) {
	conn := h.redis.Get()
	defer conn.Close()

	hash, err := redis.Bytes(conn.Do("GET", h.dataKey(userID, data)))
	if err == nil {
		err = h.update(userID, string(hash))
	}
	return string(hash), err
}

func (h *Iom) update(userID string, hash string) error {
	conn := h.redis.Get()
	defer conn.Close()
	_, err := conn.Do("ZADD", h.timeKey(userID), time.Now().UnixNano(), hash)
	return err
}

func (h *Iom) findLatestHash(userID string) (string, error) {
	conn := h.redis.Get()
	defer conn.Close()
	reply, err := redis.Values(conn.Do("ZRANGEBYSCORE", h.timeKey(userID), "-inf", "+inf", "LIMIT", "0", "1"))
	if err != nil {
		return "", err
	}
	var ret string
	if _, err := redis.Scan(reply, &ret); err != nil {
		return "", err
	}
	return ret, nil
}

func (h *Iom) make(userID string, data string) (string, error) {
	conn := h.redis.Get()
	defer conn.Close()
	data64 := base64.URLEncoding.EncodeToString([]byte(data))
	count, err := redis.Int64(conn.Do("ZCOUNT", h.timeKey(userID), -Inf, +Inf))
	if err != nil {
		return "", err
	}

	var hash string
	if count > int64(MaxSize) {
		hash, err = h.findLatestHash(userID)
		if err != nil {
			return "", err
		}
		data, err = h.grab(userID, hash)
		if err != nil {
			return "", err
		}
		conn.Do("DEL", h.dataKey(userID, data))
	} else {
		hash, err = HashFromCount(count)
		if err != nil {
			return "", nil
		}
	}
	if _, err = conn.Do("SET", h.hashKey(userID, hash), data64); err != nil {
		return "", err
	}
	if _, err = conn.Do("SET", h.dataKey(userID, data), hash); err != nil {
		return "", err
	}
	err = h.update(userID, hash)
	return hash, err
}

func (h *Iom) hashKey(userID, hash string) string {
	return fmt.Sprintf("%s:hash:%s:%s", h.prefix, userID, hash)
}

func (h *Iom) dataKey(userID, data string) string {
	url64 := base64.URLEncoding.EncodeToString([]byte(data))
	return fmt.Sprintf("%s:hash_url:%s:%s", h.prefix, userID, url64)
}

func (h *Iom) timeKey(userID string) string {
	return fmt.Sprintf("%s:hash_time:%s", h.prefix, userID)
}
