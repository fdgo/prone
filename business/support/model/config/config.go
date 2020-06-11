package config

import (
	"business/support/libraries/psql"
	"business/support/libraries/redis"
	"errors"
)

var (
	dbPool       *psql.Client
	redisPool    *redis.RedisPool
	exchangePool *psql.Client
)

func InitConfigDB(db *psql.Client) {
	dbPool = db
}

func InitRedisDB(db *redis.RedisPool) {
	redisPool = db
}

func InitExchangePool(db *psql.Client) {
	exchangePool = db
}

func AutoMigrate() error {
	if dbPool == nil {
		return errors.New("Config DB is nil")
	}

	return dbPool.NewConn().AutoMigrate(
		new(GlobalConfig),
		new(WithDrawalConfig),
		new(Area),
		new(SubArea),
	).Error
}
