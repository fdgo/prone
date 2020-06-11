package adaptor

import (
	"business/support/config"
	"business/support/dbpool"
	emailSender "business/support/libraries/email"
	"business/support/libraries/loggers"
	"business/support/libraries/message"
	"business/support/libraries/psql"
	"business/support/libraries/redis"
	cfg "business/support/model/config"
	"business/support/rpc"
	"time"
)

var (
	dbPool         *psql.Client
	dbActivityPool *psql.Client
	dbRecordsPool  *psql.Client
	confDBPool     *psql.Client
	redisPool      *redis.RedisPool
	accountRPC     *rpc.AccountRPC
	snsClient      *message.SNSClient
	sqsClient      *message.SQSClient
	localTime      *time.Location
)

func init() {
	var err error
	dbPool, err = dbpool.InitDB("ExchangeRead")
	if err != nil {
		panic(err)
	}
	dbActivityPool, err = dbpool.InitDB("ActivitiesRead")
	if err != nil {
		panic(err)
	}
	dbRecordsPool, err = dbpool.InitDB("RecordsRead")
	if err != nil {
		panic(err)
	}
	confDBPool, err = dbpool.InitDB("ConfigRead")
	if err != nil {
		panic(err)
	}
	redisPool, err = dbpool.InitRedis()
	if err != nil {
		panic(err)
	}

	cfg.InitConfigDB(confDBPool)
	cfg.InitRedisDB(redisPool)

	snsClient, err = message.NewSNSClient(config.Conf.ActionSNSConf)
	if err != nil {
		loggers.Error.Panicf("new sns client error:%s", err.Error())
	}
	sqsClient, err = message.NewSQSClient(config.Conf.NotifySQSConf)
	if err != nil {
		loggers.Error.Panicf("new sns client error:%s", err.Error())
	}
	localTime, err = time.LoadLocation("Asia/Chongqing")
	if err != nil {
		loggers.Error.Panicf("load local time error:%s", err.Error())
	}

	emailSender.Init(redisPool)
	accountRPC = rpc.NewAccountRPC()
}

func AccountRPC() *rpc.AccountRPC {
	return accountRPC
}
