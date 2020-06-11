package errors

import (
	"net/http"
)

var (
	InvalidSettleType      = NewError(http.StatusOK, "INVALID_SETTLE_TYPE", "Invalid settle type")
	InvalidSettleStatus    = NewError(http.StatusOK, "INVALID_SETTLE_STATUS", "Invalid settle status")
	NoAccountId            = NewError(http.StatusOK, "NO_ACCOUNT_ID", "No account id")
	NoCoinCode             = NewError(http.StatusOK, "NO_COIN_CODE", "No coin code")
	VolumeIsZero           = NewError(http.StatusOK, "VOLUME_IS_ZERO", "volume is zero")
	VolumeIsTooSmall       = NewError(http.StatusOK, "VOLUME_IS_TOOL_SMALL", "volume is tool small")
	AccountNotFound        = NewError(http.StatusOK, "ACCOUNT_NOT_FOUND", "account not found")
	NoWithdrawAddress      = NewError(http.StatusOK, "NO_WITHDRAW_ADDRESS", "No withdraw address")
	WithdrawAddressExisted = NewError(http.StatusOK, "WITHDRAW_ADDRESS_EXISTED", "You cannot add duplicate withdrawal addresses")
	InsufficientBalance    = NewError(http.StatusOK, "INSUFFICIENT_BALANCE", "insufficient balance")
	InvalidKYCStatus       = NewError(http.StatusOK, "INVALID_KYC_STATUS", "Invalid KYC status")
)
