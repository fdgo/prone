package main

import (
	"access/middleware/inject"
	"access/pkg/setting"
	"access/routers"
	"fmt"
	"github.com/fvbock/endless"
	"log"
	"syscall"
)

// @title access
// @version 1.2.2
// @description  access

// @contact.name fred
// @contact.url https://github.com/fred
// @contact.email fred@sina.com

// @license.name MIT
// @license.url  https://access/blob/master/LICENSE

// @host   127.0.0.1:8000
// @BasePath
func main() {
	err := inject.LoadCasbinPolicyData()
	if err != nil {
		panic("加载casbin策略数据发生错误:  " + err.Error())
	}
	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20
	endless.DefaultReadTimeOut = readTimeout
	endless.DefaultWriteTimeOut = writeTimeout
	endless.DefaultMaxHeaderBytes = maxHeaderBytes
	server := endless.NewServer(endPoint, routersInit)
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}
