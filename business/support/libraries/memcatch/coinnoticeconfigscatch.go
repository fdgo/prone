package memcatch

import (
	"business/support/domain"
	"sync"
	"time"
)

type CoinNoticeConfigsMemCatch struct {
	sync.RWMutex
	expiration int64
	configs    map[string][]*domain.CoinNoticeConfig
}

var CoinNoticeConfigsCatch *CoinNoticeConfigsMemCatch

const CoinNoticeConfigCatchTimeOut int64 = 5 * 60

func InitCoinNoticeConfigCatch() {
	CoinNoticeConfigsCatch = &CoinNoticeConfigsMemCatch{
		configs: make(map[string][]*domain.CoinNoticeConfig),
	}
}

func (this *CoinNoticeConfigsMemCatch) GetConfigs(local string) []*domain.CoinNoticeConfig {
	curTime := time.Now().UTC().Unix()
	this.RLock()
	defer this.RUnlock()
	if curTime > this.expiration {
		this.configs = make(map[string][]*domain.CoinNoticeConfig)
		return nil
	} else if config, exist := this.configs[local]; exist {
		return config
	} else {
		return nil
	}
}

func (this *CoinNoticeConfigsMemCatch) SetConfigs(local string, configs []*domain.CoinNoticeConfig) {
	if nil == configs {
		return
	}
	expiration := time.Now().UTC().Unix() + CoinNoticeConfigCatchTimeOut
	this.Lock()
	defer this.Unlock()
	this.configs[local] = configs
	this.expiration = expiration
}
