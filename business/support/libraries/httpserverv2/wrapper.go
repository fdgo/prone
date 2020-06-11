package httpserver

import (
	"business/support/libraries/errors"
	"business/support/libraries/loggers"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/Tomasen/realip"
)

// handlerWrapper Every http api hundler will be processed by this wrapper
func HandlerWrapper(f HandleFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			response *Response
		)
		request := NewRequest(r)
		startTime := time.Now()
		defer func() {
			if r := recover(); r != nil {
				recoverError()
				response = NewResponseWithError(errors.InternalServerError)
				response.Write(w)
			}
			loggers.Info.Printf(`%3d %s %s %s %s %v request:%s`,
				response.HTTPCode,
				response.Code,
				request.Method,
				getPathAndQuery(r.URL),
				realip.RealIP(r),
				time.Since(startTime),
				request.reqStr(),
			)
		}()
		response = f(request)
		response.ID = request.ID
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
