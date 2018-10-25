package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"newsWeb/models"
	"database/sql"
	"encoding/base64"
)

type UserController struct {
	beego.Controller
}

func (this *UserController)ShowRegister()  {

	this.TplName="register.html"
}
func (this *UserController)HandleRegister()  {
	//1.接受数据
	userName:=this.GetString("userName")
	pwd:=this.GetString("password")
	//2.校验数据
	if userName==""||pwd==""{
		this.Data["errmsg"]="用户名密码不能为空"
		this.TplName="register.html"
		return
	}
	//查询
	/*
	conn,err :=sql.Open("mysql","root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	defer conn.Close()
	if err !=nil{
		beego.Error("数据库链接失败",err)
	}
	res ,err :=conn.Query("select name from user")
	if err != nil{
		beego.Error("查询错误",err)
	}
	var name string

	for res.Next(){
		res.Scan(&name)
		beego.Info(name)
	}*/
	conn,err:=sql.Open("mysql","root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	defer conn.Close()
	if err!=nil{
		beego.Error("数据库连接失败",err)
	}
	res,err:=conn.Query("select user_name from user")
	if err!=nil{
		beego.Error("查询错误",err)
	}
	var username string
	for res.Next(){
		res.Scan(&username)
		if userName==username{
			this.Data["errmsg"]="已经存在该用户"
			this.TplName="register.html"
			return
		}
	}
	//插入数据
	o:=orm.NewOrm()
	//获取插入对象
	var user models.User
	user.UserName=userName
	user.Pwd=pwd
	_,err =o.Insert(&user)
	if err!=nil{
		this.Data["errmsg"] = "注册失败，请重新注册！"
		this.TplName = "register.html"
		return
	}
	//返回数据
	this.Redirect("/login",302)
}
func (this *UserController)ShowLogin()  {
	dec:=this.Ctx.GetCookie("userName")
	userName,_:=base64.StdEncoding.DecodeString(dec)
	if string(userName)!=""{
		this.Data["userName"]=string(userName)
		this.Data["checked"]="checked"
	}else {
		this.Data["userName"]=""
		this.Data["checked"]=""
	}
	this.TplName="login.html"
}
func (this *UserController)HandleLogin() {
	userName := this.GetString("userName")
	pwd := this.GetString("password")
	if userName == "" || pwd == "" {
		this.Data["errmsg"] = "用户名密码不能为空"
		this.TplName = "login.html"
		return
	}
	//查询
	o := orm.NewOrm()
	var user models.User
	user.UserName=userName
	err:=o.Read(&user,"UserName")
	if err!=nil{
		this.Data["errmsg"]="用户名不存在"
		this.TplName="login.html"
		return
	}
	if user.Pwd!=pwd{
		this.Data["errmsg"]="用户名密码不正确"
		this.TplName="login.html"
		return
	}
	//是否记住登陆名
	remember:=this.GetString("remember")
	beego.Info(remember)
	if remember=="on"{
		enc:=base64.StdEncoding.EncodeToString([]byte(userName))
		this.Ctx.SetCookie("userName",enc,3600*1)
	}else {
		this.Ctx.SetCookie("userName",userName,-1)
	}
	this.SetSession("userName",userName)
	this.Redirect("/article/articleList",302)
}
func (this*UserController)Logout()  {
	//删除session
	this.DelSession("userName")
	//返回页面
	this.Redirect("/login",302)
}