package AutoManager

import (
	"git.ntmc.tech/root/MinoIC-PE/Controllers"
	"git.ntmc.tech/root/MinoIC-PE/MinoConfigure"
	"git.ntmc.tech/root/MinoIC-PE/MinoDatabase"
	"git.ntmc.tech/root/MinoIC-PE/MinoKey"
	"git.ntmc.tech/root/MinoIC-PE/PterodactylAPI"
	"github.com/astaxie/beego"
	"time"
)

func LoopTasksManager() {
	go func() {
		interval, err := MinoConfigure.GetConf().Int("AutoTaskInterval")
		if err != nil {
			interval = 10
			beego.Error("cant get AutoTaskInterval ,set it to 10sec as default")
		}
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		for {
			select {
			case <-ticker.C:
				go PterodactylAPI.CheckServers()
			case <-ticker.C:
				go PterodactylAPI.CacheNeededEggs()
			case <-ticker.C:
				go PterodactylAPI.CacheNeededServers()
			case <-ticker.C:
				go Controllers.RefreshWareInfo()
			case <-ticker.C:
				go MinoKey.DeleteOutdatedKeys()
			case <-ticker.C:
				go func() {
					DB := MinoDatabase.GetDatabase()
					beego.Info("DB_OpenConnections: ", DB.DB().Stats().OpenConnections, " - ",
						DB.DB().Stats().WaitCount)
				}()
			}
		}
	}()
}