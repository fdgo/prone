package handler

import (
	"business/services/base/ifaccount/adaptor"
	"business/support/domain"
	"business/support/libraries/auth"
	"business/support/libraries/errors"
	"business/support/libraries/loggers"
	"fmt"
)

const ApiKeyBindIPMaxLimit = 5

func checkVerifyCode(account *domain.Account, code string) *errors.Error {
	if account.GAKey != "" {
		gaCode, err := auth.CurGACode(account.GAKey)
		if err != nil {
			loggers.Error.Printf("CheckVerifyCode Invalid GAKey:%s", account.GAKey)
			return errors.InternalServerError
		}
		if fmt.Sprintf("%06d", gaCode) != code {
			loggers.Error.Printf("CheckVerifyCode invalid GACode")
			return errors.IncorretGACode
		}
	} else if account.Phone.FullNumber() != "" {
		if err := adaptor.VerifySMSCode(account.Phone, code, domain.VERIFYCODE_TYPE_API_KEY); err != nil {
			loggers.Warn.Printf("CheckVerifyCode uid:%d verify code error", account.AccountId)
			return errors.SMSVerifyCodeError
		}
	} else {
		if err := adaptor.VerifyEmailCode(account.Email, code, domain.VERIFYCODE_TYPE_API_KEY); err != nil {
			loggers.Warn.Printf("AddApiKey uid:%d verify code error", account.AccountId)
			return errors.EmailVerifyCodeError
		}
	}
	return nil
}
