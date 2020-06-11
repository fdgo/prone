package domain

import (
	"github.com/shopspring/decimal"
	"time"
)

type SEND_CHECK_TYPE int

const (
	SEND_TYPE_NONE = iota
	SEND_TYPE_CHECK_LOGIN
	SEND_TYPE_CHECK_UN_REGISTER
	SEND_TYPE_CHECK_REGISTER
)

type SMSTemplate struct {
	ID                      int64           `json:"-" gorm:"primary_key"`
	Name                    string          `json:"name"`
	Body                    string          `json:"body"`
	Language                string          `json:"language"`
	TemplateID              int             `json:"template_id"`
	TemplateInternationalID int             `json:"template_international_id"`
	Service                 string          `json:"service"`
	SendType                SEND_CHECK_TYPE `json:"send_type" gorm:"default:0"`
	CreatedAt               *time.Time      `json:"-"`
	UpdatedAt               *time.Time      `json:"-"`
}

func (*SMSTemplate) TableName() string {
	return "sms_templates"
}

type SMSMessage struct {
	PhoneNumer   PHONE       `json:"phone_numer"`
	Service      string      `json:"service"`
	Source       interface{} `json:"source"`
	Language     string      `json:"language"`
	TemplateID   int         `json:"template_id"`
	TemplateName string      `json:"template_name"`
}

type SMSService struct {
	ID       int64  `json:"-" gorm:"primary_key"`
	AreaCode string `json:"area_code"`
	Service  string `json:"service"`
}

const (
	SMS_SERVIEC_SENDCLOUD = "sendcloud"
	SMS_SERVICE_AWSSNS    = "awssns"
	SMS_SERVICE_NEC       = "nec"
	SMS_SERVICE_DEFAULT   = "*"
)

func (*SMSService) TableName() string {
	return "sms_services"
}

type CoinDepositNotifyVolLimit struct {
	ID       int64           `json:"-" gorm:"primary_key"`
	CoinCode string          `json:"coin_code"`
	LimitVol decimal.Decimal `json:"limit_vol" gorm:"type:decimal(36,18)"`
}
