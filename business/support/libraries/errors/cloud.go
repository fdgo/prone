package errors

import (
	"net/http"
)

var (
	InvalidTimestamp          = NewError(http.StatusOK, "INVALID_TIMESTAMP", "invalid timestamp")
	InvalidSignature          = NewError(http.StatusOK, "INVALID_SIGNATURE", "invalid signature")
	CloudAccountNotFound      = NewError(http.StatusOK, "CLOUD_ACCOUNT_NOT_FOUND", "cloud account not found")
	CloudAppNotFound          = NewError(http.StatusOK, "CLOUD_APP_NOT_FOUND", "cloud app not found")
	CloudTradeNotFound        = NewError(http.StatusOK, "CLOUD_TRADE_NO_NOT_FOUND", "out_trade_no not found")
	CloudTradeExisted         = NewError(http.StatusOK, "CLOUD_TRADE_NO_EXISTED", "out_trade_no already existed")
	CloudAppNotActived        = NewError(http.StatusOK, "CLOUD_APP_NOT_ACTIVED", "cloud app not actived")
	CloudAppDisabled          = NewError(http.StatusOK, "CLOUD_APP_DISABLE", "cloud app disable")
	CloudAccountCountLimit    = NewError(http.StatusOK, "CLOUD_ACCOUNT_COUNT_LIMIT", "cloud account count limit")
	OriginUIDHadExisted       = NewError(http.StatusOK, "ORIGIN_UID_HAD_EXISTED", "ORIGIN_UID_HAD_EXISTED")
	CloudAppNoPublicKey       = NewError(http.StatusOK, "CLOUD_APP_NO_PUBLICK_KEY", "cloud app no publick key")
	CloudApiKeyNotFound       = NewError(http.StatusOK, "CLOUD_API_KEY_NOT_FOUND", "cloud api key not found")
	InvalidExpiredTs          = NewError(http.StatusOK, "INVALID_EXPIRED_TS", "invalid expired ts")
	CloudApiKeyExpired        = NewError(http.StatusOK, "CLOUD_API_KEY_EXPIRED", "cloud api key expired")
	TransferVolPrecisionError = NewError(http.StatusOK, "TRANSFER_VOL_PRECISION_ERROR", "transfer vol precision error")
)
