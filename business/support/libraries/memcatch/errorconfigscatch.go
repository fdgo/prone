package memcatch

import (
	"business/support/domain"
	"sync"
	"time"
)

type LocalError struct {
	Local string
	Erros map[string]string
}

type ErrorConfigsMemCatch struct {
	sync.RWMutex
	expiration int64
	configs    map[string]*LocalError
}

// Data map[string]*QuotesNode
var ErrorConfigsCatch *ErrorConfigsMemCatch

const errorConfigsCatchCatchTimeOut int64 = 30 * 60

func InitErrorConfigsCatch() {
	ErrorConfigsCatch = new(ErrorConfigsMemCatch)
	ErrorConfigsCatch.configs = make(map[string]*LocalError)
}

func (this *ErrorConfigsMemCatch) GetConfigs(local string) (map[string]string, bool) {
	curTime := time.Now().UTC().Unix()
	this.RLock()
	defer this.RUnlock()
	if curTime > this.expiration {
		return nil, true
	} else {
		localErrors, ok := this.configs[local]
		if ok && nil != localErrors {
			return localErrors.Erros, false
		} else {
			return nil, false
		}
	}
}

func (this *ErrorConfigsMemCatch) SetConfigs(configs []*domain.ErrorConfig) {
	if nil == configs {
		return
	}
	configMap := make(map[string]*LocalError)
	for _, item := range configs {
		if localErrors, ok := configMap[item.Local]; ok {
			localErrors.Erros[item.ErrCode] = item.ErrMsg
		} else {
			localErrors := new(LocalError)
			localErrors.Erros = make(map[string]string)
			localErrors.Local = item.Local
			localErrors.Erros[item.ErrCode] = item.ErrMsg
			configMap[item.Local] = localErrors
		}
	}
	this.Lock()
	defer this.Unlock()
	this.configs = configMap
	this.expiration = time.Now().UTC().Unix() + errorConfigsCatchCatchTimeOut
}
