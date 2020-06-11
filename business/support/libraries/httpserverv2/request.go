package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/Tomasen/realip"
	uuid "github.com/satori/go.uuid"
)

var requestID uint32

type Request struct {
	*http.Request
	ID string `json:"id,omitempty"`
}

func init() {
	requestID = 0
}

func NewRequest(r *http.Request) *Request {
	id, _ := uuid.NewV4()
	request := Request{
		Request: r,
		ID:      id.String(),
	}
	return &request
}

func (r *Request) ClientIP() string {
	return realip.RealIP(r.Request)
}

func (r *Request) reqStr() string {
	req := map[string]string{
		"client_ip": r.ClientIP(),
		"id":        r.ID,
	}
	for k := range r.Form {
		req[k] = r.Form.Get(k)
	}
	logStr, _ := json.Marshal(req)
	return string(logStr)
}
