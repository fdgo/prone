package memcatch

import (
	"business/support/domain"
	"sync"
	"time"
)

type PhoneCodesMemCatch struct {
	sync.RWMutex
	expiration int64
	codes      *domain.PhoneCode
}

var PhoneCodesCatch *PhoneCodesMemCatch

const phoneCodesCatchTimeOut int64 = 3 * 60

func InitPhoneCodesCatch() {
	PhoneCodesCatch = new(PhoneCodesMemCatch)
}

func (this *PhoneCodesMemCatch) GetPhoneCode() *domain.PhoneCode {
	curTime := time.Now().UTC().Unix()
	this.RLock()
	defer this.RUnlock()
	if curTime > this.expiration {
		return nil
	} else {
		return this.codes
	}
}

func (this *PhoneCodesMemCatch) SetPhoneCode(codes *domain.PhoneCode) {
	if nil == codes {
		return
	}
	expiration := time.Now().UTC().Unix() + phoneCodesCatchTimeOut
	this.Lock()
	defer this.Unlock()
	this.codes = codes
	this.expiration = expiration
}
