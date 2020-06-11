//package models
//
//import (
//	"github.com/astaxie/beego/orm"
//	_ "github.com/go-sql-driver/mysql"
//	"time"
//)
//
//type User struct {
//	Id int `orm:"pk;auto"`
//	Name string `orm:"unique"`
//	Pwd string
//}
//
////------------------------------
//type Article struct {
//	Id int `orm:"pk;auto"`
//	ArtiName string `orm:"size(32);default(一篇好文章)"`
//	Atime time.Time `orm:"auto_now;type(datetime)"`
//	Acount int `orm:"default(0);null"`
//	Acontent string `orm:"null"`  //默认不能为NULL, 这样设置之后可以为NULL
//	Aimg string
//	BeginTime time.Time `orm:"auto_now_add;type(datetime)"` //第一次保存时才设置时间
//	UpdateTime time.Time `orm:"auto_now;type(date)"`        //保存时都会对时间自动更新
//}
//type ArticleType struct {
//	Id int
//	TypeName string `orm:"size(16)"`
//}
////--------------------------------
//type Testtable struct {
//	Id int `orm:"pk;auto"`
//	Money float64 `orm:"digits(12);decimals(2);default(0.00)"` //总共12位,小数点2位
//}
//
//func init()  {
//	orm.RegisterDataBase("default","mysql","root:000000@tcp(120.27.239.127:3306)/haixian?charset=utf8&loc=Local")
//	orm.RegisterModel(new(User),new(Article),new(Testtable))
//	orm.RunSyncdb("default",true,true)//   verbose是否可见创建过程
//}



package models

import (
"github.com/astaxie/beego/orm"
_ "github.com/go-sql-driver/mysql"
"time"
)
//----------------------多对多begin--------
type User struct {
	Id int `orm:"pk;auto"`
	Name string `orm:"unique"`
	Pwd string
	Article []*Article `orm:"rel(m2m)"`
}
//User,Article 多对多
//----------------------多对多end--------


//----------------------一对多begin--------
type Article struct {
	Id int `orm:"pk;auto"`
	ArtiName string `orm:"size(32);default(一篇好文章)"`
	Atime time.Time `orm:"auto_now;type(datetime)"`
	Acount int `orm:"default(0);null"`
	Acontent string `orm:"null"`  //默认不能为NULL, 这样设置之后可以为NULL
	Aimg string
	User []*User `orm:"reverse(many)"`                      //多对多
	ArticleType *ArticleType `orm:"rel(fk)"`                //一对多
	BeginTime time.Time `orm:"auto_now_add;type(datetime)"` //第一次保存时才设置时间
	UpdateTime time.Time `orm:"auto_now;type(date)"`        //保存时都会对时间自动更新
}
type ArticleType struct {
	Id int
	TypeName string `orm:"size(16)"`
	Article []*Article `orm:"reverse(many)"`
}
//-----------------------一对多end---------
type Testtable struct {
	Id int `orm:"pk;auto"`
	Money float64 `orm:"digits(12);decimals(2);default(0.00)"` //总共12位,小数点2位
}

func init()  {
	orm.RegisterDataBase("default","mysql","root:000000@tcp(120.27.239.127:3306)/haixian?charset=utf8&loc=Local")
	orm.RegisterModel(new(User),new(Article),new(Testtable),new(ArticleType))
	orm.RunSyncdb("default",false,true)//   verbose是否可见创建过程
}