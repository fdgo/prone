package main

import (
	"github.com/astaxie/beego"
	_ "gitee.com/fredgo/back/beego-haixian/routers"
	_ "gitee.com/fredgo/back/beego-haixian/models"
	"strconv"
)

func main() {
	beego.AddFuncMap("ShowPrePage",handlePrePage)
	beego.AddFuncMap("ShowNextPage",handleNextPage)
	beego.Run()
}

func handlePrePage(data int) string  {
	pageIndex := data -1
	return strconv.Itoa(pageIndex)
}
func handleNextPage(data int ) string  {
	pageIndex := data +1
	return strconv.Itoa(pageIndex)
}
