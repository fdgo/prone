package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type COIN_FEE_CONFIG_STATUS int

const (
	COIN_FEE_CONFIG_STATUS_UNKOWN  COIN_FEE_CONFIG_STATUS = iota
	COIN_FEE_CONFIG_STATUS_ENABLE                         // 可以用
	COIN_FEE_CONFIG_STATUS_DISABLE                        // 禁用中
)

type COIN_FEE_CONFIG_TYPE int

const (
	COIN_FEE_TYPE_CONFIG_UNKOWN COIN_FEE_CONFIG_TYPE = iota // value --> 0
	COIN_FEE_TYPE_CONFIG_BUY                                // 只收买的       // value --> 1
	COIN_FEE_TYPE_CONFIG_SELL                               // 只收卖的      // value --> 2
	COIN_FEE_TYPE_CONFIG_ALL                                // 两边都收
)

type CoinFeeConfig struct {
	ID        int64                  `json:"-" gorm:"primary_key"`
	CoinCode  string                 `json:"coin_code"`
	StockCode string                 `json:"stock_code"`
	FeeRatio  decimal.Decimal        `json:"fee_ratio" gorm:"type:decimal(36,18)"`
	Status    COIN_FEE_CONFIG_STATUS `json:"status"`
	Type      COIN_FEE_CONFIG_TYPE   `json:"type"`
	CreatedAt *time.Time             `json:"-"`
	UpdatedAt *time.Time             `json:"-"`
}
