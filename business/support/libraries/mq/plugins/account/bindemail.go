package account

import (
	"business/support/domain"
	"business/support/libraries/mq"
	"encoding/json"
)

type BindEmailReq struct {
	mq.Request
	domain.Account
}

func (*BindEmailReq) TubeName() string {
	return "AccountCenter.BindEmail.Request"
}

func (r *BindEmailReq) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *BindEmailReq) New() mq.Message {
	return &BindEmailReq{}
}

type BindEmailResp struct {
	mq.Response
}

func (*BindEmailResp) TubeName() string {
	return "AccountCenter.BindEmail.Response"
}

func (r *BindEmailResp) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *BindEmailResp) New() mq.Message {
	return &BindEmailResp{}
}
