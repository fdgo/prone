package adaptor

import (
	"business/support/config"
	"business/support/domain"
	"business/support/libraries/loggers"
	"business/support/libraries/psql"
	"business/support/libraries/redis"
	"errors"
)

var (
	confDB    *psql.Client
	exDB      *psql.Client
	recordDB  *psql.Client
	redisPool *redis.RedisPool
)

func init() {
	dbConf, ok := config.Conf.DBS["ConfigRead"]
	if !ok {
		loggers.Error.Printf("adaptor init error no db config")
		panic(errors.New("no db config"))
	}
	confDB = psql.NewClient(&dbConf, config.IsDebug)

	dbConf, ok = config.Conf.DBS["ExchangeRead"]
	if !ok {
		loggers.Error.Printf("adaptor init error no db config")
		panic(errors.New("no db config"))
	}
	exDB = psql.NewClient(&dbConf, config.IsDebug)

	dbConf, ok = config.Conf.DBS["RecordsRead"]
	if !ok {
		loggers.Error.Printf("adaptor init error no db config")
		panic(errors.New("no db config"))
	}
	recordDB = psql.NewClient(&dbConf, config.IsDebug)

	if err := confDB.NewConn().AutoMigrate(new(domain.CoinDepositNotifyVolLimit), new(domain.SMSTemplate), new(domain.EmailTemplate)).Error; err != nil {
		loggers.Error.Printf("AutoMigrate error:%s", err.Error())
		panic(errors.New("AutoMigrate failed"))
	}

	redisConf, ok := config.Conf.DBS["Redis"]
	if !ok {
		loggers.Error.Printf("adaptor init error no db config")
		panic(errors.New("no db config"))
	}
	redisPool = redis.NewPool(&redisConf)
}
