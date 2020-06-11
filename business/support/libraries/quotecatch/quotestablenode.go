package quotecatch

import (
	"business/support/domain"
	"sync"
	"time"
)

type QuotesMap struct {
	Data map[string]*QuotesNode
}
type QuotesMapNode struct {
	sync.RWMutex
	Data map[string]*QuotesNode
}

type QuotesTableNode struct {
	QuotesNodes   QuotesMapNode
	FreshTimeSpan int64
	Timeout       int // 单位是天,几天的超时时间
}

func (this *QuotesMapNode) getQuotes(coinCode string) *QuotesNode {
	this.RLock()
	defer this.RUnlock()
	quotesNode, ok := this.Data[coinCode]
	if ok {
		return quotesNode
	}
	return nil
}

// freshTimeSpan
// 单位秒
// 如果最后一条数据的时间距离当前时间没有超过freshTimeSpan,则不需要从数据库加载新数据了

// timeout
// 单位是天,几天的超时时间,一旦缓存的数据的时间超过你timeout指定的天数
// 将缓存的数据删除掉
func NewQuotesTable(
	freshTimeSpan int64,
	timeout int) *QuotesTableNode {
	node := new(QuotesTableNode)
	node.QuotesNodes.Data = make(map[string]*QuotesNode)
	node.FreshTimeSpan = freshTimeSpan
	node.Timeout = timeout
	return node
}

func (this *QuotesTableNode) NewFreshTime(coinCode string) int64 {
	quotesNode := this.QuotesNodes.getQuotes(coinCode)
	if nil != quotesNode {
		return quotesNode.newFreshTime()
	}
	currentTime := time.Now().UTC().Unix()
	newFreshTime := currentTime - int64(this.Timeout)*24*60*60
	if this.Timeout < 1 {
		newFreshTime = 0
	}
	return newFreshTime
}

func (this *QuotesTableNode) QueryQuotes(
	coinCode string,
	startTime int64,
	endTime int64) []*domain.Quote {
	quotesNode := this.QuotesNodes.getQuotes(coinCode)
	if nil != quotesNode {
		quotes := quotesNode.getQuotes(startTime, endTime)
		if nil == quotes {
			return nil
		}
		return quotes
	} else {
		return nil
	}
}

func (this *QuotesTableNode) AddQuotes(
	coinCode string,
	quotes []*domain.Quote) {
	if len(coinCode) <= 0 ||
		nil == quotes ||
		len(quotes) <= 0 {
		return
	}
	this.QuotesNodes.Lock()
	if existQuotes, ok := this.QuotesNodes.Data[coinCode]; ok {
		this.QuotesNodes.Unlock()
		existQuotes.addQuotes(quotes)
		return
	}
	newNode := new(QuotesNode)
	newNode.quotes = quotes
	newNode.timeout = this.Timeout
	newNode.freshTimeSpan = this.FreshTimeSpan
	this.QuotesNodes.Data[coinCode] = newNode
	this.QuotesNodes.Unlock()
}
