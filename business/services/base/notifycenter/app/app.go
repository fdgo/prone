package app

import (
	"business/services/base/notifycenter/worker"
	"business/support/libraries/loggers"
	"time"
)

func Start(addr string) {
	loggers.Debug.Printf("OK")
	sender := worker.Sender{Interval: time.Second * 5}

	go sender.Start()
	loggers.Info.Printf("NotifyCenter service start [\033[0;32;1mOK\t%+v\033[0m]", addr)
	select {}
}
