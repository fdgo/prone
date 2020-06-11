package httpserver

import (
	"business/support/domain"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/Tomasen/realip"

	"github.com/gorilla/mux"
)

type AUTH_TYPE int

const (
	AUTH_TYPE_UNKOWN AUTH_TYPE = 0
	AUTH_TYPE_NORMAL AUTH_TYPE = 1 // 普通用户认证
	AUTH_TYPE_ACCESS AUTH_TYPE = 2 // accessKey认证
	AUTH_TYPE_CLOUD  AUTH_TYPE = 4 // cloud认证
)

var requestID uint32

type Request struct {
	ID             uint32              `json:"id,omitempty"`
	Method         string              `json:"method,omitempty"`
	Body           io.Reader           `json:"-"`
	BodyBuff       *bytes.Buffer       `json:"body,omitempty"`
	RemoteAddr     string              `json:"remoteAddr,omitempty"`
	QueryParams    url.Values          `json:"queryPrams,omitempty"`
	UrlParams      map[string]string   `json:"URLParams,omitempty"`
	URL            *url.URL            `json:"URL,omitempty"`
	Header         http.Header         `json:"header,omitempty"`
	Uid            int64               `json:"-"`
	Ts             int64               `json:"-"`
	Dev            string              `json:"-"`
	Ver            string              `json:"-"`
	Sign           string              `json:"-"`
	Ssid           string              `json:"-"`
	ActionSign     string              `json:"-"`
	Address        domain.ETHADDRESS   `json:"-"`
	AccessKey      string              `json:"-"`
	ExpiredTs      int64               `json:"-"`
	Language       string              `json:"-"`
	IDFA           string              `json:"-"`
	Model          string              `json:"-"`
	Udid           string              `json:"-"`
	ConnectionType string              `json:"-"`
	Platform       domain.APP_PLATFORM `json:"-"`
	Software       string              `json:"-"`
	Session        *domain.Session     `json:"-"`
	MultipartForm  *multipart.Form     `json:"-"`
	authType       AUTH_TYPE           `json:"-"`
}

func init() {
	requestID = 0
}

// NewRequest New a Request from http.Request
func NewRequest(r *http.Request) *Request {
	var request Request
	request.ID = atomic.AddUint32(&requestID, 1)
	request.Method = r.Method
	request.Body = r.Body
	request.BodyBuff = new(bytes.Buffer)
	request.BodyBuff.ReadFrom(r.Body)
	request.RemoteAddr = realip.RealIP(r)
	request.Header = r.Header
	request.UrlParams = mux.Vars(r)
	request.QueryParams = r.URL.Query()
	request.URL = r.URL
	request.Ver = r.Header.Get("bifund-Ver")
	request.Dev = r.Header.Get("bifund-Dev")
	request.Ssid = r.Header.Get("bifund-Ssid")
	request.Sign = r.Header.Get("bifund-Sign")
	request.ActionSign = r.Header.Get("bifund-Actionsign")
	request.AccessKey = r.Header.Get("bifund-AccessKey")
	request.Language = r.Header.Get("bifund-Language")
	request.IDFA = r.Header.Get("bifund-idfa")
	request.Model = r.Header.Get("bifund-model")
	request.Udid = r.Header.Get("bifund-udid")
	request.ConnectionType = r.Header.Get("bifund-Connection-Type")
	request.MultipartForm = r.MultipartForm
	tsStr := r.Header.Get("bifund-Ts")
	if tsStr != "" {
		request.Ts, _ = strconv.ParseInt(tsStr, 10, 64)

	}
	uidStr := r.Header.Get("bifund-Uid")
	if uidStr != "" {
		request.Uid, _ = strconv.ParseInt(uidStr, 10, 64)
	}
	expiredTsStr := r.Header.Get("bifund-ExpiredTs")
	if expiredTsStr != "" {
		request.ExpiredTs, _ = strconv.ParseInt(expiredTsStr, 10, 64)
	}
	if len(request.Dev) > 0 {
		dev := strings.ToLower(request.Dev)
		dev = strings.TrimSpace(dev)
		request.Dev = dev
		if dev == "android" {
			request.Platform = domain.APP_PLATFORM_ANDOID
		} else if dev == "ios" {
			request.Platform = domain.APP_PLATFORM_IOS
		} else if dev == "web" {
			request.Platform = domain.APP_PLATFORM_WEB
		}
	}
	request.Software = r.Header.Get("bifund-Software")
	return &request
}

// Parse Parse the JSON-encoded Request.Body store the result in the value pointed to by v
// If cache is true, cache the body in BodyBuff
func (self *Request) Parse(v interface{}) error {
	return json.Unmarshal(self.BodyBuff.Bytes(), v)
}

func (self *Request) ParseCache(v interface{}) error {
	return json.Unmarshal(self.BodyBuff.Bytes(), v)
}

func (self *Request) ParseByXML(v interface{}) error {
	return xml.Unmarshal(self.BodyBuff.Bytes(), v)
}

func (self *Request) ClientIP() string {
	return self.RemoteAddr
}

func (self *Request) API() string {
	if self.URL.RawQuery == "" {
		return self.URL.Path
	} else {
		rawQuery, _ := url.QueryUnescape(self.URL.RawQuery)
		return self.URL.Path + "?" + rawQuery
	}
}

func (self *Request) AuthType() AUTH_TYPE {
	if self.authType != AUTH_TYPE_UNKOWN {
		return self.authType
	}
	if len(self.Sign) < 32 {
		return AUTH_TYPE_UNKOWN
	}
	if len(self.AccessKey) <= 0 {
		self.authType = AUTH_TYPE_NORMAL
	} else if self.Ts < 1 {
		self.authType = AUTH_TYPE_UNKOWN
	} else if self.ExpiredTs < 1 {
		self.authType = AUTH_TYPE_ACCESS
	} else {
		self.authType = AUTH_TYPE_CLOUD
	}
	return self.authType
}

func (self *Request) IsSupportAuth(at int) bool {
	authType := int(self.AuthType())
	if authType == at&authType {
		return true
	}
	return false
}
