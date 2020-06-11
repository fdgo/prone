package sequence

import (
	"business/support/libraries/redis"
	"errors"
)

var (
	pool *redis.RedisPool
	key  string
)

func init() {
	pool = nil
	key = "BASE.SEQUENCE.CURRENT"
}

func Init(redisPool *redis.RedisPool) {
	pool = redisPool
}

func GenSequence() (int64, error) {
	if pool == nil {
		return 0, errors.New("RedisPool Nil")
	}
	conn := pool.NewConn()
	return conn.Incr(key).Result()
}
