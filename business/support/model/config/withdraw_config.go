package config

import (
	"business/support/domain"
	"business/support/libraries/loggers"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

type WITHDRAWAL_CONFIG_STATUS int32

const (
	WITHDRAWAL_CONFIG_STATUS_UNKOWN  WITHDRAWAL_CONFIG_STATUS = iota
	WITHDRAWAL_CONFIG_STATUS_ENABLE                           // 可以用
	WITHDRAWAL_CONFIG_STATUS_DISABLE                          // 禁用中
)

type WithDrawalConfig struct {
	ID                 int64                    `json:"-" gorm:"primary_key"`
	CoinName           string                   `json:"coin_name"`
	Fee                decimal.Decimal          `json:"fee" gorm:"type:decimal(36,18)"`
	FeeCoinName        string                   `json:"fee_coin_name"`
	MinVol             decimal.Decimal          `json:"min_vol" gorm:"type:decimal(36,18)"`
	DepositMinVol      decimal.Decimal          `json:"deposit_min_vol" gorm:"type:decimal(36,18)"`
	KYCMaxVol          decimal.Decimal          `json:"kyc_max_vol" gorm:"type:decimal(36,18)"` // 没有通过KYC认证的最大提现限制
	Status             WITHDRAWAL_CONFIG_STATUS `json:"-"`
	CoinPrecision      int64                    `json:"coin_precision"`
	WalletMaxVol       decimal.Decimal          `json:"wallet_max_vol" gorm:"type:decimal(36,18);DEFAULT:0"`        //向大钱包充值的最大值
	MiddleWalletMaxVol decimal.Decimal          `json:"middle_wallet_max_vol" gorm:"type:decimal(36,18);DEFAULT:0"` //中间归集钱包(OKEX钱包)的最大值
	CreatedAt          *time.Time               `json:"-"`
	UpdatedAt          *time.Time               `json:"-"`
}

type CoinDetailConfig struct {
	WithDrawalConfig
	CoinGroup     string `json:"coin_group"`
	TokenAddress  string `json:"token_address"`
	TokenDecimals int32  `json:"token_decimals"`
}
type WithdrawConfigResp struct {
	Total   int32               `json:"total"`
	Configs []*CoinDetailConfig `json:"configs"`
}

func (this *WithDrawalConfig) TableName() string {
	return "coin_withdrawal_configs"
}

func (c *WithDrawalConfig) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}

func (c WithDrawalConfig) MarshalBinary() (data []byte, err error) {
	return json.Marshal(c)
}

func (*WithDrawalConfig) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("created_at", time.Now().UTC())
	scope.SetColumn("updated_at", time.Now().UTC())
	return nil
}
func (*WithDrawalConfig) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("updated_at", time.Now().UTC())
	return nil
}

func AddWithdrawConfig(c *WithDrawalConfig) error {
	if c.CoinName == "" {
		return errors.New("No coin name")
	}
	return dbPool.NewConn().Create(c).Error
}

func UpdateWithdrawConfig(c *WithDrawalConfig) error {
	if c.CoinName == "" {
		return errors.New("No coin name")
	}
	if dbPool == nil {
		return errors.New("Config DB is nil")
	}
	updates := map[string]interface{}{
		"coin_name":     c.CoinName,
		"fee":           c.Fee,
		"fee_coin_name": c.FeeCoinName,
		"min_vol":       c.MinVol,
		"kyc_max_vol":   c.KYCMaxVol,
		"status":        c.Status,
	}
	return dbPool.NewConn().Model(c).Where("coin_name = ?", c.CoinName).Updates(updates).Error
}

func GetWithdrawConfig(coinName string) (*WithDrawalConfig, error) {
	conf := new(WithDrawalConfig)
	if dbPool == nil {
		return nil, errors.New("Config DB is nil")
	}

	if err := dbPool.NewConn().Where("coin_name  = ?", coinName).First(conf).Error; err != nil {
		return nil, err
	}
	return conf, nil
}

func GetAndCacheWithdrawConfig(coinName string) (*WithDrawalConfig, error) {
	var (
		redisKey = fmt.Sprintf("CONFIG.WITHDRAW_CONFIG.%s", coinName)
		conf     = new(WithDrawalConfig)
		err      error
	)
	if dbPool == nil {
		return nil, errors.New("Config DB is nil")
	}
	if redisPool == nil {
		return nil, errors.New("Redis DB is nil")
	}
	if err = redisPool.NewConn().Get(redisKey).Scan(conf); err == nil {
		return conf, nil
	}
	if conf, err = GetWithdrawConfig(coinName); err != nil {
		return nil, err
	}
	if err = redisPool.NewConn().Set(redisKey, conf, time.Hour*2).Err(); err != nil {
		loggers.Warn.Printf("GetAndCacheGlobal cache global config [%s] error %s", redisKey, err.Error())
	}
	return conf, nil
}

func GetWithdrawConfigs() ([]*WithDrawalConfig, error) {
	var confs []*WithDrawalConfig
	if dbPool == nil {
		return nil, errors.New("Config DB is nil")
	}

	if err := dbPool.NewConn().Where("status = ?", WITHDRAWAL_CONFIG_STATUS_ENABLE).Find(&confs).Error; err != nil {
		return nil, err
	}
	return confs, nil
}

func GetAndCacheWithdrawConfigs() ([]*WithDrawalConfig, error) {
	var (
		redisKey = fmt.Sprintf("CONFIG.WITHDRAW_CONFIG.ALL")
		confs    []*WithDrawalConfig
		err      error
	)
	if dbPool == nil {
		return nil, errors.New("Config DB is nil")
	}
	if redisPool == nil {
		return nil, errors.New("Redis DB is nil")
	}

	v, err := redisPool.NewConn().Get(redisKey).Bytes()
	if err == nil {
		if err := json.Unmarshal(v, &confs); err == nil {
			return confs, nil
		}
	}

	confs, err = GetWithdrawConfigs()
	if err != nil {
		return nil, err
	}

	if v, err = json.Marshal(confs); err == nil {
		if err := redisPool.NewConn().Set(redisKey, v, time.Hour*2).Err(); err != nil {
			loggers.Warn.Printf("GetAndCacheGlobal cache global config [%s] error %s", redisKey, err.Error())
		}
	}

	return confs, nil
}

func GetWithdrawSettleConfig(coin string) (*CoinDetailConfig, error) {
	configs, err := GetWithdrawSettlesConfigs()
	if err != nil {
		return nil, err
	}

	for _, v := range configs {
		if v.CoinName == coin {
			return v, nil
		}
	}
	return nil, errors.New("config null")
}

func GetWithdrawSettlesConfigs() ([]*CoinDetailConfig, error) {
	s := getConfig()
	if s != nil {
		return s, nil
	}

	var wconfigs []*WithDrawalConfig
	db := dbPool.NewConn()
	err := db.Find(&wconfigs).Error
	if err != nil {
		loggers.Error.Println("GetWithdrawConfig ", err)
		return nil, err
	}

	if len(wconfigs) == 0 {
		loggers.Error.Println("GetWithdrawConfig query withdrawconfig ret invalid")
		return nil, errors.New("WithDrawalConfig data null")
	}

	tconfigs := make(map[string]*domain.ContractAddress)
	var tokenConfigs []*domain.ContractAddress
	err = db.Find(&tokenConfigs).Error
	if err != nil {
		loggers.Error.Println("GetWithdrawConfig ", err)
		return nil, err
	}

	if len(tokenConfigs) == 0 {
		loggers.Error.Println("GetWithdrawConfig query contractaddress ret invalid")
		return nil, errors.New("ContractAddress data null")
	}

	for _, t := range tokenConfigs {
		tconfigs[t.CoinCode] = t
	}

	cconfigs := make(map[string]*domain.Coin)

	var coins []*domain.Coin
	err = exchangePool.NewConn().Find(&coins).Error
	if err != nil {
		loggers.Error.Println("GetWithdrawConfig ", err)
		return nil, err
	}

	if len(coins) == 0 {
		loggers.Error.Println("GetWithdrawConfig query coins ret invalid")
		return nil, errors.New("coin data null")
	}

	for _, c := range coins {
		cconfigs[c.Name] = c
	}

	cfg := make([]*CoinDetailConfig, 0)
	for _, w := range wconfigs {
		var c CoinDetailConfig
		c.WithDrawalConfig = *w
		d, ok := cconfigs[w.CoinName]
		if ok {
			c.CoinGroup = d.CoinGroup
		}
		e, o := tconfigs[w.CoinName]
		if o {
			c.TokenAddress = e.Address
			c.TokenDecimals = e.Decimals
		}

		cfg = append(cfg, &c)
	}

	saveConfig(cfg)
	return cfg, nil
}

func saveConfig(config []*CoinDetailConfig) {
	redisKey := domain.KConfig

	data, e := json.Marshal(config)
	if e == nil {
		_, e = redisPool.NewConn().Set(redisKey.Key, data, redisKey.Timeout).Result()
		if e != nil {
			loggers.Warn.Println("SaveConfig", e)
		}
	}
}

func getConfig() []*CoinDetailConfig {
	var redisdata []*CoinDetailConfig
	redisKey := domain.KConfig
	str, e := redisPool.NewConn().Get(redisKey.Key).Result()
	if e == nil {
		e = json.Unmarshal([]byte(str), &redisdata)
		if e == nil {
			return redisdata
		}
	}

	return nil
}
