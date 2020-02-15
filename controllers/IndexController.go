package controllers

import (
	"github.com/astaxie/beego"
)

type IndexController struct {
	beego.Controller
}

func (this *IndexController) Get() {
	this.TplName = "Loading.html"
	handleNavbar(&this.Controller)
	this.Data["u"] = 0
}
