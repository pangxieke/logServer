package storage

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/config"
	"strings"
)

var (
	producer sarama.SyncProducer
)

func InitKafka(cfg config.Configer) (err error) {
	addr := cfg.String("kafka::addrs")
	if addr == "" {
		return errors.New("kafka::addrs is empty")
	}

	addrs := strings.Split(addr, ",")
	if len(addrs) == 0 {

	}

	producer, err = sarama.NewSyncProducer(addrs, nil)
	if err != nil {
		return err
	}

	return nil
}
