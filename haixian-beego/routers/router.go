package routers

import (
	"gitee.com/fredgo/back/beego-haixian/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/Article/*",beego.BeforeRouter,FiltFunc)
	beego.Router("/Register", &controllers.UserController{},"get:ShowReg;post:HandleReg")
	beego.Router("/Login", &controllers.UserController{},"get:ShowLogin;post:HandleLogin")
	beego.Router("/Article/ShowArticle",&controllers.MainController{},"get:ShowArticleList;post:HandleSelect")
	beego.Router("/Article/AddArticle",&controllers.MainController{},"get:ShowAdd;post:HandleAdd")
	beego.Router("/Article/ArticleContent",&controllers.MainController{},"get:ShowContent")
	beego.Router("/Article/DeleteArticle",&controllers.MainController{},"get:HandleDelete")
	beego.Router("/Article/UpdateArticle",&controllers.MainController{},"get:ShowUpdate;post:HandleUpdate")
    beego.Router("/Article/AddArticleType",&controllers.MainController{},"get:ShowAddType;post:HandleAddType")
	beego.Router("/Logout",&controllers.UserController{},"get:Logout")

}
var FiltFunc = func(ctx *context.Context)  {
	name := ctx.Input.Session("name")
	if name == nil{
		ctx.Redirect(302,"/Login")
	}
}