package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"github.com/pangxieke/logServer/storage"
	"github.com/pangxieke/logServer/util"
	"time"
)

type LogCtl struct {
	BaseController

	producer sarama.SyncProducer
}

func (this *LogCtl) Prepare() {
	this.BaseController.Prepare()

}

type LogRequest struct {
	MacAddr         string `json:"mac_addr"`
	ClientIp        string `json:"client_ip"`
	Version         string `json:"version"`
	EventId         string `json:"event_id"`
	EventTime       int64  `json:"event_time"`
	EventTimeFormat time.Time
	RequestId       string      `json:"request_id"`
	Payload         interface{} `json:"payload"`
}

func (this *LogCtl) Upload() {
	var req LogRequest
	var err error

	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	if err != nil {
		logs.Error("upload, body = %s, err = %v", string(this.Ctx.Input.RequestBody), err)

		util.OnError(this.Ctx, 0, "error")
		return
	}

	req.ClientIp = this.XRealIp

	req.EventTimeFormat = time.Unix(0, req.EventTime*int64(time.Millisecond))

	body, _ := json.Marshal(req)
	fmt.Println(string(body))

	// 日志写入storage
	storage.SendMsg(string(body))

	util.OnSuccess(this.Ctx, nil)

}
