package errorconfig

import (
	"business/support/dbpool"
	"business/support/domain"
	"business/support/libraries/psql"
	"fmt"
)

var dbPool *psql.Client

func init() {
	var err error
	dbPool, err = dbpool.InitDB("ExchangeRead")
	if err != nil {
		panic(fmt.Sprintf("ErrorConfig init db error:%s", err.Error()))
	}
	//memcatch.InitErrorConfigsCatch()
}

func loadAllErrorConfigs() []*domain.ErrorConfig {
	db := dbPool.NewConn()
	if nil == db {
		return nil
	}
	var errConfigs []*domain.ErrorConfig
	dbResult := db.Find(&errConfigs)
	if dbResult.RecordNotFound() {
		return nil
	}
	if dbResult.Error != nil {
		return nil
	}
	return errConfigs
}

func QueryErrMsg(local string, errCode string) (string, bool) {
	//localErrors, isNeedReload := memcatch.ErrorConfigsCatch.GetConfigs(local)
	//if nil == localErrors && isNeedReload {
	//	if dbLocalErrors := loadAllErrorConfigs(); nil != dbLocalErrors {
	//		memcatch.ErrorConfigsCatch.SetConfigs(dbLocalErrors)
	//		localErrors, _ = memcatch.ErrorConfigsCatch.GetConfigs(local)
	//	}
	//}
	//if nil == localErrors {
	//	return "", false
	//} else {
	//	if errMsg, ok := localErrors[errCode]; ok {
	//		return errMsg, true
	//	}
	//	return "", false
	//}
	return "", true
}
