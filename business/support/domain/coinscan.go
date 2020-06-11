package domain

type COIN_SCAN_STATUS int

const (
	COIN_SCAN_STATUS_UNKOWN  COIN_SCAN_STATUS = iota
	COIN_SCAN_STATUS_ENABLE                   // 可以用
	COIN_SCAN_STATUS_DISABLE                  // 禁用中
)

type CoinScan struct {
	ID        int64            `json:"-" gorm:"primary_key"`
	CoinGroup string           `json:"coin_group"`
	Address   string           `json:"address"`
	Tx        string           `json:"tx"`
	Regex     string           `json:"regex"`
	Status    COIN_SCAN_STATUS `json:"-"`
}

func (*CoinScan) TableName() string {
	return "coin_scans"
}
