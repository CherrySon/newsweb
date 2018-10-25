package main

import (
	_ "newsWeb/routers"
	"github.com/astaxie/beego"
	_"newsWeb/models"
)

func main() {
	//在beego.run之前把两个函数对应起来
	beego.AddFuncMap("PrePage",PrePageIndex)
	beego.AddFuncMap("NextPage",NextPage)
	beego.AddFuncMap("Index",Index)
	beego.Run()

}
//第二步 在代码里定义一个函数
func PrePageIndex(pageIndex int)int  {
	prePage:=pageIndex-1
	if prePage<1{
		prePage=1
	}
	return prePage
}
//定一个函数
func NextPage(pageIndex int,pageCount float64)int  {
	if pageIndex+1>int(pageCount){
		return pageIndex
	}
	return pageIndex+1
}
func Index(Index int)int  {
	Index=Index+1
	return Index
}

