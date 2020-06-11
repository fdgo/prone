package domain

import (
	"time"
)

// CoinNoticeConfig 货币提示配置
type CoinNoticeConfig struct {
	ID        int64      `json:"-" gorm:"primary_key;auto_increment"`
	CoinCode  string     `json:"coin_code" gorm:"size:64"`
	Local     string     `json:"-" gorm:"size:32"`
	Deposit   string     `json:"deposit"`
	Withdraw  string     `json:"withdraw"`
	CreatedAt *time.Time `json:"-"`
	UpdatedAt *time.Time `json:"-"`
}
