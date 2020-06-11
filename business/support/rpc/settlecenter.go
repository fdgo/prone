package rpc

import (
	"business/support/domain"
	"context"

	"github.com/smallnest/rpcx/client"
)

type BankRPC struct {
	Client client.XClient
}

func NewBankRPC() *BankRPC {
	name := "BankRPC"
	return &BankRPC{
		Client: NewClient(name),
	}
}

func (r *BankRPC) WithdrawReview(ctx context.Context, args *domain.Settle, reply *domain.Settle) error {
	return r.Client.Call(ctx, "WithdrawReview", args, reply)
}

func (r *BankRPC) WithdrawUnfreeze(ctx context.Context, args *domain.Settle, reply *domain.Settle) error {
	return r.Client.Call(ctx, "WithdrawUnfreeze", args, reply)
}

func (r *BankRPC) WithdrawSuccess(ctx context.Context, args *domain.Settle, reply *domain.Settle) error {
	return r.Client.Call(ctx, "WithdrawSuccess", args, reply)
}

func (r *BankRPC) WithdrawFailed(ctx context.Context, args *domain.Settle, reply *domain.Settle) error {
	return r.Client.Call(ctx, "WithdrawFailed", args, reply)
}

func (r *BankRPC) Deposit(ctx context.Context, args *domain.Settle, reply *domain.Settle) error {
	return r.Client.Call(ctx, "Deposit", args, reply)
}

func (r *BankRPC) Destroy(ctx context.Context, args *domain.Settle, reply *domain.Settle) error {
	return r.Client.Call(ctx, "Destroy", args, reply)
}
