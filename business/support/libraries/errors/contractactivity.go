package errors

import "net/http"

var (
	KYCNoAuth                 = NewError(http.StatusOK, "KYC_NO_AUTH", "Need to do KYC authentication")
	NoBindPhone               = NewError(http.StatusOK, "NO_BIND_PHONE", "Need to bind phone")
	NotContractAccount        = NewError(http.StatusOK, "NOT_CONTRACT_ACCOUNT", "Please open a contract account")
	ActivityNotStart          = NewError(http.StatusOK, "ACTIVITY_NOT_START", "Activity is not being started")
	ActivityHadStop           = NewError(http.StatusOK, "ACTIVITY_HAD_STOP", "Activity has been suspended")
	ActivityHadEnd            = NewError(http.StatusOK, "ACTIVITY_HAD_END", "Activity has been over")
	ActivityNotFound          = NewError(http.StatusOK, "ACTIVITY_NOT_FOUND", "Activity is not found")
	HadParticipatedActivity   = NewError(http.StatusOK, "HAD_PARTICIPATED_ACTIVITY", "Had been participated this activity")
	CumulationDepositVolLimit = NewError(http.StatusOK, "CUMULATION_DEPOSIT_VOL_LIMIT", "Deposit total balance is insufficient, please deposit first.")
)
