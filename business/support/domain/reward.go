package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

//type REWARD_TYPE int32
//
//const (
//	REWARD_TYPE_UNKOWN   REWARD_TYPE = iota
//	REWARD_TYPE_SIGN                 // 注册
//	REWARD_TYPE_RELATION             // 师徒
//)

type REWARD_STATUS int32

const (
	REWARD_STATUS_UNKOWN   REWARD_STATUS = iota
	REWARD_STATUS_APPROVAL               // 审批中
	REWARD_STATUS_CONFIRM                // 确认
	REWARD_STATUS_REJECT                 // 审批不通过
	REWARD_STATUS_TIMEOUT                // 超时
)

//type Reward struct {
//	ID        int64           `json:"-" gorm:"primary_key"`
//	RewardId  int64           `json:"-" gorm:"type:int8;not null"`
//	AccountId int64           `json:"account_id" gorm:"type:int8;not null"`
//	PlanId    int64           `json:"plan_id" gorm:"type:int8;not null"`
//	Reward    decimal.Decimal `json:"reward" gorm:"type:decimal(36,18);not null"`
//	CoinCode  string          `json:"coin_code,omitempty" gorm:"size:64;"`
//	Type      REWARD_TYPE     `json:"type" gorm:"type:int4;DEFAULT:1; not null"`
//	Status    REWARD_STATUS   `json:"status" gorm:"type:int4;DEFAULT:1; not null"`
//	CreatedAt *time.Time      `json:"created_at,omitempty;not null"`
//	UpdatedAt *time.Time      `json:"-"`
//}
//
//func (this *Reward) TableName() string {
//	return "rewards"
//}
//
//func (this *Reward) IsValid() bool {
//	if nil == this {
//		return false
//	}
//	if this.RewardId < 1 ||
//		this.AccountId < 1 ||
//		this.Reward.Sign() <= 0 ||
//		len(this.CoinCode) < 1 ||
//		this.Type == REWARD_TYPE_UNKOWN {
//		return false
//	}
//	return true
//}

type USER_REWARD_TYPE int

const (
	USER_REWARD_TYPE_UNKOWN         USER_REWARD_TYPE = iota
	USER_REWARD_TYPE_AIRDROP                         // 空投
	USER_REWARD_TYPE_ACTIVITY_POINT                  // 活动积分
	USER_REWARD_TYPE_TRADE_POINT                     // 交易积分
	USER_REWARD_TYPE_DEPOSIT_POINT                   // 充值积分
	USER_REWARD_TYPE_REBATE                          // 返佣
	USER_REWARD_TYPE_BONUS                           // 赠送
	USER_REWARD_TYPE_QD                              // 渠道
	USER_REWARD_TYPE_CONTRACT_TASK                   // 合约任务
	USER_REWARD_TYPE_SPOT_TASK                       // 现货任务
)

type USER_REWARD_DISTRIBUTE_TYPE int

const (
	USER_REWARD_DISTRIBUTE_TYPE_UNKOWN USER_REWARD_DISTRIBUTE_TYPE = iota
	USER_REWARD_DISTRIBUTE_MANUL                                   // 需要管理员确认
	USER_REWARD_DISTRIBUTE_AUTO                                    // 自动
	USER_REWARD_DISTRIBUTE_BATCH                                   // 批处理
)

type UserReward struct {
	RewardId     int64            `json:"-" gorm:"primary_key;type:int8;not null"`
	AccountId    int64            `json:"account_id,omitempty" gorm:"type:int8;not null"`
	ActivityId   int64            `json:"activity_id,omitempty" gorm:"type:int8;not null"`
	ReferId      string           `json:"refer_id,omitempty" gorm:"size:64;"` // 在返佣奖励中是关联的用户ID
	OrigId       string           `json:"orig_id,omitempty" gorm:"size:64;"`  // 在返佣奖励中是交易ID
	Reward       decimal.Decimal  `json:"reward,omitempty" gorm:"type:decimal(36,18);not null"`
	RewardType   USER_REWARD_TYPE `json:"reward_type,omitempty" gorm:"type:int4;DEFAULT:1; not null"`
	Status       REWARD_STATUS    `json:"status,omitempty" gorm:"type:int4;DEFAULT:1; not null"`
	Auditor      string           `json:"auditor,omitempty" gorm:"size:255"`
	RejectReason string           `json:"reject_reason,omitempty" gorm:"size:255"`
	CreatedAt    *time.Time       `json:"created_at,omitempty" gorm:"not null"`
	ConfirmedAt  *time.Time       `json:"-"`
	ExpirationAt *time.Time       `json:"-"`
}

func (this *UserReward) IsValid() bool {
	if nil == this {
		return false
	}
	if this.RewardId < 1 ||
		this.AccountId < 1 ||
		this.Reward.Sign() <= 0 ||
		this.RewardType == USER_REWARD_TYPE_UNKOWN {
		return false
	}
	return true
}

////////////////////////////

type UserPointReward struct {
	UserReward
}

func (this *UserPointReward) TableName() string {
	return "user_point_rewards"
}

type UserAssetReward struct {
	UserReward
	RewardCoin string `json:"reward_coin,omitempty" gorm:"size:64"`
}

func (this *UserAssetReward) TableName() string {
	return "user_asset_rewards"
}
