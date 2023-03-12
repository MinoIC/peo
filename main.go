package main

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/minoic/peo/api"
	"github.com/minoic/peo/internal/controllers"
	"github.com/minoic/peo/internal/cron"
)

const Version = "v0.2.0"

var (
	BUILD_TIME string
	GO_VERSION string
)

func main() {
	controllers.BuildTime = BUILD_TIME
	controllers.GoVersion = GO_VERSION
	api.InitRouter()
	go cron.LoopTasksManager()
	web.Run()
}

// todo: add code comments
