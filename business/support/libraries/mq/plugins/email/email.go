package email

import (
	"business/support/libraries/mq"
	"encoding/json"
	"time"
)

type EmailType int

const (
	_ EmailType = iota
	EmailTypeHTML
	EmailTypeTemplate
)

type EmailMessage struct {
	header       mq.Header          `json:"-"`
	From         string             `json:"from"`
	FromName     string             `json:"fromName"`
	To           string             `json:"to"`
	Cc           string             `json:"cc,omitempty"`
	Subject      string             `json:"subject"`
	EmailType    EmailType          `json:"emailType"`
	Html         string             `json:"html,omitempty"`
	Attachments  *map[string][]byte `json:"attachments,omitempty"`
	TemplateName string             `json:"templateName,omitempty"`
	TemplateBody interface{}        `json:"templateBody,omtiempty"`
}

var (
	TubeName string
	SendTTR  time.Duration
)

func init() {
	TubeName = "EMAIL"
	SendTTR = time.Minute * 1
}

func (eMsg *EmailMessage) Marshal() ([]byte, error) {
	return json.Marshal(eMsg)
}

func (eMsg *EmailMessage) GetHeader() *mq.Header {
	return &eMsg.header
}

func (eMsg *EmailMessage) SetHeader(header *mq.Header) {
	eMsg.header = *header
}

func (eMsg *EmailMessage) SetDelay(delay time.Duration) {
	eMsg.header.Delay = delay
}

func (eMsg *EmailMessage) SetPri(pri uint32) {
	eMsg.header.Pri = pri
}

func (eMsg *EmailMessage) ID() uint64 {
	return eMsg.header.ID
}
