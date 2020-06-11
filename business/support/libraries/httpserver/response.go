package httpserver

import (
	"business/support/domain"
	"business/support/libraries/errors"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	HTTPCode     int           `json:"-"`
	Code         string        `json:"errno"`
	Msg          string        `json:"message"`
	Header       http.Header   `json:"-"`
	Data         interface{}   `json:"data,omitempty"`
	IsGZip       bool          `json:"-"`
	IsResetToken bool          `json:"-"`
	MsgData      []interface{} `json:"-"`
}

type XMLResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

// NewResponse
func NewResponse() *Response {
	resp := &Response{
		HTTPCode: http.StatusOK,
		Code:     "OK",
		Msg:      "Success",
		Header:   make(http.Header),
	}
	resp.Header.Set("Content-Type", "application/json")
	return resp
}

// NewResponse
func NewGZipResponse() *Response {
	resp := NewResponse()
	resp.Header.Set("Content-Encoding", "gzip")
	resp.IsGZip = true
	return resp
}

// NewResponseWithError
func NewResponseWithError(err *errors.Error) *Response {
	resp := &Response{
		HTTPCode: err.HttpCode(),
		Code:     err.Error(),
		Msg:      err.Msg(),
		Header:   make(http.Header),
	}
	resp.Header.Set("Content-Type", "application/json")
	return resp
}

// NewResponseWithFormatError 构造格式化参数的error
func NewResponseWithFormatError(err *errors.Error, extra []interface{}) *Response {
	resp := &Response{
		HTTPCode: err.HttpCode(),
		Code:     err.Error(),
		Msg:      err.Msg(),
		MsgData:  extra,
		Header:   make(http.Header),
	}
	resp.Header.Set("Content-Type", "application/json")
	return resp
}

// NewResponseForRedirect
func NewResponseForRedirect(location string) *Response {
	resp := &Response{
		HTTPCode: http.StatusFound,
		Header:   make(http.Header),
	}
	resp.Header.Set("Location", location)
	return resp
}

// Write
func (r *Response) Write(w http.ResponseWriter) {
	for k, vv := range r.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(r.HTTPCode)
	body, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	if !r.IsGZip {
		w.Write(body)
	} else {
		zipper := gzip.NewWriter(w)
		zipper.Write(body)
		zipper.Close()
	}
}

func (r *Response) SetSession(s *domain.Session) {
	uid := fmt.Sprintf("%d", s.Account.AccountId)
	r.Header.Set("bifund-Ssid", s.Ssid)
	r.Header.Set("bifund-Token", s.Token)
	r.Header.Set("bifund-Uid", uid)
}
