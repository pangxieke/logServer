package server

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"sync"
	"time"
)

var topic string

var t time.Duration = time.Second * 10

var storageClient read

type read interface {
	Read()
	SendMsg(s interface{})
}

func Init(cfg config.Configer) {
	topic = cfg.String("topic")
	syncTime := cfg.DefaultInt("time::sync_time", 10)
	t = time.Duration(syncTime) * time.Second
	fmt.Println("time:", t)

	// 存储类型
	s := cfg.String("storage")
	if s == "kafka" {
		storageClient = NewKafka()
	} else {
		storageClient = NewRedis()
	}
}

// 定时任务
func SyncTimer(sig chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	timer := time.NewTicker(t)
	defer timer.Stop()

Loop:
	for {
		select {
		case <-timer.C:
			logs.Info("it's time for read")
			storageClient.Read()

		case <-sig:
			break Loop
		}
	}

	logs.Info("exit timer for SyncTimer")
}
