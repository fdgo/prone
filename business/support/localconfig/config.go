package localconfig

import (
	"business/support/domain"
	"business/support/libraries/loggers"
	"flag"
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/jinzhu/configor"
)

type Config struct {
	Services map[string]domain.ServerConf
	Consuls  []string
	LogLevel string
	Ip       string
}

var (
	Conf           Config
	IsDebug        bool
	ConfigFilePath string
)

var (
	configFilePath = flag.String("c", "/data/etc/business/localconfig.yml", "config file path")
	isDebug        = flag.Bool("d", false, "debug mode")
)

func init() {
	decimal.DivisionPrecision = 19
	flag.Parse()
	if isDebug != nil {
		IsDebug = *isDebug
	}

	ConfigFilePath = *configFilePath
	if err := configor.Load(&Conf, *configFilePath); err != nil {
		fmt.Printf("localconfig Read config file:%s error:%s", *configFilePath, err.Error())
	}

	if !IsDebug {
		loggers.InitLogLevel(Conf.LogLevel)
	}
}
