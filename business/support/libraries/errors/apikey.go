package errors

import "net/http"

var (
	NoPassphrase         = NewError(http.StatusOK, "NO_PASSPHRASE", "No passphrase")
	IncorretPassphrase   = NewError(http.StatusOK, "INCORRET_PASSPHRASE", "Incorret passphrase")
	InvalidIPFormat      = NewError(http.StatusOK, "InvalidIPFormat", "invallid ip format")
	NoExpireTime         = NewError(http.StatusOK, "NO_EXPIRE_TIME", "No expire time")
	ApiKeyNotFound       = NewError(http.StatusOK, "API_KEY_NOT_FOUND", "Api key not found")
	ApiKeyBindIPMaxLimit = NewError(http.StatusOK, "API_KEY_BIND_IP_MAX_LIMIT", "Api key bind ip max limit")
	NoAccessKey          = NewError(http.StatusOK, "NO_ACCESS_KEY", "No access key")
)
