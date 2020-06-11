package httpserver

import (
	"business/support/libraries/errorconfig"
	"business/support/libraries/errors"
	"business/support/libraries/loggers"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

// handlerWrapper Every http api hundler will be processed by this wrapper
func HandlerWrapper(f HandleFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var response *Response
		request := NewRequest(r)
		startTime := time.Now()
		defer func() {
			if r := recover(); r != nil {
				recoverError()
				response = NewResponseWithError(errors.InternalServerError)
				response.Write(w)
			}
			ts := time.Since(startTime)
			if ts < time.Second && response.Code == "OK" {
				loggers.Info.Printf(`%3d %s %s %s %s %s %s %d %v`,
					response.HTTPCode,
					response.Code,
					request.Method,
					getPathAndQuery(r.URL),
					request.ClientIP(),
					request.Dev,
					request.Ver,
					request.Uid,
					ts,
				)
			} else {
				loggers.Warn.Printf(`%3d %s %s %s %s %s %s %d %v`,
					response.HTTPCode,
					response.Code,
					request.Method,
					getPathAndQuery(r.URL),
					request.ClientIP(),
					request.Dev,
					request.Ver,
					request.Uid,
					ts,
				)
			}
		}()
		response = f(request)
		if 200 != response.HTTPCode || "OK" != response.Code {
			// 多语言转化错误
			if errMsg, ok := errorconfig.QueryErrMsg(request.Language, response.Code); ok {
				if len(response.MsgData) > 0 {
					errMsg = fmt.Sprintf(errMsg, response.MsgData[0:]...)
				}
				response.Msg = errMsg
			}
		}
		response.Write(w)

	}
}

// getPathAndQuery Combile path and raw query
func getPathAndQuery(u *url.URL) string {
	if u.RawQuery == "" {
		return u.Path
	} else {
		rawQuery, _ := url.QueryUnescape(u.RawQuery)
		return u.Path + "?" + rawQuery
	}
}

func recoverError() {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	loggers.Error.Printf("%s", buf[:n])
	loggers.Error.Printf("End here")
}
