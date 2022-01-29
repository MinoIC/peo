package main

import (
	"github.com/astaxie/beego"
	"github.com/minoic/peo/api"
	"github.com/minoic/peo/internal/cron"
)

const Version = "v0.1.0"

func main() {
	beego.LoadAppConfig("ini", "app.conf")
	beego.BConfig.WebConfig.Session.SessionDisableHTTPOnly = true
	api.InitRouter()
	go cron.LoopTasksManager()
	beego.Run()
}

// todo: add code comments