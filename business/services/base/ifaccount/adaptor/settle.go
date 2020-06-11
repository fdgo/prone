package adaptor

import (
	"business/support/domain"
	"business/support/libraries/errors"
	"business/support/libraries/loggers"
	"business/support/rpc"
	"context"
	"regexp"
	"time"

	"github.com/shopspring/decimal"
)

func GetUserBindAddrss(accountId int64, coinGroup string) (string, error) {
	db := dbPool.NewConn()
	var addr domain.UserDepositAddress
	dbResult := db.Select(`deposit_address`).
		Where("account_id = ? AND coin_group = ?", accountId, coinGroup).
		First(&addr)
	if dbResult.Error != nil {
		return "", dbResult.Error
	}
	if dbResult.RecordNotFound() {
		return "", nil
	}
	return addr.DepositAddress, nil
}

func getBindAddrss(address *domain.UserDepositAddress) (*domain.UserDepositAddress, error) {
	db := dbPool.NewConn()
	var addr domain.UserDepositAddress
	dbResult := db.Where("account_id = ? AND coin_group = ?", address.AccountId, address.CoinGroup).First(&addr)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	if dbResult.RecordNotFound() {
		return nil, nil
	}
	return &addr, nil
}

func BindAddress(accountId int64, coinCode string) (*domain.UserDepositAddress, error) {
	var address domain.UserDepositAddress
	address.AccountId = accountId
	coinGroup, err := GetCoinGroup(coinCode)
	if err != nil {
		return nil, errors.InvalidCoinCode
	}
	address.CoinGroup = coinGroup
	if addr, _ := getBindAddrss(&address); nil != addr {
		return addr, nil
	}
	var reply rpc.BindSettleAddressReply
	if err := accountRPC.BindSettleAddress(context.Background(), &address, &reply); err != nil {
		return nil, errors.InternalServerError
	}
	if nil != reply.Error {
		return nil, reply.Error
	} else if nil == reply.Address {
		return nil, errors.SysError
	}
	return reply.Address, nil
}

func Withdraw(settle *domain.Settle) error {
	var reply rpc.WithdrawReply
	if err := accountRPC.Witdraw(context.Background(), settle, &reply); err != nil {
		return errors.InternalServerError
	}
	if nil != reply.Error {
		loggers.Warn.Printf("Withdraw accountId:%d coin_code:%s error:%s", settle.AccountId, settle.CoinCode, reply.Error.Error())
		return reply.Error
	}
	return nil
}

func GetSettles(types []domain.SETTLE_TYPE, coin string, uid int64, limit, offset int) (*[]domain.Settle, error) {
	db := dbPool.NewConn()
	var settles []domain.Settle
	if types != nil {
		db = db.Where("type in (?)", types)
	}
	if coin != "" {
		db = db.Where("coin_code = ?", coin)
	}
	dbResult := db.Where("account_id = ?", uid).Order("updated_at desc").Limit(limit).Offset(offset).Find(&settles)
	if dbResult.Error != nil && !dbResult.RecordNotFound() {
		return nil, dbResult.Error
	}
	return &settles, nil
}

func GetTotalSettles(types []domain.SETTLE_TYPE, coin string, uid int64) (int64, error) {
	db := dbPool.NewConn()
	var total int64
	if types != nil {
		db = db.Where("type in (?)", types)
	}
	if coin != "" {
		db = db.Where("coin_code = ?", coin)
	}
	if err := db.Model(new(domain.Settle)).Where("account_id = ?", uid).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// 只返回充值,提现记录
func GetSettlesV2(typ domain.SETTLE_TYPE, coin string, uid int64, limit, offset int) (*[]domain.Settle, error) {
	db := dbPool.NewConn()
	var settles []domain.Settle
	if typ != 0 {
		db = db.Where("type = ?", typ)
	} else {
		db = db.Where("type = ? OR type = ?", domain.SETTLE_TYPE_DEPOSIT, domain.SETTLE_TYPE_WITHDRAW)
	}
	if coin != "" {
		db = db.Where("coin_code = ?", coin)
	}
	dbResult := db.Where("account_id = ?", uid).Order("updated_at desc").Limit(limit).Offset(offset).Find(&settles)
	if dbResult.Error != nil && !dbResult.RecordNotFound() {
		return nil, dbResult.Error
	}
	return &settles, nil
}

func GetTotalSettlesV2(typ domain.SETTLE_TYPE, coin string, uid int64) (int64, error) {
	db := dbPool.NewConn()
	var total int64
	if typ != 0 {
		db = db.Where("type = ?", typ)
	} else {
		db = db.Where("type = ? OR type = ?", domain.SETTLE_TYPE_DEPOSIT, domain.SETTLE_TYPE_WITHDRAW)
	}
	if coin != "" {
		db = db.Where("coin_code = ?", coin)
	}
	if err := db.Model(new(domain.Settle)).Where("account_id = ?", uid).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func GetCoinGroup(coinCode string) (string, error) {
	db := dbPool.NewConn()
	var coin domain.Coin
	dbResult := db.Select(`coin_group`).
		Where("name = ?", coinCode).
		First(&coin)
	if dbResult.RecordNotFound() {
		return "", errors.NotFound
	}
	if dbResult.Error != nil {
		return "", dbResult.Error
	}
	return coin.CoinGroup, nil
}

func CheckAddress(coinCode, address string) (string, error) {
	coinGroup, err := GetCoinGroup(coinCode)
	if err != nil {
		return "", errors.ParameterError
	}
	var (
		ok      bool
		fmtAddr domain.ETHADDRESS
	)

	switch coinGroup {
	case "ETH":
		fmtAddr, ok = domain.FormatETHAddress(address)
		if !ok {
			return "", errors.InvalidAddress
		}
		address = string(fmtAddr)
	case "BTC", "USDT":
		ok, err = regexp.MatchString("^[a-km-zA-HJ-NP-Z1-9]{25,34}$", address)
		if err != nil || !ok {
			return "", errors.InvalidAddress
		}
	case "EOS":
		ok, err = regexp.MatchString("^[1-5\\.a-z]{1,12}$", address)
		if err != nil || !ok {
			return "", errors.InvalidAddress
		}
	}
	return address, nil
}

func AddWithdrawAddress(addr *domain.UserWithdrawAddress) (*domain.UserWithdrawAddress, error) {
	var reply domain.UserWithdrawAddress
	if err := accountRPC.AddWithdrawAddress(context.Background(), addr, &reply); err != nil {
		if errors.Equal(err, errors.WithdrawAddressExisted) {
			return nil, errors.WithdrawAddressExisted
		}
		return nil, errors.InternalServerError
	}
	return &reply, nil
}

func DeleteWithdrawAddress(addr *domain.UserWithdrawAddress) error {
	var reply domain.UserWithdrawAddress
	if err := accountRPC.DeleteWithdrawAddress(context.Background(), addr, &reply); err != nil {
		if errors.Equal(err, errors.AddressNotFound) {
			return errors.AddressNotFound
		}
		return errors.InternalServerError
	}
	return nil
}

func UpdateWithdrawAddress(addr *domain.UserWithdrawAddress) (*domain.UserWithdrawAddress, error) {
	var reply domain.UserWithdrawAddress
	if err := accountRPC.UpdateWithdrawAddress(context.Background(), addr, &reply); err != nil {
		if errors.Equal(err, errors.AddressNotFound) {
			return nil, errors.AddressNotFound
		}
		return nil, errors.InternalServerError
	}
	return &reply, nil
}

func GetWithdrawAddress(addr *domain.UserWithdrawAddress) ([]domain.UserWithdrawAddress, error) {
	if addr.AccountId == 0 {
		return nil, errors.ParameterError
	}
	var withdrawAddresses []domain.UserWithdrawAddress
	db := dbPool.NewConn()
	if addr.CoinCode != "" {
		db = db.Where("coin_code = ?", addr.CoinCode)
	}
	if addr.Address != "" {
		db = db.Where("address = ?", addr.Address)
	}
	if addr.Type != 0 {
		db = db.Where("type = ?", addr.Type)
	}
	if err := db.Where("account_id = ?", addr.AccountId).Order("created_at DESC").Find(&withdrawAddresses).Error; err != nil {
		return nil, err
	}
	return withdrawAddresses, nil
}

func GetWithdrawTotalDay(coinCode string, accountId int64) (*decimal.Decimal, error) {
	var (
		settle domain.Settle
		now, _ = time.Parse("2006-01-02", time.Now().UTC().Format("2006-01-02"))
	)
	if err := dbPool.NewConn().
		Select("SUM(vol) as vol").
		Where("type = ?", domain.SETTLE_TYPE_WITHDRAW).
		Where("status != ?", domain.SETTLE_STATUS_REJECTED).
		Where("status != ?", domain.SETTLE_STATUS_FAILED).
		Where("account_id = ?", accountId).
		Where("coin_code = ?", coinCode).
		Where("created_at >= ?", now).Take(&settle).Error; err != nil {
		return nil, err
	}
	return &settle.Vol, nil

}

func GetWithdrawCountHour(accountId int64) (int64, error) {
	var (
		settle = new(domain.Settle)
		result = struct{ Count int64 }{}
		start  = time.Now().UTC().Add(time.Hour * -1)
	)
	if err := dbPool.NewConn().
		Table(settle.TableName()).
		Select("COUNT(*) as count").
		Where("type = ?", domain.SETTLE_TYPE_WITHDRAW).
		Where("status != ?", domain.SETTLE_STATUS_REJECTED).
		Where("status != ?", domain.SETTLE_STATUS_FAILED).
		Where("account_id = ?", accountId).
		Where("created_at >= ?", start).Take(&result).Error; err != nil {
		return 0, err
	}
	return result.Count, nil
}
