package domain

import "time"

type COIN_STATUS int

const (
	COIN_STATUS_UNKOWN   COIN_STATUS = iota
	COIN_STATUS_ENABLE               // 可以用
	COIN_STATUS_DISABLE              // 禁用中
	COIN_STATUS_APPROVAL             // 审批中
	COIN_STATUS_REJECT               // 审批不通过
)

// Coin 货币
type Coin struct {
	ID        int64       `json:"-" gorm:"primary_key"`
	Name      string      `json:"name"`
	CoinGroup string      `json:"coin_group"`
	Rank      int64       `json:"rank,omitempty"`
	Big       string      `json:"big,omitempty"`
	Small     string      `json:"small,omitempty"`
	Gray      string      `json:"gray,omitempty"`
	Status    COIN_STATUS `json:"-"`
	CreatedAt *time.Time  `json:"-"`
	UpdatedAt *time.Time  `json:"-"`
}

// TableName 表名
func (*Coin) TableName() string {
	return "coins"
}

// ContractAddress ERC20代币信息
type ContractAddress struct {
	ID        int64      `json:"-" gorm:"primary_key"`
	CoinCode  string     `json:"coin_code,required" gorm:"size:16;unique_index"`
	Address   string     `json:"address,required"`
	Decimals  int32      `json:"decimals"`
	CreatedAt *time.Time `json:"-"`
	UpdatedAt *time.Time `json:"-"`
	CoinGroup string     `json:"coin_group"`
}

// TableName 表名
func (*ContractAddress) TableName() string {
	return "contract_addresss"
}
