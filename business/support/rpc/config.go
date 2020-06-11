package rpc

import (
	"business/support/domain"
	"context"

	"github.com/smallnest/rpcx/client"
)

type ConfigRPC struct {
	Client client.XClient
}

func NewConfigRPC() *ConfigRPC {
	name := "ConfigRPC"
	return &ConfigRPC{
		Client: NewClient(name),
	}
}

func (r *ConfigRPC) GetGlobalConfig(ctx context.Context, args *domain.ConfigRequest, reply *domain.ConfigResp) error {
	return r.Client.Call(ctx, "GetGlobalConfig", args, reply)
}
