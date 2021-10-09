package server

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/pangxieke/logServer/storage"
)

type redis struct {
}

func NewRedis() *redis {
	return new(redis)
}

func (r *redis) Read() {

	lLen := storage.CacheClient.LLen(topic)
	if lLen == 0 {
		return
	}
	val := storage.CacheClient.Pop(topic)

	fmt.Println("redis read:", val)

	var data interface{}
	err := json.Unmarshal([]byte(val), &data)
	if err != nil {
		logs.Info("json.Unmarshal err=%s", err)
	}

	fmt.Printf("json.Unmarshal %+v", data)
	//
	//res, _ := json.Marshal(data)
	//logs.Info("json.Marshal: %s", res)

	r.SendMsg(data)
}

func (r *redis) SendMsg(s interface{}) {
	logs.Info("SendMsg es, read from redis")
	storage.ESLogHandler.SendMsg(s)
}
