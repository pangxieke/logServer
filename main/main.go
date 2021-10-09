package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/pangxieke/logServer/routers"
	"github.com/pangxieke/logServer/server"
	"github.com/pangxieke/logServer/storage"
	"github.com/pangxieke/logServer/util"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	util.InitLogger(beego.AppConfig)
	storage.Init(beego.AppConfig)

	//beego.BeeLogger.DelLogger("console")

	logs.Info("logproxy server is running")

	beego.LoadAppConfig("ini", "../conf/app.conf")

	// 开启2个服务
	s, _ := beego.AppConfig.Bool("time::sync")

	var wg sync.WaitGroup
	var sig chan bool
	wg.Add(1)

	if s {
		// 定时器
		server.Init(beego.AppConfig)
		go server.SyncTimer(sig, &wg)
	}

	beego.Run()

	close(sig)
	wg.Wait()

}
