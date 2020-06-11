package errors

import "net/http"

const STOCK_OFFLINE string = "STOCK_OFFLINE"
const STOCK_PAUSE string = "STOCK_PAUSE"

var (
	StockOfflineError = NewError(http.StatusOK, STOCK_OFFLINE, "This stock is offline!")
	StockPauseError   = NewError(http.StatusOK, STOCK_PAUSE, "This stock's exchange has been paused!")
)
