package controllers

import (
	"gitee.com/fredgo/back/beego-haixian/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"path"
	"strconv"
	"time"
)

type MainController struct {
	beego.Controller
}

//---------------------------------------------------

func (c *MainController)ShowIndex()  {

	//pageIndex := 1
	pageIndex1 := c.GetString("pageIndex")
	pageIndex,err :=  strconv.Atoi(pageIndex1)
	if err != nil{
		pageIndex = 1
	}

	o := orm.NewOrm()
	var articles []models.Article
	qs := o.QueryTable("Article")
	//qs.All(&articles)
	count,err :=qs.Count()
	if err!=nil{
		beego.Info("查询所有文章信息出错！")
		return
	}
	pageSize := 1
	qs.Limit( pageSize, pageSize* (pageIndex-1) ).All(&articles)
	pageCount  := math.Ceil( float64(count)/float64(pageSize) )

	c.Data["articles"] = articles
	c.Data["count"] = count
	c.Data["pageCount"] = pageCount
	c.Data["pageIndex"] = pageIndex
	c.TplName = "index.html"
}
func (c *MainController)ShowAdd()  {
	o := orm.NewOrm()
	var types [] models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	c.Data["types"] = types
	c.TplName = "add.html"
}
func (c *MainController)HandleAdd()  {
	articleName := c.GetString("articleName")
	articleContent := c.GetString("content")
	f,h,err := c.GetFile("uploadname")
	defer f.Close()


	//1.要限制格式
	fileext := path.Ext(h.Filename)
	if fileext !=".jpg" && fileext!= ".png" && fileext!= ".icon"{
		beego.Info("上传文件格式错误！")
		return
	}
	//2.限制大小
	if h.Size > 1024*1024*3{
		beego.Info("上传文件太大！")
		return
	}
	//3.需要对文件重命名
	filename := time.Now().Format("2006-01-02 15:04:05") + fileext

	if err!=nil{
		beego.Info("上传失败!")
		return
	}else {
		c.SaveToFile("uploadname","./static/img/" + filename )
	}
	if articleContent == "" || articleName == ""{
		beego.Info("添加文章数据错误！")
		return
	}
	o := orm.NewOrm()
	arti := models.Article{}
	arti.ArtiName = articleName
	arti.Acontent = articleContent
	arti.Aimg = "/static/img/"+filename


	typeName :=c.GetString("select")
	if typeName == ""{
		beego.Info("下拉框数据错误")
		return
	}
	var artiType models.ArticleType
	artiType.TypeName = typeName
	err =o.Read(&artiType,"TypeName")
	if err != nil{
		beego.Info("获取类型错误！")
	}
	arti.ArticleType = &artiType

	_,err = o.Insert(&arti)
	if err !=nil{
		beego.Info("插入数据库错误!",err.Error())
		return
	}
	c.Redirect("/Article/ShowArticle",302)

}
func (c *MainController)ShowContent()  {
	id2 := c.GetString("id")
	beego.Info(id2)
	id,_ := strconv.Atoi(id2)

	o := orm.NewOrm()
	arti := models.Article{Id:id}
	err := o.Read(&arti)
	if err !=nil{
		beego.Info("查询错误!",err.Error())
		return
	}
	arti.Acount +=1
	m2m := o.QueryM2M(&arti,"user")
	name := c.GetSession("name")
	user := models.User{}
	user.Name = name.(string)
	o.Read(&user,"name")
	_, err = m2m.Add(&user)
	if err!=nil{
		beego.Info("插入失败！")
		return
	}
	o.Update(&arti)
	//o.LoadRelated(&arti,"user")
	var xxx []models.Article
	o.QueryTable("Article").RelatedSel("User").Filter("User__User__Name",name.(string)).Distinct().Filter("Id",id).One(&xxx)


	c.Data["article"] = arti
	c.Layout = "layout.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["contentHead"]="head.html"
	c.TplName = "content.html"
}
//显示编辑界面
func (c *MainController)ShowUpdate()  {

	id, err := c.GetInt("id")
	if err !=nil{
		beego.Info("获取文章Id错误！",err)
		return
	}
	o := orm.NewOrm()
	arti := models.Article{Id:id}
	err = o.Read(&arti)
	if err !=nil{
		beego.Info("查询错误!",err.Error())
		return
	}
	c.Data["article"] = arti

	c.TplName = "update.html"
}
func (c *MainController)HandleUpdate()  {
	id,_ := c.GetInt("id")

	artiName := c.GetString("articleName")
	content := c.GetString("content")

	f,h,err := c.GetFile("uploadname")

	var fileName string

	if err!=nil{
		beego.Info("上传失败!")
		return
	}else {

		defer f.Close()
		//1.要限制格式
		fileext := path.Ext(h.Filename)
		if fileext !=".jpg" && fileext!= ".png"{
			beego.Info("上传文件格式错误！")
			return
		}
		//2.限制大小
		if h.Size > 1024*1024*3{
			beego.Info("上传文件太大！")
			return
		}
		//3.需要对文件重命名
		fileName = time.Now().Format("2006-01-02 15:04:05") + fileext
		c.SaveToFile("uploadname","./static/img/" + fileName )
	}
	if artiName == "" || content == "" {
		beego.Info("更新数据获取失败!")
		return
	}

	o := orm.NewOrm()
	arti := models.Article{Id:id}
	err = o.Read(&arti)
	if err != nil{
		beego.Info("查询数据错误！")
		return
	}
	arti.ArtiName = artiName
	arti.Acontent = content
	arti.Aimg = "./static/img/" + fileName
	_,err = o.Update(&arti,"ArtiName","Acontent","Aimg")
	if err != nil {
		beego.Info("更新数据显示错误！")
		return
	}
	c.Redirect("/Article/ShowArticle",302)

}
func (c *MainController) HandleDelete()  {
	id, err := c.GetInt("id")
	if err!= nil{
		beego.Info("获取id数据错误！")
		return
	}
	o := orm.NewOrm()
	arti := models.Article{Id:id}
	err = o.Read(&arti)
	if err != nil{
		beego.Info("查询错误！")
		return
	}
	o.Delete(&arti)
	c.Redirect("/Article/ShowArticle",302)
}
func (c *MainController)HandleSelect()  {
	typeName := c.GetString("select")
	if typeName == ""{
		beego.Info("下拉框传递数据失败！")
		return
	}
	o := orm.NewOrm()
	var articles[]models.Article
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
	//beego.Info(articles)

}
func (c *MainController)ShowArticleList()  {
	//pageIndex := 1
	pageIndex1 := c.GetString("pageIndex")
	pageIndex,err :=  strconv.Atoi(pageIndex1)
	if err != nil{
		pageIndex = 1
	}

	o := orm.NewOrm()
	var articles []models.Article
	qs := o.QueryTable("Article")
	//
	typeName := c.GetString("select")
	//
	count,err :=qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()
	if err!=nil{
		beego.Info("查询所有文章信息出错！")
		return
	}
	pageSize := 2
	qs.Limit( pageSize, pageSize* (pageIndex-1) ).RelatedSel("ArticleType").All(&articles)
	pageCount  := math.Ceil( float64(count)/float64(pageSize) )

	FirstPage := false
	if pageIndex == 1 {
		FirstPage = true
	}
	LastPage := false
	if pageIndex == int(pageCount){
		LastPage = true
	}

	var types [] models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	c.Data["types"] = types



	var articleswithtype []models.Article
	if typeName == ""{
		beego.Info("下拉框传递数据失败！")
		qs.Limit( pageSize, pageSize* (pageIndex-1) ).RelatedSel("ArticleType").All(&articleswithtype)
	}else {
		qs.Limit(pageSize,pageSize* (pageIndex-1)).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articleswithtype)
	}
	c.Data["TypeName"]=typeName
	c.Data["FirstPage"] = FirstPage
	c.Data["LastPage"] = LastPage
	c.Data["count"] = count
	c.Data["pageCount"] = pageCount
	c.Data["pageIndex"] = pageIndex
	c.Data["articles"] = articleswithtype

	c.Layout = "layout.html"
	c.TplName = "index.html"
}

func (c *MainController)ShowAddType()  {
	o := orm.NewOrm()
	var artiTypes []models.ArticleType
	_,err := o.QueryTable("ArticleType").All(&artiTypes)
	if err != nil{
		beego.Info("查询类型错误！")
		return
	}
	c.Data["types"] = artiTypes
	c.TplName = "addType.html"
}
func (c *MainController)HandleAddType()  {
	typename := c.GetString("typeName")
	if typename == ""{
		beego.Info("添加类型数据为空!")
		return
	}
	o := orm.NewOrm()
	var artiType models.ArticleType
	artiType.TypeName = typename
	_,err := o.Insert(&artiType)
	if err != nil{
		beego.Info("插入类型数据失败！")
		return
	}
	c.Redirect("/Article/AddArticleType",302)
}
// 1XX  请求已经被接受，需要继续发送请求   100
// 2XX  请求成功  200
// 3XX  请求资源被转移，请求被转接 302
// 4xx  请求失败                404
// 5xx  服务器错误              500