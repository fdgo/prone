package adaptor

import (
	"business/support/domain"
	"business/support/libraries/errors"
	"business/support/libraries/memcatch"
)

func getCoinScansInDB() ([]*domain.CoinScan, error) {
	db := dbCfPool.NewConn()
	var coinScans []*domain.CoinScan
	dbResult := db.Where(`status = ?`, domain.COIN_SCAN_STATUS_ENABLE).
		Find(&coinScans)
	if dbResult.RecordNotFound() {
		return nil, errors.NotFound
	}
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	if len(coinScans) <= 0 {
		return nil, errors.NotFound
	}
	return coinScans, nil
}

func GetCoinScans() ([]*domain.CoinScan, error) {
	coinScans := memcatch.CoinScansCatch.GetCoinScans()
	if nil != coinScans {
		return coinScans, nil
	}
	var err error
	coinScans, err = getCoinScansInDB()
	if nil != err {
		return nil, err
	}
	if nil == coinScans {
		return nil, errors.NotFound
	}
	memcatch.CoinScansCatch.SetCoinScans(coinScans)
	return coinScans, nil
}
