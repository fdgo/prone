package memcatch

import (
	"business/support/domain"
	"sync"
	"time"
)

type CoinRateBaseMemCatch struct {
	sync.RWMutex
	expiration int64
	bases      []*domain.CoinRateBase
}

var CoinRateBaseCatch *CoinRateBaseMemCatch

const coinRateBaseCatchTimeOut int64 = 5 * 60

func InitCoinRateBaseCatch() {
	CoinRateBaseCatch = new(CoinRateBaseMemCatch)
}

func (this *CoinRateBaseMemCatch) GetBases() []*domain.CoinRateBase {
	curTime := time.Now().UTC().Unix()
	this.RLock()
	defer this.RUnlock()
	if curTime > this.expiration {
		return nil
	} else {
		return this.bases
	}
}

func (this *CoinRateBaseMemCatch) SetBases(bases []*domain.CoinRateBase) {
	if nil == bases {
		return
	}
	expiration := time.Now().UTC().Unix() + coinFeeConfigCatchTimeOut
	this.Lock()
	defer this.Unlock()
	this.bases = bases
	this.expiration = expiration
}
