package domain

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type MessageBase struct {
	Action        string     `json:"action,omitempty" gorm:"column:base_action"`      // Message由什么操作产生
	Type          string     `json:"type,omitempty" gorm:"column:base_type"`          // Message类型
	From          string     `json:"from,omitempty" gorm:"column:base_from"`          // 来源(User/Manager/ServiceName)
	FromType      string     `json:"from_type,omitempty" gorm:"column:base_from_type` // 来源类型(User/Manager/System)
	ReceiptHandle string     `json:"-" gorm:"-"`
	LogTime       *time.Time `json:"log_time"`
}

func (m MessageBase) String() string {
	return fmt.Sprintf(`{Action:%s,Type:%s,From:%s,FromType:%s,ReceiptHandle:%s,LogTime:%s`,
		m.Action, m.Type, m.From, m.FromType,
		m.ReceiptHandle, m.LogTime)
}

// Action 类型
const (
	MESSAGE_ACTION_CREATED                             = "CREATED"
	MESSAGE_ACTION_APPROVED                            = "APPROVED"
	MESSAGE_ACTION_PASSED                              = "PASSED"
	MESSAGE_ACTION_REJECTED                            = "REJECTED"
	MESSAGE_ACTION_SIGNED                              = "SIGNED"
	MESSAGE_ACTION_BROADCASTED                         = "BROADCASTED"
	MESSAGE_ACTION_SUCCESSED                           = "SUCCESSED"
	MESSAGE_ACTION_CANCELED                            = "CANCELED"
	MESSAGE_ACTION_FAILED                              = "FAILED"
	MESSAGE_ACTION_UNFREEZE_FUNDS                      = "UNFREEZE_FUNDS"
	MESSAGE_ACTION_TRADED                              = "TRADED"
	MESSAGE_ACTION_REGISTER                            = "REGISTER"
	MESSAGE_ACTION_LOGIN                               = "LOGIN"
	MESSAGE_ACTION_RESETPASSWD                         = "RESET_PASSWORD"
	MESSAGE_ACTION_ASSET_PASSWORD_ADD                  = "ADD_ASSET_PASSWORD"
	MESSAGE_ACTION_ASSET_PASSWORD_RESET                = "RESET_ASSET_PASSWORD"
	MESSAGE_ACTION_ASSET_PASSWORD_EFFECTIVE_TIME_RESET = "RESET_ASSET_PASSWORD_EFFECTIVE_TIME"
	MESSAGE_ACTION_GAKEY_ADD                           = "GAKEY_ADD"
	MESSAGE_ACTION_GAKEY_DELETE                        = "GAKEY_DELETE"
	MESSAGE_ACTION_ANTI_FISHING_TEXT_SET               = "ANTGI_FISHIING_TEXT_SET"
	MESSAGE_ACTION_ACTIVE_ACCOUNT                      = "ACCOUNT_ACTIVE"
	MESSAGE_ACTION_BIND_EMAIL                          = "BIND_EMAIL"
	MESSAGE_ACTION_BIND_PHONE                          = "BIND_PHONE"
	MESSAGE_ACTION_KYC_UPDATED                         = "KYC_UPDATED"
	MESSAGE_ACTION_KYC_SUBMITED                        = "KYC_SUBMITED"
	MESSAGE_ACTION_KYC_REJECTED                        = "KYC_REJECTED"
	MESSAGE_ACTION_KYC_APPROVED                        = "KYC_APPROVED"
	MESSAGE_ACTION_CONTRACT_REGISTER_REWARD            = "CONTRACT_REGISTER_REWARD"
)

// Message 类型
const (
	MESSAGE_TYPE_ORDER          = "ORDER"
	MESSAGE_TYPE_DEPOSIT        = "DEPOSIT"
	MESSAGE_TYPE_WITHDRAW       = "WITHDRAW"
	MESSAGE_TYPE_ACCOUNT        = "ACCOUNT"
	MESSAGE_TYPE_EMAIL          = "EMAIL"
	MESSAGE_TYPE_SMS            = "SMS"
	MESSAGE_TYPE_CONTRACT_TRADE = "CONTRACT_TRADE"
	MESSAGE_TYPE_TRANSFER_FUNDS = "TRANSFER_FUNDS"
	MESSAGE_TYPE_NOTIFY         = "NOTIFY"
	MESSAGE_TYPE_OTC_TRADE      = "OTC_TRADE"
	MESSAGE_TYPE_SPOT_TRADE     = "SPOT_TRADE"
)

const (
	FROM_TYPE_MANAGER = "MANAGER"
	FROM_TYPE_USER    = "USER"
	FROM_TYPE_SYSTEM  = "SYSTEM"
)

type Message struct {
	MessageBase
	MessageID     string         `json:"-"`
	Settle        *Settle        `json:"settle,omitempty"`
	Account       *AccountRecord `json:"account,omitempty"`
	EmailMessage  *EmailMessage  `json:"email_message,omitempty"`
	SMSMessage    *SMSMessage    `json:"sms_message,omitempty"`
	NotifyMessage *NotifyMessage `json:"notify_message,omitempty"`
	Assets        []AssetRecord  `json:"assets,omitempty"`
}

func (m *Message) String() string {
	//if nil == m {
	//	return "nil"
	//}
	//transferFundsRecords := "["
	//for i, item := range m.TransferFundsRecords {
	//	transferFundsRecords = fmt.Sprintf(`%s%d:%s,`, transferFundsRecords, i, item)
	//}
	//transferFundsRecords = fmt.Sprintf(`%s]`, transferFundsRecords)
	//orders := "["
	//for i, item := range m.Orders {
	//	orders = fmt.Sprintf(`%s%d:%s,`, orders, i, &item)
	//}
	//orders = fmt.Sprintf(`%s]`, orders)
	//assets := "["
	//for i, item := range m.Assets {
	//	assets = fmt.Sprintf(`%s%d:%s,`, assets, i, &item)
	//}
	//assets = fmt.Sprintf(`%s]`, assets)
	//contractTrades := "["
	//for i, item := range m.ContractTrades {
	//	contractTrades = fmt.Sprintf(`%s%d:%s,`, contractTrades, i, item)
	//}
	//contractTrades = fmt.Sprintf(`%s]`, contractTrades)
	//contractCashBooks := "["
	//for i, item := range m.ContractCashBooks {
	//	contractCashBooks = fmt.Sprintf(`%s%d:%s,`, contractCashBooks, i, item)
	//}
	//contractCashBooks = fmt.Sprintf(`%s]`, contractCashBooks)
	//return fmt.Sprintf(`{MessageBase:%s,MessageID:%s,Settle:%s,Trade:%s,Account:%s,TransferFundsRecords:%s,Orders:%s,Assets:%s,ContractTrades:%s,ContractCashBooks:%s}`,
	//	m.MessageBase, m.MessageID, m.Settle, m.Trade, m.Account,
	//	transferFundsRecords, orders, assets, contractTrades, contractCashBooks)
	return ""
}

// 用户资产记录
type AssetRecord struct {
	MessageBase
	SourceID           int64           `json:"source_id,omitempty"`
	FreezeVolModify    decimal.Decimal `json:"freeze_vol_motify" gorm:"type:decimal(36,18)"`
	AvailableVolModify decimal.Decimal `json:"available_vol_modify" gorm:"type:decimal(36,18)"`
	UserAssets
}

func (a *AssetRecord) String() string {
	if nil == a {
		return "nil"
	}
	return fmt.Sprintf(`{MessageBase:%s,SourceID:%d,FreezeVolModify:%s,AvailableVolModify:%s,UserAssets:%s}`,
		a.MessageBase, a.SourceID, a.FreezeVolModify, a.AvailableVolModify, a.UserAssets.Desc())
}

func (*AssetRecord) TableName() string {
	return "asset_records"
}

// 用户订单记录
type OrderRecord struct {
	MessageBase
	TradeID int64
	//Order
}

func (*OrderRecord) TableName() string {
	return "order_records"
}

// 用户充值/提现记录
type SettleRecord struct {
	MessageBase
	Settle
}

func (*SettleRecord) TableName() string {
	return "settle_records"
}

// AccountRecord 用户账户记录
type AccountRecord struct {
	MessageBase
	Account
	Device    string  `json:"device,omitempty"`
	ClientIP  string  `json:"client_ip,omitempty"`
	InviterId int64   `json:"inviter_id,omitempty"`
	MarkCode  *string `json:"mark_code,omitempty"` // 渠道带过来的用户的标记
}

func (a *AccountRecord) String() string {
	if nil == a {
		return "nil"
	}
	return fmt.Sprintf(`{MessageBase:%s,Account:%#v,Device:%s,ClientIP:%s,InviterId:%d}`,
		a.MessageBase, a.Account, a.Device, a.ClientIP, a.InviterId)
}

func (*AccountRecord) TableName() string {
	return "account_records"
}

type NotifyMessage struct {
	AccountID    int64       `json:"account_id,omitempty"`
	TemplateName string      `json:"template_name"`
	Source       interface{} `json:"source"`
}
