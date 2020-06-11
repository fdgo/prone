package dbpool

import (
	"business/support/config"
	"business/support/libraries/loggers"
	"business/support/libraries/psql"
	"business/support/libraries/redis"
	"errors"
	"sync"
)

var (
	lock      sync.Mutex
	redisPool *redis.RedisPool
	dbs       = make(map[string]*psql.Client)
)

func InitDB(name string) (*psql.Client, error) {
	lock.Lock()
	defer lock.Unlock()
	if client, ok := dbs[name]; ok {
		return client, nil
	}
	dbConf, ok := config.Conf.DBS[name]
	if !ok {
		return nil, errors.New("DB Config not found")
	}
	pool := psql.NewClient(&dbConf, config.IsDebug)
	dbs[name] = pool
	loggers.Info.Printf("dbpool init db:%s OK", name)

	return pool, nil
}

func InitRedis() (*redis.RedisPool, error) {
	lock.Lock()
	defer lock.Unlock()
	if redisPool != nil {
		return redisPool, nil
	}
	dbConf, ok := config.Conf.DBS["Redis"]
	if !ok {
		return nil, errors.New("DB Config not found")
	}
	redisPool = redis.NewPool(&dbConf)
	return redisPool, nil
}
