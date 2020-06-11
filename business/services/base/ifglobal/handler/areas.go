package handler

import (
	"business/services/base/ifglobal/adaptor"
	"business/support/libraries/errors"
	"business/support/libraries/httpserver"
	"business/support/libraries/loggers"
)

func Areas(r *httpserver.Request) *httpserver.Response {
	language := r.Header.Get("bifund-Language")
	if language != "zh-cn" && language != "en" {
		language = "en"
	}
	areas, err := adaptor.GetAreas(language)
	if err != nil {
		loggers.Error.Printf("Get areas error %s", err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}

	resp := httpserver.NewResponse()
	resp.Data = areas
	return resp
}
