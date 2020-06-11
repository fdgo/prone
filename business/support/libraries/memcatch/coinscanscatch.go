package memcatch

import (
	"business/support/domain"
	"sync"
	"time"
)

type CoinScansMemCatch struct {
	sync.RWMutex
	expiration int64
	coinScans  []*domain.CoinScan
}

var CoinScansCatch *CoinScansMemCatch

const coinScansCatchTimeOut int64 = 3 * 60

func InitCoinScansCatch() {
	CoinScansCatch = new(CoinScansMemCatch)
}

func (this *CoinScansMemCatch) GetCoinScans() []*domain.CoinScan {
	curTime := time.Now().UTC().Unix()
	this.RLock()
	defer this.RUnlock()
	if curTime > this.expiration {
		return nil
	} else {
		return this.coinScans
	}
}

func (this *CoinScansMemCatch) SetCoinScans(coinScans []*domain.CoinScan) {
	if nil == coinScans {
		return
	}
	expiration := time.Now().UTC().Unix() + coinScansCatchTimeOut
	this.Lock()
	defer this.Unlock()
	this.coinScans = coinScans
	this.expiration = expiration
}
