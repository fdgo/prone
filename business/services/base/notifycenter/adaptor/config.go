package adaptor

import (
	"business/support/domain"
	"business/support/libraries/loggers"
	"fmt"
	"github.com/shopspring/decimal"
)

func GetDepositNotifyVolLimit(coinCode string) decimal.Decimal {
	var conf domain.CoinDepositNotifyVolLimit
	redis := redisPool.NewConn()
	key := fmt.Sprintf("%s.%s", domain.KDepositVolLimit.Key, coinCode)
	v, err := redis.Get(key).Int64()
	if err == nil {
		loggers.Debug.Printf("Get %s from redis error:%s", domain.KDepositVolLimit.Key)
		return decimal.New(v, 0)
	}
	if err := confDB.NewConn().Where("coin_code = ?", coinCode).First(&conf).Error; err != nil {
		loggers.Debug.Printf("Get %s deposit notify limit vol from db error:%s", coinCode, err.Error())
		return decimal.New(0, 0)
	}
	if err := redis.Set(key, conf.LimitVol, domain.KDepositVolLimit.Timeout).Err(); err != nil {
		loggers.Debug.Printf("Set %s to redis error:%s", domain.KDepositVolLimit.Key)
	}
	return conf.LimitVol
}
