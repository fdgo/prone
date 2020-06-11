package main

import (
	"business/support/domain"
	"business/support/libraries/loggers"
	"business/support/libraries/mq"
	"business/support/libraries/mq/plugins/order"
	"time"
)

func main() {
	timeout := time.Second * 30
	mqConf := domain.MQConf{
		Host:    "192.168.5.141",
		Port:    11300,
		Network: "tcp",
	}
	accessor, err := mq.NewAcessor(&mqConf, timeout)
	if err != nil {
		loggers.Error.Printf("NewProcessor error:%s", err.Error())
	}
	accessor.Register(&order.SubmitOrderReq{}, handleSubmitOrder, timeout)
	for {
		time.Sleep(timeout)
	}
}

func handleSubmitOrder(msg mq.Message) mq.Message {
	req, ok := msg.(*order.SubmitOrderReq)
	if !ok {
		loggers.Error.Printf("Invalid req type")
		return &order.SubmitOrderResp{}
	}
	loggers.Debug.Printf("OrderID:%d AccountID:%d ", req.OrderId, req.AccountId)
	resp := order.SubmitOrderResp{}
	resp.ErrCode = ""
	return &resp
}
