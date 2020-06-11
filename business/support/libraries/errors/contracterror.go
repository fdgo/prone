package errors

import "net/http"

const CONTRACT_OFFLINE string = "CONTRACT_OFFLINE"
const CONTRACT_PAUSE string = "CONTRACT_PAUSE"
const CONTRACT_ORDER_WILL_LIQUIDATE string = "CONTRACT_ORDER_WILL_LIQUIDATE"
const CONTRACT_ORDER_HAVE_REVERSE_WAY string = "CONTRACT_ORDER_HAVE_REVERSE_WAY"
const CONTRACT_POSITION_IS_CLOSED string = "CONTRACT_POSITION_IS_CLOSED"
const CONTRACT_POSITION_IS_DELEGATING string = "CONTRACT_POSITION_IS_DELEGATING"
const CONTRACT_POSITION_VOLUME_NOT_ENOUGH string = "CONTRACT_POSITION_VOLUME_NOT_ENOUGH"
const CONTRACT_POSITION_NOT_EXSIT string = "CONTRACT_POSITION_NOT_EXSIT"
const CONTRACT_POSITION_NOT_ISOLATED string = "CONTRACT_POSITION_NOT_ISOLATED"
const CONTRACT_POSITION_LIQUIDATE_WHEN_SUB_MARGIN string = "CONTRACT_POSITION_LIQUIDATE_WHEN_SUB_MARGIN"
const CONTRACT_POSITION_BASE_LIMIT_SUB_MARGIN string = "CONTRACT_POSITION_BASE_LIMIT_SUB_MARGIN"
const CONTRACT_POSITION_WARN_WHEN_SUB_MARGIN string = "CONTRACT_POSITION_WARN_WHEN_SUB_MARGIN"
const CONTRACT_USER_ASSET_IN_DELEGATE string = "CONTRACT_USER_ASSET_IN_DELEGATE"
const CONTRACT_USER_ASSET_NOT_ENOUGH string = "CONTRACT_USER_ASSET_NOT_ENOUGH"
const CONTRACT_PLAN_ORDER_COUNT_IS_FULL string = "CONTRACT_PLAN_ORDER_COUNT_IS_FULL"
const CONTRACT_ORDER_LEVERAGE_TOO_LARGE string = "CONTRACT_ORDER_LEVERAGE_TOO_LARGE"
const CONTRACT_ORDER_LEVERAGE_TOO_SMALL string = "CONTRACT_ORDER_LEVERAGE_TOO_SMALL"
const CONTRACT_ORDER_PRICE_WOULD_NOT_TRIGGER string = "CONTRACT_ORDER_PRICE_WOULD_NOT_TRIGGER"
const CONTRACT_PLAN_ORDER_LIFECYCLE_TOO_LONG string = "CONTRACT_PLAN_ORDER_LIFECYCLE_TOO_LONG"
const CONTRACT_PLAN_ORDER_LIFECYCLE_TOO_SHORT string = "CONTRACT_PLAN_ORDER_LIFECYCLE_TOO_SHORT"

var (
	ContractOfflineError                   = NewError(http.StatusOK, CONTRACT_OFFLINE, "This contract is offline!")
	ContractPauseError                     = NewError(http.StatusOK, CONTRACT_PAUSE, "This contract's exchange has been paused!")
	ContractOrderWillLiquidateError        = NewError(http.StatusOK, CONTRACT_ORDER_WILL_LIQUIDATE, "This order would trigger user position liquidate!")
	ContractOrderHaveReverseWayError       = NewError(http.StatusOK, CONTRACT_ORDER_HAVE_REVERSE_WAY, "It is not possible to open and close simultaneously in the same position!")
	ContractPositionIsClosedError          = NewError(http.StatusOK, CONTRACT_POSITION_IS_CLOSED, "Your position is closed!")
	ContractPositionIsDelegatingError      = NewError(http.StatusOK, CONTRACT_POSITION_IS_DELEGATING, "Your position is in liquidation delegating!")
	ContractPositionVolumeNotEnoughError   = NewError(http.StatusOK, CONTRACT_POSITION_VOLUME_NOT_ENOUGH, "Your position  volume is not enough!")
	ContractPositionNotExsit               = NewError(http.StatusOK, CONTRACT_POSITION_NOT_EXSIT, "The position is not exsit")
	ContractPositionNotIsolated            = NewError(http.StatusOK, CONTRACT_POSITION_NOT_ISOLATED, "The position is not isolated")
	ContractPositionLiquidateWhenSubMargin = NewError(http.StatusOK, CONTRACT_POSITION_LIQUIDATE_WHEN_SUB_MARGIN, "The position would liquidate when sub margin")
	ContractPositionWarnWhenSubMargin      = NewError(http.StatusOK, CONTRACT_POSITION_WARN_WHEN_SUB_MARGIN, "The position would be warnning of liquidation when sub margin")
	ContractPositionBaseLimitWhenSubMargin = NewError(http.StatusOK, CONTRACT_POSITION_BASE_LIMIT_SUB_MARGIN, "The position’s margin shouldn’t be lower than the base limit")
	ContractUserAssetInDelegate            = NewError(http.StatusOK, CONTRACT_USER_ASSET_IN_DELEGATE, "You cross margin position is in liquidation delegating.")
	ContractUserAssetNotEnough             = NewError(http.StatusOK, CONTRACT_USER_ASSET_NOT_ENOUGH, "You contract account available balance not enough.")
	// 您计划委托的数量已经达到系统的设置
	ContractPlanOrderCountIsFull = NewError(http.StatusOK, CONTRACT_PLAN_ORDER_COUNT_IS_FULL, "Your plan order's count is more than system maximum limit.")
	// 订单的杠杆大小超过系统设置
	ContractOrderLeverage2Large = NewError(http.StatusOK, CONTRACT_ORDER_LEVERAGE_TOO_LARGE, "The order's leverage is too large.")
	// 订单的杠杆小于系统设置
	ContractOrderLeverage2Small = NewError(http.StatusOK, CONTRACT_ORDER_LEVERAGE_TOO_SMALL, "The order's leverage is too small.")
	// 触发价格与当前价格偏差过大
	ContractOrderWouldNotTrigger = NewError(http.StatusOK, CONTRACT_ORDER_PRICE_WOULD_NOT_TRIGGER, "The deviation between current price and trigger price is too large.")
	// 订单的委托周期超过系统设置
	ContractPlanOrderLifecycle2Long = NewError(http.StatusOK, CONTRACT_PLAN_ORDER_LIFECYCLE_TOO_LONG, "The plan order's life cycle is too long.")
	// 订单的委托周期小于系统设置
	ContractPlanOrderLifecycle2Short = NewError(http.StatusOK, CONTRACT_PLAN_ORDER_LIFECYCLE_TOO_SHORT, "The plan order's life cycle is too short.")
)
