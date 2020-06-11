package adaptor

import (
	"business/support/domain"
	"business/support/libraries/errors"
	"business/support/libraries/loggers"
	"sync"
)

type AddressCatch struct {
	sync.Mutex
	coinCode  string
	addresses []*domain.DepositAddress
}

func newAddressCatch(coinCode string) *AddressCatch {
	this := new(AddressCatch)
	this.coinCode = coinCode
	this.Lock()
	this.load()
	this.Unlock()
	return this
}

func (this *AddressCatch) load() error {
	var addresses []*domain.DepositAddress
	db := dbPool.NewConn()
	dbResult := db.Where("status = ?", domain.DEPOSIT_ADDRESS_FREE).
		Where("coin_code = ?", this.coinCode).
		Limit(1000).
		Find(&addresses)
	if nil != dbResult.Error {
		return dbResult.Error
	}
	if dbResult.RecordNotFound() {
		loggers.Warn.Printf("AddressCatch load unused list not found.coinCode:%s", this.coinCode)
		return errors.NotFound
	}
	this.addresses = addresses
	return nil
}

func (this *AddressCatch) QueryOne() *domain.DepositAddress {
	this.Lock()
	defer this.Unlock()
	size := len(this.addresses)
	if size < 1 {
		if err := this.load(); nil != err {
			return nil
		} else {
			size = len(this.addresses)
			if size < 1 {
				return nil
			}
		}
	}
	one := this.addresses[size-1]
	// 内存泄露,直到下次load时才会释放
	this.addresses = this.addresses[:size-1]
	return one
}

type AddressCatchMap struct {
	sync.Mutex
	addrMap map[string]*AddressCatch
}

func NewAddressCatchMap() *AddressCatchMap {
	this := new(AddressCatchMap)
	this.addrMap = make(map[string]*AddressCatch)
	this.Lock()
	this.load()
	this.Unlock()
	return this
}

func (this *AddressCatchMap) QueryOne(coinCode string) *domain.DepositAddress {
	this.Lock()
	catch, ok := this.addrMap[coinCode]
	if nil == catch || !ok {
		catch = newAddressCatch(coinCode)
		if nil != catch {
			this.addrMap[coinCode] = catch
		}
	}
	this.Unlock()
	if nil != catch {
		return catch.QueryOne()
	}
	return nil
}
func (this *AddressCatchMap) load() {
	type ACoinCode struct {
		CoinCode string
	}
	var coinCodes []ACoinCode
	db := dbPool.NewConn()
	dbResult := db.Table("deposit_addresss").
		Select("coin_code").
		Group("coin_code").
		Find(&coinCodes)
	if nil != dbResult.Error {
		loggers.Warn.Printf("AddressCatchMap load coinCodes failed:%#v", dbResult.Error)
		return
	}
	if dbResult.RecordNotFound() {
		loggers.Warn.Printf("AddressCatchMap load coinCodes not found.")
		return
	}
	for _, item := range coinCodes {
		catch := newAddressCatch(item.CoinCode)
		if nil != catch {
			this.addrMap[item.CoinCode] = catch
		}
	}
	return
}
