package main

import (
	"business/support/domain"
	"business/support/rpc"
	"context"
	"fmt"
)

type AccountRPC int

func (*AccountRPC) BindEmail(ctx context.Context, args *domain.Account, reply *domain.Account) error {
	reply.AccountName = "BindEmail"
	return nil
}

func (*AccountRPC) BindSettleAddress(ctx context.Context, args *domain.Account, reply *domain.Account) error {
	reply.AccountName = "BindSettleAddress"
	return nil
}
func (*AccountRPC) Create(ctx context.Context, args *domain.Account, reply *domain.Account) error {
	reply.AccountName = "Create"
	return nil
}
func (*AccountRPC) Witdraw(ctx context.Context, args *domain.Account, reply *domain.Account) error {
	reply.AccountName = "Withdraw"
	return nil
}

func main() {
	if err := rpc.Serve(":8188", "AccountRPCTest", new(AccountRPC)); err != nil {
		fmt.Println("Start account rpc server error")
	}
}
