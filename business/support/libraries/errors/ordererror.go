package errors

import "net/http"

const ORDER_SYS_IS_BUSY string = "ORDER_SYS_IS_BUSY"
const ORDER_USER_ASSETS_NO_ENOUGH string = "ORDER_USER_ASSETS_NO_ENOUGH"
const ORDER_ORDER_INVALID string = "ORDER_ORDER_INVALID"
const ORDER_ORDER_IS_FINISHED string = "ORDER_ORDER_IS_FINISHED" // ToDo
const ORDER_FEE_COIN_RATIO_NOT_FOUND string = "ORDER_FEE_COIN_RATIO_NOT_FOUND"
const ORDER_ACCOUNT_NO_ACTIVATED string = "ORDER_ACCOUNT_NO_ACTIVATED"

var (
	OrderSysBusy                 = NewError(http.StatusOK, ORDER_SYS_IS_BUSY, "Order system busy")
	OrderUserAssetsError         = NewError(http.StatusOK, ORDER_USER_ASSETS_NO_ENOUGH, "Lack of balance")
	OrderInvalidError            = NewError(http.StatusOK, ORDER_ORDER_INVALID, "Invalid order")
	OrderIsFinishedError         = NewError(http.StatusOK, ORDER_ORDER_IS_FINISHED, "order is finished")
	OrderFeeCoinRatioError       = NewError(http.StatusOK, ORDER_FEE_COIN_RATIO_NOT_FOUND, "Invalid code of order fee")
	OrderAccountNoActivatedError = NewError(http.StatusOK, ORDER_ACCOUNT_NO_ACTIVATED, "Account not activated")
)
