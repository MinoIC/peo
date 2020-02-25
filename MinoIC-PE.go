package main

import (
	"git.ntmc.tech/root/MinoIC-PE/controllers"
	"git.ntmc.tech/root/MinoIC-PE/models/AutoManager"
	"github.com/astaxie/beego"
)

func main() {
	beego.BConfig.WebConfig.Session.SessionDisableHTTPOnly = true
	beego.Router("/", &controllers.WareSellerController{})
	beego.Router("/index", &controllers.IndexController{})
	beego.Router("/reg", &controllers.RegController{})
	beego.Router("/reg/confirm/:key", &controllers.RegController{}, "get:MailConfirm")
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/admin-console", &controllers.AdminConsoleController{})
	beego.Router("/admin-console/delete-confirm/:entityID", &controllers.AdminConsoleController{}, "get:DeleteConfirm")
	beego.Router("/admin-console/new-key", &controllers.AdminConsoleController{}, "get:NewKey")
	beego.Router("/new-ware", &controllers.NewWareController{})
	beego.Router("/new-order", &controllers.OrderCreateController{})
	beego.Router("/new-pack", &controllers.NewPackController{})
	beego.Router("/order/:orderID", &controllers.OrderInfoController{})
	beego.Router("/delay", &controllers.DelayController{})
	beego.Router("/user-settings", &controllers.UserSettingsController{})
	beego.Router("/user-settings/change-password", &controllers.UserSettingsController{}, "post:UpdateUserPassword")
	beego.Router("/user-settings/change-email", &controllers.UserSettingsController{}, "post:UpdateUserEmail")
	beego.Router("/user-settings/change-email/:email", &controllers.UserSettingsController{}, "get:SendCaptcha")
	beego.Router("/user-settings/create-pterodactyl-user", &controllers.UserSettingsController{}, "get:CreatePterodactylUser")
	beego.Router("/user-messages", &controllers.UserMessagesController{})
	beego.Router("/user-console", &controllers.UserConsoleController{})
	beego.Router("/user-console/renew/:entityID/:key", &controllers.UserConsoleController{}, "get:Renew")
	beego.Router("/user-console/reinstall/:entityID/:packID", &controllers.UserConsoleController{}, "get:Reinstall")
	beego.Router("/user-terms", &controllers.UserTermsController{})
	beego.Router("/forget-password", &controllers.ForgetPasswordController{})
	beego.Router("/forget-password-mail/:email", &controllers.ForgetPasswordController{}, "get:SendMail")
	beego.ErrorController(&controllers.ErrorController{})
	AutoManager.LoopTasksManager()
	beego.Run()
}

//todo: add alipay/wxpay or more payment method
//todo: add code comments
