package adaptor

import (
	"business/support/domain"
	"business/support/libraries/errors"
	"business/support/libraries/loggers"
	"business/support/libraries/memcatch"
)

func GetCoinRateBases() ([]*domain.CoinRateBase, error) {
	catchs := memcatch.CoinRateBaseCatch.GetBases()
	if nil != catchs || len(catchs) > 0 {
		return catchs, nil
	}
	db := dbExPool.NewConn()
	if nil == db {
		loggers.Warn.Printf("GetCoins new conn failed")
		return nil, errors.New("new db conn failed")
	}
	var bases []*domain.CoinRateBase
	dbResult := db.Where("status = ?", domain.COIN_RATE_BASE_STATUS_ENABLE).
		Find(&bases)
	if dbResult.RecordNotFound() {
		return nil, errors.NotFound
	}
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	if len(bases) <= 0 {
		return nil, errors.NotFound
	}
	memcatch.CoinRateBaseCatch.SetBases(bases)
	return bases, nil
}
