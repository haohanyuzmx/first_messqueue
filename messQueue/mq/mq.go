package mq

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"messQueue/model"
	"sync"
)


func Joinmg(order model.Order)  {
	cf:=sarama.NewConfig()
	cf.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	cf.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	cf.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, cf)
	if err!=nil {
		log.Println("producer closed, err:", err)
		return
	}
	defer client.Close()
	msg:=&sarama.ProducerMessage{}
	msg.Topic="order"
	byo,err:=json.Marshal(order)
	if err!=nil {
		log.Println(err)
		return
	}
	msg.Value=sarama.ByteEncoder(byo)
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		log.Println("send msg failed, err:", err)
		return
	}
	log.Printf("pid:%v offset:%v\n", pid, offset)
}


func Dealmq()  {
	var sy sync.WaitGroup
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		log.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions("order") // 根据topic取到所有的分区
	if err != nil {
		log.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println(partitionList)
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition("order", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		// 异步从每个分区消费信息
		sy.Add(1)
		go func(sarama.PartitionConsumer) {
			defer sy.Done()
			var o model.Order
			for msg := range pc.Messages() {
				errs:=json.Unmarshal(msg.Value,&o)
				fmt.Println(o)
				if errs != nil {
					log.Println(err)
					return
				}
				model.DealOrder(o)
			}
		}(pc)
	}
	sy.Wait()
	log.Println("结束处理")
}