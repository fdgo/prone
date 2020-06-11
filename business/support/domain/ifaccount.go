package domain

import "time"

// Login request
type LoginReq struct {
	Account
	Validate string `json:"validate,omitempty"`
	Nonce    int64  `json:"nonce,omitempty"`
}

type LoginResp struct {
	Account
}

// Login record
type LoginRecord struct {
	LoginTime *time.Time `json:"login_time"`
	LoginIP   string     `json:"login_ip"`
	Device    string     `json:"device"`
}

// Get user self info response
type UserMeResp struct {
	Account
}

// Create account request
type CreateAccountReq struct {
	Account
	Nonce int64 `json:"nonce,omitempty"`
}

type CreateAccountResp struct {
	Account
}

// Register account
type RegisterAccountReq struct {
	Account
	Code       string `json:"code,omitempty"`
	InviterId  int64  `json:"inviter_id,omitempty"`
	Validate   string `json:"validate,omitempty"`
	Passphrase string `json:"passphrase,omitempty"`
}

type RegtisterAccountResp struct {
	Account
	//ApiKey *ApiKey `json:"api_key,omitempty"`
}

// Active account
type ActiveAccountReq struct {
	Code string `json:"code,omitempty"`
}

type ActiveAccountResp struct {
	Account
}

type EmailVerifyCodeReq struct {
	EMAIL `json:"email"`
}

type EmailVerifyCodeResp struct {
	Code string `json:"code"`
}

// Bind address request
type BindAddressReq struct {
	Account
	Nonce int64 `json:"nonce,omitempty"`
}

// Bind address response
type BindAddressResp struct {
	Account
}

// Depost address response
type DepositAddressResp struct {
	UserDepositAddress
}

// Send bind email verify code request
type EmailVerifyReq struct {
	Account
	Nonce int64 `json:"nonce,omitempty"`
}

// Bind email request
type BindEmailReq struct {
	Account
	EmailCode string `json:"email_code,omitempty"`
	SMSCode   string `json:"sms_code,omitempty"`
	Nonce     int64  `json:"nonce,omitempty"`
}

// Bind email request
type BindPhoneReq struct {
	Account
	EmailCode string `json:"email_code,omitempty"`
	SMSCode   string `json:"sms_code,omitempty"`
	Nonce     int64  `json:"nonce,omitempty"`
}

// Withdraw request
type WithdrawReq struct {
	Settle
	VerifyCode string `json:"verify_code,omitempty"`
	EmailCode  string `json:"email_code,omitempty"`
	SMSCode    string `json:"sms_code,omitempty"`
	GACode     uint32 `json:"ga_code,omitempty"`
	Nonce      int64  `json:"nonce,omitempty"`
}

// Reset password request
type ResetPasswordReq struct {
	Account
	Code   string `json:"code,omitempty"`
	GACode uint32 `json:"ga_code,omitempty"`
}

// Check Account exist response
type CheckAccountExistResp struct {
	Exist bool `json:"exist"`
}

// Get settles response
type SettlesResp struct {
	Total   int64    `json:"total"`
	Settles []Settle `json:"settles,omitempty" gorm:"-"`
}

//get depositaddress resp
type AddressResp struct {
	UserDepositAddress
	Memo    string `json:"memo"`
	Account string `json:"account"`
}

type ResetAssetPasswordReq struct {
	Account
	EmailCode string `json:"email_code,omitempty"`
	SMSCode   string `json:"sms_code,omitempty"`
	GACode    uint32 `json:"ga_code"`
}

// Get google authenticator key response
type GetGAKeyResp struct {
	GAKey     string `json:"ga_key"`
	LoginName string `json:"login_name"`
}

// Set google authenticator key request
type SetGAKeyReq struct {
	GACode uint32 `json:"ga_code"`
}

// Delete google authenticator key request
type DeleteGAKeyReq struct {
	GACode uint32 `json:"ga_code"`
}
