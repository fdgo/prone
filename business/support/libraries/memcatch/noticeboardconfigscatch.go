package memcatch

import (
	"business/support/domain"
	"sync"
	"time"
)

type NoticeboardConfigsMemCatch struct {
	sync.RWMutex
	expiration int64
	configs    []*domain.NoticeBoardConfig
}

var NoticeboardConfigsCatch *NoticeboardConfigsMemCatch

const noticeboardConfigCatchTimeOut int64 = 8 * 60

func InitNoticeboardConfigsCatch() {
	NoticeboardConfigsCatch = new(NoticeboardConfigsMemCatch)
}

func (this *NoticeboardConfigsMemCatch) GetConfigs() []*domain.NoticeBoardConfig {
	curTime := time.Now().UTC().Unix()
	this.RLock()
	defer this.RUnlock()
	if curTime > this.expiration {
		return nil
	} else {
		return this.configs
	}
}

func (this *NoticeboardConfigsMemCatch) SetConfigs(configs []*domain.NoticeBoardConfig) {
	if nil == configs {
		return
	}
	expiration := time.Now().UTC().Unix() + noticeboardConfigCatchTimeOut
	this.Lock()
	defer this.Unlock()
	this.configs = configs
	this.expiration = expiration
}
