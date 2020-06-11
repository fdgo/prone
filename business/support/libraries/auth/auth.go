package auth

import (
	"business/support/dbpool"
	"business/support/libraries/loggers"
	"business/support/libraries/psql"
	"business/support/libraries/redis"
	"fmt"
	"math/rand"
	"time"
)

var (
	redisPool *redis.RedisPool
	dbPool    *psql.Client
)

func init() {
	var err error
	dbPool, err = dbpool.InitDB("ExchangeRead")
	if err != nil {
		panic(fmt.Sprintf("Auth init db:ExchangeRead error:%s", err.Error()))
	}

	redisPool, err = dbpool.InitRedis()
	if err != nil {
		panic(fmt.Sprintf("Auth init db:Redis error:%s", err.Error()))
	}

}

func Nonce() int64 {
	rand.Seed(time.Now().UnixNano())
	return 100000000000 + rand.Int63n(89999999999)
}

func IncorrectAssetPasswordIncr(uid int64) (int64, error) {
	key := fmt.Sprintf("AUTH.INCORRECT_ASSET_PASSWORD.%d", uid)
	cnt, err := redisPool.NewConn().Incr(key).Result()
	if redisPool.IsNil(err) {
		return 0, err
	}

	redisPool.NewConn().Expire(key, time.Hour*1)

	return cnt, err
}

func DisableAssetOP(uid int64) {
	key := fmt.Sprintf("AUTH.DISABLE_ASSET_OP.%d", uid)
	redisPool.NewConn().Set(key, time.Now().Unix(), time.Hour*24)
}

func IsDisableAssetOP(uid int64) bool {
	key := fmt.Sprintf("AUTH.DISABLE_ASSET_OP.%d", uid)
	_, err := redisPool.NewConn().Get(key).Int64()
	if redisPool.IsNil(err) {
		return false
	}
	if err != nil {
		loggers.Warn.Printf("Check is asset op disable error %s", err.Error())
		return true
	}
	return true
}

func GetDisableAssetOPTime(uid int64) int64 {
	key := fmt.Sprintf("AUTH.DISABLE_ASSET_OP.%d", uid)
	v, err := redisPool.NewConn().Get(key).Int64()
	if redisPool.IsNil(err) {
		return 0
	}

	return v
}
