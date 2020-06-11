package quotecatch

import (
	"business/support/domain"
	"business/support/libraries/utils"
	"sync"
	"time"
)

// 带这mutex的对象,为了防止拷贝构造,
// 该对象必须以指针的方式传递
type QuotesNode struct {
	sync.RWMutex
	quotes             []*domain.Quote
	timeout            int // 单位是天,几天的超时时间
	freshTimeSpan      int64
	lastFreshTimestamp int64
}

func (this *QuotesNode) getQuotes(startTime int64, endTime int64) []*domain.Quote {
	if nil == this {
		return nil
	}
	this.RLock()
	defer this.RUnlock()
	if nil == this.quotes || len(this.quotes) < 1 {
		return nil
	}
	var quotes []*domain.Quote
	for _, item := range this.quotes {
		if item.Timestamp >= startTime && item.Timestamp <= endTime {
			quotes = append(quotes, item)
		}
	}
	return quotes
}

func (this *QuotesNode) addQuotes(quotes []*domain.Quote) {
	if nil == this || nil == quotes || len(quotes) <= 0 {
		return
	}
	this.Lock()
	defer this.Unlock()
	this.clear()
	if nil == this.quotes {
		this.quotes = quotes
		this.setLastFreshTimestamp()
		return
	}
	thisSize := len(this.quotes)
	if thisSize <= 0 {
		size := len(quotes)
		this.quotes = make([]*domain.Quote, size, size)
		copy(this.quotes, quotes)
		this.setLastFreshTimestamp()
		return
	}
	this.appendQuotes(quotes)
	this.setLastFreshTimestamp()
}

func (this *QuotesNode) appendQuotes(quotes []*domain.Quote) {
	thisSize := len(this.quotes)
	newSize := len(quotes)
	for i := 0; i < newSize; i++ {
		for j := thisSize - 1; j >= 0; j-- {
			if quotes[i].Timestamp > this.quotes[j].Timestamp {
				this.quotes = append(this.quotes, quotes[i:]...)
				return
			}
			if quotes[i].Timestamp == this.quotes[j].Timestamp {
				this.quotes[j] = quotes[i]
				break
			}
		}
	}
}

func (this *QuotesNode) newFreshTime() int64 {
	currentTime := time.Now().UTC().Unix()
	newFreshTime := currentTime - int64(this.timeout)*24*60*60
	if this.timeout < 1 {
		newFreshTime = 0
	}
	if this.lastFreshTimestamp < 1 {
		return newFreshTime
	}
	if (currentTime - this.lastFreshTimestamp) >= this.freshTimeSpan {
		dayDis := utils.SimpleDisDays(currentTime, this.lastFreshTimestamp)
		if this.timeout > 0 && dayDis > this.timeout {
			this.Lock()
			this.quotes = nil
			this.Unlock()
			return newFreshTime
		} else {
			return this.getLastNodeTimestamp()
		}
	}
	return -1
}

func (this *QuotesNode) clear() {
	// 删除超时的数据
	if this.timeout <= 0 {
		return
	}
	if nil == this.quotes || len(this.quotes) < 1 {
		return
	}
	currentTime := time.Now().UTC().Unix()
	for i, item := range this.quotes {
		dayDis := utils.SimpleDisDays(currentTime, item.Timestamp)
		if dayDis <= this.timeout {
			if i > 0 {
				this.quotes = this.quotes[i:]
			}
			return
		}
	}
	this.quotes = nil
}

func (this *QuotesNode) setLastFreshTimestamp() {
	thisSize := len(this.quotes)
	if thisSize < 1 {
		return
	}
	this.lastFreshTimestamp = time.Now().UTC().Unix()
}

func (this *QuotesNode) getLastNodeTimestamp() int64 {
	this.RLock()
	defer this.RUnlock()
	thisSize := len(this.quotes)
	if thisSize < 1 {
		return this.lastFreshTimestamp - 1
	} else {
		lt := this.quotes[thisSize-1].Timestamp
		return lt - 1
	}
}
