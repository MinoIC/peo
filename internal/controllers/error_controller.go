package controllers

import (
	"github.com/astaxie/beego"
	"github.com/minoic/peo/internal/configure"
)

type ErrorController struct {
	beego.Controller
}

func (this *ErrorController) Error400() {
	DelayRedirect(DelayInfo{
		URL:    configure.WebHostName,
		Detail: "请求参数有误",
		Title:  "400 Bad Request",
	}, &this.Controller, 400)
	this.TplName = "Delay.html"
}

func (this *ErrorController) Error401() {
	DelayRedirect(DelayInfo{
		URL:    configure.WebHostName,
		Detail: "未经授权，请求要求验证身份",
		Title:  "401 Unauthorized",
	}, &this.Controller, 401)
	this.TplName = "Delay.html"
}

func (this *ErrorController) Error403() {
	DelayRedirect(DelayInfo{
		URL:    configure.WebHostName,
		Detail: "服务器拒绝请求",
		Title:  "403 Forbidden",
	}, &this.Controller, 403)
	this.TplName = "Delay.html"
}

func (this *ErrorController) Error404() {
	DelayRedirect(DelayInfo{
		URL:    configure.WebHostName,
		Detail: "找不到指定页面: " + configure.WebHostName + this.Ctx.Request.URL.String(),
		Title:  "404 Not Found 😭",
	}, &this.Controller, 404)
	this.TplName = "Delay.html"
}

func (this *ErrorController) Error405() {
	DelayRedirect(DelayInfo{
		URL:    configure.WebHostName,
		Detail: "不被允许的方法: " + this.Ctx.Request.Method,
		Title:  "405 Method not Allowed",
	}, &this.Controller, 405)
	this.TplName = "Delay.html"
}

func (this *ErrorController) Error500() {
	DelayRedirect(DelayInfo{
		URL:    this.Ctx.Request.Referer(),
		Detail: "服务器遇到了一个未曾预料的状况",
		Title:  "500 Internal Server Error",
	}, &this.Controller, 500)
	this.TplName = "Delay.html"
}

func (this *ErrorController) Error502() {
	DelayRedirect(DelayInfo{
		URL:    this.Ctx.Request.Referer(),
		Detail: "从上游服务器接收到无效的响应",
		Title:  "502 Bad Gateway",
	}, &this.Controller, 502)
	this.TplName = "Delay.html"
}

func (this *ErrorController) Error503() {
	DelayRedirect(DelayInfo{
		URL:    this.Ctx.Request.Referer(),
		Detail: "临时的服务器维护或者过载",
		Title:  "503 Service Unavailable",
	}, &this.Controller, 503)
	this.TplName = "Delay.html"
}

func (this *ErrorController) Error504() {
	DelayRedirect(DelayInfo{
		URL:    this.Ctx.Request.Referer(),
		Detail: "未能及时从上游服务器或者辅助服务器收到响应",
		Title:  "504 Gateway Timeout",
	}, &this.Controller, 504)
	this.TplName = "Delay.html"
}
