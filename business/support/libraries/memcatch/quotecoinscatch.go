package memcatch

import (
	"business/support/domain"
	"sync"
	"time"
)

type QuoteCoinsMemCatch struct {
	sync.RWMutex
	expiration int64
	coins      []*domain.Coin
}

var QuoteCoinsCatch *QuoteCoinsMemCatch

const quoteCoinsCatchTimeOut int64 = 3 * 60

func InitQuoteCoinsCatch() {
	QuoteCoinsCatch = new(QuoteCoinsMemCatch)
}

func (this *QuoteCoinsMemCatch) GetCoins() []*domain.Coin {
	curTime := time.Now().UTC().Unix()
	this.RLock()
	defer this.RUnlock()
	if curTime > this.expiration {
		return nil
	} else {
		return this.coins
	}
}

func (this *QuoteCoinsMemCatch) SetCoins(coins []*domain.Coin) {
	if nil == coins {
		return
	}
	expiration := time.Now().UTC().Unix() + quoteCoinsCatchTimeOut
	this.Lock()
	defer this.Unlock()
	this.coins = coins
	this.expiration = expiration
}
