package adaptor

import (
	"business/support/config"
	"business/support/domain"
	"business/support/libraries/loggers"
	"business/support/libraries/psql"
	"errors"
)

var (
	dbPool *psql.Client
)

func init() {
	if dbConf, ok := config.Conf.DBS["RecordsWrite"]; ok {
		dbPool = psql.NewClient(&dbConf, config.IsDebug)
		if dbPool == nil {
			panic(errors.New("RecordWrite db new client failed!"))
		}
	} else {
		loggers.Error.Printf("adaptor init RecordWrite db error no db config")
		panic(errors.New("no db config"))
	}

	if err := dbPool.NewConn().AutoMigrate(
		&domain.AssetRecord{},
		&domain.OrderRecord{},
		&domain.AccountRecord{},
		&domain.SettleRecord{}).Error; err != nil {
		loggers.Error.Printf("adaptor auto migrate domain error:%s", err.Error())
		panic(err)
	}

}
