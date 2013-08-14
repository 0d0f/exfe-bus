package qutimer

import (
	"github.com/garyburd/redigo/redis"
)

type Task struct {
	pool     *redis.Pool
	prefix   string
	key      string
	number   int
	b        []interface{}
	complete *redis.Script
	release  *redis.Script
}

func newTask(pool *redis.Pool, prefix, key string, number int, b []interface{}) *Task {
	return &Task{
		pool:     pool,
		prefix:   prefix,
		key:      key,
		number:   number,
		b:        b,
		complete: redis.NewScript(2, completeScript),
		release:  redis.NewScript(2, releaseScript),
	}
}

func (t *Task) Data() []interface{} {
	return t.b
}

const completeScript = `
local prefix = KEYS[1]
local key    = KEYS[2]
local number  = ARGV[1]
local timer = prefix..":timer"
local queue = prefix..":queue:"..key
local queue_overwrite = queue..":overwrite"
local queue_locker    = queue..":locker"
local queue_data      = queue..":data"

redis.call("DEL", queue_overwrite)
redis.call("LTRIM", queue, number, -1)
if redis.call("LLEN", queue) == 0 then
	redis.call("DEL", queue_data)
	redis.call("ZREM", timer, key)
end
redis.call("DEL", queue_locker)`

func (t *Task) Complete() error {
	conn := t.pool.Get()
	defer conn.Close()
	if _, err := t.complete.Do(conn, t.prefix, t.key, t.number); err != nil {
		return err
	}
	return nil
}

const releaseScript = `
local prefix = KEYS[1]
local key    = KEYS[2]
local ontime = ARGV[1]
local timer = prefix..":timer"
local queue = prefix..":queue:"..key
local queue_locker    = queue..":locker"

redis.call("DEL", queue_locker)
redis.call("ZADD", timer, ontime, key)`

func (t *Task) Release(ontime int64) error {
	conn := t.pool.Get()
	defer conn.Close()
	if _, err := t.release.Do(conn, t.prefix, t.key, ontime); err != nil {
		return err
	}
	return nil
}
