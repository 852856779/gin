package services

import (
    "context"
    "fmt"
    "github.com/IBM/sarama"
    "log"
    "math"
    "os"
    // "os/signal"
    "sync"
    // "syscall"
)

type synckakfaService struct{}


var SyncKakfaService = new(synckakfaService)
// 更多参考:https://kpretty.tech/archives/gokafkaclient

var addrs = []string{"192.168.31.110:9092"}
// var Topic = []string{"sun", "topic2", "topic3"} 
var Topic = "sun";

func (synckakfaService *synckakfaService) Producer(jsonData string ,topic string) {
    //    生产者配置
    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForAll          // ACK
    config.Producer.Partitioner = sarama.NewRandomPartitioner // 分区
    // 异步回调(两个channel, 分别是成功和错误)
    config.Producer.Return.Successes = true // 确认
    config.Producer.Return.Errors = true

    sarama.Logger = log.New(os.Stdout, "[Sarama]", log.LstdFlags)

    // 连接kafka
    // 同步
    //client, err := sarama.NewSyncProducer(addrs, config)
    // 异步
    client, err := sarama.NewAsyncProducer(addrs, config)
    if err != nil {
        fmt.Println("producer error", err)
        return
    }

    defer func() {
        _ = client.Close()
    }()

    // 封装消息
    msg := &sarama.ProducerMessage{
        Topic: topic,
        Value: sarama.StringEncoder(jsonData),
    }
	client.Input() <- msg   //异步发送消息
    // pid, offset, err := client.SendMessage(msg)
    // if err != nil {
    //     fmt.Println("send failed", err)
    //     return
    // }

    //fmt.Printf("pid:%v offset:%v \n", pid, offset)
}


func (synckakfaService *synckakfaService) ProducerSync(msg sarama.ProducerMessage) error {
    //    生产者配置
    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForAll          // ACK
    config.Producer.Partitioner = sarama.NewRandomPartitioner // 分区
    // 异步回调(两个channel, 分别是成功和错误)
    config.Producer.Return.Successes = true // 确认
    config.Producer.Return.Errors = true

    sarama.Logger = log.New(os.Stdout, "[Sarama]", log.LstdFlags)

    // 连接kafka
    // 同步
    client, err := sarama.NewSyncProducer(addrs, config)
    if err != nil {
        fmt.Println("producer error", err)
        return err
    }

    defer func() {
        _ = client.Close()
    }()

    pid, offset, err := client.SendMessage(&msg)
    if err != nil {
        fmt.Println("send failed", err)
        return err
    }

    fmt.Printf("pid:%v offset:%v \n", pid, offset)
    return err
}

func groupConsumer() {
    groupId := "sarama-consumer"
    config := sarama.NewConfig()
    // 关闭自动提交 和 初始化策略(oldest|newest)
    config.Consumer.Offsets.AutoCommit.Enable = false
    config.Consumer.Offsets.Initial = sarama.OffsetOldest
    sarama.NewConsumerGroup(addrs, groupId, config)
}

// func main() {
//     // 生产者
//     //producer()
//     // 消费者(只能连续读取,中断期间会丢失数据)
//     //consumer()

//     // 消费者
//     //groupConsumer()
//     SimpleConsumer()
// }

var groupID = "sarama-consumer"
var asyncOffset chan struct{}
var wg sync.WaitGroup

const defaultOffsetChannelSize = math.MaxInt

func (synckakfaService *synckakfaService) SimpleConsumer() {
    fmt.Printf("初始化111");
    // MongoDBService.Connect();
    brokers := addrs
    // 消费者配置
    config := sarama.NewConfig()
    // 关闭自动提交
    config.Consumer.Offsets.AutoCommit.Enable = false
    config.Consumer.Offsets.Initial = sarama.OffsetOldest
    // 开启日志
    logger := log.New(os.Stdout, "[Sarama] ", log.LstdFlags)
    sarama.Logger = logger
    consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
    if err != nil {
        panic(err)
    }
    defer func() { _ = consumer.Close() }()
    // 搞一个上下文用于终止消费者
    // ctx, cancelFunc := context.WithCancel(context.Background())
    // 监听终止信号
    //wg.Add(1)
    go func() {
        logger.Println("monitor signal")
        // quit := make(chan os.Signal, 1)
        // signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
        // <-quit
        logger.Println("stop consumer")
        //cancelFunc()
        //wg.Done()
    }()
   
    // 消费数据
    err = consumer.Consume(context.Background(), []string{Topic}, &Consumer{})
    if err != nil {
        panic(err)
    }
    // 等待所有偏移量都提交完毕再退出
    logger.Println("当前存在未提交的偏移量")
    // wg.Wait()
    logger.Println("结束")
    // os.Exit(0)
    // logger.Println("结束")
}

type Consumer struct{}

func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error {
    // 初始化异步提交的channel
    fmt.Printf("初始化");
    asyncOffset = make(chan struct{}, defaultOffsetChannelSize)
    wg.Add(1)
    // 异步提交偏移量
    go func() {
        for range asyncOffset {
            session.Commit()
        }
        wg.Done()
    }()
    return nil
}

func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error {
    // 关闭通道
    fmt.Printf("关闭通道");
    close(asyncOffset)
    return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    fmt.Printf("标记");
EXIT:
    for {
        select {
        case message := <-claim.Messages():
            fmt.Printf(string(message.Value));
            //log.Printf("Message claimed: key= %s, value = %s, timestamp = %v, Topic = %s", string(message.Key), string(message.Value), message.Timestamp, message.Topic)
            // 标记消息，并不是提交偏移量
            session.MarkMessage(message, "")
            // 异步提交
            asyncOffset <- struct{}{}
        case <-session.Context().Done():
            log.Println("cancel consumer")
            break EXIT
        }
    }
    return nil
}


