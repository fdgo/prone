package memcatch

import (
	"sync"
	"time"
)

type catchItem struct {
	expiration int64
	obj        interface{}
}

type MemCatch struct {
	sync.RWMutex
	data map[string]*catchItem
}

func NewMemCatch() *MemCatch {
	catch := new(MemCatch)
	catch.data = make(map[string]*catchItem)
	return catch
}

func (this *MemCatch) Update(key string, obj interface{}, timeout int64) {
	if len(key) < 1 {
		return
	}
	if nil == obj {
		return
	}
	var expiration int64
	if timeout > 0 {
		expiration = time.Now().UTC().Unix() + timeout
	} else {
		expiration = 0
	}
	item := &(catchItem{expiration, obj})
	this.Lock()
	defer this.Unlock()
	this.data[key] = item
}

func (this *MemCatch) Query(key string) interface{} {
	if len(key) < 1 {
		return nil
	}
	this.RLock()
	defer this.RUnlock()
	if item, ok := this.data[key]; ok {
		return item
	}
	return nil
}
