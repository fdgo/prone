package handler

import (
	"business/services/base/ifaccount/adaptor"
	baseconfig "business/support/config"
	"business/support/domain"
	"business/support/libraries/auth"
	"business/support/libraries/errors"
	"business/support/libraries/httpserver"
	"business/support/libraries/loggers"
	"business/support/model/config"
	"strconv"

	"github.com/shopspring/decimal"
)

func Withdraw(r *httpserver.Request) *httpserver.Response {
	if r.Session == nil {
		loggers.Error.Printf("%s need auth", r.API())
		return httpserver.NewResponseWithError(errors.Forbidden)
	}

	var (
		req domain.WithdrawReq
		err error
	)
	if err := r.Parse(&req); err != nil {
		loggers.Warn.Printf("Withdraw parse body error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.ParameterError)
	}

	if len(req.CoinCode) < 1 ||
		!domain.CheckVol(req.Vol) ||
		req.Nonce == 0 {
		loggers.Warn.Printf("Withdraw invalid param")
		return httpserver.NewResponseWithError(errors.ParameterError)
	}

	if req.EmailCode == "" && req.SMSCode == "" && req.VerifyCode == "" {
		loggers.Warn.Printf("Withdraw no verify code input")
		return httpserver.NewResponseWithError(errors.NoVerifyCode)
	}

	// 校验提现地址格式
	switch req.Settle.Type {
	case domain.SETTLE_TYPE_TRANSFER_OUT:
		toAccountID, err := strconv.ParseInt(req.ToAddress, 10, 64)
		if err != nil {
			loggers.Warn.Printf("Withdraw check address error:%s", err.Error())
			return httpserver.NewResponseWithError(errors.InvalidAddress)
		}
		_, err = adaptor.GetAccountByUid(toAccountID)
		if err != nil {
			loggers.Warn.Printf("Withdraw check address error:%s", err.Error())
			if errors.NotFound.Equal(err) {
				return httpserver.NewResponseWithError(errors.UidNotFound)
			}
			return httpserver.NewResponseWithError(errors.InternalServerError)
		}
	default:
		req.ToAddress, err = adaptor.CheckAddress(req.CoinCode, req.ToAddress)
		if err != nil {
			loggers.Warn.Printf("Withdraw check address error:%s", err.Error())
			return httpserver.NewResponseWithError(errors.InvalidAddress)
		}
	}

	account := r.Session.Account

	// 如果设置了Google验证器,校验Google验证码是否正确
	if account.GAKey != "" {
		gaCode, err := auth.CurGACode(account.GAKey)
		if err != nil {
			loggers.Error.Printf("Invalid GAKey %s", account.GAKey)
			return httpserver.NewResponseWithError(errors.InternalServerError)
		}
		if req.GACode != gaCode {
			loggers.Warn.Printf("Withdraw ip:%s incorret GACode key:%s %d %d", r.RemoteAddr, account.GAKey, req.GACode, gaCode)
			return httpserver.NewResponseWithError(errors.IncorretGACode)
		}
	}

	// 判断账号是否被禁用
	if account.AccountStatus == domain.ACCOUNT_STATUS_DISABLE {
		loggers.Warn.Printf("Withdraw ip:%s email:%s account disable", r.RemoteAddr, account.Email)
		return httpserver.NewResponseWithError(errors.AccountHadDisabled)
	}

	// 校验短信/邮件验证码
	if req.VerifyCode != "" {
		if err := adaptor.VerifyEmailCode(account.Email, req.VerifyCode, domain.VERIFYCODE_TYPE_WITHDRAW); err != nil {
			loggers.Warn.Printf("Withdraw ip:%s account email:%s verify code erro", r.RemoteAddr, account.Email)
			return httpserver.NewResponseWithError(errors.VerifyCodeError)
		}
	}
	if req.EmailCode != "" {
		if err := adaptor.VerifyEmailCode(account.Email, req.EmailCode, domain.VERIFYCODE_TYPE_WITHDRAW); err != nil {
			loggers.Warn.Printf("Withdraw ip:%s account email:%s verify code erro", r.RemoteAddr, account.Email)
			return httpserver.NewResponseWithError(errors.VerifyCodeError)
		}
	}
	if req.SMSCode != "" {
		if err := adaptor.VerifySMSCode(account.Phone, req.SMSCode, domain.VERIFYCODE_TYPE_WITHDRAW); err != nil {
			loggers.Warn.Printf("Withdraw ip:%s account phone:%s verify code erro", r.RemoteAddr, account.Phone)
			return httpserver.NewResponseWithError(errors.VerifyCodeError)
		}
	}

	// 如果设置交易密码,并且安全设置要求检验,校验交易密码
	//assetPassword := r.Header.Get("bifund-AssetPassword")
	//if account.AssetPasswordEffectiveTime >= 0 && assetPassword != account.AssetPassword {
	//	loggers.Warn.Printf("Withdraw ip:%s account asset password code incorret", r.RemoteAddr)
	//	return httpserver.NewResponseWithError(errors.IncorretAccountOrPassword)
	//}

	// 检验提现配置
	conf, err := config.GetAndCacheWithdrawConfig(req.Settle.CoinCode)
	if err != nil {
		loggers.Error.Printf("Withdarw get withdraw config error %s", err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}
	if conf.Status == config.WITHDRAWAL_CONFIG_STATUS_DISABLE {
		loggers.Error.Printf("Withdraw %s disable now", req.Settle.CoinCode)
		return httpserver.NewResponseWithError(errors.WithdrawDisable)
	}
	if req.Settle.Vol.LessThan(conf.MinVol) {
		return httpserver.NewResponseWithError(errors.VolumeIsTooSmall)
	}
	if account.KYCStatus != domain.KYC_STATUS_APPROVED && conf.KYCMaxVol.IsPositive() {
		vol, err := adaptor.GetWithdrawTotalDay(req.Settle.CoinCode, r.Uid)
		if err != nil {
			loggers.Error.Printf("Withdraw get withdraw total error %s", err.Error())
			return httpserver.NewResponseWithError(errors.InternalServerError)
		}
		if vol.Add(req.Settle.Vol).GreaterThan(conf.KYCMaxVol) {
			return httpserver.NewResponseWithError(errors.MaxLimitVolNoKYC)
		}
	}

	// 校验账号类型是否允许提现
	var ownerTypes []int
	if err := config.GetAndParseCacheGlobal("DISABLE_WITHDRAW_OWNER_TYPES", &ownerTypes); err != nil {
		ownerTypes = []int{}
	}
	for i := range ownerTypes {
		if ownerTypes[i] == r.Session.Account.OwnerType {
			loggers.Error.Printf("Withdraw owner_type:%d disable now", r.Session.Account.OwnerType)
			return httpserver.NewResponseWithError(errors.WithdrawDisable)
		}
	}

	// 内部转账存在次数限制时,校验转账次数
	if maxCount, err := config.GetAndCacheGlobal("SETTLE.TRANSFER_MAX_COUNT").Int64(); err == nil && req.Type == domain.SETTLE_TYPE_TRANSFER_OUT {
		count, err := adaptor.GetWithdrawCountHour(r.Uid)
		if err != nil {
			loggers.Error.Printf("Witdraw get latest hour count error %s", err.Error())
			return httpserver.NewResponseWithError(errors.InternalServerError)
		}
		if count >= maxCount {
			loggers.Debug.Printf("Withdraw %s max limit count internal transfer in one hour")
			return httpserver.NewResponseWithError(errors.MaxLimitCountTransfer)
		}
	}

	settle := req.Settle
	settle.AccountId = account.AccountId
	switch settle.Type {
	case domain.SETTLE_TYPE_TRANSFER_OUT:
		settle.Fee = decimal.Zero
		settle.FeeCoinCode = ""
	default:
		settle.Type = domain.SETTLE_TYPE_WITHDRAW
		settle.Fee = conf.Fee
		settle.FeeCoinCode = conf.FeeCoinName
	}
	if err2 := adaptor.Withdraw(&settle); err2 != nil {
		loggers.Warn.Printf("Withdraw ip:%s error:%s", r.RemoteAddr, err2.Error())
		return httpserver.NewResponseWithError(err2.(*errors.Error))
	}

	return httpserver.NewResponse()
}

func WithdrawTest(r *httpserver.Request) *httpserver.Response {
	var req domain.WithdrawReq
	if err := r.Parse(&req); err != nil {
		loggers.Warn.Printf("Withdraw parse body error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.ParameterError)
	}

	if len(req.CoinCode) < 1 ||
		!domain.CheckVol(req.Vol) ||
		req.Nonce == 0 ||
		len(req.ToAddress) < 1 {
		loggers.Warn.Printf("Withdraw invalid param")
		return httpserver.NewResponseWithError(errors.ParameterError)
	}
	account, err := adaptor.GetAccountByUid(r.Uid)
	if err == errors.NotFound {
		loggers.Warn.Printf("Withdraw ip:%s uid:%d account not found", r.RemoteAddr, r.Uid)
		return httpserver.NewResponseWithError(errors.NotFound)
	}
	if account.AccountStatus == domain.ACCOUNT_STATUS_DISABLE {
		loggers.Warn.Printf("Withdraw ip:%s email:%s account disable", r.RemoteAddr, account.Email)
		return httpserver.NewResponseWithError(errors.AccountHadDisabled)
	}
	settle := req.Settle
	settle.AccountId = account.AccountId
	settle.Type = domain.SETTLE_TYPE_WITHDRAW
	if err2 := adaptor.Withdraw(&settle); err2 != nil {
		loggers.Warn.Printf("Withdraw ip:%s error:%s", r.RemoteAddr, err2.Error())
		return httpserver.NewResponseWithError(err2.(*errors.Error))
	}

	return httpserver.NewResponse()
}

func DepositAddress(r *httpserver.Request) *httpserver.Response {
	if r.Session == nil {
		loggers.Error.Printf("%s need auth", r.API())
		return httpserver.NewResponseWithError(errors.Forbidden)
	}

	coinCode := r.QueryParams.Get("coin")
	if coinCode == "" {
		loggers.Warn.Printf("DepositAddress no coin code input")
		return httpserver.NewResponseWithError(errors.NoCoinCode)
	}

	account := r.Session.Account

	if account.AccountStatus == domain.ACCOUNT_STATUS_DISABLE {
		loggers.Warn.Printf("DepositAddress ip:%s uid:%d account disable", r.RemoteAddr, r.Uid)
		return httpserver.NewResponseWithError(errors.AccountHadDisabled)
	}
	addr, err := adaptor.BindAddress(r.Uid, coinCode)
	if nil != err {
		return httpserver.NewResponseWithError(err.(*errors.Error))
	}
	if addr == nil {
		return httpserver.NewResponseWithError(errors.InternalServerError)
	} else {
		addr.CreatedAt = nil
	}

	loggers.Debug.Printf("UserDepositAddress %s %d %s", addr.CoinGroup, r.Uid, addr.DepositAddress)
	resp := httpserver.NewResponse()
	data := domain.AddressResp{UserDepositAddress: *addr, Memo: ""}
	if addr.CoinGroup == "EOS" {
		BankConfig, ok := baseconfig.Conf.Services["Bank"]
		if ok {
			data.Memo = data.DepositAddress
			data.DepositAddress = ""
			data.Account = BankConfig.EOSAccount
		} else {
			data.DepositAddress = ""
			data.Account = ""
			data.Memo = ""
		}
	}
	resp.Data = data
	return resp
}

func Settles(r *httpserver.Request) *httpserver.Response {
	typ := r.QueryParams.Get("type")
	coin := r.QueryParams.Get("coin")
	if coin == "" {
		loggers.Warn.Printf("Settles invalid query params no coin code input")
		return httpserver.NewResponseWithError(errors.BadRequest)
	}
	var settleTypes []domain.SETTLE_TYPE
	switch typ {
	case "withdraw":
		settleTypes = append(settleTypes, domain.SETTLE_TYPE_WITHDRAW)
		settleTypes = append(settleTypes, domain.SETTLE_TYPE_TRANSFER_OUT)
		settleTypes = append(settleTypes, domain.SETTLE_TYPE_OTC_OUT)
	case "deposit":
		settleTypes = append(settleTypes, domain.SETTLE_TYPE_DEPOSIT)
		settleTypes = append(settleTypes, domain.SETTLE_TYPE_TRANSFER_IN)
		settleTypes = append(settleTypes, domain.SETTLE_TYPE_OTC_IN)
	case "bouns":
		settleTypes = append(settleTypes, domain.SETTLE_TYPE_BONUS)
	case "airdrop":
		settleTypes = append(settleTypes, domain.SETTLE_TYPE_AIRDROP)
	case "transferE2C":
		settleTypes = append(settleTypes, domain.SETTLE_TYPE_TRANSFER_E2C)
	case "transferC2E":
		settleTypes = append(settleTypes, domain.SETTLE_TYPE_TRANSFER_C2E)
	case "OTCOUT":
		settleTypes = append(settleTypes, domain.SETTLE_TYPE_OTC_OUT)
	case "OTCIN":
		settleTypes = append(settleTypes, domain.SETTLE_TYPE_OTC_IN)
	}
	limit, err := strconv.Atoi(r.QueryParams.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(r.QueryParams.Get("offset"))
	if err != nil || offset <= 0 {
		offset = 0
	}

	settles, err := adaptor.GetSettles(settleTypes, coin, r.Uid, limit, offset)
	if err != nil {
		loggers.Error.Printf("Settles GetSettles error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}
	total, err := adaptor.GetTotalSettles(settleTypes, coin, r.Uid)
	if err != nil {
		loggers.Error.Printf("Settles GetSettles error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}
	result := domain.SettlesResp{
		Total:   total,
		Settles: *settles,
	}
	resp := httpserver.NewResponse()
	resp.Data = result
	return resp
}

func SettlesV2(r *httpserver.Request) *httpserver.Response {
	typ := r.QueryParams.Get("type")
	coin := r.QueryParams.Get("coin")
	if coin == "" {
		loggers.Warn.Printf("Settles invalid query params no coin code input")
		return httpserver.NewResponseWithError(errors.BadRequest)
	}
	var settleType domain.SETTLE_TYPE
	if typ == "withdraw" {
		settleType = domain.SETTLE_TYPE_WITHDRAW
	} else if typ == "deposit" {
		settleType = domain.SETTLE_TYPE_DEPOSIT
	} else {
		settleType = 0
	}
	limit, err := strconv.Atoi(r.QueryParams.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(r.QueryParams.Get("offset"))
	if err != nil || offset <= 0 {
		offset = 0
	}

	settles, err := adaptor.GetSettlesV2(settleType, coin, r.Uid, limit, offset)
	if err != nil {
		loggers.Error.Printf("Settles GetSettles error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}

	total, err := adaptor.GetTotalSettlesV2(settleType, coin, r.Uid)
	if err != nil {
		loggers.Error.Printf("Settles GetSettles error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}
	result := domain.SettlesResp{
		Total:   total,
		Settles: *settles,
	}
	resp := httpserver.NewResponse()
	resp.Data = result
	return resp
}

func RechargeAmount(r *httpserver.Request) *httpserver.Response {
	var req domain.RechargeAmount
	if err := r.Parse(&req); err != nil || r.Uid < 1 {
		loggers.Warn.Printf("RechargeAmount parse body error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.ParameterError)
	}
	req.AccountId = r.Uid
	if err := adaptor.AddRechargeAmount(&req); nil != err {
		return httpserver.NewResponseWithError(errors.SysError)
	} else {
		return httpserver.NewResponse()
	}
}

func AddWithdrawAddress(r *httpserver.Request) *httpserver.Response {
	var (
		req   domain.UserWithdrawAddress
		reply *domain.UserWithdrawAddress
		err   error
	)
	if err = r.Parse(&req); err != nil {
		loggers.Warn.Printf("AddWithdrawAddress parse body error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.ParameterError)
	}
	if req.Address == "" {
		loggers.Warn.Printf("AddWithdrawAddress no address input")
		return httpserver.NewResponseWithError(errors.NoAccountId)
	}
	if req.Type == domain.WITHDRAW_ADDRESS_TYPE_CHAIN || req.Type == 0 {
		req.Type = domain.WITHDRAW_ADDRESS_TYPE_CHAIN
		if req.CoinCode == "" {
			loggers.Warn.Printf("AddWithdrawAddress no coin code input")
			return httpserver.NewResponseWithError(errors.NoCoinCode)
		}
		req.Address, err = adaptor.CheckAddress(req.CoinCode, req.Address)
		if err != nil {
			loggers.Warn.Printf("AddWithdrawAddress check address error:%s", err.Error())
			return httpserver.NewResponseWithError(errors.InvalidAddress)
		}
	} else {
		_, err := strconv.ParseInt(req.Address, 10, 64)
		if err != nil {
			loggers.Warn.Printf("AddWithdrawAddress check address error:%s", err.Error())
			return httpserver.NewResponseWithError(errors.InvalidAddress)
		}
	}

	req.AccountId = r.Uid
	reply, err = adaptor.AddWithdrawAddress(&req)
	if err != nil {
		loggers.Warn.Printf("AddWithdrawAddress error:%s", err.Error())
		return httpserver.NewResponseWithError(err.(*errors.Error))
	}
	resp := httpserver.NewResponse()

	resp.Data = reply
	return resp
}

func GetWithdrawAddress(r *httpserver.Request) *httpserver.Response {
	var (
		req domain.UserWithdrawAddress
		err error
	)
	if err = r.Parse(&req); err != nil {
		loggers.Warn.Printf("UpdateWithdrawAddress parse body error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.ParameterError)
	}
	req.AccountId = r.Uid
	addresss, err := adaptor.GetWithdrawAddress(&req)
	if err != nil {
		loggers.Error.Printf("GetWithdrawAddress error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}
	resp := httpserver.NewResponse()
	resp.Data = addresss
	return resp
}

func UpdateWithdrawAddress(r *httpserver.Request) *httpserver.Response {
	var (
		req   domain.UserWithdrawAddress
		reply *domain.UserWithdrawAddress
		err   error
	)
	if err = r.Parse(&req); err != nil {
		loggers.Warn.Printf("UpdateWithdrawAddress parse body error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.ParameterError)
	}
	if req.Address == "" {
		loggers.Warn.Printf("UpdateWithdrawAddress no address input")
		return httpserver.NewResponseWithError(errors.NoAccountId)
	}
	if req.CoinCode == "" {
		loggers.Warn.Printf("UpdateWithdrawAddress no coin code input")
		return httpserver.NewResponseWithError(errors.NoCoinCode)
	}
	req.Address, err = adaptor.CheckAddress(req.CoinCode, req.Address)
	if err != nil {
		loggers.Warn.Printf("UpdateWithdrawAddress check address error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.InvalidAddress)
	}
	req.AccountId = r.Uid
	reply, err = adaptor.UpdateWithdrawAddress(&req)
	if err != nil {
		loggers.Warn.Printf("UpdateWithdrawAddress error:%s", err.Error())
		return httpserver.NewResponseWithError(err.(*errors.Error))
	}
	resp := httpserver.NewResponse()
	resp.Data = reply
	return resp
}

func DeleteWithdrawAddress(r *httpserver.Request) *httpserver.Response {
	var (
		req domain.UserWithdrawAddress
		err error
	)
	if err = r.Parse(&req); err != nil {
		loggers.Warn.Printf("DeleteWithdrawAddress parse body error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.ParameterError)
	}
	if req.Address == "" {
		loggers.Warn.Printf("DeleteWithdrawAddress no address input")
		return httpserver.NewResponseWithError(errors.NoAccountId)
	}
	if req.CoinCode == "" {
		loggers.Warn.Printf("DeleteWithdrawAddress no coin code input")
		return httpserver.NewResponseWithError(errors.NoCoinCode)
	}
	req.AccountId = r.Uid
	if err := adaptor.DeleteWithdrawAddress(&req); err != nil {
		loggers.Warn.Printf("DeleteWithdrawAddress error:%s", err.Error())
		return httpserver.NewResponseWithError(err.(*errors.Error))
	}
	return httpserver.NewResponse()
}
