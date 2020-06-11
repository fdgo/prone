package main

import (
	"business/services/base/notifycenter/app"
	"business/support/config"
	"errors"
)

func main() {
	serverConf, ok := config.Conf.Services["NotifyCenter"]
	if !ok {
		panic(errors.New("No NotifyCenter server config"))
	}
	app.Start(serverConf.ListenAddress)
}
