package errors

import (
	"net/http"
)

var (
	NoSignStr  = NewError(http.StatusOK, "NO_SIGNSTR", "No signstr")
	NoSettleId = NewError(http.StatusOK, "NO_SETTLE_ID", "No settle id")
	NoTxHash   = NewError(http.StatusOK, "NO_TX_HASH", "No TxHash")
)
