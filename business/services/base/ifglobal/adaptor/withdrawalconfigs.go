package adaptor

import (
	"business/support/model/config"
)

//func getWithdrawalConfigInDB() ([]*domain.WithDrawalConfig, error) {
//	db := dbExPool.NewConn()
//	var configs []*domain.WithDrawalConfig
//	dbResult := db.Where(`status = ?`, domain.WITHDRAWAL_CONFIG_STATUS_ENABLE).Find(&configs)
//	if dbResult.RecordNotFound() {
//		return nil, errors.NotFound
//	}
//	if dbResult.Error != nil {
//		return nil, dbResult.Error
//	}
//	// loggers.Debug.Printf("getWithdrawalConfigInDB,configs:%#v", configs)
//	if len(configs) <= 0 {
//		return nil, errors.NotFound
//	}
//	return configs, nil
//}
//
//func GetWithdrawalConfigs() ([]*domain.WithDrawalConfig, error) {
//	configs := memcatch.WithDrawalConfigsCatch.GetConfigs()
//	if nil != configs {
//		return configs, nil
//	}
//	var err error
//	configs, err = getWithdrawalConfigInDB()
//	if nil != err {
//		return nil, err
//	}
//	if nil == configs {
//		return nil, errors.NotFound
//	}
//	memcatch.WithDrawalConfigsCatch.SetConfigs(configs)
//	return configs, nil
//}

func GetWithdrawalConfigs() ([]*config.WithDrawalConfig, error) {
	return config.GetAndCacheWithdrawConfigs()
}
