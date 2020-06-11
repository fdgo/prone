package domain

import (
	"fmt"
	"time"
)

type RedisKey struct {
	Key     string
	Timeout time.Duration
	Size    int64
}

// Bank
var (
	KMaxSettledBlockNum   = RedisKey{"Bank.MAX_SETTLE_BLOCK_NUM", time.Second * 0, 0}
	KLatestAttentionID    = RedisKey{"Bank.LATEST_ATTENTION_ID", time.Second * 0, 0}
	KWaitCheckTransaction = RedisKey{"Bank.WAIT_CHECK_TRANSACTIONS", time.Second * 0, 0}
	KWaitApproveSettles   = RedisKey{"Bank.WAIT_APPROVE_SETTLES", time.Second * 0, 0}
	KSupportToken         = RedisKey{"Bank.SUPPORT_TOKEN", time.Hour, 0}
	KConfig               = RedisKey{"Bank.CONFIG", time.Second * 0, 0}
)

// stock
var (
	KSpotStocks      = RedisKey{"SPOTCENTER.STOCKS", 0, 0}
	KWSSpotDepth     = RedisKey{"SPOTCENTER.WS_DEPTH", 0, 0}
	KWSSpotTick      = RedisKey{"SPOTCENTER.WS_TICKER", 0, 0}
	KWSSpotTrade     = RedisKey{"SPOTCENTER.WS_TRADES", 0, 0}
	KWSSpotAliveUser = RedisKey{"SPOTCENTER.ALIVE_USER", time.Minute * 5, 0}
	KWSSpotUnicast   = RedisKey{"SPOTCENTER.UNICAST", 0, 0}
)

func GetSpotDepthBuyRedisKey(stockCode string) RedisKey {
	key := fmt.Sprintf("SPOTCENTER.DEPTH_BUY_%s", stockCode)
	return RedisKey{key, 0, 0}
}

func GetSpotDepthSellRedisKey(stockCode string) RedisKey {
	key := fmt.Sprintf("SPOTCENTER.DEPTH_SELL_%s", stockCode)
	return RedisKey{key, 0, 0}
}

func GetDepositAddresssKey(coinGroup string) RedisKey {
	key := fmt.Sprintf("Bank.%s_DEPOSIT_ADDRESSS", coinGroup)
	return RedisKey{key, time.Second * 0, 0}
}
func GetSettledBlockNumKey(coinGroup string) RedisKey {
	key := fmt.Sprintf("Bank.MAX_%s_BLOCK_NUM", coinGroup)
	return RedisKey{key, time.Second * 0, 0}
}

func GetSignedWithdrawsKey(coinGroup string) RedisKey {
	key := fmt.Sprintf("Bank.%s_SIGNED", coinGroup)
	return RedisKey{key, time.Second * 0, 0}
}

func GetWithdrawsKey(coinGroup string) RedisKey {
	key := fmt.Sprintf("Bank.%s_WITHDRAWS", coinGroup)
	return RedisKey{key, time.Second * 0, 0}
}

func GetSystemAddressKey(coinGroup string) RedisKey {
	key := fmt.Sprintf("Bank.%s_SYSTEM_ADDRESS", coinGroup)
	return RedisKey{key, time.Hour, 0}
}

func GetWalletAddressKey(coinGroup string) RedisKey {
	key := fmt.Sprintf("Bank.%s_WALLET_ADDRESS", coinGroup)
	return RedisKey{key, time.Hour, 0}
}

// notifycenter
var (
	KWaitSendEmails  = RedisKey{"NOTIFYCENTER.WAIT_SEND_EMAILS", time.Second * 0, 0}
	KEmailServies    = RedisKey{"NOTIFYCENTER.EMAIL_SERVICES", time.Minute * 5, 0}
	KEmailTemplate   = RedisKey{"NOTIFYCENTER.EMAIL_TEMPLATE", time.Hour * 2, 0}
	KSMSServies      = RedisKey{"NOTIFYCENTER.SMS_SERVICES", time.Minute * 5, 0}
	KSMSTemplate     = RedisKey{"NOTIFYCENTER.SMS_TEMPLATE", time.Hour * 2, 0}
	KDepositVolLimit = RedisKey{"NOTIFYCENTER.DEPOSIT_VOL_LIMIT", time.Hour * 2, 0}
)

// usercenter
var (
	KFreeDepositAddress = RedisKey{"USERCENTER.FREE_DEPOSIT_ADDRESS", time.Second * 0, 0}
)

// ifaccount
var (
	KVerifyCode      = RedisKey{"IFACCOUNT.VERIFYCODE", time.Minute * 10, 0}
	KVerifyCodeTimes = RedisKey{"IFACCOUNT.VERIFYCODE_TIMES", time.Second * 58, 0}
	KLoginRetryTimes = RedisKey{"IFACCOUNT.LOGIN_RETRY_TIMES", time.Minute * 10, 0}
	KDisableLogin    = RedisKey{"IFACCOUNT.DISABLE_LOGIN", time.Hour * 8, 0}
	KGAKey           = RedisKey{"IFACCOUNT.GAKEY", time.Minute * 15, 0}
	KWithdrawConfig  = RedisKey{"IFACCOUNT.WITHDRAW_CONFIG", time.Hour * 2, 0}
	KCaptchLimit     = RedisKey{"IFACCOUNT.CAPTCH_LIMIT", time.Hour * 2, 0}
)

// quotes
var (
	KQuotesExchangeRate    = RedisKey{"QUOTES.EXCHANGE_RATE", 0, 0}
	KQuotesCoinMarketPrice = RedisKey{"QUOTES.COIN_MARKET_PRICE", time.Second * 60 * 30, 0}
)

// ifmarkets
var (
	KSpotCoins        = RedisKey{"SPOTS.COINS", time.Second * 30, 0}
	KSpotCoinicons    = RedisKey{"SPOTS.COIN_ICONS", time.Second * 30, 0}
	KSpotCoinRateBase = RedisKey{"SPOTS.COIN_RATE_BASE", time.Second * 30, 0}
)

// backstage
var (
	KDashboardState      = RedisKey{"BACKSTAGE.DASHBOARD_STATE", time.Minute * 15, 0}
	KStockState          = RedisKey{"BACKSTAGE.STOCK_COUNT", time.Hour, 0}
	KCoinState           = RedisKey{"BACKSTAGE.COIN_COUNT", time.Hour, 0}
	KContractState       = RedisKey{"BACKSTAGE.CONTRACT_COUNT", time.Hour, 0}
	KContractTradeStatis = RedisKey{"BACKSTAGE.CONTRACT_TRADE_NAME", time.Hour, 0}
)

// ifglobal
var (
	KAreas = RedisKey{"IFGLOBAL.AREAS", time.Hour * 4, 0}
)

// contract
var (
	KContracts       = RedisKey{"CONTRACTCENTER.CONTRACTS", 0, 0}
	KContractRate    = RedisKey{"CONTRACTCENTER.RATE", 0, 0}
	KQuotesTimeUnit  = RedisKey{"CONFIG.QUOTES_TIME_UNIT", time.Hour * 4, 0}
	KWSContractTrade = RedisKey{"CONTRACTCENTER.WS_TRADES", 0, 0}
	KWSPNLS          = RedisKey{"CONTRACTCENTER.WS_PNLS", 0, 0}
	KWSContractDepth = RedisKey{"CONTRACTCENTER.WS_DEPTH", 0, 0}
	KWSContractTick  = RedisKey{"CONTRACTCENTER.WS_TICKERS", 0, 0}
	KWSRateLimit     = RedisKey{"CONTRACTCENTER.RATE_LIMIT", time.Hour, 0}
	KWSAliveUser     = RedisKey{"CONTRACTCENTER.ALIVE_USER", time.Minute * 5, 0}
	KWSUnicast       = RedisKey{"CONTRACTCENTER.UNICAST", 0, 0}
	KWSPLIQB         = RedisKey{"CONTRACTCENTER.WS_PLIQB", 0, 0}
)

// OTCCenter
var (
	KCancelOTCOrderLimit = RedisKey{"OTCCENTER.CANCEL_OTC_ORDER_LIMIT", time.Hour * 24, 0}
)

// CloudCenter
var (
	KCloudTransferSettle = RedisKey{"CLOUDCENTER.TRANSFER_SETTLE_LIST", 0, 0}
	KCloudApp            = RedisKey{"IFCLOUD.CLOUD_APP", time.Hour * 2, 0}
	KCloudApiKey         = RedisKey{"AUTH.CLOUD_API_KEY", time.Hour * 2, 0}
)

// 通用的
var (
	KDelayWhiteAccount = RedisKey{"COMMON.DELAY_WHITE_ACCOUNT", 0, 0}
	KDelayBlackAccount = RedisKey{"COMMON.DELAY_BLACK_ACCOUNT", 0, 0}
	KDelayConfig       = RedisKey{"COMMON.DELAY_CONFIG", 0, 0}
)

//func GetOrderMatchingStatusRedisKey(stockCode string) RedisKey {
//	sc := NewCoinStockCode(stockCode)
//	key := fmt.Sprintf("ORDERCENTER.ORDER_MATCHING_STATUS_%s_%s",
//		sc.LCode, sc.RCode)
//	return RedisKey{key, time.Second * 0, 0}
//}
//
//func GetOrderUserOperationLimitRedisKey(accountId int64) RedisKey {
//	key := fmt.Sprintf("BASE.ACCOUNT_%d_OPER_LIMIT", accountId)
//	return RedisKey{key, time.Second * 5, 0}
//}
//
//func GetUserStocksRedisKey() RedisKey {
//	return RedisKey{"ORDERCENTER.USER_STOCKS", 0, 0}
//}
//
//func GetTradesRedisKey(stockCode string) RedisKey {
//	sc := NewCoinStockCode(stockCode)
//	if nil == sc {
//		return RedisKey{"", 0, 0}
//	}
//	key := fmt.Sprintf("ORDERCENTER.%s_%s_TRADES", sc.LCode, sc.RCode)
//	return RedisKey{key, 0, 360}
//}
//
//func GetOrdersRedisKey(stockCode string) RedisKey {
//	sc := NewCoinStockCode(stockCode)
//	if nil == sc {
//		return RedisKey{"", 0, 0}
//	}
//	key := fmt.Sprintf("ORDERCENTER.%s_%s_ORDERS", sc.LCode, sc.RCode)
//	return RedisKey{key, 0, 360}
//}

//func GetStockRealtimeTickerRedisKey(stockCode string) RedisKey {
//	sc := NewCoinStockCode(stockCode)
//	key := fmt.Sprintf("%s_%s_STOCK_REAL_TIME_TICKER", sc.LCode, sc.RCode)
//	return RedisKey{key, time.Second * 0, 0}
//}
//
////func GetStock24HourTickerRedisKey(stockCode string) RedisKey {
////	sc := NewCoinStockCode(stockCode)
////	key := fmt.Sprintf("%s_%s_STOCK.24_HOUR_TICKER", sc.LCode, sc.RCode)
////	return RedisKey{key, time.Second * 0, 0}
////}
////
////func GetStockYesterdayRealtimeTickerRedisKey(stockCode string) RedisKey {
////	sc := NewCoinStockCode(stockCode)
////	key := fmt.Sprintf("%s_%s_STOCK_YESTERDAY_REAL_TIME_TICKER", sc.LCode, sc.RCode)
////	return RedisKey{key, time.Minute * 7, 0}
////}
////
////func GetQuoteRealtimeTickerRedisKey(coinCode string) RedisKey {
////	key := fmt.Sprintf("%s_REAL_TIME_PRICE", coinCode)
////	return RedisKey{key, time.Second * 0, 0}
////}
//
//func GetQuoteMarketRedisKey(coinCode string) RedisKey {
//	key := fmt.Sprintf("QUOTES.%s_MARKET", coinCode)
//	return RedisKey{key, time.Second * 0, 0}
//}
//
//// supply
//func GetQuoteCoinSupplyRedisKey(coinCode string) RedisKey {
//	key := fmt.Sprintf("QUOTES.%s_SUPPLY_VOLUME", coinCode)
//	return RedisKey{key, time.Hour * 12, 0}
//}
//
func GetRechargeAmountRedisKey(coinCode string) RedisKey {
	key := fmt.Sprintf("IFACCOUNT.RECHARGE_AMOUNT_%s", coinCode)
	return RedisKey{key, 0, 0}
}

//
//// GetBlockBalanceRedisKey 链上资产
//func GetBlockBalanceRedisKey(coinCode string) RedisKey {
//	key := fmt.Sprintf("BACKSTAGE.BLOCK_BALANCE_%s", coinCode)
//	return RedisKey{key, 0, 0}
//}
//
//func GetApiLimitOfUnitTimeKey(accountID int64) RedisKey {
//	key := fmt.Sprintf("BASE.API_LIMIT_OF_UNIT_TIME.%d", accountID)
//	return RedisKey{key, time.Second * 20, 0}
//}
//
//// GetCoinBriefRedisKey 货币描述
//func GetCoinBriefRedisKey(coinCode string) RedisKey {
//	key := fmt.Sprintf("IFGLOBAL.COIN_BRIEF_%s", coinCode)
//	return RedisKey{key, time.Hour, 0}
//}
//
//// GetCoinBriefRedisKey 货币描述
//func GetCoinRedisKey(coinCode string) RedisKey {
//	key := fmt.Sprintf("IFGLOBAL.COIN_%s", coinCode)
//	return RedisKey{key, time.Hour, 0}
//}
//
//// GetContractMatchingStatusRedisKey 合约的撮合状态
//func GetContractMatchingStatusRedisKey(contractId int64) RedisKey {
//	key := fmt.Sprintf("CONTRACTCENTER.ORDER_MATCHING_STATUS_%d",
//		contractId)
//	return RedisKey{key, time.Second * 16, 0}
//}
//
//// GetContractPNLRedisKey 合约的PNL
//func GetContractPNLRedisKey(contractId int64, positionType POSITION_TYPE) RedisKey {
//	key := fmt.Sprintf("CONTRACTCENTER.PNL_%d_%d",
//		contractId,
//		positionType)
//	return RedisKey{key, time.Second * 0, 0}
//}
//
//func GetUserContractsRedisKey() RedisKey {
//	return RedisKey{"CONTRACTCENTER.USER_CONTRACTS", 0, 0}
//}
//
//func GetContractRealtimeTickerRedisKey(contractID int64) RedisKey {
//	key := fmt.Sprintf("CONTRACTCENTER.REAL_TIME_TICKER_%d", contractID)
//	return RedisKey{key, time.Second * 0, 0}
//}
//
//func GetContract24HourTickerRedisKey(contractID int64) RedisKey {
//	key := fmt.Sprintf("CONTRACTCENTER.24_HOUR_TICKER_%d", contractID)
//	return RedisKey{key, time.Second * 0, 0}
//}
//
//func GetContractHistoryTradesRedisKey(contractID int64) RedisKey {
//	key := fmt.Sprintf("CONTRACTCENTER.TRADES_HISTORY_%d", contractID)
//	return RedisKey{key, 0, 660}
//}
//
//func GetFundingRateRedisKey(contractID int64, ts int64) RedisKey {
//	key := fmt.Sprintf("CONTRACTCENTER.RATE_%d_%d", contractID, ts)
//	return RedisKey{key, 0, 0}
//}
//
//func GetContractPriceRedisKey(contractID int64) RedisKey {
//	key := fmt.Sprintf("CONTRACTCENTER.PRICE_%d", contractID)
//	return RedisKey{key, time.Hour * 24, 0}
//}
//
//func GetContractTradesRedisKey(contractID int64) RedisKey {
//	key := fmt.Sprintf("CONTRACTCENTER.TRADES_%d", contractID)
//	return RedisKey{key, 0, 0}
//}
//
//func GetContractDepthRedisKey(contractID int64) RedisKey {
//	key := fmt.Sprintf("CONTRACTCENTER.DEPTH_%d", contractID)
//	return RedisKey{key, 0, 0}
//}
//
//func GetContractDepthBuyRedisKey(contractID int64) RedisKey {
//	key := fmt.Sprintf("CONTRACTCENTER.DEPTH_BUY_%d", contractID)
//	return RedisKey{key, 0, 0}
//}
//
//func GetContractDepthSellRedisKey(contractID int64) RedisKey {
//	key := fmt.Sprintf("CONTRACTCENTER.DEPTH_SELL_%d", contractID)
//	return RedisKey{key, 0, 0}
//}
//
//func GetContractPositionSizeRedisKey(contractID int64) RedisKey {
//	key := fmt.Sprintf("CONTRACTCENTER.POSITION_SIZE_%d", contractID)
//	return RedisKey{key, time.Hour * 48, 0}
//}
//
//func GetContractPlayPositionSizeRedisKey(contractID int64) RedisKey {
//	key := fmt.Sprintf("CONTRACTCENTER.PLAY_POSITION_SIZE_%d", contractID)
//	return RedisKey{key, time.Hour * 48, 0}
//}
//
//func GetContractPlanOrderConfigRedisKey() RedisKey {
//	return RedisKey{"CONTRACTCENTER.PLAN_ORDER_CONFIGS", 0, 0}
//}
//
//func GetContractMarketOrderConfigRedisKey(contractID int64) RedisKey {
//	key := fmt.Sprintf("CONTRACTCENTER.MARKET_ORDER_CONFIG_%d", contractID)
//	return RedisKey{key, 0, 0}
//}
//
//func GetStatisTaskTopRedisKey(coincode string) RedisKey {
//	key := fmt.Sprintf("STATISTASK.TOP.%s", coincode)
//
//	return RedisKey{key, time.Hour * 4, 0}
//}
//
//func GetQDLimitOfUnitTimeKey(qduid int64) RedisKey {
//	key := fmt.Sprintf("BASE.QD_LIMIT_OF_UNIT_TIME.%d", qduid)
//	return RedisKey{key, time.Minute * 10, 0}
//}
//
//func GetTradeStatisticsLimitOfUnitTimeKey(accountID int64, contractID int64) RedisKey {
//	key := fmt.Sprintf("CONTRACTCENTER.TRADE_RIVAL_FEE_LIMIT_OF_UNIT_TIME.%d", accountID)
//	return RedisKey{key, time.Minute * 6, 0}
//}
//
func GetContractAddressKey(coin string) RedisKey {
	key := fmt.Sprintf("Bank.CONTRACTADDRESS.GROUP_%s", coin)
	return RedisKey{key, 0, 0}
}

//
//func GetContractPermissionCacheKey(contractID int64) RedisKey {
//	key := fmt.Sprintf("API_PERMISSION.CONTRACT_%d", contractID)
//	return RedisKey{key, 0, 0}
//}
