package account

import (
	"business/support/domain"
	"business/support/libraries/mq"
	"encoding/json"
)

type WithdrawReq struct {
	mq.Request
	domain.Settle
}

func (*WithdrawReq) TubeName() string {
	return "AccountCenter.Withdraw.Request"
}

func (r *WithdrawReq) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (*WithdrawReq) New() mq.Message {
	return &WithdrawReq{}
}

type WithdrawResp struct {
	mq.Response
}

func (*WithdrawResp) TubeName() string {
	return "AccountCenter.Withdraw.Response"
}

func (r *WithdrawResp) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (*WithdrawResp) New() mq.Message {
	return &WithdrawResp{}
}
