package handler

import (
	"business/services/base/ifglobal/adaptor"
	"business/support/libraries/errors"
	"business/support/libraries/httpserver"
)

func PhoneCode(r *httpserver.Request) *httpserver.Response {
	codes, _ := adaptor.GetPhoneCode()
	if nil == codes {
		return httpserver.NewResponseWithError(errors.NotFound)
	}
	resp := httpserver.NewGZipResponse()
	resp.Data = codes
	return resp
}
