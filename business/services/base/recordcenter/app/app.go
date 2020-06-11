package app

import (
	"business/services/base/recordcenter/worker"
	"business/support/libraries/httpserver"
	"business/support/libraries/loggers"
	"time"
)

func init() {
	loggers.Info.Printf("RecordCenter service init  [\033[0;32;1mOK\033[0m]")
}

func Start(addr string) {
	r := httpserver.NewRouter()
	recoder := worker.Recorder{
		Interval: time.Second * 5,
	}
	go recoder.Start()
	go r.ListenAndServe(addr)
	loggers.Info.Printf("RecordCenter service start [\033[0;32;1mOK\t%+v\033[0m]", addr)
	select {}
}
