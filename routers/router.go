package routers

import (
	"github.com/astaxie/beego"
	"github.com/pangxieke/logServer/controllers"
)

func init() {
	beego.Router("/", &controllers.LogCtl{}, "post:Upload")
}
