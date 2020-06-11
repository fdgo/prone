package errors

import (
	"net/http"
)

var (
	AuthNoDev                  = NewError(http.StatusOK, "AUTH_NO_DEV", "No device info")
	AuthAccountDisable         = NewError(http.StatusOK, "AUTH_ACCOUNT_DISABLE", "Account had disable")
	AuthVerifySignFailed       = NewError(http.StatusOK, "AUTH_VERIFY_SIGN_FAILED", "Signature verification failure")
	AuthMaxLoginRetry          = NewError(http.StatusOK, "AUTH_MAX_LOGIN_RETRY", "Too many password errors, please try again after 8 hours.")
	AuthMaxLoginRetryCountDown = NewError(http.StatusOK, "AUTH_MAX_LOGIN_RETRY_COUNTDOWN", "Too many password errors, please try again after %s hours.")
	AssetPermissionDenied      = NewError(http.StatusOK, "ASSERT_PERMISSION_DENIED", "Asset Permission denied")
	IncorretAssetPassword      = NewError(http.StatusOK, "INCORRET_ASSET_PASSWORD", "Incorret asset password")
	DisableAssetOP             = NewError(http.StatusOK, "DISABLE_ASSET_OPERATION", "Disable asset opration")
	DisableAssetOPCountDown    = NewError(http.StatusOK, "DISABLE_ASSET_OPERATION_COUNTDOWN", "Disable asset opration")
)
