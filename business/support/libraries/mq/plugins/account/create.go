package account

import (
	"business/support/domain"
	"business/support/libraries/mq"
	"encoding/json"
)

type CreateAccountReq struct {
	mq.Request
	domain.Account
}

func (*CreateAccountReq) TubeName() string {
	return "AccountCenter.Create.Request"
}

func (r *CreateAccountReq) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (*CreateAccountReq) New() mq.Message {
	return &CreateAccountReq{}
}

type CreateAccountResp struct {
	domain.Account
	mq.Response
}

func (*CreateAccountResp) TubeName() string {
	return "AccountCenter.Create.Response"
}

func (r *CreateAccountResp) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (*CreateAccountResp) New() mq.Message {
	return &CreateAccountResp{}
}
