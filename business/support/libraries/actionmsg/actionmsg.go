package actionmsg

import (
	"business/support/domain"
	"business/support/libraries/msg"

	"github.com/shopspring/decimal"
)

type ACTION_MSG_ID int32

const (
	ACTION_MSG_ID_UNKOWN  ACTION_MSG_ID = iota
	ACTION_MSG_ID_REGIST                // 注册
	ACTION_MSG_ID_DEPOSIT               // 充值
	ACTION_MSG_ID_TRADE                 // 交易
)

type RegistMsg struct {
	AccountId int64 `json:"account_id,omitempty"`
	InviterId int64 `json:"inviter_id,omitempty"`
}

type DepositMsg struct {
	SettleId  int64           `json:"settle_id,omitempty"`
	AccountId int64           `json:"account_id,omitempty"`
	CoinCode  string          `json:"coin_code,omitempty"`
	Vol       decimal.Decimal `json:"vol,omitempty"`
}

type TradeMsg struct {
	domain.DealTradeInfo
}

func Receive(waitTimeSeconds int64) ([]*msg.Msg, error) {
	return receiver.Receive(waitTimeSeconds)
}

func DeleteMsgs(msgs []*msg.Msg) {
	receiver.DeleteMsgs(msgs)
}

func DeleteMsg(msg *msg.Msg) {
	receiver.DeleteMsg(msg)
}

func SendRegistMsg(accountId int64, inviterId int64) {
	v := RegistMsg{AccountId: accountId, InviterId: inviterId}
	if nil != sender {
		go sender.Send(int32(ACTION_MSG_ID_REGIST), &v)
	} else if nil != broadcaster {
		go broadcaster.Send(int32(ACTION_MSG_ID_REGIST), &v)
	}
}

func SendDepositMsg(
	settleId int64,
	accountId int64,
	coinCode string,
	vol decimal.Decimal) {
	v := DepositMsg{
		SettleId:  settleId,
		AccountId: accountId,
		CoinCode:  coinCode,
		Vol:       vol}
	if nil != sender {
		go sender.Send(int32(ACTION_MSG_ID_DEPOSIT), &v)
	} else if nil != broadcaster {
		go broadcaster.Send(int32(ACTION_MSG_ID_DEPOSIT), &v)
	}
}

func SendTradeMsg(
	sellOrder *domain.Order,
	buyOrder *domain.Order,
	trade *domain.Trade) {
	dealTradeInfo := domain.DealTradeInfo{
		SellOrderId:     sellOrder.OrderID,
		SellAccountId:   sellOrder.AccountID,
		SellFee:         trade.SellFee,
		SellFeeCoinCode: sellOrder.FeeCoinCode,
		BuyOrderId:      buyOrder.OrderID,
		BuyAccountId:    buyOrder.AccountID,
		BuyFee:          trade.BuyFee,
		BuyFeeCoinCode:  buyOrder.FeeCoinCode,
		TradeId:         trade.TradeID,
		StockCode:       trade.StockCode,
		Way:             trade.Way,
		Fluctuation:     trade.Fluctuation,
	}
	if nil != sender {
		go sender.Send(int32(ACTION_MSG_ID_TRADE), &dealTradeInfo)
	} else if nil != broadcaster {
		go broadcaster.Send(int32(ACTION_MSG_ID_TRADE), &dealTradeInfo)
	}
}
