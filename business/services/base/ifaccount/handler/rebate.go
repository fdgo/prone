package handler

import (
	"business/services/base/ifaccount/adaptor"
)

type RebateInfo struct {
	Rewards         []*adaptor.RebateReward `json:"rewards"`
	ChildCount      int                     `json:"child_count"`
	TradeChildCount int                     `json:"trade_child_count"`
}

type RebateResponse struct {
	Rebate *RebateInfo `json:"rebate"`
}
