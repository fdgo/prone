package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// Coin 货币
type CoinRate struct {
	ID        int64           `json:"-" gorm:"primary_key"`
	CoinName  string          `json:"coin_name"`
	Rate      decimal.Decimal `json:"rate" gorm:"type:decimal(36,18);not null"`
	CreatedAt *time.Time      `json:"-"`
	UpdatedAt *time.Time      `json:"-"`
}

func (*CoinRate) TableName() string {
	return "coin_rates"
}
