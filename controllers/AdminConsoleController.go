package controllers

import (
	"git.ntmc.tech/root/MinoIC-PE/models/MinoConfigure"
	"git.ntmc.tech/root/MinoIC-PE/models/MinoDatabase"
	"git.ntmc.tech/root/MinoIC-PE/models/MinoKey"
	"git.ntmc.tech/root/MinoIC-PE/models/MinoSession"
	"git.ntmc.tech/root/MinoIC-PE/models/PterodactylAPI"
	"github.com/astaxie/beego"
	"html/template"
	"strconv"
	"sync"
)

type AdminConsoleController struct {
	beego.Controller
}

func (this *AdminConsoleController) Prepare() {
	this.TplName = "AdminConsole.html"
	this.Data["u"] = 4
	handleNavbar(&this.Controller)
	sess := this.StartSession()
	if !MinoSession.SessionIslogged(sess) {
		DelayRedirect(DelayInfo{
			URL:    MinoConfigure.WebHostName + "/login",
			Detail: "正在跳转到登录",
			Title:  "您还没有登录",
		}, &this.Controller)
	} else if !MinoSession.SessionIsAdmin(sess) {
		DelayRedirect(DelayInfo{
			URL:    MinoConfigure.WebHostName,
			Detail: "正在跳转到主页",
			Title:  "您不是管理员",
		}, &this.Controller)
	}
	DB := MinoDatabase.GetDatabase()
	/* delete confirm */
	var (
		dib           []MinoDatabase.DeleteConfirm
		deleteServers []struct {
			ServerName            string
			ServerConsoleHostName template.URL
			ServerIdentifier      string
			DeleteURL             template.URL
			ServerOwner           string
			ServerEXP             string
			ServerHostName        string
		}
	)
	DB.Find(&dib)
	for i, d := range dib {
		var entity MinoDatabase.WareEntity
		if DB.Where("id = ?", d.WareID).First(&entity).RecordNotFound() {
			DB.Delete(&d)
		} else {
			pteServer := PterodactylAPI.GetServer(PterodactylAPI.ConfGetParams(), entity.ServerExternalID)
			deleteServers = append(deleteServers, struct {
				ServerName            string
				ServerConsoleHostName template.URL
				ServerIdentifier      string
				DeleteURL             template.URL
				ServerOwner           string
				ServerEXP             string
				ServerHostName        string
			}{
				ServerName:            pteServer.Name,
				ServerConsoleHostName: template.URL(PterodactylAPI.PterodactylGethostname(PterodactylAPI.ConfGetParams()) + "/server/" + pteServer.Identifier),
				ServerIdentifier:      pteServer.Identifier,
				DeleteURL:             template.URL(MinoConfigure.WebHostName + "/admin-console/delete-confirm/" + strconv.Itoa(int(entity.ID))),
				ServerOwner:           entity.UserExternalID,
				ServerEXP:             entity.ValidDate.Format("2006-01-02"),
				ServerHostName:        entity.HostName,
			})
			if deleteServers[i].ServerName == "" {
				deleteServers[i].ServerName = "无法获取服务器名称"
			}
			if deleteServers[i].ServerIdentifier == "" {
				deleteServers[i].ServerIdentifier = "无法获取编号"
			}
		}
	}
	//beego.Debug(deleteServers )
	this.Data["deleteServers"] = deleteServers
	/* panel stats*/
	var (
		specs    []MinoDatabase.WareSpec
		entities []MinoDatabase.WareEntity
		users    []MinoDatabase.User
		packs    []MinoDatabase.Pack
		keys     []MinoDatabase.WareKey
		orders   []MinoDatabase.Order
		wg       sync.WaitGroup
	)
	wg.Add(6)
	go func() {
		DB.Find(&specs)
		wg.Done()
	}()
	go func() {
		DB.Find(&entities)
		wg.Done()
	}()
	go func() {
		DB.Find(&users)
		wg.Done()
	}()
	go func() {
		DB.Find(&packs)
		wg.Done()
	}()
	go func() {
		DB.Find(&keys)
		wg.Done()
	}()
	go func() {
		DB.Where("confirmed = ?", true).Find(&orders)
		wg.Done()
	}()
	wg.Wait()
	this.Data["specAmount"] = len(specs)
	this.Data["entityAmount"] = len(entities)
	this.Data["userAmount"] = len(users)
	this.Data["packAmount"] = len(packs)
	this.Data["keyAmount"] = len(keys)
	this.Data["orderAmount"] = len(orders)
	this.Data["specs"] = specs
}

func (this *AdminConsoleController) Get() {}

func (this *AdminConsoleController) DeleteConfirm() {
	entityID := this.Ctx.Input.Param(":entityID")
	entityIDint, err := strconv.Atoi(entityID)
	if err != nil {
		_, _ = this.Ctx.ResponseWriter.Write([]byte("FAILED"))
	}
	if err := PterodactylAPI.ConfirmDelete(uint(entityIDint)); err != nil {
		_, _ = this.Ctx.ResponseWriter.Write([]byte("FAILED"))
	}
	_, _ = this.Ctx.ResponseWriter.Write([]byte("SUCCESS"))
}

func (this *AdminConsoleController) NewKey() {
	if !this.CheckXSRFCookie() {
		_, _ = this.Ctx.ResponseWriter.Write([]byte("XSRF 验证失败"))
		return
	}
	keyAmount, err := this.GetInt("key_amount", 1)
	if err != nil || keyAmount <= 0 || keyAmount >= 100 {
		_, _ = this.Ctx.ResponseWriter.Write([]byte("输入不合理的 KEY 数量"))
		return
	}
	validDuration, err := this.GetInt("valid_duration", 60)
	if err != nil || validDuration <= 0 {
		_, _ = this.Ctx.ResponseWriter.Write([]byte("输入不合理的有效期"))
		return
	}
	DB := MinoDatabase.GetDatabase()
	specID, err := this.GetInt("spec_id")
	if err != nil || DB.Where("id = ?", specID).First(&MinoDatabase.WareSpec{}).RecordNotFound() {
		_, _ = this.Ctx.ResponseWriter.Write([]byte("选择了无效的商品"))
		return
	}
	err = MinoKey.GeneKeys(keyAmount, uint(specID), validDuration, 20)
	if err != nil {
		_, _ = this.Ctx.ResponseWriter.Write([]byte("在数据库中创建 Key 失败"))
		return
	}
	_, _ = this.Ctx.ResponseWriter.Write([]byte("SUCCESS"))
}

func (this *AdminConsoleController) CheckXSRFCookie() bool {
	if !this.EnableXSRF {
		return true
	}
	token := this.Ctx.Input.Query("_xsrf")
	if token == "" {
		token = this.Ctx.Request.Header.Get("X-Xsrftoken")
	}
	if token == "" {
		token = this.Ctx.Request.Header.Get("X-Csrftoken")
	}
	if token == "" {
		return false
	}
	if this.XSRFToken() != token {
		return false
	}
	return true
}
