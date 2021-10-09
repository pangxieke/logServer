package storage

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
)

type storage interface {
	SendMsg(string)
}

var (
	client storage
	//同步标准
	sync bool
)

func Init(cfg config.Configer) {
	// 存储
	initStorage(cfg)

	sync, _ = cfg.Bool("sync")

}

func initStorage(cfg config.Configer) {
	s := cfg.String("storage")
	topic := cfg.String("topic")
	if s == "" {
		panic("storage is empty")
	}
	if topic == "" {
		panic("topic is empty")
	}

	switch s {
	case "redis":
		if err := InitRedis(cfg); err != nil {
			panic(err)
		}
		client = &Redis{key: topic}
	case "kafka":
		if err := InitKafka(cfg); err != nil {
			panic(err)
		}
		client = &Kafka{topic: topic}
	//case "elastic":
	//	if err := InitElastic(cfg); err != nil {
	//		panic(err)
	//	}
	//	client = newESClient(topic)
	default:
		panic("storage type error")
	}

	if err := InitElastic(cfg); err != nil {
		panic(err)
	}
	ESLogHandler = newESClient(topic)
}

func SendMsg(str string) {
	if !sync {
		go client.SendMsg(str)
	} else {
		client.SendMsg(str)
	}
	fmt.Println("storage success")
}

type Redis struct {
	key string
}

func (r *Redis) SendMsg(s string) {
	key := r.key
	CacheClient.Push(key, s)
}

type Kafka struct {
	topic string
}

func (k *Kafka) SendMsg(s string) {

	hashkey := fmt.Sprintf("%s", s)
	var key = sarama.StringEncoder(hashkey)
	var value sarama.ByteEncoder = []byte(s)

	msg := sarama.ProducerMessage{
		Topic: k.topic,
		Key:   key,
		Value: value,
	}
	p, offset, err := producer.SendMessage(&msg)
	logs.Info("kafka SendMessage topic=%s, %d, %d, err=%s", k.topic, p, offset, err)
}
