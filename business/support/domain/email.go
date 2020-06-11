package domain

import (
	"time"
)

type Email struct {
	Sender    string `json:"sender"`
	Service   string `json:"services"`
	Recipient EMAIL  `json:"recipient"`
	Subject   string `json:"subject"`
	HtmlBody  string `json:"html_body"`
	TextBody  string `json:"text_body"`
	CharSet   string `json:"-"`
}

type EmailMessage struct {
	Sender       string      `json:"sender"`
	Recipient    EMAIL       `json:"recipient"`
	Subject      string      `json:"subject"`
	TemplateName string      `json:"template_name"`
	Language     string      `json:"language"`
	Source       interface{} `json:"source"`
}

type EmailTemplate struct {
	ID        int64           `json:"-" gorm:"primary_key"`
	Name      string          `json:"name"`
	Subject   string          `json:"subject"`
	Body      string          `json:"body"`
	SendType  SEND_CHECK_TYPE `json:"send_type" gorm:"default:0"`
	CreatedAt *time.Time      `json:"-"`
	UpdatedAt *time.Time      `json:"-"`
}

type EmailService struct {
	ID      int64  `json:"-" gorm:"primary_key"`
	Domain  string `json:"domain"`
	Service string `json:"service"`
}

const (
	EMAIL_SERVIEC_SENDCLOUD = "sendcloud"
	EMAIL_SERVICE_AWSSES    = "awsses"
	EMAIL_SERVICE_DEFAULT   = "*"
)

func (*EmailService) TableName() string {
	return "email_services"
}

const (
	VERIFYCODE_TYPE_ACTIVE_ACCOUNT       = "ActiveVerifyCode"
	VERIFYCODE_TYPE_RESET_PASSWORD       = "ResetPasswordVerifyCode"
	VERIFYCODE_TYPE_WITHDRAW             = "WithdrawVerifyCode"
	VERIFYCODE_TYPE_REGISTER             = "RegisterVerifyCode"
	VERIFYCODE_TYPE_BIND_PHONE           = "BindPhoneVerifyCode" // bind email or phone number
	VERIFYCODE_TYPE_BIND_EMAIL           = "BindEmailVerifyCode" // bind email or phone number
	VERIFYCODE_TYPE_RESET_ASSET_PASSWORD = "ResetAssetPasswordVerifyCode"
	VERIFYCODE_TYPE_API_KEY              = "APIKeyVerifyCode"
)
