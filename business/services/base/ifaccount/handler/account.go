package handler

import (
	"business/services/base/ifaccount/adaptor"
	"business/support/domain"
	"business/support/libraries/auth"
	"business/support/libraries/errors"
	"business/support/libraries/httpserver"
	"business/support/libraries/loggers"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func Login(r *httpserver.Request) *httpserver.Response {
	var (
		req     domain.LoginReq
		account *domain.Account
		err     error
		//ip      = r.ClientIP()
	)

	if err = r.Parse(&req); err != nil {
		if err == errors.InvalidEmailFormat {
			return httpserver.NewResponseWithError(errors.InvalidEmailFormat)
		}
		loggers.Debug.Printf("Login parse body error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.ParameterError)
	}
	if req.Email == "" && req.Phone == "" {
		loggers.Debug.Printf("Login no email input")
		return httpserver.NewResponseWithError(errors.NoEmailOrPhone)
	}

	if req.Nonce == 0 {
		loggers.Debug.Printf("Login no nonce input")
		return httpserver.NewResponseWithError(errors.BadRequest)
	}

	if !auth.CheckTs(req.Nonce) {
		loggers.Warn.Printf("Login ip:%s,email:%s,phone:%s,dev:%s,version:%s,invalid nonce:%d", r.RemoteAddr, req.Email, req.Phone, r.Dev, r.Ver, req.Nonce)
		return httpserver.NewResponseWithError(errors.ClientTimeInvalidError)
	}

	//if adaptor.CheckLoginCaptchLimit(ip) {
	if req.Validate == "" {
		return httpserver.NewResponseWithError(errors.NeedCaptch)
	}
	if !adaptor.CheckCaptchValidate(req.Validate) {
		return httpserver.NewResponseWithError(errors.InvalidValidate)
	}
	//}

	if req.Email != "" {
		account, err = adaptor.GetAccountByEmail(req.Email)
		if err == errors.NotFound {
			loggers.Warn.Printf("Login ip:%s email:%s account not found", r.RemoteAddr, req.Email)
			return httpserver.NewResponseWithError(errors.IncorretAccountOrPassword)
		}
	} else {
		account, err = adaptor.GetAccountByPhone(req.Phone)
		if err == errors.NotFound {
			loggers.Warn.Printf("Login ip:%s phone:%s account not found", r.RemoteAddr, req.Phone)
			return httpserver.NewResponseWithError(errors.IncorretAccountOrPassword)
		}
	}

	if account.AccountStatus == domain.ACCOUNT_STATUS_PASSWORD_SAFE_WARN {
		loggers.Warn.Printf("Login ip:%s uid:%d is safe warn", r.RemoteAddr, account.AccountId)
		return httpserver.NewResponseWithError(errors.AccountSafeWarn)
	}

	if auth.IsDisabledAccountLogin(account.AccountId) {
		ts := auth.GetAccountDisabledTime(account.AccountId)
		t := math.Floor(float64(time.Now().Unix()-ts) / 3600)
		f := fmt.Sprintf("%g", 8-t)
		loggers.Warn.Printf("Login account:%s disable login,Login ip:%s,last:%s", account.String(), r.RemoteAddr, f)
		if ts <= 1 {
			return httpserver.NewResponseWithError(errors.AuthMaxLoginRetry)
		}
		return httpserver.NewResponseWithFormatError(errors.AuthMaxLoginRetryCountDown, []interface{}{f})
	}

	md5Ctx := md5.New()
	ctx := fmt.Sprintf("%s%d", account.Password, req.Nonce)
	md5Ctx.Write([]byte(ctx))
	cipherStr := md5Ctx.Sum(nil)
	if strings.ToLower(hex.EncodeToString(cipherStr)) != strings.ToLower(req.Password) {
		cnt := auth.LoginRetry(account.AccountId)
		l := 5 - cnt
		loggers.Warn.Printf("Login account:%s login retry failed,Login ip:%s,lastcnt:%d", account.String(), r.RemoteAddr, l)
		if l <= 0 {
			//账号或密码错误，请在8小时后重试
			return httpserver.NewResponseWithError(errors.AuthMaxLoginRetry)
		}

		//账户或密码输入错误，您还有x次尝试机会
		return httpserver.NewResponseWithFormatError(errors.IncorretPassword, []interface{}{l})
	}

	if account.AccountStatus == domain.ACCOUNT_STATUS_DISABLE {
		loggers.Warn.Printf("Login ip:%s email:%s account disable", r.RemoteAddr, req.Email)
		return httpserver.NewResponseWithError(errors.AccountHadDisabled)
	}

	// 系统账户禁止登陆
	if account.OwnerType == int(domain.OWNER_TYPE_SYSTEM) {
		loggers.Warn.Printf("Login ip:%s email:%s account system", r.RemoteAddr, account.String())
		return httpserver.NewResponseWithError(errors.AccountHadDisabled)
	}

	resp := httpserver.NewResponse()
	session, err := auth.NewSession(account, r.Ver, r.Dev, r.Ts)
	if err != nil {
		loggers.Error.Printf("Login NewSession error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}
	resp.SetSession(session)

	// Broadcast 登陆消息
	publishAccountRecord(account, domain.MESSAGE_ACTION_LOGIN, r.ClientIP(), "", 0)

	loginResp := domain.LoginResp{
		Account: *account.UserInfo(),
	}
	resp.Data = loginResp
	return resp
}

func Logout(r *httpserver.Request) *httpserver.Response {
	if err := auth.DeleteSession(r.Uid, r.Dev); err != nil {
		loggers.Error.Printf("Logout delete session error:%s", err.Error())
	}
	return httpserver.NewResponse()
}

func UserMe(r *httpserver.Request) *httpserver.Response {
	if r.Session == nil {
		loggers.Error.Printf("%s need auth", r.API())
		return httpserver.NewResponseWithError(errors.Forbidden)
	}
	account := r.Session.Account
	assets, err := adaptor.GetUserAssets(account.AccountId)
	if err != nil && err != errors.NotFound {
		loggers.Error.Printf("GetAccount GetUserAssets uid:%d error:%s", r.Uid, err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}
	for _, asset := range *assets {
		if asset.FreezeVol.IsPositive() || asset.AvailableVol.IsPositive() {
			account.Assets = append(account.Assets, asset)
		}
	}
	resp := httpserver.NewResponse()
	resp.Data = &domain.UserMeResp{Account: *account.UserInfo()}
	return resp
}

func Register(r *httpserver.Request) *httpserver.Response {
	var (
		req     domain.RegisterAccountReq
		account *domain.Account
		err     error
	)

	if err = r.Parse(&req); err != nil {
		if err == errors.InvalidEmailFormat {
			return httpserver.NewResponseWithError(errors.InvalidEmailFormat)
		}
		loggers.Debug.Printf("Register parse body error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.ParameterError)
	}
	if req.Email == "" && req.Phone == "" {
		loggers.Debug.Print("Register no email or phone number input")
		return httpserver.NewResponseWithError(errors.NoEmailOrPhone)
	}

	if req.Password == "" {
		loggers.Debug.Print("Register no Password input")
		return httpserver.NewResponseWithError(errors.NoPassword)
	}

	if isRegisterIpLimit(r.RemoteAddr) {
		loggers.Warn.Printf("Register %s is limit", r.RemoteAddr)
		err := errors.NewError(200, "IP_REGISTER_RATE_LIMIT", "ip register rate limit")
		return httpserver.NewResponseWithError(err)
	}
	req.AccountStatus = domain.ACCOUNT_STATUS_INACTIVATED
	// Create account
	if req.Email != "" {
		if req.Code == "" {
			loggers.Warn.Printf("Register email:%s no verify code input", req.Phone)
			return httpserver.NewResponseWithError(errors.NoVerifyCode)
		}
		account, err = adaptor.GetAccountByEmail(req.Email)
		if err != nil && err != errors.NotFound {
			loggers.Error.Printf("Register get account email:%s error:%s", req.Email, err.Error())
			return httpserver.NewResponseWithError(errors.InternalServerError)
		}
		if account != nil {
			loggers.Debug.Printf("Register email:%s exist", req.Email)
			return httpserver.NewResponseWithError(errors.EmailHadExisted)
		}
		req.AccountType = domain.ACCOUNT_TYPE_EMAIL
		if err = adaptor.VerifyEmailCode(req.Email, req.Code, domain.VERIFYCODE_TYPE_REGISTER); err != nil {
			return httpserver.NewResponseWithError(errors.VerifyCodeError)
		}
		req.AccountStatus = domain.ACCOUNT_STATUS_ACTIVATED
	} else {
		if req.Code == "" {
			loggers.Warn.Printf("Register phone:%s no verify code input", req.Phone)
			return httpserver.NewResponseWithError(errors.NoVerifyCode)
		}
		account, err = adaptor.GetAccountByPhone(req.Phone)
		if err != nil && err != errors.NotFound {
			loggers.Error.Printf("Register get account phone:%s error:%s", req.Phone, err.Error())
			return httpserver.NewResponseWithError(errors.InternalServerError)
		}
		if account != nil {
			loggers.Debug.Printf("Register phone:%s exist", req.Phone)
			return httpserver.NewResponseWithError(errors.PhoneHadExisted)
		}
		req.AccountType = domain.ACCOUNT_TYPE_PHONE
		if err = adaptor.VerifySMSCode(req.Phone, req.Code, domain.VERIFYCODE_TYPE_REGISTER); err != nil {
			return httpserver.NewResponseWithError(errors.VerifyCodeError)
		}
		req.AccountStatus = domain.ACCOUNT_STATUS_ACTIVATED

		loggers.Debug.Printf("Register inviterId:%d", req.InviterId)
	}

	req.RegisterIp = r.ClientIP()
	req.QD = r.QueryParams.Get("qd")
	var markCode *string
	mc := r.QueryParams.Get("markcode")
	if len(mc) > 0 && len(mc) <= 64 {
		markCode = &mc
	}
	ac, err1 := adaptor.CreateAccount(&req.Account, req.InviterId)
	if err1 != nil {
		loggers.Error.Printf("Register create account email:%s error:%s", req.Email, err1.Error())
		return httpserver.NewResponseWithError(err1.(*errors.Error))
	}

	session, err := auth.NewSession(ac, r.Ver, r.Dev, r.Ts)
	if err != nil {
		loggers.Error.Printf("Register NewSession error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}

	publishResigerAccountRecord(ac, domain.MESSAGE_ACTION_REGISTER, r.ClientIP(), "", req.InviterId, markCode)

	resp := httpserver.NewResponse()
	result := domain.RegtisterAccountResp{
		Account: *ac.UserInfo(),
	}
	resp.Data = result
	resp.SetSession(session)
	return resp
}

func ActiveAccount(r *httpserver.Request) *httpserver.Response {
	if r.Session == nil {
		loggers.Error.Printf("%s need auth", r.API())
		return httpserver.NewResponseWithError(errors.Forbidden)
	}

	var req domain.ActiveAccountReq
	if err := r.Parse(&req); err != nil {
		loggers.Warn.Printf("ActiveAccount parse req error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.BadRequest)
	}

	account := r.Session.Account

	if req.Code == "" {
		loggers.Error.Printf("ActiveAccount no verify code input")
		return httpserver.NewResponseWithError(errors.NoVerifyCode)
	}
	if err := adaptor.VerifyEmailCode(account.Email, req.Code, domain.VERIFYCODE_TYPE_ACTIVE_ACCOUNT); err != nil {
		return httpserver.NewResponseWithError(errors.VerifyCodeError)
	}

	if account.AccountStatus == domain.ACCOUNT_STATUS_DISABLE {
		loggers.Warn.Printf("ActiveAccount account:%d had disable", r.Uid)
		return httpserver.NewResponseWithError(errors.AccountHadDisabled)
	}

	if account.AccountStatus == domain.ACCOUNT_STATUS_ACTIVATED {
		loggers.Warn.Printf("ActiveAccount account:%d had actived", r.Uid)
		return httpserver.NewResponseWithError(errors.AccountHadActived)
	}

	if err2 := adaptor.ActiveAccount(account); err2 != nil {
		loggers.Error.Printf("ActiveAccount error:%s", err2.Error())
		return httpserver.NewResponseWithError(err2.(*errors.Error))
	}

	publishAccountRecord(account, domain.MESSAGE_ACTION_ACTIVE_ACCOUNT, r.ClientIP(), "", 0)

	account.AccountStatus = domain.ACCOUNT_STATUS_ACTIVATED
	resp := httpserver.NewResponse()
	resp.Data = account.UserInfo()

	return resp
}

func ResetPassword(r *httpserver.Request) *httpserver.Response {
	var (
		req     domain.ResetPasswordReq
		account *domain.Account
		err     error
	)

	if err = r.Parse(&req); err != nil {
		if err == errors.InvalidEmailFormat {
			return httpserver.NewResponseWithError(errors.InvalidEmailFormat)
		}
		loggers.Warn.Printf("ResetPassword parse req error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.BadRequest)
	}
	if req.Code == "" {
		loggers.Warn.Printf("ResetPassword no verify code input")
		return httpserver.NewResponseWithError(errors.NoVerifyCode)
	}
	if req.Email == "" && req.Phone == "" {
		loggers.Warn.Printf("ResetPassword no email input")
		return httpserver.NewResponseWithError(errors.NoEmailOrPhone)
	}
	if req.Password == "" {
		loggers.Warn.Printf("ResetPassword no password input")
		return httpserver.NewResponseWithError(errors.NoPassword)
	}
	if req.Email != "" {
		account, err = adaptor.GetAccountByEmail(req.Email)
		if nil == account {
			loggers.Warn.Printf("ResetPassword invalid email")
			return httpserver.NewResponseWithError(errors.AccountNotFound)
		}
		if err := adaptor.VerifyEmailCode(req.Email, req.Code, domain.VERIFYCODE_TYPE_RESET_PASSWORD); err != nil {
			loggers.Warn.Printf("ResetPassword verify code error")
			return httpserver.NewResponseWithError(errors.VerifyCodeError)
		}
	} else {
		account, err = adaptor.GetAccountByPhone(req.Phone)
		if nil == account {
			loggers.Warn.Printf("ResetPassword invalid phone")
			return httpserver.NewResponseWithError(errors.AccountNotFound)
		}
		if err := adaptor.VerifySMSCode(req.Phone, req.Code, domain.VERIFYCODE_TYPE_RESET_PASSWORD); err != nil {
			loggers.Warn.Printf("ResetPassword verify code error")
			return httpserver.NewResponseWithError(errors.VerifyCodeError)
		}
	}

	account.Password = req.Password
	if err2 := adaptor.ResetPassword(account); err2 != nil {
		loggers.Error.Printf("ResetPassword error:%s", err2.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}
	auth.EnableLogin(account.AccountId)

	httpserver.CleanSession(account.AccountId)
	session, err := auth.NewSession(account, r.Ver, r.Dev, r.Ts)
	if err != nil {
		loggers.Error.Printf("Login NewSession error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}

	publishAccountRecord(account, domain.MESSAGE_ACTION_RESETPASSWD, r.ClientIP(), "", 0)

	resp := httpserver.NewResponse()
	resp.SetSession(session)
	loginResp := domain.LoginResp{
		Account: *account.UserInfo(),
	}
	resp.Data = loginResp

	return resp
}

func BindEmail(r *httpserver.Request) *httpserver.Response {
	if r.Session == nil {
		loggers.Error.Printf("%s need auth", r.API())
		return httpserver.NewResponseWithError(errors.Forbidden)
	}

	var req domain.BindEmailReq
	if err := r.Parse(&req); err != nil {
		loggers.Warn.Printf("BindEmail parse req error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.BadRequest)
	}
	if req.SMSCode == "" || req.EmailCode == "" {
		loggers.Warn.Printf("BindEmail no verify code")
		return httpserver.NewResponseWithError(errors.NoVerifyCode)
	}
	if req.Email == "" {
		loggers.Warn.Printf("BindEmail no email input")
		return httpserver.NewResponseWithError(errors.NoEmail)
	}

	if _, err := adaptor.GetAccountByEmail(req.Email); err == nil {
		loggers.Warn.Printf("BindEmail email:%s already exist", req.Email)
		return httpserver.NewResponseWithError(errors.EmailHadExisted)
	}

	account := r.Session.Account

	if account.Email != "" {
		loggers.Warn.Printf("BindEmail uid:%d had bind email:%s", r.Uid, account.Email)
		return httpserver.NewResponseWithError(errors.EmailAlreadyBind)
	}
	if account.Phone == "" {
		loggers.Warn.Printf("BindEmail uid:%d phone is null", r.Uid)
		return httpserver.NewResponseWithError(errors.SMSVerifyCodeError)
	}
	if err := adaptor.VerifySMSCode(account.Phone, req.SMSCode, domain.VERIFYCODE_TYPE_BIND_EMAIL); err != nil {
		loggers.Warn.Printf("BindEmail uid:%d verify code error", r.Uid)
		return httpserver.NewResponseWithError(errors.SMSVerifyCodeError)
	}
	if err := adaptor.VerifyEmailCode(req.Email, req.EmailCode, domain.VERIFYCODE_TYPE_BIND_EMAIL); err != nil {
		loggers.Warn.Printf("BindEmail uid:%d verify code error", r.Uid)
		return httpserver.NewResponseWithError(errors.EmailVerifyCodeError)
	}
	account.Email = req.Email
	if rErr := adaptor.BindEmail(account); rErr != nil {
		loggers.Warn.Printf("BindEmail uid:%d error:%s", r.Uid, rErr.Error())
		return httpserver.NewResponseWithError(rErr.(*errors.Error))
	}

	publishAccountRecord(account, domain.MESSAGE_ACTION_BIND_EMAIL, r.ClientIP(), "", 0)

	return httpserver.NewResponse()
}

func BindPhone(r *httpserver.Request) *httpserver.Response {
	if r.Session == nil {
		loggers.Error.Printf("%s need auth", r.API())
		return httpserver.NewResponseWithError(errors.Forbidden)
	}

	var req domain.BindPhoneReq
	if err := r.Parse(&req); err != nil {
		loggers.Warn.Printf("BindPhone parse req error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.BadRequest)
	}
	if req.SMSCode == "" || req.EmailCode == "" {
		loggers.Warn.Printf("BindPhone no verify code")
		return httpserver.NewResponseWithError(errors.NoVerifyCode)
	}
	if req.Phone == "" {
		loggers.Warn.Printf("BindPhone no email input")
		return httpserver.NewResponseWithError(errors.NoEmail)
	}

	if _, err := adaptor.GetAccountByPhone(req.Phone); err == nil {
		loggers.Warn.Printf("BindPhone phone:%s already exist", req.Phone)
		return httpserver.NewResponseWithError(errors.PhoneHadExisted)
	}

	account := r.Session.Account
	if account.AccountStatus == domain.ACCOUNT_STATUS_INACTIVATED {
		loggers.Error.Printf("BindPhone uid:%d not active", r.Uid)
		return httpserver.NewResponseWithError(errors.AccountNotActived)
	}
	if account.Phone != "" {
		loggers.Warn.Printf("BindPhone uid:%d had bind email:%s", r.Uid, account.Email)
		return httpserver.NewResponseWithError(errors.PhoneAlreadyBind)
	}
	if account.Email == "" {
		loggers.Warn.Printf("BindPhone uid:%d phone is null", r.Uid)
		return httpserver.NewResponseWithError(errors.EmailVerifyCodeError)
	}
	if err := adaptor.VerifySMSCode(req.Phone, req.SMSCode, domain.VERIFYCODE_TYPE_BIND_PHONE); err != nil {
		loggers.Warn.Printf("BindPhone uid:%d verify code error", r.Uid)
		return httpserver.NewResponseWithError(errors.SMSVerifyCodeError)
	}
	if err := adaptor.VerifyEmailCode(account.Email, req.EmailCode, domain.VERIFYCODE_TYPE_BIND_PHONE); err != nil {
		loggers.Warn.Printf("BindPhone uid:%d verify code error", r.Uid)
		return httpserver.NewResponseWithError(errors.EmailVerifyCodeError)
	}
	account.Phone = req.Phone
	if rErr := adaptor.BindPhone(account); rErr != nil {
		loggers.Warn.Printf("BindPhone uid:%d error:%s", r.Uid, rErr.Error())
		return httpserver.NewResponseWithError(rErr.(*errors.Error))
	}

	publishAccountRecord(account, domain.MESSAGE_ACTION_BIND_PHONE, r.ClientIP(), "", 0)

	return httpserver.NewResponse()
}

func CheckAccountExist(r *httpserver.Request) *httpserver.Response {
	email := r.QueryParams.Get("email")
	phone := r.QueryParams.Get("phone")
	if email == "" && phone == "" {
		return httpserver.NewResponseWithError(errors.ParameterError)
	}

	resp := httpserver.NewResponse()
	if email != "" {
		if !domain.CheckEmailFormat(email) {
			return httpserver.NewResponseWithError(errors.InvalidEmailFormat)
		}
		_, err := adaptor.GetAccountByEmail(domain.FormatEmail(email))
		if err == nil {
			resp.Data = domain.CheckAccountExistResp{true}
			return resp
		}
		resp.Data = domain.CheckAccountExistResp{false}
		return resp
	}

	if !domain.CheckPhoneFormat(phone) {
		return httpserver.NewResponseWithError(errors.InvalidPhoneFormat)
	}
	fmtPhone := domain.FormatPhone(phone)
	if fmtPhone == "" {
		return httpserver.NewResponseWithError(errors.InvalidPhoneFormat)
	}
	_, err := adaptor.GetAccountByPhone(fmtPhone)
	if err == nil {
		resp.Data = domain.CheckAccountExistResp{true}
		return resp
	}
	resp.Data = domain.CheckAccountExistResp{false}
	return resp
}

func AddAssetPassword(r *httpserver.Request) *httpserver.Response {
	if r.Session == nil {
		loggers.Error.Printf("%s need auth", r.API())
		return httpserver.NewResponseWithError(errors.Forbidden)
	}

	var (
		req domain.Account
		err error
	)

	if err = r.Parse(&req); err != nil {
		loggers.Warn.Printf("AddAssetPassword parse body error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.ParameterError)
	}

	if req.AssetPassword == "" {
		loggers.Warn.Printf("AddAssetPassword error no asset password input")
		return httpserver.NewResponseWithError(errors.NoAssetPassword)
	}

	account := r.Session.Account

	if account.AssetPassword != "" {
		loggers.Warn.Printf("AddAssetPassword asset password had existed")
		return httpserver.NewResponseWithError(errors.AssetPasswordHadExisted)
	}

	account.AssetPassword = req.AssetPassword
	account.AssetPasswordEffectiveTime = 7200
	if err := adaptor.AddAssetPassword(account); err != nil {
		loggers.Warn.Printf("AddAssetPassword error:%s", err.Error())
		return httpserver.NewResponseWithError(err.(*errors.Error))
	}

	publishAccountRecord(account, domain.MESSAGE_ACTION_ASSET_PASSWORD_ADD, r.ClientIP(), "", 0)

	resp := httpserver.NewResponse()
	resp.Data = account.UserInfo()
	return resp
}

func ResetAssetPassword(r *httpserver.Request) *httpserver.Response {
	if r.Session == nil {
		loggers.Error.Printf("%s need auth", r.API())
		return httpserver.NewResponseWithError(errors.Forbidden)
	}

	var (
		req domain.ResetAssetPasswordReq
		err error
	)
	if err = r.Parse(&req); err != nil {
		loggers.Warn.Printf("ResetAssetPassword parse body error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.ParameterError)
	}
	if req.AssetPassword == "" {
		loggers.Warn.Printf("ResetAssetPassword error no asset password input")
		return httpserver.NewResponseWithError(errors.NoAssetPassword)
	}
	if req.EmailCode == "" && req.SMSCode == "" {
		loggers.Warn.Printf("ResetAssetPassword error no verify code input")
		return httpserver.NewResponseWithError(errors.NoVerifyCode)
	}

	account := r.Session.Account
	if req.EmailCode != "" {
		if err := adaptor.VerifyEmailCode(account.Email, req.EmailCode, domain.VERIFYCODE_TYPE_RESET_ASSET_PASSWORD); err != nil {
			loggers.Warn.Printf("ResetAssetPassword ip:%s account email:%s verify code error", r.RemoteAddr, account.Email)
			return httpserver.NewResponseWithError(errors.VerifyCodeError)
		}
	}
	if req.SMSCode != "" {
		if err := adaptor.VerifySMSCode(account.Phone, req.SMSCode, domain.VERIFYCODE_TYPE_RESET_ASSET_PASSWORD); err != nil {
			loggers.Warn.Printf("ResetAssetPassword ip:%s account phone:%s verify code error", r.RemoteAddr, account.Phone)
			return httpserver.NewResponseWithError(errors.VerifyCodeError)
		}
	}
	if account.GAKey != "" {
		gaCode, err := auth.CurGACode(account.GAKey)
		if err != nil {
			loggers.Error.Printf("Invalid GAKey %s", account.GAKey)
			return httpserver.NewResponseWithError(errors.InternalServerError)
		}
		if req.GACode != gaCode {
			loggers.Warn.Printf("Withdraw ip:%s incorret ga code", r.RemoteAddr)
			return httpserver.NewResponseWithError(errors.IncorretGACode)
		}
	}
	account.AssetPassword = req.AssetPassword
	if err := adaptor.ResetAssetPassword(account); err != nil {
		loggers.Warn.Printf("ResetAssetPassword error:%s", err.Error())
		return httpserver.NewResponseWithError(err.(*errors.Error))
	}

	publishAccountRecord(account, domain.MESSAGE_ACTION_ASSET_PASSWORD_RESET, r.ClientIP(), "", 0)

	return httpserver.NewResponse()
}

func ResetAssetPasswordEffectiveTime(r *httpserver.Request) *httpserver.Response {
	if r.Session == nil {
		loggers.Error.Printf("%s need auth", r.API())
		return httpserver.NewResponseWithError(errors.Forbidden)
	}

	var (
		req domain.Account
		err error
	)
	if err = r.Parse(&req); err != nil {
		loggers.Warn.Printf("ResetAssetPasswordEffectiveTime parse body error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.ParameterError)
	}
	if req.AssetPassword == "" {
		loggers.Warn.Printf("ResetAssetPasswordEffectiveTime error no asset password input")
		return httpserver.NewResponseWithError(errors.NoAssetPassword)
	}

	if req.AssetPasswordEffectiveTime < 0 && req.AssetPasswordEffectiveTime != domain.ASSET_PASSWORD_EFFECTIVE_TIME_PERMANENT {
		loggers.Warn.Printf("ResetAssetPasswordEffectiveTime invalid asset password effective time:%d", req.AssetPasswordEffectiveTime)
		return httpserver.NewResponseWithError(errors.InvalidEffectiveTime)
	}

	account := r.Session.Account
	if auth.IsDisableAssetOP(account.AccountId) {
		ts := auth.GetDisableAssetOPTime(account.AccountId)
		t := math.Floor(float64(time.Now().Unix()-ts) / 3600)
		f := fmt.Sprintf("%g", 24-t)
		loggers.Warn.Printf("ResetAssetPasswordEffectiveTime account %d optime %d countdown %s", account.AccountId, ts, f)
		if ts <= 1 {
			return httpserver.NewResponseWithError(errors.DisableAssetOP)
		}
		return httpserver.NewResponseWithFormatError(errors.DisableAssetOPCountDown, []interface{}{f})
	}
	if account.AssetPassword != req.AssetPassword {
		count, err := auth.IncorrectAssetPasswordIncr(account.AccountId)
		if err != nil {
			loggers.Error.Printf("IncorrectAssetPasswordIncr error %s", err.Error())
		}
		if count >= 6 {
			auth.DisableAssetOP(account.AccountId)
		}
		l := 6 - count
		loggers.Warn.Printf("ResetAssetPasswordEffectiveTime asset password error cnt %d", l)
		if l <= 0 {
			//账号或密码错误，请在24小时后重试
			return httpserver.NewResponseWithError(errors.DisableAssetOP)
		}

		//账户或密码输入错误，您还有x次尝试机会
		return httpserver.NewResponseWithFormatError(errors.IncorretAssetPassword, []interface{}{l})
	}

	account.AssetPasswordEffectiveTime = req.AssetPasswordEffectiveTime
	if err := adaptor.ResetAssetPasswordEffectiveTime(account); err != nil {
		loggers.Warn.Printf("ResetAssetPasswordEffectiveTime error %s", err.Error())
		return httpserver.NewResponseWithError(err.(*errors.Error))
	}

	publishAccountRecord(account, domain.MESSAGE_ACTION_ASSET_PASSWORD_EFFECTIVE_TIME_RESET, r.ClientIP(), "", 0)

	resp := httpserver.NewResponse()
	resp.Data = account.UserInfo()
	return resp
}

func SetAntiFishingText(r *httpserver.Request) *httpserver.Response {
	if r.Session == nil {
		loggers.Error.Printf("%s need auth", r.API())
		return httpserver.NewResponseWithError(errors.Forbidden)
	}
	var (
		req domain.Account
		err error
	)
	if err = r.Parse(&req); err != nil {
		loggers.Warn.Printf("SetAntiFishingText parse body error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.ParameterError)
	}
	r.Session.Account.AntiFishingText = req.AntiFishingText
	if err := adaptor.SetAntiFishingText(r.Session.Account); err != nil {
		loggers.Warn.Printf("SetAntiFishingText error:%s", err.Error())
		return httpserver.NewResponseWithError(err.(*errors.Error))
	}

	publishAccountRecord(r.Session.Account, domain.MESSAGE_ACTION_ANTI_FISHING_TEXT_SET, r.ClientIP(), "", 0)

	resp := httpserver.NewResponse()
	resp.IsResetToken = true
	return resp
}

func GetGAKey(r *httpserver.Request) *httpserver.Response {
	if r.Session == nil {
		loggers.Error.Printf("%s need auth", r.API())
		return httpserver.NewResponseWithError(errors.Forbidden)
	}

	if r.Session.Account.GAKey != "" {
		return httpserver.NewResponseWithError(errors.GAKeyHadExisted)
	}
	gaKey, err := adaptor.GenGAKey(r.Session.Account)
	if err != nil {
		loggers.Error.Printf("GetGAKey error %s", err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}
	gaKeyResp := domain.GetGAKeyResp{
		GAKey:     gaKey,
		LoginName: r.Session.Account.String(),
	}

	resp := httpserver.NewResponse()
	resp.Data = gaKeyResp
	return resp
}

func SetGAKey(r *httpserver.Request) *httpserver.Response {
	if r.Session == nil {
		loggers.Error.Printf("%s need auth", r.API())
		return httpserver.NewResponseWithError(errors.Forbidden)
	}

	var (
		req domain.SetGAKeyReq
		err error
	)
	if err = r.Parse(&req); err != nil {
		loggers.Warn.Printf("SetGAkey parse body error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.ParameterError)
	}
	if r.Session.Account.GAKey != "" {
		return httpserver.NewResponseWithError(errors.GAKeyHadExisted)
	}

	gaKey, err := adaptor.GetCacheGAKey(r.Session.Account)
	if err != nil {
		loggers.Warn.Printf("SetGAkey error %s", err.Error())
		return httpserver.NewResponseWithError(errors.GAKeyExpired)
	}

	gaCode, err := auth.CurGACode(gaKey)
	if err != nil {
		loggers.Error.Printf("Invalid GAKey %s", r.Session.Account.GAKey)
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}
	if req.GACode != gaCode {
		loggers.Debug.Printf("SetGAKey invalid GACode key:%s %d:%d", gaKey, req.GACode, gaCode)
		return httpserver.NewResponseWithError(errors.IncorretGACode)
	}

	r.Session.Account.GAKey = gaKey
	if err := adaptor.SetGAKey(r.Session.Account); err != nil {
		loggers.Warn.Printf("SetGAkey error %s", err.Error())
		return httpserver.NewResponseWithError(err.(*errors.Error))
	}

	publishAccountRecord(r.Session.Account, domain.MESSAGE_ACTION_GAKEY_ADD, r.ClientIP(), "", 0)

	return httpserver.NewResponse()
}

func DeleteGAKey(r *httpserver.Request) *httpserver.Response {
	if r.Session == nil {
		loggers.Error.Printf("%s need auth", r.API())
		return httpserver.NewResponseWithError(errors.Forbidden)
	}

	var (
		req domain.DeleteGAKeyReq
		err error
	)
	if err = r.Parse(&req); err != nil {
		loggers.Warn.Printf("DeleteGAKey parse body error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.ParameterError)
	}

	gaCode, err := auth.CurGACode(r.Session.Account.GAKey)
	if err != nil {
		loggers.Error.Printf("DeleteGAKey Invalid GAKey %s", r.Session.Account.GAKey)
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}
	if req.GACode != gaCode {
		loggers.Debug.Printf("DeleteGAKey invalid GACode key:%s %d:%d", r.Session.Account.GAKey, req.GACode, gaCode)
		return httpserver.NewResponseWithError(errors.IncorretGACode)
	}

	r.Session.Account.GAKey = ""
	if err := adaptor.SetGAKey(r.Session.Account); err != nil {
		loggers.Warn.Printf("SetGAkey error %s", err.Error())
		return httpserver.NewResponseWithError(err.(*errors.Error))
	}

	publishAccountRecord(r.Session.Account, domain.MESSAGE_ACTION_GAKEY_DELETE, r.ClientIP(), "", 0)

	return httpserver.NewResponse()
}

func GetLoginRecord(r *httpserver.Request) *httpserver.Response {
	if r.Session == nil {
		loggers.Error.Printf("%s need auth", r.API())
		return httpserver.NewResponseWithError(errors.Forbidden)
	}

	limit := 10
	v := r.QueryParams.Get("limit")
	if v != "" {
		if n, err := strconv.Atoi(v); err != nil && n > 0 && n < 50 {
			limit = n
		}
	}
	loginRecords, err := adaptor.GetLoginRecord(r.Uid, int64(limit))
	if err != nil {
		loggers.Error.Printf("GetLoginRecord error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}
	resp := httpserver.NewResponse()
	resp.Data = loginRecords
	return resp
}

func SetAccountName(r *httpserver.Request) *httpserver.Response {
	if r.Session == nil {
		loggers.Error.Printf("%s need auth", r.API())
		return httpserver.NewResponseWithError(errors.Forbidden)
	}
	var account domain.Account

	if err := r.Parse(&account); err != nil {
		loggers.Warn.Printf("%s parse body error", r.API())
		return httpserver.NewResponseWithError(errors.ParameterError)
	}

	if account.AccountName == "" {
		loggers.Warn.Printf("%s no account name input", r.API())
		return httpserver.NewResponseWithError(errors.ParameterError)
	}

	if err := domain.CheckAccountName(account.AccountName); err != nil {
		return httpserver.NewResponseWithError(err)
	}

	account.AccountId = r.Session.Account.AccountId
	reply, err := adaptor.SetAccountName(&account)
	if err != nil {
		loggers.Warn.Printf("%s SetAccountName error %s", r.API(), err.Error())
		switch {
		case errors.AccountNameHadExisted.Equal(err):
			return httpserver.NewResponseWithError(errors.AccountNameHadExisted)
		case errors.AccountNameHadBeenUsed.Equal(err):
			return httpserver.NewResponseWithError(errors.AccountNameHadBeenUsed)
		default:
			return httpserver.NewResponseWithError(errors.InternalServerError)
		}
	}

	resp := httpserver.NewResponse()
	resp.Data = reply.UserInfo()

	return resp
}

func UpdateAccountAvatar(r *httpserver.Request) *httpserver.Response {
	if r.Session == nil {
		loggers.Error.Printf("%s need auth", r.API())
		return httpserver.NewResponseWithError(errors.Forbidden)
	}
	var account domain.Account

	if err := r.Parse(&account); err != nil {
		loggers.Warn.Printf("%s parse body error", r.API())
		return httpserver.NewResponseWithError(errors.ParameterError)
	}

	if account.Avatar == "" {
		loggers.Warn.Printf("%s no account name input", r.API())
		return httpserver.NewResponseWithError(errors.ParameterError)
	}
	if !domain.CheckURL(account.Avatar) {
		return httpserver.NewResponseWithError(errors.BadRequest)
	}

	account.AccountId = r.Session.Account.AccountId
	reply, err := adaptor.UpdateAccountAvatar(&account)
	if err != nil {
		loggers.Warn.Printf("%s UpdateAccountAvatar error %s", r.API(), err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}

	resp := httpserver.NewResponse()
	resp.Data = reply.UserInfo()

	return resp
}

func publishAccountRecord(account *domain.Account, action, clientIP, device string, inviterId int64) {
	record := domain.Message{
		MessageBase: domain.MessageBase{
			Action:   action,
			Type:     domain.MESSAGE_TYPE_ACCOUNT,
			From:     fmt.Sprintf("%d", account.AccountId),
			FromType: domain.FROM_TYPE_USER,
		},
		Account: &domain.AccountRecord{
			Account:   *account,
			ClientIP:  clientIP,
			Device:    device,
			InviterId: inviterId,
		},
	}
	record.Account.Password = ""
	record.Account.AssetPassword = ""
	record.Account.GAKey = ""
	adaptor.PublishAccountRecord(&record)
}

func publishResigerAccountRecord(account *domain.Account, action, clientIP, device string, inviterId int64, markCode *string) {
	record := domain.Message{
		MessageBase: domain.MessageBase{
			Action:   action,
			Type:     domain.MESSAGE_TYPE_ACCOUNT,
			From:     fmt.Sprintf("%d", account.AccountId),
			FromType: domain.FROM_TYPE_USER,
		},
		Account: &domain.AccountRecord{
			Account:   *account,
			ClientIP:  clientIP,
			Device:    device,
			InviterId: inviterId,
			MarkCode:  markCode,
		},
	}
	record.Account.Password = ""
	record.Account.AssetPassword = ""
	record.Account.GAKey = ""
	adaptor.PublishAccountRecord(&record)
}

func isRegisterIpLimit(ip string) bool {
	conf := adaptor.GetRegisterIpLimitConfig()
	if conf == nil {
		return false
	}
	lastTime := time.Now().UTC().Add(-(time.Second * time.Duration(conf.LimitTime)))

	if adaptor.GetAccountByRegisterIp(ip, &lastTime) >= conf.Count {
		return true
	}
	return false
}
