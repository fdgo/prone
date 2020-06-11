package main

import (
	"business/support/domain"
	"business/support/libraries/loggers"
	"business/support/libraries/mq"
	"business/support/libraries/mq/plugins/order"
	"time"
)

func main() {
	timeout := time.Second * 5
	mqConf := domain.MQConf{
		Host:    "192.168.5.141",
		Port:    11300,
		Network: "tcp",
	}
	processor, err := mq.NewProcessor(&mqConf)
	if err != nil {
		loggers.Error.Printf("NewProcessor error:%s", err.Error())
	}
	processor.Register(&order.SubmitOrderResp{}, timeout)

	for i := 0; i < 10; i++ {
		val, err := processor.Put(&req, timeout)
		if err != nil {
			loggers.Error.Printf("Put error:%s", err.Error())
			continue
		}
		resp, ok := val.(*order.SubmitOrderResp)
		if !ok {
			loggers.Error.Printf("Type Not Match")
			break
		}
		loggers.Debug.Printf("errno:%d errmsg:%s", resp.ErrCode, resp.ErrMsg)
	}

	for {
		time.Sleep(time.Second * 10)
	}
}
