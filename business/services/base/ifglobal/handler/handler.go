package handler

import (
	"business/support/config"
	"business/support/domain"
	"business/support/libraries/loggers"
)

var ifglobalConfig domain.ServerConf

func init() {
	ifglobalConfig = config.Conf.Services["IFGlobal"]
	if len(ifglobalConfig.ListenAddress) <= 0 {
		loggers.Warn.Printf("Load IFGlobal failed \n")
		panic(88)
	}
}

func GetServiceConfig() *domain.ServerConf {
	return &ifglobalConfig
}
