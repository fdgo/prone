package memcatch

import (
	"business/support/domain"
	"sync"
	"time"
)

type CoinFeeConfigsMemCatch struct {
	sync.RWMutex
	expiration int64
	configs    []*domain.CoinFeeConfig
}

var CoinFeeConfigCatch *CoinFeeConfigsMemCatch

const coinFeeConfigCatchTimeOut int64 = 5 * 60

func InitCoinFeeConfigCatch() {
	CoinFeeConfigCatch = new(CoinFeeConfigsMemCatch)
}

func (this *CoinFeeConfigsMemCatch) GetConfigs() []*domain.CoinFeeConfig {
	curTime := time.Now().UTC().Unix()
	this.RLock()
	defer this.RUnlock()
	if curTime > this.expiration {
		return nil
	} else {
		return this.configs
	}
}

func (this *CoinFeeConfigsMemCatch) SetConfigs(configs []*domain.CoinFeeConfig) {
	if nil == configs {
		return
	}
	expiration := time.Now().UTC().Unix() + coinFeeConfigCatchTimeOut
	this.Lock()
	defer this.Unlock()
	this.configs = configs
	this.expiration = expiration
}
