package worker

import (
	"business/services/base/usercenter/adaptor"
	"business/support/domain"
	"business/support/libraries/errors"
	"business/support/libraries/loggers"
	"business/support/rpc"
	"context"
)

func (*AccountRPC) Create(ctx context.Context, args *domain.Account, reply *rpc.CreateAccountReply) error {
	if nil == reply {
		return errors.InternalServerError
	}
	if args.AccountType == domain.ACCOUNT_TYPE_EMAIL {
		eAc, err := adaptor.GetAccountByEmail(args.Email)
		if err != nil && err != errors.NotFound {
			loggers.Debug.Printf("handleCreateAccount error:%s", err.Error())
			reply.Error = errors.InternalServerError
			return nil
		}

		if nil != eAc {
			loggers.Debug.Printf("handleCreateAccount email:%s", args.Email)
			reply.Error = errors.AccountHadExisted
			return nil
		}
	} else {
		eAc, err := adaptor.GetAccountByPhone(args.Phone)
		if err != nil && err != errors.NotFound {
			loggers.Debug.Printf("handleCreateAccount error:%s", err.Error())
			reply.Error = errors.InternalServerError
			return nil
		}

		if nil != eAc {
			loggers.Debug.Printf("handleCreateAccount phone:%s", args.Phone)
			reply.Error = errors.AccountHadExisted
			return nil
		}
	}
	account, err := adaptor.AddAccount(args)
	if err != nil {
		loggers.Debug.Printf("handleCreateAccount error:%s", err.Error())
		reply.Error = errors.InternalServerError
		return nil
	}
	reply.Account = account
	return nil
}

func (*AccountRPC) ActiveAccount(ctx context.Context, args *domain.Account, reply *rpc.ActiveAccountReply) error {
	if nil == reply {
		return errors.InternalServerError
	}
	if args.AccountId == 0 {
		loggers.Debug.Printf("handleActiveAccount no account id")
		reply.Error = errors.NoAccountId
		return nil
	}
	if err := adaptor.ActiveAccount(args.AccountId); err != nil {
		loggers.Debug.Printf("handleActiveAccount error:%s", err.Error())
		reply.Error = errors.InternalServerError
		return nil
	}
	return nil
}

func (*AccountRPC) ResetPassword(ctx context.Context, args *domain.Account, reply *rpc.ResetPasswordReply) error {
	if nil == reply {
		return errors.InternalServerError
	}
	if args.AccountId == 0 {
		loggers.Debug.Printf("handleResetPassword no account id")
		reply.Error = errors.NoAccountId
		return nil
	}
	if args.Password == "" {
		loggers.Debug.Printf("handleResetPassword no password")
		reply.Error = errors.NoPassword
		return nil
	}
	if err := adaptor.ResetPassword(args); err != nil {
		loggers.Debug.Printf("handleResetPassword error:%s", err.Error())
		reply.Error = errors.InternalServerError
		return nil
	}
	reply.Account = args
	return nil
}

func (r *AccountRPC) AddUserAssets(ctx context.Context, args *domain.UserAssets, reply *rpc.AddUserAssetsReply) error {
	loggers.Debug.Printf("rpc AddUserAssets,arg:%s,reply:%#v", args.Desc(), reply)
	if nil == args || nil == reply || !domain.CheckVol(args.AvailableVol) {
		return rpc.SysError
	}
	ok, _ := adaptor.AddUserAsserts(args)
	if !ok {
		reply.IsOk = false
		reply.Error = errors.NoAccountId
	} else {
		reply.IsOk = true
	}
	return nil
}

func (r *AccountRPC) BindEmail(ctx context.Context, args *domain.Account, reply *rpc.BindEmailReply) error {
	if nil == args || nil == reply {
		return rpc.SysError
	}
	if args.Email == "" {
		reply.Error = errors.NoEmail
		return nil
	}
	if args.AccountId == 0 {
		reply.Error = errors.NoAccountId
		return nil
	}
	update := domain.Account{
		Email: args.Email,
	}
	if err := adaptor.UpdateAccountByUid(args.AccountId, &update); err != nil {
		reply.Error = errors.InternalServerError
		return nil
	}
	return nil
}

func (r *AccountRPC) BindPhone(ctx context.Context, args *domain.Account, reply *rpc.BindPhoneReply) error {
	if nil == args || nil == reply {
		return rpc.SysError
	}
	if args.Email == "" {
		reply.Error = errors.NoEmail
		return nil
	}
	if args.AccountId == 0 {
		reply.Error = errors.NoAccountId
		return nil
	}
	update := domain.Account{
		Phone: args.Phone,
	}
	if err := adaptor.UpdateAccountByUid(args.AccountId, &update); err != nil {
		reply.Error = errors.InternalServerError
		return nil
	}
	return nil
}

func (r *AccountRPC) AddAssetPassword(ctx context.Context, args *domain.Account, reply *domain.Account) error {
	if args == nil || reply == nil {
		return errors.ParameterError
	}
	if args.AccountId == 0 {
		return errors.NoAccountId
	}
	if args.AssetPassword == "" {
		return errors.NoAssetPassword
	}
	if err := adaptor.AddAssetPassword(args); err != nil {
		loggers.Warn.Printf("AddAssetPassword error %s", err.Error())
		if err == errors.AccountNotFound || err == errors.AssetPasswordHadExisted {
			return err
		}
		return errors.InternalServerError
	}
	reply = args
	return nil
}

func (r *AccountRPC) ResetAssetPassword(ctx context.Context, args *domain.Account, reply *domain.Account) error {
	if args == nil || reply == nil {
		return errors.ParameterError
	}
	if args.AccountId == 0 {
		return errors.NoAccountId
	}
	if args.AssetPassword == "" {
		return errors.NoAssetPassword
	}
	if err := adaptor.ResetAssetPassword(args); err != nil {
		loggers.Warn.Printf("ResetAssetPassword error %s", err.Error())
		if err == errors.AccountNotFound {
			return err
		}
		return errors.InternalServerError
	}
	reply = args
	return nil
}

func (r *AccountRPC) ResetAssetPasswordEffectiveTime(ctx context.Context, args *domain.Account, reply *domain.Account) error {
	if args == nil || reply == nil {
		return errors.ParameterError
	}
	if args.AccountId == 0 {
		return errors.NoAccountId
	}
	if err := adaptor.ResetAssetPasswordEffectiveTime(args); err != nil {
		loggers.Warn.Printf("ResetAssetPasswordEffectiveTime error %s", err.Error())
		if err == errors.AccountNotFound {
			return err
		}
		return errors.InternalServerError
	}
	reply = args
	return nil
}

func (r *AccountRPC) SetAntiFishingText(ctx context.Context, args *domain.Account, reply *domain.Account) error {
	if args == nil || reply == nil {
		return errors.ParameterError
	}
	if args.AccountId == 0 {
		return errors.NoAccountId
	}
	if err := adaptor.SetAntiFishingText(args); err != nil {
		loggers.Warn.Printf("SetAntiFishingText error :%s", err.Error())
		if err == errors.AccountNotFound {
			return err
		}
		return errors.InternalServerError
	}
	reply = args
	return nil
}

func (r *AccountRPC) SetGAKey(ctx context.Context, args *domain.Account, reply *domain.Account) error {
	if args == nil || reply == nil {
		return errors.ParameterError
	}

	if args.AccountId == 0 {
		return errors.NoAccountId
	}

	// 如果GAKey 为空，表示取消google验证码
	if err := adaptor.SetGAKey(args); err != nil {
		loggers.Warn.Printf("SetGAKey error %s", err.Error())
		if err == errors.AccountNotFound {
			return err
		}
		return errors.InternalServerError
	}
	reply = args
	return nil
}

func (r *AccountRPC) KYCUpdate(ctx context.Context, args *domain.Account, reply *domain.Account) error {
	if args == nil || reply == nil {
		return errors.ParameterError
	}
	if args.AccountId == 0 {
		return errors.NoAccountId
	}
	if args.KYCStatus == domain.KYC_STATUS_APPROVED {
		if args.IDNo == "" {
			return errors.ParameterError
		}
		ok, err := adaptor.CheckKYCIDNo(args.IDNo)
		if err != nil {
			return errors.InternalServerError
		}
		if !ok {
			return errors.KYCIDNoHadExisted
		}
	}

	if err := adaptor.KYCInfoUpdate(args); err != nil {
		loggers.Warn.Printf("KYCUpdate error %s", err.Error())
		if err == errors.AccountNotFound {
			return err
		}
		if err == errors.InvalidKYCStatus {
			return err
		}
		return errors.InternalServerError
	}
	reply = args
	return nil
}

func (r *AccountRPC) SetAccountName(ctx context.Context, args *domain.Account, reply *domain.Account) error {
	if args == nil || reply == nil || args.AccountName == "" || args.AccountId == 0 {
		return errors.ParameterError
	}
	account, err := adaptor.GetAccountByUid(args.AccountId)
	if err != nil {
		return errors.InternalServerError
	}
	if account.AccountName != "" {
		return errors.AccountNameHadExisted
	}
	if !adaptor.CheckAccountName(args.AccountName) {
		return errors.AccountNameHadBeenUsed
	}

	if err := adaptor.UpdateAccountName(account.AccountId, args.AccountName); err != nil {
		return err
	}
	account.AccountName = args.AccountName

	reply = account.UserInfo()
	return nil
}

func (r *AccountRPC) UpdateAccountAvatar(ctx context.Context, args *domain.Account, reply *domain.Account) error {
	if args == nil || reply == nil || args.Avatar == "" || args.AccountId == 0 {
		return errors.ParameterError
	}
	if err := adaptor.UpdateAccountAvatar(args.AccountId, args.Avatar); err != nil {
		return err
	}

	account, err := adaptor.GetAccountByUid(args.AccountId)
	if err != nil {
		return errors.InternalServerError
	}

	reply = account.UserInfo()
	return nil
}
