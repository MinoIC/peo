package controllers

import (
	"NTOJ/models"
	"github.com/astaxie/beego"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Get() {
	this.TplName = "login.html"
}

func (this *LoginController) Post() {
	this.TplName = "login.html"
	DB := models.GetDatabase()
	defer DB.Close()
	loginEOU := this.GetString("loginEOU")
	loginPass := this.GetString("loginPass")
	var user models.User
	if !DB.Where("Email = ?", loginEOU).Or("Name = ?", loginEOU).First(&user).RecordNotFound() {
		if loginPass == user.Password {
			this.Data["loginReturnData"] = "logged in!"
			this.SetSession("LST", models.GeneToken(user.Name))
			this.SetSession("UN", user.Name)
		} else {
			this.Data["loginReturnData"] = "login failed!!"
		}
	}
	return
}
