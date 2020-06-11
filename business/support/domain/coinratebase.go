package domain

import (
	"time"
)

type COIN_RATE_BASE_STATUS int

const (
	COIN_RATE_BASE_STATUS_UNKOWN  COIN_RATE_BASE_STATUS = iota
	COIN_RATE_BASE_STATUS_ENABLE                        // 可以用
	COIN_RATE_BASE_STATUS_DISABLE                       // 禁用中
)

type CoinRateBase struct {
	ID        int64                 `json:"-" gorm:"primary_key"`
	Name      string                `json:"name"`
	Status    COIN_RATE_BASE_STATUS `json:"-"`
	CreatedAt time.Time             `json:"-"`
	UpdatedAt time.Time             `json:"-"`
}
