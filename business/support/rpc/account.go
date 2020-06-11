package rpc

import (
	"business/support/domain"
	"business/support/libraries/errors"
	"context"

	"github.com/smallnest/rpcx/client"
)

type AccountRPC struct {
	Client client.XClient
}

func NewAccountRPC() *AccountRPC {
	name := "AccountRPC"
	return &AccountRPC{
		Client: NewClient(name),
	}
}

type BindSettleAddressReply struct {
	Address *domain.UserDepositAddress `json:"address"`
	Error   *errors.Error              `json:"error"`
}

type CreateAccountReply struct {
	Account *domain.Account `json:"account"`
	Error   *errors.Error   `json:"error"`
}

type WithdrawReply struct {
	Settle *domain.Settle `json:"settle"`
	Error  *errors.Error  `json:"error"`
}

type ActiveAccountReply struct {
	Account *domain.Account `json:"account"`
	Error   *errors.Error   `json:"error"`
}

type ResetPasswordReply struct {
	Account *domain.Account `json:"account"`
	Error   *errors.Error   `json:"error"`
}

type AddUserAssetsReply struct {
	IsOk  bool          `json:"is_ok"`
	Error *errors.Error `json:"error"`
}

type BindEmailReply struct {
	Account *domain.Account `json:"account"`
	Error   *errors.Error   `json:"error"`
}

type BindPhoneReply struct {
	Account *domain.Account `json:"account"`
	Error   *errors.Error   `json:"error"`
}

func (r *AccountRPC) BindSettleAddress(ctx context.Context, args *domain.UserDepositAddress, reply *BindSettleAddressReply) error {
	return r.Client.Call(ctx, "BindSettleAddress", args, reply)
}

func (r *AccountRPC) Create(ctx context.Context, args *domain.Account, reply *CreateAccountReply) error {
	return r.Client.Call(ctx, "Create", args, reply)
}

func (r *AccountRPC) Witdraw(ctx context.Context, args *domain.Settle, reply *WithdrawReply) error {
	return r.Client.Call(ctx, "Withdraw", args, reply)
}

func (r *AccountRPC) ActiveAccount(ctx context.Context, args *domain.Account, reply *ActiveAccountReply) error {
	return r.Client.Call(ctx, "ActiveAccount", args, reply)
}

func (r *AccountRPC) ResetPassword(ctx context.Context, args *domain.Account, reply *ResetPasswordReply) error {
	return r.Client.Call(ctx, "ResetPassword", args, reply)
}

func (r *AccountRPC) AddUserAssets(ctx context.Context, args *domain.UserAssets, reply *AddUserAssetsReply) error {
	return r.Client.Call(ctx, "AddUserAssets", args, reply)
}

func (r *AccountRPC) BindEmail(ctx context.Context, args *domain.Account, reply *BindEmailReply) error {
	return r.Client.Call(ctx, "BindEmail", args, reply)
}

func (r *AccountRPC) BindPhone(ctx context.Context, args *domain.Account, reply *BindPhoneReply) error {
	return r.Client.Call(ctx, "BindPhone", args, reply)
}

func (r *AccountRPC) AddWithdrawAddress(ctx context.Context, args *domain.UserWithdrawAddress, reply *domain.UserWithdrawAddress) error {
	return r.Client.Call(ctx, "AddWithdrawAddress", args, reply)
}

func (r *AccountRPC) DeleteWithdrawAddress(ctx context.Context, args *domain.UserWithdrawAddress, reply *domain.UserWithdrawAddress) error {
	return r.Client.Call(ctx, "DeleteWithdrawAddress", args, reply)
}

func (r *AccountRPC) UpdateWithdrawAddress(ctx context.Context, args *domain.UserWithdrawAddress, reply *domain.UserWithdrawAddress) error {
	return r.Client.Call(ctx, "UpdateWithdrawAddress", args, reply)
}

func (r *AccountRPC) AddAssetPassword(ctx context.Context, args *domain.Account, reply *domain.Account) error {
	return r.Client.Call(ctx, "AddAssetPassword", args, reply)
}

func (r *AccountRPC) ResetAssetPassword(ctx context.Context, args *domain.Account, reply *domain.Account) error {
	return r.Client.Call(ctx, "ResetAssetPassword", args, reply)
}

func (r *AccountRPC) ResetAssetPasswordEffectiveTime(ctx context.Context, args *domain.Account, reply *domain.Account) error {
	return r.Client.Call(ctx, "ResetAssetPasswordEffectiveTime", args, reply)
}

func (r *AccountRPC) SetAntiFishingText(ctx context.Context, args *domain.Account, reply *domain.Account) error {
	return r.Client.Call(ctx, "SetAntiFishingText", args, reply)
}

func (r *AccountRPC) SetGAKey(ctx context.Context, args *domain.Account, reply *domain.Account) error {
	return r.Client.Call(ctx, "SetGAKey", args, reply)
}

func (r *AccountRPC) KYCUpdate(ctx context.Context, args *domain.Account, reply *domain.Account) error {
	return r.Client.Call(ctx, "KYCUpdate", args, reply)
}

func (r *AccountRPC) UpdateAccountAvatar(ctx context.Context, args *domain.Account, reply *domain.Account) error {
	return r.Client.Call(ctx, "UpdateAccountAvatar", args, reply)
}

func (r *AccountRPC) SetAccountName(ctx context.Context, args *domain.Account, reply *domain.Account) error {
	return r.Client.Call(ctx, "SetAccountName", args, reply)
}
