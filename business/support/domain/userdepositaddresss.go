package domain

import (
	"time"
)

type UserDepositAddress struct {
	ID             int64      `json:"-" gorm:"primary_key"`
	AccountId      int64      `json:"account_id,omitempty"`
	CoinGroup      string     `json:"coin_group,omitempty"`
	DepositAddress string     `json:"deposit_address"`
	CreatedAt      *time.Time `json:"-" gorm:"type:timestamp(6);not null"`
	UpdatedAt      *time.Time `json:"-"`
}

func (this *UserDepositAddress) TableName() string {
	return "user_deposit_addresss_new"
}
