package adaptor

import (
	"business/support/config"
	"business/support/domain"
	"business/support/libraries/loggers"
	"business/support/libraries/message"
	"business/support/libraries/mq"
	"business/support/libraries/psql"
	"business/support/libraries/redis"
	"business/support/libraries/sequence"
	"errors"
	"time"
)

var (
	dbPool         *psql.Client
	dbActivityPool *psql.Client
	dbCfPool       *psql.Client
	redisPool      *redis.RedisPool
	processor      *mq.Processor
	mqConn         *mq.Conn
	putTimeout     time.Duration
	addrCatchs     *AddressCatchMap
	snsClient      *message.SNSClient
)

func initConfigDB() {
	dbConf, ok := config.Conf.DBS["ConfigRead"]
	if !ok {
		loggers.Warn.Printf("Load Config database failed \n")
		panic(44)
	}
	dbCfPool = psql.NewClient(&dbConf, config.IsDebug)
	if nil == dbCfPool {
		loggers.Warn.Printf("Load Config database new client failed \n")
		panic(55)
	}
}

func initActivityDB() {
	dbConf, ok := config.Conf.DBS["ActivitiesRead"]
	if !ok {
		loggers.Warn.Printf("Load ActivitiesRead database failed \n")
		panic(44)
	}
	dbActivityPool = psql.NewClient(&dbConf, config.IsDebug)
	if nil == dbActivityPool {
		loggers.Warn.Printf("Load ActivitiesRead database new client failed \n")
		panic(55)
	}
}

func init() {
	var err error
	dbConf, ok := config.Conf.DBS["ExchangeWrite"]
	if !ok {
		loggers.Error.Printf("adaptor init error no db config")
		panic(errors.New("no db config"))
	}
	dbPool = psql.NewClient(&dbConf, config.IsDebug)

	if err := dbPool.NewConn().AutoMigrate(
		&domain.UserWithdrawAddress{},
		&domain.Account{},
	).Error; err != nil {
		panic(err)
	}

	redisConf, ok := config.Conf.DBS["Redis"]
	if !ok {
		loggers.Error.Printf("adaptor init error no db config")
		panic(errors.New("no db config"))
	}

	if snsClient, err = message.NewSNSClient(config.Conf.ActionSNSConf); err != nil {
		loggers.Error.Panic("new record sender error:%s", err.Error())
	}

	redisPool = redis.NewPool(&redisConf)
	initConfigDB()
	initActivityDB()
	sequence.Init(redisPool)
	//memcatch.InitWithDrawalConfigCatch()
	//memcatch.InitUserLevelConfigCatch()
	addrCatchs = NewAddressCatchMap()
}
