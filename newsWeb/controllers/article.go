package controllers

import (
	"github.com/astaxie/beego"
	"time"
	"path"
	"github.com/astaxie/beego/orm"
	"newsWeb/models"
	"math"

	"strconv"
	"database/sql"
)

type ArticleController struct {
	beego.Controller
}
//展示文章列表页
func (this *ArticleController)ShowArticleList()  {
	userName := this.GetSession("userName")
	if userName == nil{
		this.Redirect("/login",302)
		return
	}
	this.Data["userName"] = userName.(string)
	//查询数据库，拿出数据，传递给视图
	//获取orm对象
	o := orm.NewOrm()
	//获取查询对象
	var articles []models.Article
	//查询
	//queryseter  高级查询使用的数据类型
	qs := o.QueryTable("Article")
	//查询所有的文章
	//qs.All(&articles)//select * from article

	//实现分页
	//获取总记录数和总页数
	count,_:= qs.Count()

	pageSize := int64(5)

	pageCount := float64(count) / float64(pageSize)

	pageCount = math.Ceil(pageCount)

	//向上取整
	//把数据传递给视图
	this.Data["count"] = count
	this.Data["pageCount"] = pageCount

	//获取首页末页数据
	pageIndex ,err := this.GetInt("pageIndex")
	if err != nil{
		pageIndex = 1
	}
	//获取分页的数据
	start := pageSize * (int64(pageIndex)  -1 )
	//RelatedSel 一对多关系表查询中，用来指定另外一张表的函数
	qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articles)  //queryseter


	//根据传递的类型获取相应的文章
	//获取数据
	typeName := this.GetString("select")
	this.Data["typeName"] = typeName

	qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)  //queryseter

	//获取所有类型
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)
	this.Data["articleTypes"] = articleTypes
	this.Data["pageIndex"] = pageIndex
	this.Data["articles"] = articles
	this.Layout="layout.html"
	this.TplName = "index.html"
}

//展示添加文章页面
func (this *ArticleController)ShowAddArticle()  {
	//在文章列表中展示session
	userName:=this.GetSession("userName")
	if userName==nil{
		this.Redirect("/login",302)
		return
	}
	this.Data["userName"]=userName.(string)
	o:=orm.NewOrm()
	var articleTypes []models.ArticleType
	qs:=o.QueryTable("ArticleType")
	qs.All(&articleTypes)
	this.Data["articleTypes"]=articleTypes
	this.Layout="layout.html"
	this.TplName="add.html"
}
//处理添加文章业务
func (this *ArticleController)HandleAddArticle()  {
	//接受数据
	articleName:=this.GetString("articleName")
	content:=this.GetString("content")
	//校验数据
	if articleName==""||content==""{
		this.Data["errmsg"]="文章标题内容不能为空"
		this.TplName="add.html"
		return
	}
	//接受图片
	file,head,err:=this.GetFile("uploadname")
	if err!=nil{
		this.Data["errmsg"]="获取文件失败"
		this.TplName="add.html"
		return
	}
	defer file.Close()
	//判断大小
	if head.Size>500000{
		this.Data["errmsg"]="图片太大，上传失败"
		this.TplName="add.html"
		return
	}
	//判断格式
	fileExt:=path.Ext(head.Filename)
	if fileExt!=".jpg" && fileExt!=".png"{
		this.Data["errmsg"]="图片格式不正确"
		this.TplName="add.html"
		return
	}
	//防止文件重名
	fileName:=time.Now().Format("2006-01-02-15-04-05")+fileExt
	this.SaveToFile("uploadname","./static/image/"+fileName)

	//处理数据
	//数据库的插入操作
	//1.获取orm对象
	o:=orm.NewOrm()
	//2.获取插入对象
	var article models.Article
	//3.给插入对象赋值
	article.Title=articleName
	article.Concent=content
	article.Image="/static/image/"+fileName

	var articleType models.ArticleType
	typeName:=this.GetString("select")
	articleType.TypeName = typeName
	o.Read(&articleType,"TypeName")
	article.ArticleType = &articleType
	//插入
	_,err=o.Insert(&article)
	if err!=nil{
		beego.Error(err)
		this.Data["errmsg"]="添加文章失败，清重新添加"
		this.TplName="add.html"
		return
	}
	//添加成功返回列表页面
	this.Redirect("/article/articleList",302)
}
//实现ShowArticleDetial函数
func (this *ArticleController)ShowArticleDetial()  {
	//在文章列表中展示session
	userName:=this.GetSession("userName")
	if userName==nil{
		this.Redirect("/login",302)
		return
	}
	this.Data["userName"]=userName.(string)
	//获取传递过来的ID
	id,err:=this.GetInt("id")
	//数据校验
	if err!=nil{
		this.Data["errmsg"]="数据请求错误"
		this.Redirect("/article/ShowArticleDetial",302)
		return
	}
	//根据id查询文章信息
	o:=orm.NewOrm()
	//获取对象
	var article models.Article
	//article.Id=id
	//err=o.Read(&article)
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("Id",id).One(&article)
	if err!=nil{
		this.Data["errmsg"]="请求路径错误"
		return
	}
	//获取多对多对象
	m2m:=o.QueryM2M(&article,"Users")
	//获取要插入的数据
	var user models.User
	user.UserName=userName.(string)
	o.Read(&user,"UserName")
	//插入多对多关系
	m2m.Add(user)
	//多对多关系查询
	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id",id).Distinct().All(&users)
	this.Data["users"]=users
	//查询文章后阅读此书加1
	article.ReadCount+=1
	o.Update(&article)
	//获取数据之后，指定视图，并给视图传递数据
	this.Data["article"]=article
	//var articleType models.ArticleType
	//o.Read(&articleType,"TypeName")
	//beego.Info("articleType.TypeName:",articleType.TypeName)
	//this.Data["typeName"]=articleType.TypeName
	this.Layout="layout.html"
	this.TplName="content.html"

}
func (this *ArticleController)ShowUpdateArticle()  {
	//在文章列表中展示session
	userName:=this.GetSession("userName")
	if userName==nil{
		this.Redirect("/login",302)
		return
	}
	this.Data["userName"]=userName.(string)
	//获取文章id
	articleId ,err:=this.GetInt("id")
	errmsg:=this.GetString("errmsg")
	if errmsg!=""{
		this.Data["errmsg"]=errmsg
	}
	if err!=nil{
		this.Redirect("/article/articleList?errmsg",302)
		return
	}
	o:=orm.NewOrm()
	var article models.Article
	article.Id=articleId
	o.Read(&article)
	//传递给视图
	this.Data["article"]=article
	this.Layout="layout.html"
	this.TplName="update.html"

}
//文件上传函数
func UpLoadFile(filepath string, this *ArticleController)string {
	file,head,err:=this.GetFile(filepath)
	defer file.Close()
	if err!=nil{
		this.Data["errmsg"]="上传图片错误，清重新添加"
		this.TplName="index.html"
		return ""
	}
	fileExt:=path.Ext(head.Filename)
	if fileExt!=".jpg"&&fileExt!=".png"{
		this.Data["errmsg"]="上传图片格式不正确，清重新添加"
		this.TplName="update.html"
		return ""
	}
	fileName:=time.Now().Format("2006-01-02-15:04:05")+fileExt
	this.SaveToFile("uploadname","./static/image/"+fileName)
	return "/static/image/"+fileName
}
func (this *ArticleController)HandleUpdate()  {
	//获取数据
	articleId,err:=this.GetInt("id")
	articleName:=this.GetString("articleName")
	content:=this.GetString("content")
	img:=UpLoadFile("uploadname",this)
	if err!=nil || articleName==""||content==""||img==""{
		errmsg:="编辑数据不完整"
		this.Redirect("/UpdateArticle?id="+strconv.Itoa(articleId)+"&errmsg="+errmsg,302)
		return
	}
	//更新数据
	o:=orm.NewOrm()
	var article models.Article
	article.Id=articleId
	if err=o.Read(&article);err!=nil{
		errmsg:="传递文章id错误"
		this.Redirect("/UpdateArticle?id="+strconv.Itoa(articleId)+"&errmsg="+errmsg,302)
	}
	article.Title=articleName
	article.Concent=content
	article.Image=img
	o.Update(&article)
	//返回视图
	this.Redirect("/article/articleList",302)
}
//删除文章
func (this *ArticleController)DeleteArticle()  {
	id ,err:=this.GetInt("id")
	if err!=nil{
		errmsg:="删除错误"
		this.Redirect("/article/ShowArticleList?errmsg"+errmsg,302)
		return
	}
	//delete article
	o:=orm.NewOrm()
	var article models.Article
	article.Id=id
	o.Delete(&article)
	//getback
	this.Redirect("/article/articleList",302)
}
func (this *ArticleController)AddArticleType()  {
	//在文章列表中展示session
	userName:=this.GetSession("userName")
	if userName==nil{
		this.Redirect("/login",302)
		return
	}
	this.Data["userName"]=userName.(string)
	o:=orm.NewOrm()
	var articleTypes []models.ArticleType
	qs:=o.QueryTable("ArticleType")
	qs.All(&articleTypes)
	this.Data["articleTypes"]=articleTypes
	this.Layout="layout.html"
	this.TplName="addType.html"
}
func (this *ArticleController)HandleUpdateArticle()  {
	typeName:=this.GetString("typeName")
	if typeName==""{
		errmsg:="增加类型不允许为空"
		this.Redirect("/article/addArticleType?errmsg"+errmsg,302)
	}
	o:=orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName=typeName
	if _,err:=o.Insert(&articleType);err!=nil{
		errmsg:="添加数据失败"
		this.Redirect("/article/addArticleType?errmsg"+errmsg,302)
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
	defer  conn.Close()
	if err!=nil{
		beego.Info(err)
		this.TplName="addType.html"
	}
	res,_:=conn.Query("select type_name from article_type")
	var typename string
	for res.Next(){
		res.Scan(&typename)
		if typeName==typename{
			beego.Info(err)
			this.TplName="addType.html"
		}
	}
	this.Redirect("/article/addArticleType",302)
}
func (this * ArticleController)DeleteType()  {
	typeId,err:=this.GetInt("id")
	if err!=nil{
		errmsg:="删除错误"
		this.Redirect("/article/DeleteType?errmsg"+errmsg,302)
		return
	}
	o:=orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id=typeId
	_,err = o.Delete(&articleType)
	if err != nil{
		errmsg:="删除错误"
		this.Redirect("/article/DeleteType?errmsg"+errmsg,302)
		return
	}
	this.Redirect("/article/addArticleType",302)
}
