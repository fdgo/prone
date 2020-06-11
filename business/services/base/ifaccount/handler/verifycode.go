package handler

import (
	"business/services/base/ifaccount/adaptor"
	"business/support/domain"
	"business/support/libraries/auth"
	"business/support/libraries/errors"
	"business/support/libraries/httpserver"
	"business/support/libraries/loggers"
)

var BindTypes = map[string]bool{
	"BindEmailVerifyCode": true,
	"BindPhoneVerifyCode": true,
}

func VerifyCodePost(r *httpserver.Request) *httpserver.Response {
	var (
		req domain.VerifyCodeReq
		//typ = r.QueryParams.Get("type")
	)
	if err := r.Parse(&req); err != nil {
		loggers.Warn.Printf("VerifyCodePost parse body error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.BadRequest)
	}
	//if adaptor.IncrVerifyCodeCaptchLimit(typ, r.ClientIP()) {
	if req.Validate == "" {
		return httpserver.NewResponseWithError(errors.NeedCaptch)
	}
	if !adaptor.CheckCaptchValidate(req.Validate) {
		return httpserver.NewResponseWithError(errors.InvalidValidate)
	}
	//}

	return VerifyCode(r)
}

func VerifyCode(r *httpserver.Request) *httpserver.Response {
	var (
		email    = r.QueryParams.Get("email")
		phone    = r.QueryParams.Get("phone")
		fmtEmail domain.EMAIL
		fmtPhone domain.PHONE
		typ      = r.QueryParams.Get("type")
	)
	if email == "" && phone == "" {
		loggers.Warn.Printf("VerifyCode no email input")
		return httpserver.NewResponseWithError(errors.NoEmailOrPhone)
	}
	if email != "" {
		if !domain.CheckEmailFormat(email) {
			return httpserver.NewResponseWithError(errors.InvalidEmailFormat)
		}
		fmtEmail = domain.FormatEmail(email)
	}
	if phone != "" {
		if !domain.CheckPhoneFormat(phone) {
			loggers.Debug.Printf("VerifyCode %s invalid phone format", phone)
			return httpserver.NewResponseWithError(errors.InvalidPhoneFormat)
		}
		fmtPhone = domain.FormatPhone(phone)
	}

	if typ == "" {
		loggers.Warn.Printf("VerifyCode no type input")
		return httpserver.NewResponseWithError(errors.NoVerifyCodeType)
	}
	language := r.Language
	if "" == language {
		language = "zh-cn"
	}

	tempalte, err := adaptor.GetSMSTemplate(typ, language)
	if err != nil {
		loggers.Warn.Printf("SendVerifyCode get template:%s language:%s error:%s", typ, language, err.Error())
		return httpserver.NewResponseWithError(errors.InvalidVerifyCodeType)
	}
	switch tempalte.SendType {
	case domain.SEND_TYPE_NONE:
		break
	case domain.SEND_TYPE_CHECK_LOGIN:
		{
			session, err := auth.SAuth(r.Uid, r.Sign, r.Ver, r.Dev, r.Ts)
			if err != nil {
				loggers.Warn.Printf("SendVerifyCode session auth error:%s", err.Error())
				return httpserver.NewResponseWithError(err.(*errors.Error))
			}

			if _, ok := BindTypes[typ]; ok {
				if fmtEmail != "" && session.Account.Email != "" && fmtEmail != session.Account.Email {
					loggers.Warn.Printf("SendVerifyCode input email not session email")
					return httpserver.NewResponseWithError(errors.BadRequest)
				}
				if fmtPhone.FullNumber() != "" && session.Account.Phone.FullNumber() != "" && fmtPhone.FullNumber() != session.Account.Phone.FullNumber() {
					loggers.Warn.Printf("SendVerifyCode input phone not session phone")
					return httpserver.NewResponseWithError(errors.BadRequest)
				}
			} else {
				if fmtEmail != "" && fmtEmail != session.Account.Email {
					loggers.Warn.Printf("SendVerifyCode input email not session email")
					return httpserver.NewResponseWithError(errors.BadRequest)
				}
				if fmtPhone.FullNumber() != "" && fmtPhone.FullNumber() != session.Account.Phone.FullNumber() {
					loggers.Warn.Printf("SendVerifyCode input phone not session phone")
					return httpserver.NewResponseWithError(errors.BadRequest)
				}
			}
			break
		}
	case domain.SEND_TYPE_CHECK_REGISTER:
		if fmtEmail != "" {
			account, err := adaptor.GetAccountByEmail(fmtEmail)
			if err != nil {
				loggers.Warn.Printf("SendVerifyCode get account by email:%s error:%s", email, err.Error())
				return httpserver.NewResponseWithError(errors.AccountNotFound)
			}
			if account.AccountStatus == domain.ACCOUNT_STATUS_DISABLE {
				return httpserver.NewResponseWithError(errors.AccountHadDisabled)
			}
			break
		}
		if fmtPhone.FullNumber() != "" {
			account, err := adaptor.GetAccountByPhone(fmtPhone)
			if err != nil {
				loggers.Warn.Printf("SendVerifyCode get account by email:%s error:%s", email, err.Error())
				return httpserver.NewResponseWithError(errors.AccountNotFound)
			}
			if account.AccountStatus == domain.ACCOUNT_STATUS_DISABLE {
				return httpserver.NewResponseWithError(errors.AccountHadDisabled)
			}
			break
		}
	case domain.SEND_TYPE_CHECK_UN_REGISTER:
		if fmtEmail != "" {
			account, _ := adaptor.GetAccountByEmail(fmtEmail)
			if account != nil {
				loggers.Warn.Printf("SendVerifyCode type:%s email:%s account had existed", typ, email)
				return httpserver.NewResponseWithError(errors.AccountHadExisted)
			}
			break
		}
		if fmtPhone.FullNumber() != "" {
			account, _ := adaptor.GetAccountByPhone(fmtPhone)
			if account != nil {
				loggers.Warn.Printf("SendVerifyCode type:%s phone:%s account had existed", typ, phone)
				return httpserver.NewResponseWithError(errors.AccountHadExisted)
			}
			break
		}
	default:
		break
	}

	if fmtEmail != "" {
		if err := EmailVerifyCode(fmtEmail, typ, language); err != nil {
			return httpserver.NewResponseWithError(err.(*errors.Error))
		}
	}
	if fmtPhone.FullNumber() != "" {
		if err := SMSVerifyCode(fmtPhone, typ, language); err != nil {
			return httpserver.NewResponseWithError(err.(*errors.Error))
		}
	}
	return httpserver.NewResponse()
}

func EmailVerifyCode(email domain.EMAIL, typ, language string) error {
	if err := adaptor.SendEmailVerifyCode(email, typ, language); err != nil {
		loggers.Warn.Printf("EmailVerifyCode get verify code type:%s error:%s", typ, err.Error())
		if err == errors.InvalidVerifyCodeType {
			return errors.InvalidVerifyCodeType
		}
		if err == errors.RequestTooOfften {
			return errors.RequestTooOfften
		}
		return errors.InternalServerError
	}
	return nil
}

func SMSVerifyCode(phone domain.PHONE, typ, language string) error {
	if err := adaptor.SendSMSVerifyCode(phone, typ, language); err != nil {
		loggers.Warn.Printf("PhoneVerifyCode get verify code type:%s error:%s", typ, err.Error())
		if err == errors.InvalidVerifyCodeType {
			return errors.InvalidVerifyCodeType
		}
		if err == errors.RequestTooOfften {
			return errors.RequestTooOfften
		}
		return errors.InternalServerError
	}
	return nil
}
