package app

import (
	"business/services/base/usercenter/worker"
	"business/support/libraries/loggers"
)

func init() {
	loggers.Info.Printf("UserCenter service init  [\033[0;32;1mOK\033[0m]")
}

func Start(addr string) {
	go worker.Start()
	loggers.Info.Printf("UserCenter service start [\033[0;32;1mOK\t%+v\033[0m]", addr)
	select {}
}
