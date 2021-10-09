package util

import (
	beecontext "github.com/astaxie/beego/context"
)

func OnError(ctx *beecontext.Context, code int, msg string) {
	data := make(map[string]interface{})

	data["success"] = false
	data["code"] = code
	data["msg"] = msg
	data["data"] = make(map[string]interface{})

	ctx.Output.JSON(data, false, false)
}

func OnSuccess(ctx *beecontext.Context, response interface{}) {
	data := make(map[string]interface{})

	data["success"] = true
	data["code"] = 0
	data["msg"] = "success"
	data["data"] = response

	ctx.Output.JSON(data, false, false)
}

func OnReturnJson(ctx *beecontext.Context, data interface{}) {
	ctx.Output.JSON(data, false, false)
}
