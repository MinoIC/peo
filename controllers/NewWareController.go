package controllers

import (
	"github.com/astaxie/beego"
)

var WareInfo map[string]InfoDetail

type NewWareController struct {
	beego.Controller
}
type InfoDetail struct {
	FriendlyName string
	Description  string
	Type         string
}

func init() {
	WareInfo = map[string]InfoDetail{
		"cpu": {
			FriendlyName: "CPU 限制 (%)",
			Description:  "每100个CPU限制数值表示可以占用一个CPU线程(Thread)",
			Type:         "int",
		},
		"disk": {
			FriendlyName: "磁盘限制 (MB)",
			Description:  "服务器的磁盘限制",
			Type:         "int",
		},
		"memory": {
			FriendlyName: "内存限制 (MB)",
			Description:  "服务器的内存限制",
			Type:         "int",
		},
		"swap": {
			FriendlyName: "SWAP内存限制 (MB)",
			Description:  "SWAP内存，即虚拟内存，映射到磁盘中",
			Type:         "int",
		},
		"location_id": {
			FriendlyName: "地区ID",
			Description:  "翼龙面板中的地区(Location)的ID",
			Type:         "int",
		},
		"dedicated_ip": {
			FriendlyName: "专用IP",
			Description:  "为服务器设置专用IP (可选)",
			Type:         "bool",
		},
		"nest_id": {
			FriendlyName: "Nest ID",
			Description:  "服务器使用的Nest的ID",
			Type:         "int",
		},
		"io": {
			FriendlyName: "Block IO 大小",
			Description:  "Block IO 大小 (10-1000) (默认500)",
			Type:         "int",
		},
		"egg_id": {
			FriendlyName: "Egg ID",
			Description:  "服务器使用的Egg的ID",
			Type:         "int",
		},
		"pack_id": {
			FriendlyName: "Pack ID",
			Description:  "服务器使用的Pack的ID",
			Type:         "int",
		},
		"port_range": {
			FriendlyName: "配备给服务器的端口范围",
			Description:  "端口范围，以逗号分隔以分配给服务器（例如：25565-25570、25580-25590）（可选）",
			Type:         "text",
		},
		"startup": {
			FriendlyName: "启动命令",
			Description:  "定制启动命令以分配给创建的服务器（可选）",
			Type:         "text",
		},
		"image": {
			FriendlyName: "镜像",
			Description:  "自定义Docker映像以分配给创建的服务器（可选）",
			Type:         "text",
		},
		"database": {
			FriendlyName: "数据库数量",
			Description:  "客户端将能够为其服务器创建此数量的数据库（可选）",
			Type:         "int",
		},
		"server_name": {
			FriendlyName: "服务器名字",
			Description:  "面板上显示的服务器名称（可选）",
			Type:         "text",
		},
		"oom_disabled": {
			FriendlyName: "关闭 OOM Killer",
			Description:  "是否应禁用“内存不足杀手”（可选）",
			Type:         "bool",
		},
	}
}
func (this *NewWareController) Get() {
	beego.Info(WareInfo["server_name"])
	this.TplName = "NewWare.html"
}