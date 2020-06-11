package adaptor

import (
	"business/support/domain"
	"business/support/libraries/errors"
	"time"

	"github.com/shopspring/decimal"
)

type RebateReward struct {
	CoinCode string          `json:"coin_code" gorm:"type:varchar(64);"`
	Vol      decimal.Decimal `json:"vol"`
}

func GetAccountRebateRewards(
	accountId int64) ([]*RebateReward, error) {
	currentTime := time.Now().UTC()
	var rewards []*RebateReward
	db := dbPool.NewConn()
	dbResult := db.Model(&domain.UserAssetReward{}).
		Select(`sum(reward) as vol,reward_coin as coin_code`).
		Where(`account_id = ? AND
			status = ? AND
			reward_type = ? AND
			(expiration_at IS NULL OR expiration_at >= ?)`,
			accountId,
			domain.REWARD_STATUS_CONFIRM,
			domain.USER_REWARD_TYPE_REBATE,
			&currentTime).Group(`reward_coin`).
		Scan(&rewards)
	if dbResult.RecordNotFound() {
		return nil, errors.NotFound
	}
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return rewards, nil
}

func GetUserAssetRewards(
	rtype domain.USER_REWARD_TYPE,
	coin string,
	accountId int64, limit, offset int) ([]*domain.UserAssetReward, error) {
	currentTime := time.Now().UTC()
	var rewards []*domain.UserAssetReward
	db := dbPool.NewConn()
	if rtype != 0 {
		db = db.Where("reward_type = ?", rtype)
	}
	if coin != "" {
		db = db.Where("reward_coin = ?", coin)
	}
	dbResult := db.Select(`reward,reward_type,reward_coin,status,reject_reason,created_at`).
		Where(`account_id = ? AND
			status != ? AND
			(expiration_at IS NULL OR expiration_at >= ?)`,
			accountId,
			domain.REWARD_STATUS_APPROVAL,
			&currentTime).
		Order("reward_id desc").
		Limit(limit).
		Offset(offset).
		Find(&rewards)
	if dbResult.RecordNotFound() {
		return nil, errors.NotFound
	}
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return rewards, nil
}

func GetAccountRebateUserCount(
	accountId int64) (int, error) {
	currentTime := time.Now().UTC()
	var count []int
	db := dbActivityPool.NewConn()
	dbResult := db.Model(&domain.UserAssetReward{}).
		Where(`account_id = ? AND
			reward_type = ? AND
			(expiration_at IS NULL OR expiration_at >= ?)`,
			accountId,
			domain.USER_REWARD_TYPE_REBATE,
			&currentTime).
		Pluck(`COUNT(DISTINCT(refer_id))`, &count)
	if dbResult.RecordNotFound() {
		return 0, errors.NotFound
	}
	if dbResult.Error != nil {
		return 0, dbResult.Error
	}
	if len(count) < 1 {
		return 0, nil
	}
	return count[0], nil
}
