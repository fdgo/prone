package worker

import (
	"business/support/config"
	"business/support/domain"
	"business/support/libraries/loggers"
	"business/support/rpc"
)

var (
	RPCName    string
	ServerConf domain.ServerConf
)

func init() {
	RPCName = "AccountRPC"
	conf, ok := config.Conf.Services["UserCenter"]
	if !ok {
		panic("No UserCenter config")
	}
	ServerConf = conf
}

type AccountRPC int

func Start() {
	loggers.Info.Printf("%s listen %s", RPCName, ServerConf.RPCAddress)
	if err := rpc.Serve(ServerConf.RPCAddress, RPCName, new(AccountRPC)); err != nil {
		panic(err.Error())
	}
}
