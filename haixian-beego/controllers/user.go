package controllers

import (
	"gitee.com/fredgo/back/beego-haixian/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)
type UserController struct {
	beego.Controller
}
func (c *UserController)ShowReg()  {
	c.TplName = "register.html"
}
func (c *UserController)HandleReg()  {
	userName := c.GetString("userName")
	pwd :=c.GetString("pwd")
	beego.Info(userName,pwd)
	if userName == "" || pwd == ""{
		beego.Info("userName,pwd can not be empty!")
		c.Redirect("/Register",302)
		return
	}
	o := orm.NewOrm()
	user := models.User{}
	user.Name = userName
	user.Pwd = pwd
	_,err := o.Insert(&user)
	if err!=nil{
		beego.Info("insert error!",err)
		c.Redirect("/Register",302)
		return
	}
	//c.Ctx.WriteString("注册成功！")

	c.TplName = "login.html"
	c.Redirect("/Login",302)
}

func  (c * UserController)ShowLogin()  {
	name :=c.Ctx.GetCookie("name")
	if name != ""{
		c.Data["name"]=name
		c.Data["check"]="checked"
	}
	c.TplName = "login.html"
}

func (c * UserController)HandleLogin()  {
	name := c.GetString("userName")
	pwd := c.GetString("pwd")
	check := c.GetString("remember")
	o := orm.NewOrm()
	user := models.User{}
	user.Name = name
	err := o.Read(&user,"name")
	if err!=nil{
		beego.Info("用户名实效！")
		c.TplName = "login.html"
		return
	}
	if user.Pwd != pwd {
		beego.Info("密码实效！")
		c.TplName = "login.html"
		return
	}
	if check=="on" {
		c.Ctx.SetCookie("name",name,time.Second*3600)
	}else {
		c.Ctx.SetCookie("name","sss",-1)
	}
	c.SetSession("name",name)
	c.Redirect("/Article/ShowArticle",302)
}

func (c * UserController)Logout()  {
	c.DelSession("name")
	c.Redirect("/Login",302)
}