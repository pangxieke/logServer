package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"time"
)

type BaseController struct {
	beego.Controller
	XRealIp string
	begin   int64
}

func (this *BaseController) Prepare() {
	this.begin = time.Now().UnixNano()

	logs.Debug("url = %s, body = %s",
		this.Ctx.Input.URL(), string(this.Ctx.Input.RequestBody))

	this.XRealIp = this.Ctx.Input.Header("X-Real-IP")
	if this.XRealIp == "" {
		this.XRealIp = this.Ctx.Input.Header("X-Forwarded-For")
	}

	logs.Info("X-Client-IP = %s", this.XRealIp)
}

func (this *BaseController) Finish() {
	elapse := (time.Now().UnixNano() - this.begin) / 1e6
	logs.Info("url = %s, use %d msec", this.Ctx.Input.URL(), elapse)
}
