package adaptor

import (
	"business/support/dbpool"
	"business/support/libraries/psql"
	"business/support/libraries/redis"
	cfg "business/support/model/config"
	"time"

	"github.com/jinzhu/gorm"
)

var (
	dbActivityPool *psql.Client
	dbSpotPool     *psql.Client
	dbQuotePool    *psql.Client
	dbExPool       *psql.Client
	dbCfPool       *psql.Client
	redisPool      *redis.RedisPool
)

func init() {
	var err error
	dbActivityPool, err = dbpool.InitDB("ActivitiesRead")
	if err != nil {
		panic(err)
	}

	dbSpotPool, err = dbpool.InitDB("SpotsRead")
	if err != nil {
		panic(err)
	}

	dbQuotePool, err = dbpool.InitDB("QuotesRead")
	if err != nil {
		panic(err)
	}

	dbExPool, err = dbpool.InitDB("ExchangeRead")
	if err != nil {
		panic(err)
	}
	dbCfPool, err = dbpool.InitDB("ConfigRead")
	if err != nil {
		panic(err)
	}

	redisPool, err = dbpool.InitRedis()
	if err != nil {
		panic(err)
	}

	cfg.InitConfigDB(dbCfPool)
	cfg.InitRedisDB(redisPool)

	gorm.NowFunc = func() time.Time {
		return time.Now().UTC()
	}

	//memcatch.InitStocksCatch()
	//memcatch.InitSpotCoinsCatch()
	//memcatch.InitQuoteCoinsCatch()
	//memcatch.InitCoinFeeConfigCatch()
	//memcatch.InitCoinRateBaseCatch()
	//memcatch.InitUsdRatesCatch()
	////memcatch.InitWithDrawalConfigCatch()
	//memcatch.InitNoticeboardConfigsCatch()
	//memcatch.InitAppbuildsCatch()
	//memcatch.InitCoinScansCatch()
	//memcatch.InitUserLevelConfigCatch()
	//memcatch.InitAppSkinPkgCatch()
	//memcatch.InitPhoneCodesCatch()
	//memcatch.InitCoinNoticeConfigCatch()
}
