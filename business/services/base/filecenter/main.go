package main

import (
	"business/services/base/filecenter/app"
	"business/support/config"
	"errors"
)

func main() {
	serverConf, ok := config.Conf.Services["FileCenter"]
	if !ok {
		panic(errors.New("No FileCenter server config"))
	}
	app.Start(serverConf.ListenAddress)
}
