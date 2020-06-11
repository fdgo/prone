package memcatch

import (
	"business/support/domain"
	"sync"
	"time"
)

type AppSkinPkgMemCatch struct {
	sync.RWMutex
	expiration int64
	skinPkg    *domain.AppSkinPkg
}

var AppSkinPkgCatch *AppSkinPkgMemCatch

const appSkinPkgCatchTimeOut int64 = 3 * 60

func InitAppSkinPkgCatch() {
	AppSkinPkgCatch = new(AppSkinPkgMemCatch)
}

func (this *AppSkinPkgMemCatch) GetAppSkinPkg() *domain.AppSkinPkg {
	curTime := time.Now().UTC().Unix()
	this.RLock()
	defer this.RUnlock()
	if curTime > this.expiration {
		return nil
	} else {
		return this.skinPkg
	}
}

func (this *AppSkinPkgMemCatch) SetAppSkinPkg(skinPkg *domain.AppSkinPkg) {
	if nil == skinPkg {
		return
	}
	expiration := time.Now().UTC().Unix() + appSkinPkgCatchTimeOut
	this.Lock()
	defer this.Unlock()
	this.skinPkg = skinPkg
	this.expiration = expiration
}
