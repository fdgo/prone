package httpserver

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"business/support/libraries/version"
)

type HandleFunc func(*Request) *Response

// Router 路由
type Router struct {
	*mux.Router
}

// NewRouter 新路由
func NewRouter() Router {
	r := Router{mux.NewRouter()}
	r.RouteAlive()
	return r
}

// RouteAlive 活跃检测
func (r Router) RouteAlive() *mux.Route {
	return r.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		response := NewResponse()
		response.Data = version.ServicesInfo()
		response.Write(w)
	})
}

// RouteHandleFunc 路由
func (r Router) RouteHandleFunc(path string, f HandleFunc) *mux.Route {
	return r.HandleFunc(path, HandlerWrapper(f))
}

// ListenAndServe This function blocks
func (r Router) ListenAndServe(addr string) error {
	r.Use(handlers.HTTPMethodOverrideHandler)
	return http.ListenAndServe(addr, r)
}

// ListenAndServeCORS This function blocks
func (r Router) ListenAndServeCORS(addr string, opts ...handlers.CORSOption) error {
	r.Use(handlers.HTTPMethodOverrideHandler)
	return http.ListenAndServe(addr, handlers.CORS(opts...)(r))
}
