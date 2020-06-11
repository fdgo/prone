package handler

import (
	"business/services/base/ifaccount/adaptor"
	"business/support/libraries/httpserver"
)

func CaptchCheck(r *httpserver.Request) *httpserver.Response {
	var (
		ip     = r.ClientIP()
		action = r.QueryParams.Get("action")
		ret    = adaptor.CheckCaptchLimit(action, ip)
	)
	resp := httpserver.NewResponse()
	resp.Data = map[string]bool{
		"need": ret,
	}
	return resp
}
