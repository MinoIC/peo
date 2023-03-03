package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/jinzhu/gorm"
	"github.com/minoic/glgf"
	"github.com/minoic/peo/internal/configure"
	"github.com/minoic/peo/internal/database"
	"github.com/minoic/peo/internal/pterodactyl"
	"github.com/minoic/peo/internal/session"
	"time"
)

var wareInfo = []InputField{
	{
		Name:           "ware_name",
		FriendlyName:   "商品名称",
		Description:    "显示在商品的标题",
		Type:           "text",
		AdditionalTags: "required",
		Required:       true,
	},
	{
		Name:         "ware_description",
		FriendlyName: "商品介绍",
		Description:  "显示在商品的介绍",
		Type:         "text",
	},
	{
		Name:           "cpu",
		FriendlyName:   "CPU 限制 (%)",
		Description:    "每100个CPU限制数值表示可以占用一个CPU线程(Thread)",
		Type:           "number",
		AdditionalTags: "required",
		Required:       true,
	},
	{
		Name:           "disk",
		FriendlyName:   "磁盘限制 (MB)",
		Description:    "服务器的磁盘限制",
		Type:           "number",
		AdditionalTags: "required",
		Required:       true,
	},
	{
		Name:           "memory",
		FriendlyName:   "内存限制 (MB)",
		Description:    "服务器的内存限制",
		Type:           "number",
		AdditionalTags: "required",
		Required:       true,
	},
	{
		Name:           "swap",
		FriendlyName:   "SWAP内存限制 (MB)",
		Description:    "SWAP内存，即虚拟内存，映射到磁盘中",
		Type:           "number",
		AdditionalTags: "required",
		Required:       true,
	},
	{
		Name:           "io",
		FriendlyName:   "Block IO 大小",
		Description:    "Block IO 大小 (10-1000) (默认填500)",
		Type:           "number",
		AdditionalTags: "required",
		Required:       true,
	},
	{
		Name:           "backups",
		FriendlyName:   "备份数量",
		Description:    "允许的备份数量 (默认填0)",
		Type:           "number",
		AdditionalTags: "required",
		Required:       true,
	},
	{
		Name:           "node_id",
		FriendlyName:   "节点ID",
		Description:    "服务器将会在这个节点上创建",
		Type:           "number",
		AdditionalTags: "required",
		Required:       true,
	},
	{
		Name:           "nest_id",
		FriendlyName:   "Nest ID",
		Description:    "服务器使用的Nest的ID",
		Type:           "number",
		AdditionalTags: "required",
		Required:       true,
	},
	{
		Name:           "egg_id",
		FriendlyName:   "Egg ID",
		Description:    "服务器使用的默认Egg的ID",
		Type:           "number",
		AdditionalTags: "required",
		Required:       true,
	},
	/*{
		Name:         "dedicated_ip",
		FriendlyName: "专用IP",
		Description:  "为服务器设置专用IP (可选)",
		Type:         "checkbox",
	},*/
	{
		Name:         "startup",
		FriendlyName: "启动命令",
		Description:  "定制启动命令以分配给创建的服务器（可选）",
		Type:         "text",
	},
	{
		Name:         "image",
		FriendlyName: "镜像",
		Description:  "自定义Docker映像以分配给创建的服务器（可选）",
		Type:         "text",
	},
	/*		{
			Name:         "database",
			FriendlyName: "数据库数量",
			Description:  "客户端将能够为其服务器创建此数量的数据库（可选）",
			Type:         "int",
		},*/
	/*		{
				Name:         "start_on_completion",
				FriendlyName: "立即启动",
				Description:  "是否在安装完成后立即启动服务器",
				Type:         "checkbox",
			},
			{
				Name:         "oom_disabled",
				FriendlyName: "开启 OOM Killer",
				Description:  "是否应开启“内存不足杀手”（推荐关闭）",
				Type:         "checkbox",
			},*/
	{
		Name:           "delete_time",
		FriendlyName:   "删除延迟（天）",
		Description:    "商品从失效暂停到被删除的时间",
		Type:           "number",
		AdditionalTags: "required",
		Required:       true,
	},
	{
		Name:           "price",
		FriendlyName:   "价格（每三十天/人民币）",
		Description:    "打折前的原价（任意正整数）",
		Type:           "number",
		AdditionalTags: "required",
		Required:       true,
	},
	{
		Name:           "discount0",
		FriendlyName:   "月付折扣",
		Description:    "付款时减去的百分比(0-100的整数)",
		Type:           "number",
		AdditionalTags: "required",
		Required:       true,
	},
	{
		Name:           "discount1",
		FriendlyName:   "季付折扣",
		Description:    "付款时减去的百分比(0-100的整数)",
		Type:           "number",
		AdditionalTags: "required",
		Required:       true,
	},
	{
		Name:           "discount2",
		FriendlyName:   "年付折扣",
		Description:    "付款时减去的百分比(0-100的整数)",
		Type:           "number",
		AdditionalTags: "required",
		Required:       true,
	},
}

type NewWareController struct {
	web.Controller
}

type InputField struct {
	Name           string
	FriendlyName   string
	Description    string
	Type           string
	AdditionalTags string
	Required       bool
	Default        interface{}
}

func (this *NewWareController) Prepare() {
	this.TplName = "NewWare.html"
	sess := this.StartSession()
	if !session.Logged(sess) {
		this.Abort("401")
	} else if !session.IsAdmin(sess) {
		this.Abort("401")
	}
	handleNavbar(&this.Controller)
	// glgf.Info(wareInfo)
	this.Data["options"] = wareInfo
	this.Data["u"] = 0
}
func (this *NewWareController) Get() {}

// todo: add nest/egg select instead of input ID
func (this *NewWareController) Post() {
	if !this.CheckXSRFCookie() {
		this.Data["hasError"] = true
		this.Data["hasErrorText"] = "XSRF 验证失败！"
		return
	}
	cli := pterodactyl.ClientFromConf()
	// formText,_:=template.ParseFiles("tpls/forms/formgroup.html")
	ware := database.WareSpec{
		Model:           gorm.Model{},
		WareName:        this.GetString("ware_name"),
		WareDescription: this.GetString("ware_description"),
		DockerImage:     this.GetString("image"),
	}
	var (
		err          error
		hasError     bool
		hasErrorText string
	)
	ware.Cpu, err = this.GetInt("cpu")
	if err != nil {
		glgf.Error(err)
		hasError = true
		hasErrorText = "POST 表单获取错误 cpu " + err.Error()
	} else if ware.Cpu < 0 {
		hasError = true
		hasErrorText = "CPU 输入值不合法"
	}
	ware.Disk, err = this.GetInt("disk")
	if err != nil {
		glgf.Error(err)
		hasError = true
		hasErrorText = "POST 表单获取错误 disk " + err.Error()
	} else if ware.Disk <= 0 {
		hasError = true
		hasErrorText = "DISK 输入值不合法"
	}
	ware.Memory, err = this.GetInt("memory")
	if err != nil {
		glgf.Error(err)
		hasError = true
		hasErrorText = "POST 表单获取错误 memory " + err.Error()
	} else if ware.Memory <= 0 {
		hasError = true
		hasErrorText = "Memory 输入值不合法"
	}
	ware.Io, err = this.GetInt("io")
	if err != nil {
		glgf.Error(err)
		hasError = true
		hasErrorText = "POST 表单获取错误 io " + err.Error()
	} else if ware.Io < 100 || ware.Io > 1000 {
		hasError = true
		hasErrorText = "Block IO Weight 输入了不建议的值"
	}
	ware.Backups, err = this.GetInt("backups")
	if err != nil {
		glgf.Error(err)
		hasError = true
		hasErrorText = "POST 表单获取错误 backups " + err.Error()
	} else if ware.Backups < 0 || ware.Backups > 1000 {
		hasError = true
		hasErrorText = "备份数量输入了不建议的值"
	}
	ware.Swap, err = this.GetInt("swap")
	if err != nil {
		glgf.Error(err)
		hasError = true
		hasErrorText = "POST 表单获取错误 swap " + err.Error()
	} else if ware.Swap < (-1) {
		hasError = true
		hasErrorText = "Swap 输入值不合法"
	}
	var discount [3]int
	if discount[0], err = this.GetInt("discount0"); err != nil {
		glgf.Error(err)
		hasError = true
		hasErrorText = "POST 表单获取错误 discount0 " + err.Error()
	} else if discount[0] > 100 || discount[0] < 0 {
		hasError = true
		hasErrorText = "Discount0 输入值不合法"
	}
	if discount[1], err = this.GetInt("discount1"); err != nil {
		glgf.Error(err)
		hasError = true
		hasErrorText = "POST 表单获取错误 discount1 " + err.Error()
	} else if discount[1] > 100 || discount[1] < 0 {
		hasError = true
		hasErrorText = "Discount1 输入值不合法"
	}
	if discount[2], err = this.GetInt("discount2"); err != nil {
		glgf.Error(err)
		hasError = true
		hasErrorText = "POST 表单获取错误 discount2 " + err.Error()
	} else if discount[2] > 100 || discount[2] < 0 {
		hasError = true
		hasErrorText = "Discount2 输入值不合法"
	}
	ware.Node, err = this.GetInt("node_id")
	_, nerr := cli.GetNode(ware.Node)
	if err != nil {
		glgf.Error(err)
		hasError = true
		hasErrorText = "POST 表单获取错误 discount " + err.Error()
	} else if ware.Node < 0 || nerr != nil {
		hasError = true
		hasErrorText = "Node ID 输入值小于 0 或找不到该节点"
	}
	ware.Nest, err = this.GetInt("nest_id")
	if err != nil {
		glgf.Error(err)
		hasError = true
		hasErrorText = "POST 表单获取错误 nest_id " + err.Error()
	}
	ware.Egg, err = this.GetInt("egg_id")
	if err != nil {
		glgf.Error(err)
		hasError = true
		hasErrorText = "POST 表单获取错误 egg_id " + err.Error()
	} else if _, err := cli.GetEgg(ware.Nest, ware.Egg); err != nil {
		hasError = true
		hasErrorText = "在翼龙面板中找不到对应的 EGG"
	}
	price, err := this.GetUint32("price", 999)
	if err != nil {
		glgf.Error(err)
		hasError = true
		hasErrorText = "POST 表单获取错误 price " + err.Error()
	} else if price < 0 {
		hasError = true
		hasErrorText = "价格不能设置为负"
	}
	e, err := this.GetInt("delete_time")
	if err != nil {
		glgf.Error(err)
		hasError = true
		hasErrorText = "POST 表单获取错误 price " + err.Error()
	} else if e < 0 {
		hasError = true
		hasErrorText = "删除延迟不能小于 0"
	} else {
		ware.DeleteDuration = time.Duration(e*24) * time.Hour
	}
	if hasError {
		this.Data["hasError"] = true
		this.Data["hasErrorText"] = hasErrorText
	} else {
		ware.PricePerMonth = uint(price)
		ware.OomDisabled = true
		ware.StartOnCompletion = true
		// todo: handle database number
		DB := database.Mysql()
		for i, d := range []time.Duration{
			30 * 24 * time.Hour,
			90 * 24 * time.Hour,
			365 * 24 * time.Hour,
		} {
			// glgf.Debug(d)
			wareTemp := ware
			wareTemp.ValidDuration = d
			wareTemp.Discount = discount[i]
			DB.Create(&wareTemp)
		}
		DelayRedirect(DelayInfo{
			URL:    configure.Viper().GetString("WebHostName") + "/new-ware",
			Detail: "正在跳转回添加页面",
			Title:  "添加商品成功！",
		}, &this.Controller)
	}
}

func (this *NewWareController) CheckXSRFCookie() bool {
	if !this.EnableXSRF {
		return true
	}
	token := this.GetString("_xsrf")
	if token == "" {
		token = this.Ctx.Input.Query("_xsrf")
	}
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
