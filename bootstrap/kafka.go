package bootstrap

import (
    // "context"
    "fmt"
    "github.com/IBM/sarama"
	"testproject/app/services"
    // "log"
    // "math"
    // "os"
    // "os/signal"
    // "sync"
    // "syscall"
)

var addrs = []string{"192.168.31.232:9092"}
var topic = "sun";
func InitializeConsumer() {
    consumer, err := sarama.NewConsumer(addrs, nil)
    if err != nil {
        fmt.Printf("fail to start consumer, err:%v \n", err)
        return
    }
    partitionList, err := consumer.Partitions(topic) // 通过topic获取所有分区
    if err != nil {
        fmt.Printf("fail to get partition list, err:%v\n", err)
        return
    }

    fmt.Println(partitionList)
    for partition := range partitionList { // 遍历所有分区
        pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
        if err != nil {
            fmt.Printf("failed to start consumer for partition %d, err:%v\n", partition, err)
            return
        }
        defer pc.AsyncClose()
        go func(sarama.PartitionConsumer) {
            for msg := range pc.Messages() {
                // 当设置了key的时候,不为空
                fmt.Println(msg.Topic);
                //执行消费代码
                services.ElasticSearchHander.Sync(string(msg.Value));
                fmt.Printf("Partition:%d Offset:%d Key:%s Value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
            }
        }(pc)
    }
    //time.Sleep(5 * time.Second)
    select {}
}

