package handler

import (
	"business/services/base/ifglobal/adaptor"
	"business/support/domain"
	"business/support/libraries/errors"
	"business/support/libraries/httpserver"
	"strconv"
)

type CheckUpdateParam struct {
	Platform domain.APP_PLATFORM
	Version  string
}

func GetCheckUpdateParam(r *httpserver.Request) *CheckUpdateParam {
	var params CheckUpdateParam
	temp := r.QueryParams.Get("platform")
	if len(temp) <= 0 {
		return nil
	}
	platform, err := strconv.ParseInt(temp, 10, 32)
	if nil != err {
		return nil
	}
	params.Platform = domain.APP_PLATFORM(platform)
	params.Version = r.QueryParams.Get("version")
	return &params
}

type AppBuildsParam struct {
	Platform domain.APP_PLATFORM
}

func GetAppBuildsParam(r *httpserver.Request) *AppBuildsParam {
	var params AppBuildsParam
	temp := r.QueryParams.Get("platform")
	if len(temp) <= 0 {
		return &params
	}
	platform, err := strconv.ParseInt(temp, 10, 32)
	if nil != err {
		return nil
	}
	params.Platform = domain.APP_PLATFORM(platform)
	return &params
}

func CheckUpdate(r *httpserver.Request) *httpserver.Response {
	params := GetCheckUpdateParam(r)
	if nil == params {
		return httpserver.NewResponseWithError(errors.ParameterError)
	}
	build, err := adaptor.QueryAppUpdate(params.Platform, params.Version, r.Language)
	if nil != err || nil == build {
		return httpserver.NewResponseWithError(errors.NotFound)
	}
	resp := httpserver.NewGZipResponse()
	resp.Data = build
	return resp
}

func AppBuilds(r *httpserver.Request) *httpserver.Response {
	params := GetAppBuildsParam(r)
	if nil == params {
		return httpserver.NewResponseWithError(errors.ParameterError)
	}
	builds, err := adaptor.GetAppPlatformBuilds(params.Platform, r.Language)
	if nil != err {
		return httpserver.NewResponseWithError(err.(*errors.Error))
	}
	resp := httpserver.NewGZipResponse()
	resp.Data = builds
	return resp
}
