package main

import (
	"business/support/domain"
	"business/support/rpc"
	"context"
	"fmt"
)

func main() {
	args := domain.Account{
		AccountName: "bifund-TEST",
	}
	reply := rpc.BindEmailReply{}

	c := rpc.NewAccountRPC()

	for i := 0; i < 10; i++ {
		err := c.BindEmail(context.Background(), &args, &reply)
		if err != nil {
			fmt.Printf("Call BindEmail error:%s\n", err.Error())
			continue
		}
		fmt.Println(reply.Account.Email)
	}
}
