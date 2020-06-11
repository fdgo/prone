package main

import (
	"business/services/base/usercenter/app"
	"business/support/config"
	"errors"
)

func main() {
	serverConf, ok := config.Conf.Services["UserCenter"]
	if !ok {
		panic(errors.New("No UserCenter server config"))
	}
	app.Start(serverConf.ListenAddress)
}
