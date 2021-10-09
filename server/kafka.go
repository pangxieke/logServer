package server

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"github.com/pangxieke/logServer/storage"
	"sync"
)

var (
	wg sync.WaitGroup
)

type kafka struct {
}

func NewKafka() *kafka {
	return new(kafka)
}

func (k *kafka) Read() {

	consumer, err := sarama.NewConsumer([]string{"112.74.187.73:9092"}, nil)
	if err != nil {
		logs.Info("fail to start consumer, err:%s", err)
		return
	}
	// 根据topic取到所有的分区
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		logs.Info("fail to get list of partition:%s", err)
		return
	}
	fmt.Println(partitionList)

	for partition := range partitionList {
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			logs.Info("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}

		defer pc.AsyncClose()
		// 异步从每个分区消费信息
		wg.Add(1)
		go func(sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%s Value:%s", msg.Partition, msg.Offset, msg.Key, msg.Value)

				k.SendMsg(string(msg.Value))
			}
			fmt.Println("kafka3")
		}(pc)
		wg.Wait()
	}
}

func (k *kafka) SendMsg(s interface{}) {
	logs.Info("SendMsg es, read from kafka")
	storage.ESLogHandler.SendMsg(s)
}
