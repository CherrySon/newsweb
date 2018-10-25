package routers

import (
	"newsWeb/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/article/*",beego.BeforeExec,funcFilter)
    beego.Router("/", &controllers.MainController{})
    beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleRegister")
    beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
    beego.Router("/article/articleList",&controllers.ArticleController{},"get:ShowArticleList")
    beego.Router("/article/addArticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArticle")
	beego.Router("/article/ShowArticleDetial",&controllers.ArticleController{},"get:ShowArticleDetial")
	beego.Router("/article/updateArticle",&controllers.ArticleController{},"get:ShowUpdateArticle;post:HandleUpdate")
	beego.Router("/article/deleteArticle",&controllers.ArticleController{},"get:DeleteArticle")
	beego.Router("/article/addArticleType",&controllers.ArticleController{},"get:AddArticleType;post:HandleUpdateArticle")
	beego.Router("/article/deleteType",&controllers.ArticleController{},"get:DeleteType")
	beego.Router("/article/logout",&controllers.UserController{},"get:Logout")
}
var funcFilter = func(ctx*context.Context) {
	userName:=ctx.Input.Session("userName")
	if userName==nil{
		ctx.Redirect(302,"/login")
	}
}