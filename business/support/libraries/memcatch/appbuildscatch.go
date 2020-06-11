package memcatch

import (
	"business/support/domain"
	"sync"
	"time"
)

type AppbuildsMemCatch struct {
	sync.RWMutex
	expiration int64
	appbuilds  []*domain.AppBuild
}

var AppbuildsCatch *AppbuildsMemCatch

const appbuildCatchTimeOut int64 = 3 * 60

func InitAppbuildsCatch() {
	AppbuildsCatch = new(AppbuildsMemCatch)
}

func (this *AppbuildsMemCatch) GetAppBuilds() []*domain.AppBuild {
	curTime := time.Now().UTC().Unix()
	this.RLock()
	defer this.RUnlock()
	if curTime > this.expiration {
		return nil
	} else {
		return this.appbuilds
	}
}

func (this *AppbuildsMemCatch) SetAppBuilds(builds []*domain.AppBuild) {
	if nil == builds {
		return
	}
	expiration := time.Now().UTC().Unix() + appbuildCatchTimeOut
	this.Lock()
	defer this.Unlock()
	this.appbuilds = builds
	this.expiration = expiration
}
