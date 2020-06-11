package domain

import (
	"time"
)

type Session struct {
	Ssid       string     `json:"ssid"`
	Ts         int64      `json:"ts"`
	Token      string     `json:"token"`
	AssetToken string     `json:"asset_token"`
	Account    *Account   `json:"-"`
	CreatedAt  *time.Time `json:"created_at"`
}
