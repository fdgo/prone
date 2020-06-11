package order

import (
	"business/support/domain"
	"business/support/libraries/mq"
	"encoding/json"
)

type SubmitOrderReq struct {
	mq.Request
	domain.Order
}

func (*SubmitOrderReq) TubeName() string {
	return "OrderCenter.SubmitOrder.Request"
}

func (r *SubmitOrderReq) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *SubmitOrderReq) New() mq.Message {
	return &SubmitOrderReq{}
}

type SubmitOrderResp struct {
	mq.Response
}

func (*SubmitOrderResp) TubeName() string {
	s := mq.Response{}
	s.ErrCode = ""
	return "OrderCenter.SubmitOrder.Response"
}

func (r *SubmitOrderResp) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *SubmitOrderResp) New() mq.Message {
	return &SubmitOrderResp{}
}
