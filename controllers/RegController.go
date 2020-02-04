package controllers

import (
	"git.ntmc.tech/root/MinoIC-PE/models/MinoConfigure"
	"git.ntmc.tech/root/MinoIC-PE/models/MinoDatabase"
	"git.ntmc.tech/root/MinoIC-PE/models/MinoEmail"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/utils/captcha"
	uuid "github.com/satori/go.uuid"
)

var cpt *captcha.Captcha

func init() {
	// use beego cache system store the captcha data
	store := cache.NewMemoryCache()
	cpt = captcha.NewWithFilter("/captcha/", store)
	cpt.StdHeight = 40
	cpt.StdWidth = 100
	cpt.ChallengeNums = 4
}

type RegController struct {
	beego.Controller
}

func (this *RegController) Get() {
	this.TplName = "Register.html"
	this.Data["webHostName"] = MinoConfigure.ConfGetHostName()
	this.Data["webApplicationName"] = MinoConfigure.ConfGetWebName()
}

func (this *RegController) Post() {
	this.TplName = "Register.html"
	//beego.Info("user posted!")
	registerEmail := this.GetString("registerEmail")
	registerPassword := this.GetString("registerPassword")
	registerPasswordConfirm := this.GetString("registerPasswordConfirm")
	registerName := this.GetString("registerName")
	cptSuccess := cpt.VerifyReq(this.Ctx.Request)
	agreement, err := this.GetBool("agreement", false)
	if err != nil {
		beego.Error(err)
	}
	DB := MinoDatabase.GetDatabase()
	defer DB.Close()
	if !cptSuccess {
		this.Data["hasError"] = true
		this.Data["hasErrorText"] = "验证码输入错误，请重试！"
	} else if !agreement {
		this.Data["hasError"] = true
		this.Data["hasErrorText"] = "您必须同意我们的用户协议才能注册！"
	} else if registerPassword != registerPasswordConfirm {
		this.Data["hasError"] = true
		this.Data["hasErrorText"] = "两次密码输入不一致，请检查！"
	} else if !DB.Where("Name = ?", registerName).RecordNotFound() {
		this.Data["hasError"] = true
		this.Data["hasErrorText"] = "您输入的用户名已被占用！"
	} else if !DB.Where("Email = ?", registerEmail).RecordNotFound() {
		this.Data["hasError"] = true
		this.Data["hasErrorText"] = "您输入的邮箱已被占用！"
	} else {
		newUuid := uuid.NewV4()
		newUser := MinoDatabase.User{
			Name:           registerName,
			Email:          registerEmail,
			Password:       registerPassword,
			UUID:           newUuid,
			IsAdmin:        false,
			EmailConfirmed: false,
		}
		DB.Create(&newUser)
		if MinoConfigure.ConfGetSMTPEnabled() {
			if err := MinoEmail.ConfirmRegister(newUser); err != nil {
				beego.Error(err)
				DelayRedirect(DelayInfo{
					URL:    MinoConfigure.ConfGetHostName() + "/reg",
					Detail: "即将跳转到注册页面",
					Title:  "邮件发送失败，请联系网站管理员！",
				}, &this.Controller)
			} else {
				DelayRedirect(DelayInfo{
					URL:    MinoConfigure.ConfGetHostName() + "/login",
					Detail: "即将跳转到登陆页面",
					Title:  "请前往您的邮箱进行验证！",
				}, &this.Controller)
			}
		} else {
			DelayRedirect(DelayInfo{
				URL:    MinoConfigure.ConfGetHostName() + "/login",
				Detail: "即将跳转到登陆页面",
				Title:  "注册成功！",
			}, &this.Controller)
		}
	}
	this.Data["webHostName"] = MinoConfigure.ConfGetHostName()
	this.Data["webApplicationName"] = MinoConfigure.ConfGetWebName()
	//todo: create Pterodactyl user at the same time
}
