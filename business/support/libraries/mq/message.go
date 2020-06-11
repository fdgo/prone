package mq

import (
	"business/support/libraries/errors"
	"time"
)

type Header struct {
	ID       uint64        `json:"-"`
	Tube     string        `json:"-"`
	State    string        `json:"-"`
	Pri      uint32        `json:"-"`
	Delay    time.Duration `json:"-"`
	TTR      time.Duration `json:"-"`
	Age      time.Duration `json:"-"`
	File     uint32        `json:"-"`
	Rservers uint32        `json:"-"`
	Timeouts uint32        `json:"-"`
	Releases uint32        `json:"-"`
	Buries   uint32        `json:"-"`
	Kicks    uint32        `json:"-"`
	TimeLeft time.Duration `json:"-"`
}

type Message interface {
	GetHeader() *Header
	SetHeader(header *Header)
	Marshal() ([]byte, error)
	TubeName() string
	ReqID() uint64
	SetReqID(id uint64)
	New() Message
}

type MessageImpl struct {
	header Header
}

// Please re-implement this function
func (*MessageImpl) TubeName() string {
	return ""
}

func (m *MessageImpl) GetHeader() *Header {
	return &m.header
}

func (m *MessageImpl) SetHeader(header *Header) {
	m.header = *header
}

// Please re-implement this function
func (m *MessageImpl) Marshal() ([]byte, error) {
	return nil, errors.New("Not implement marshal func")
}

func (m *MessageImpl) ReqID() uint64 {
	return 0
}

func (m *MessageImpl) SetReqID(id uint64) {
}

// Please re-implement this function
func (m *MessageImpl) New() Message {
	return &MessageImpl{}
}

type Request struct {
	MessageImpl
}

type Response struct {
	MessageImpl
	errors.Error
	ReqId uint64 `json:"req_id"`
}

func (m *Response) ReqID() uint64 {
	return m.ReqId
}

func (m *Response) SetReqID(id uint64) {
	m.ReqId = id
}

func (m *Response) New() Message {
	return &Response{}
}
