package rmodel

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

type BreadcrumbCacheSaver struct {
	r          *redis.Pool
	saveScript *redis.Script
	loadScript *redis.Script
}

func NewBreadcrumbCacheSaver(r *redis.Pool) *BreadcrumbCacheSaver {
	ret := &BreadcrumbCacheSaver{
		r: r,
	}
	ret.saveScript = redis.NewScript(1, `
		local user_id = KEYS[1]
		local data = ARGV[1]
		local now = ARGV[2]
		local userkey = "exfe:v3:routex:user_"..user_id
		local matchkey = "exfe:v3:routex:user_"..user_id..":cross"
		local crosses = redis.call("ZRANGEBYSCORE", matchkey, now, "+INF")
		redis.call("EXPIRE", userkey, 600)
		for i = 1, #crosses do
			local c = crosses[i]
			redis.call("SET", matchkey..":"..c, data)
		end
		return crosses
	`)
	ret.loadScript = redis.NewScript(0, `
		local cross_id = ARGV[1]
		local key_pattern = "exfe:v3:routex:*:cross:"..cross_id
		local keys = redis.call("KEYS", key_pattern)
		local ret = {}
		for i = 1, #keys do
			local key = keys[i]
			local user_id
			string.gsub(key, "exfe:v3:routex:user_(.-):cross:.*", function(p) user_id=p end)
			table.insert(ret, tonumber(user_id))
			table.insert(ret, redis.call("GET", key))
		end
		return ret
	`)
	return ret
}

func (s *BreadcrumbCacheSaver) ukey(userId int64) string {
	return fmt.Sprintf("exfe:v3:routex:user_%d", userId)
}

func (s *BreadcrumbCacheSaver) ckey(crossId, userId int64) string {
	return fmt.Sprintf("exfe:v3:routex:user_%d:cross:%d", userId, crossId)
}

func (s *BreadcrumbCacheSaver) EnableCross(userId, crossId int64, afterInSecond int) error {
	key, conn := s.ukey(userId)+":cross", s.r.Get()
	defer conn.Close()

	till := time.Now().Add(time.Duration(afterInSecond) * time.Second).Unix()
	if err := conn.Send("ZADD", key, till, crossId); err != nil {
		return err
	}
	if err := conn.Send("EXPIRE", key, 7200); err != nil {
		return err
	}
	if err := conn.Flush(); err != nil {
		return err
	}
	return nil
}

func (s *BreadcrumbCacheSaver) DisableCross(userId, crossId int64) error {
	key, conn := s.ukey(userId)+":cross", s.r.Get()
	defer conn.Close()

	if _, err := conn.Do("ZREM", key, crossId); err != nil {
		return err
	}
	return nil
}

func (s *BreadcrumbCacheSaver) LoadAllCross(crossId int64) (map[int64]SimpleLocation, error) {
	conn := s.r.Get()
	defer conn.Close()

	values, err := redis.Values(s.loadScript.Do(conn, crossId))
	if err != nil {
		return nil, err
	}
	ret := make(map[int64]SimpleLocation)
	for len(values) > 0 {
		var userId int64
		var value []byte
		if values, err = redis.Scan(values, &userId, &value); err != nil {
			continue
		}
		var l SimpleLocation
		if err := json.Unmarshal(value, &l); err != nil {
			continue
		}
		ret[userId] = l
	}
	return ret, nil
}

func (s *BreadcrumbCacheSaver) SaveCross(userId int64, l SimpleLocation) ([]int64, error) {
	b, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}
	conn := s.r.Get()
	defer conn.Close()

	reply, err := redis.Values(s.saveScript.Do(conn, userId, b, time.Now().Unix()))
	if err != nil {
		return nil, err
	}
	var ret []int64
	for len(reply) > 0 {
		var crossId int64
		reply, err = redis.Scan(reply, &crossId)
		if err != nil {
			return nil, err
		}
		ret = append(ret, crossId)
	}
	return ret, nil
}

func (s *BreadcrumbCacheSaver) LoadCross(userId, crossId int64) (SimpleLocation, bool, error) {
	key, conn := s.ckey(crossId, userId), s.r.Get()
	defer conn.Close()

	var ret SimpleLocation
	reply, err := redis.Bytes(conn.Do("GET", key))
	if err == redis.ErrNil {
		return ret, false, nil
	}
	if err != nil {
		return ret, false, err
	}
	if err := json.Unmarshal(reply, &ret); err != nil {
		return ret, false, err
	}
	return ret, true, nil
}

func (s *BreadcrumbCacheSaver) RemoveCross(userId, crossId int64) error {
	key, conn := s.ckey(crossId, userId), s.r.Get()
	defer conn.Close()

	if _, err := conn.Do("DEL", key); err != nil {
		return err
	}
	return nil
}

func (s *BreadcrumbCacheSaver) Save(userId int64, l SimpleLocation) error {
	key, conn := s.ukey(userId), s.r.Get()
	defer conn.Close()

	b, err := json.Marshal(l)
	if err != nil {
		return err
	}
	if _, err := conn.Do("SET", key, b, "EX", 600); err != nil {
		return err
	}
	return nil
}

func (s *BreadcrumbCacheSaver) Load(userId int64) (SimpleLocation, error) {
	key, conn := s.ukey(userId), s.r.Get()
	defer conn.Close()

	var ret SimpleLocation
	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return ret, err
	}
	if err := json.Unmarshal(reply, &ret); err != nil {
		return ret, err
	}
	return ret, nil
}
