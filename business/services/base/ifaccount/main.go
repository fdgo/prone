package main

import (
	"business/services/base/ifaccount/app"
	"business/support/config"
	"errors"
)

func main() {
	serverConf, ok := config.Conf.Services["IFAccount"]
	if !ok {
		panic(errors.New("No IFAccount server config"))
	}
	app.Start(serverConf.ListenAddress)
}
