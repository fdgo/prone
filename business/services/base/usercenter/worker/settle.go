package worker

import (
	"business/services/base/usercenter/adaptor"
	"business/support/domain"
	"business/support/libraries/errors"
	"business/support/libraries/loggers"
	"business/support/rpc"
	"context"
)

func (*AccountRPC) BindSettleAddress(ctx context.Context, args *domain.UserDepositAddress, reply *rpc.BindSettleAddressReply) error {
	if nil == reply {
		return errors.InternalServerError
	}
	if args.AccountId == 0 {
		loggers.Debug.Printf("HandleBindAddress no account id")
		reply.Error = errors.NoAccountId
		return nil
	}

	if args.CoinGroup == "" {
		loggers.Debug.Printf("HandleBindAddress no coin group")
		reply.Error = errors.InvalidCoinCode
		return nil
	}

	_, err := adaptor.GetAccountByUid(args.AccountId)
	if err != nil {
		loggers.Error.Printf("HandleBindAddress account_id:%d error:%s", args.AccountId, err.Error())
		if err == errors.NotFound {
			reply.Error = errors.NotFound
		}
		reply.Error = errors.InternalServerError
		return nil
	}

	userDepositAddress, err := adaptor.BindUserDepositAddress(args)
	if err != nil {
		loggers.Error.Printf("HandleBindAddress bind deposit address error:%s", err.Error())
		reply.Error = errors.InternalServerError
		return nil
	}

	if err := adaptor.AddDepositAddressToCache(userDepositAddress); err != nil {
		loggers.Error.Printf("HandleBindAddress  error:%s", err.Error())
		reply.Error = errors.InternalServerError
		return nil
	}

	reply.Address = userDepositAddress
	return nil
}

func (*AccountRPC) Withdraw(ctx context.Context, args *domain.Settle, reply *rpc.WithdrawReply) error {
	if nil == reply {
		return errors.InternalServerError
	}
	if args.AccountId == 0 {
		loggers.Debug.Printf("HandleWithdraw no account id")
		reply.Error = errors.NoAccountId
		return nil
	}
	if args.CoinCode == "" {
		loggers.Debug.Printf("HandleWithdraw no coin code")
		reply.Error = errors.NoCoinCode
		return nil
	}
	if !domain.CheckVol(args.Vol) {
		loggers.Debug.Printf("HandleWithdraw vol invalid vol:%s", args.Vol.String())
		reply.Error = errors.VolumeIsZero
		return nil
	}
	if args.Type == domain.SETTLE_TYPE_FREEZE {
		args.Status = domain.SETTLE_STATUS_REJECTED
	} else {
		if len(args.ToAddress) < 1 {
			loggers.Debug.Printf("HandleWithdraw address invalid")
			reply.Error = errors.ParameterError
			return nil
		}
		args.Status = domain.SETTLE_STATUS_CREATED
	}
	if err := adaptor.Withdraw(args); err != nil {
		loggers.Error.Printf("HandleWithdraw account_id:%d error:%s", args.AccountId, err.Error())
		if err == errors.InsufficientBalance {
			reply.Error = errors.InsufficientBalance
		}
		reply.Error = errors.InternalServerError
		return nil
	}
	reply.Settle = args
	return nil
}

func (*AccountRPC) AddWithdrawAddress(ctx context.Context, args *domain.UserWithdrawAddress, reply *domain.UserWithdrawAddress) error {
	if args == nil || reply == nil {
		return errors.ParameterError
	}
	if args.AccountId == 0 {
		loggers.Warn.Printf("Handle AddWithdrawAddress no account id input")
		return errors.NoAccountId
	}
	if args.Address == "" {
		loggers.Warn.Printf("Handle AddWithdrawAddress no address id input")
		return errors.NoAddress
	}
	if args.Type == domain.WITHDRAW_ADDRESS_TYPE_CHAIN {
		if args.CoinCode == "" {
			loggers.Warn.Printf("Handle AddWithdrawAddress no coin code input")
			return errors.NoCoinCode
		}
	}
	if err := adaptor.AddWithdrawAddress(args); err != nil {
		loggers.Error.Printf("Handle AddWithdrawAddress error:%s", err.Error())
		return err
	}
	*reply = *args
	return nil
}

func (*AccountRPC) DeleteWithdrawAddress(ctx context.Context, args *domain.UserWithdrawAddress, reply *domain.UserWithdrawAddress) error {
	if args == nil || reply == nil {
		return errors.ParameterError
	}
	if args.AccountId == 0 {
		loggers.Warn.Printf("Handle DeleteWithdrawAddress no account id input")
		return errors.NoAccountId
	}
	if args.Address == "" {
		loggers.Warn.Printf("Handle DeleteWithdrawAddress no address id input")
		return errors.NoAddress
	}
	if args.Type == domain.WITHDRAW_ADDRESS_TYPE_CHAIN {
		if args.CoinCode == "" {
			loggers.Warn.Printf("Handle DeleteWithdrawAddress no coin code input")
			return errors.NoCoinCode
		}
	}
	if err := adaptor.DeleteWithdrawAddress(args); err != nil {
		loggers.Error.Printf("Handle DeleteWithdrawAddress error:%s", err.Error())
		return err
	}
	return nil
}

func (*AccountRPC) UpdateWithdrawAddress(ctx context.Context, args *domain.UserWithdrawAddress, reply *domain.UserWithdrawAddress) error {
	if args == nil || reply == nil {
		return errors.ParameterError
	}
	if args.AccountId == 0 {
		loggers.Warn.Printf("Handle UpdateWithdrawAddress no account id input")
		return errors.NoAccountId
	}
	if args.Address == "" {
		loggers.Warn.Printf("Handle UpdateWithdrawAddress no address id input")
		return errors.NoAddress
	}
	if args.Type == domain.WITHDRAW_ADDRESS_TYPE_CHAIN {
		if args.CoinCode == "" {
			loggers.Warn.Printf("Handle UpdateWithdrawAddress no coin code input")
			return errors.NoCoinCode
		}
	}
	if err := adaptor.UpdateWithdrawAddress(args); err != nil {
		loggers.Error.Printf("Handle UpdateWithdrawAddress error:%s", err.Error())
		return err
	}
	*reply = *args
	return nil
}
