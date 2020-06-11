package account

import (
	"business/support/domain"
	"business/support/libraries/mq"
	"encoding/json"
)

type BindSettleAddressReq struct {
	mq.Request
	domain.Account
}

func (*BindSettleAddressReq) TubeName() string {
	return "AccountCenter.BindSettleAddress.Request"
}

func (r *BindSettleAddressReq) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (*BindSettleAddressReq) New() mq.Message {
	return &BindSettleAddressReq{}
}

type BindSettleAddressResp struct {
	mq.Response
	domain.Account
}

func (*BindSettleAddressResp) TubeName() string {
	return "AccountCenter.BindSettleAddress.Response"
}

func (r *BindSettleAddressResp) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (*BindSettleAddressResp) New() mq.Message {
	return &BindSettleAddressResp{}
}
