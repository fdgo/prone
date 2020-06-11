package httpserver

import (
	"business/support/libraries/errors"
	"compress/gzip"
	"encoding/json"
	"net/http"
)

var signFunc ResponseSignFunc

type ResponseSignFunc func(r *Response, body []byte)

type Response struct {
	HTTPCode     int         `json:"-"`
	Code         string      `json:"errno"`
	Msg          string      `json:"message"`
	ID           string      `json:"id"`
	Header       http.Header `json:"-"`
	Data         interface{} `json:"data,omitempty"`
	IsGZip       bool        `json:"-"`
	IsResetToken bool        `json:"-"`
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

func RegsiterSignFunc(f ResponseSignFunc) {
	signFunc = f
}

// NewResponse
func NewGZipResponse() *Response {
	resp := NewResponse()
	resp.Header.Set("Content-Encoding", "gzip")
	resp.IsGZip = true
	return resp
}

// NewResponseWithError
func NewResponseWithError(err error) *Response {
	e, ok := errors.Errors[err.Error()]
	if !ok {
		e = errors.InternalServerError
	}
	resp := &Response{
		HTTPCode: e.HttpCode(),
		Code:     e.Error(),
		Msg:      e.Msg(),
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
func (r *Response) Write(w http.ResponseWriter) string {
	body, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	signFunc(r, body)
	for k, vv := range r.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(r.HTTPCode)
	if !r.IsGZip {
		w.Write(body)
		return string(body)
	} else {
		zipper := gzip.NewWriter(w)
		zipper.Write(body)
		zipper.Close()
		return "gzip"
	}
}
