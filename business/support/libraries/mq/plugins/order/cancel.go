package order

import (
	"business/support/domain"
	"business/support/libraries/mq"
	"encoding/json"
)

type CancelOrderReq struct {
	mq.Request
	domain.Order
}

func (*CancelOrderReq) TubeName() string {
	return "OrderCenter.CancelOrder.Request"
}

func (r *CancelOrderReq) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (*CancelOrderReq) New() mq.Message {
	return &CancelOrderReq{}
}

type CancelOrderResp struct {
	mq.Response
}

func (*CancelOrderResp) TubeName() string {
	return "OrderCenter.CancelOrder.Response"
}

func (r *CancelOrderResp) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (*CancelOrderResp) New() mq.Message {
	return &CancelOrderResp{}
}
