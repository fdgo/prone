package config

import (
	"business/support/libraries/loggers"
	"business/support/libraries/utils"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type GlobalConfig struct {
	ID        int64      `json:"-" gorm:"primary_key" xlsx:"-"`
	Key       string     `json:"key" gorm:"unique_index"`
	Value     []byte     `json:"value"`
	TimeOut   int64      `json:"time_out"`
	Error     error      `json:"error" gorm:"-"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (c *GlobalConfig) TableName() string {
	return "global_configs"
}

func (c *GlobalConfig) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}

func (c GlobalConfig) MarshalBinary() (data []byte, err error) {
	return json.Marshal(c)
}

func (c *GlobalConfig) Val() string {
	return string(c.Value)
}

func (c *GlobalConfig) Result() (string, error) {
	return c.Val(), c.Error
}

func (c *GlobalConfig) Bytes() ([]byte, error) {
	return c.Value, c.Error
}

func (c *GlobalConfig) Int64() (int64, error) {
	if c.Error != nil {
		return 0, c.Error
	}
	return strconv.ParseInt(c.Val(), 10, 64)
}

func (c *GlobalConfig) Uint64() (uint64, error) {
	if c.Error != nil {
		return 0, c.Error
	}
	return strconv.ParseUint(c.Val(), 10, 64)
}

func (c *GlobalConfig) Scan(val interface{}) error {
	if c.Error != nil {
		return c.Error
	}
	return utils.Scan(c.Value, val)
}

func (c *GlobalConfig) UnmarshalJSON(val interface{}) error {
	return json.Unmarshal(c.Value, val)
}

func AddGlobalConfig(c *GlobalConfig) error {
	if c.Key == "" || c.Value == nil {
		return errors.New("No key or value")
	}
	return dbPool.NewConn().Create(c).Error
}

func UpdateGlobalConfig(c *GlobalConfig) error {
	if c.Key == "" {
		return errors.New("No key")
	}
	if dbPool == nil {
		return errors.New("Config DB is nil")
	}
	updates := map[string]interface{}{
		"key":      c.Key,
		"value":    c.Value,
		"time_out": c.TimeOut,
	}
	return dbPool.NewConn().Model(c).Where("key = ?", c.Key).Updates(updates).Error
}

func GetGlobal(key string) *GlobalConfig {
	var conf GlobalConfig
	if dbPool == nil {
		conf.Error = errors.New("Config DB is nil")
		return &conf
	}

	if err := dbPool.NewConn().Where("key = ?", key).First(&conf).Error; err != nil {
		conf.Error = err
	}
	return &conf
}

func GetAndCacheGlobal(key string) *GlobalConfig {
	var (
		redisKey = fmt.Sprintf("CONFIG.GLOBAL.%s", key)
		conf     = new(GlobalConfig)
	)
	if dbPool == nil {
		conf.Error = errors.New("Config DB is nil")
		return conf
	}
	if redisPool == nil {
		conf.Error = errors.New("Redis DB is nil")
		return conf
	}
	if err := redisPool.NewConn().Get(redisKey).Scan(conf); err == nil {
		return conf
	}
	if conf = GetGlobal(key); conf.Error != nil {
		return conf
	}
	if conf.TimeOut > 0 {
		if err := redisPool.NewConn().Set(redisKey, conf, time.Second*time.Duration(conf.TimeOut)).Err(); err != nil {
			loggers.Warn.Printf("GetAndCacheGlobal cache global config [%s] error %s", redisKey, err.Error())
		}
	}
	return conf
}

func GetAndParseCacheGlobal(key string, v interface{}) error {
	conf := GetAndCacheGlobal(key)
	if conf.Error != nil {
		return conf.Error
	}
	if err := json.Unmarshal(conf.Value, v); err != nil {
		return err
	}
	return nil
}
