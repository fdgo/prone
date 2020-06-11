package adaptor

import (
	"business/support/domain"
	"business/support/libraries/errors"
	"business/support/libraries/loggers"
	"encoding/json"
	"fmt"
)

func AddRechargeAmount(recharge *domain.RechargeAmount) error {
	conn := redisPool.NewConn()
	if nil == conn {
		return errors.SysError
	}
	coinGroup, err := GetCoinGroup(recharge.CoinCode)
	if nil != err {
		return err
	}
	if len(coinGroup) > 0 {
		userAddress, err := GetUserBindAddrss(recharge.AccountId, coinGroup)
		if nil != err {
			return err
		}
		if len(userAddress) > 0 {
			recharge.Address = userAddress
		}
	}
	v, err := json.Marshal(recharge)
	if err != nil {
		return err
	}
	redisKey := domain.GetRechargeAmountRedisKey(recharge.CoinCode)
	accountId := fmt.Sprintf("%d", recharge.AccountId)
	if err := conn.HSet(redisKey.Key, accountId, v).Err(); err != nil {
		loggers.Error.Printf("AddRechargeAmount recharge:%#v to redis error:%s", recharge, err.Error())
		return err
	}
	return nil
}
