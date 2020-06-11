package domain

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// import "fmt"

type UserAssets struct {
	ID           int64           `json:"-" gorm:"primary_key"`
	AccountId    int64           `json:"account_id"`
	CoinCode     string          `json:"coin_code"`                                // 代币代码
	FreezeVol    decimal.Decimal `json:"freeze_vol" gorm:"type:decimal(36,18)"`    // 冻结资产
	AvailableVol decimal.Decimal `json:"available_vol" gorm:"type:decimal(36,18)"` // 可用资产
	CreatedAt    *time.Time      `json:"created_at"`
	UpdatedAt    *time.Time      `json:"updated_at"`
}

func (this *UserAssets) Desc() string {
	return fmt.Sprintf(`UserAssets:{id:%d,account_id:%d,coin_code:%s,
		freeze_vol:%s,available_vol:%s,created_at:%#v,updated_at:%#v}`,
		this.ID, this.AccountId, this.CoinCode, this.FreezeVol,
		this.AvailableVol, this.CreatedAt, this.UpdatedAt)
}

func (*UserAssets) TableName() string {
	return GetUserAssetsTableName()
}

func GetUserAssetsTableName() string {
	return "user_assets"
}

func (this *UserAssets) Copy() *UserAssets {
	if nil == this {
		return nil
	}
	other := *this
	return &other
}

type UserCoin struct {
	Code string
	Vol  decimal.Decimal
}

func (this *UserCoin) Copy() *UserCoin {
	if nil == this {
		return nil
	}
	item := *this
	return &item
}

func (this *UserCoin) String() string {
	return this.Desc()
}

func (this *UserCoin) Desc() string {
	return fmt.Sprintf(`UserCoin:{Code:%s,Vol:%s}`, this.Code, this.Vol)
}

func (this *UserCoin) IsValidCoin() bool {
	if nil == this {
		return false
	}
	if len(this.Code) <= 0 {
		return false
	}
	if this.Vol.Sign() <= 0 {
		return false
	}
	return true
}

//  id	long	否		主键，自增
//  userId	long	否 	 	用户ID
//  coinCode	text	否		代币代码
//  freezevol	long	是	0	数量
//  availableVol	long	否		可用数量
//  createTime	timestamp	否		添加时间
//  updateTime	timestamp	否		更新时间
