package bootstrap

import (
    "context"
    "fmt"
    "github.com/IBM/sarama"
	"testproject/app/services"
    // "log"
    "go.uber.org/zap"
    "time"
    // "os/signal"
    // "log"
    // "math"
    // "os"
    // "os/signal"
    "sync"
    // "syscall"
    log "github.com/sirupsen/logrus"
)

var addrs = []string{"192.168.31.110:9092"}
// var topic = "sun";
// var Topic = []string{"sun", "redisKeyDel"} 
var topics = []string{"sun", "redisKeyDel"} 

func InitializeConsumer() {
    consumer, err := sarama.NewConsumer(addrs, nil)
    if err != nil {
        fmt.Printf("fail to start consumer, err:%v \n", err)
        return
    }
    for _, topic := range topics { 
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
					fmt.Println("普通轮询消费");
                    fmt.Println(msg.Topic);
                    //执行消费代码
                    setConsumerFunction(msg)
                    fmt.Printf("Partition:%d Offset:%d Key:%s Value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
                }
            }(pc)
        }
    }

    //time.Sleep(5 * time.Second)
    select {}
}

//根据topic而选择不同的消费
func setConsumerFunction(message *sarama.ConsumerMessage){
	topic := message.Topic
    switch topic {
    case "sun":
        services.ElasticSearchHander.Sync(string(message.Value));
    case "redisKeyDel"://用于删除缓存的真实队列
        // services.ElasticSearchHander.Sync(message);
    case "realTopic": //延迟队列redis
        services.RedisHander.DelCache(message);
    }
}


// type KafkaConfig struct {
// 	BrokerList []string
// 	Topic      []string
// 	GroupId    []string
// 	Cfg        *sarama.Config
// 	PemPath    string
// 	KeyPath    string
// 	CaPemPath  string
// }
//启动消费组
var ConsumerGroupReal  sarama.ConsumerGroup
var ConsumerGroupDelay sarama.ConsumerGroup
var Producer sarama.SyncProducer
var DelayTime  time.Duration = time.Minute * 5
var DelayTopic string = "delayTopic"
var RealTopic  string = "realTopic"
// var KafkaDelayQueue *KafkaDelayQueueProducer
// var KafkaDelayQueue    *KafkaDelayQueueProducer
func InitializekafkaGroup(){

	// 创建配置
	config := sarama.NewConfig()
	// 设置消费者组名称
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
    // config.Producer.RequiredAcks = sarama.WaitForAll          // ACK
    // config.Producer.Partitioner = sarama.NewRandomPartitioner // 分区
	//config.Consumer.Offsets.Initial = -2                    // 未找到组消费位移的时候从哪边开始消费
    // 异步回调(两个channel, 分别是成功和错误)
    config.Producer.Return.Successes = true // 确认
    config.Producer.Return.Errors = true
    ConsumerGroupReal, err := sarama.NewConsumerGroup(addrs, "group-1", config)
	if err != nil {
        log.Error(err)	
    }
	ConsumerGroupDelay, err = sarama.NewConsumerGroup(addrs, "group-2", config)
	if err != nil {
        log.Error(err)
	}
    ConsumerToRequestRta(ConsumerGroupReal)

    // return ConsumerGroupReal,ConsumerGroupDelay
    // KafkaDelayQueue = NewKafkaDelayQueueProducer(Producer, ConsumerGroupDelay, DelayTime, DelayTopic, RealTopic, log)
}

func GetKafkaDelayQueue() *KafkaDelayQueueProducer{
	KafkaDelayQueue := NewKafkaDelayQueueProducer(Producer, ConsumerGroupDelay, DelayTime, DelayTopic, RealTopic)
	return KafkaDelayQueue
}

//延时队列开始

// KafkaDelayQueueProducer 延迟队列生产者，包含了生产者和延迟服务
type KafkaDelayQueueProducer struct {
	producer   sarama.SyncProducer // 生产者
	delayTopic string              // 延迟服务主题
}

// DelayServiceConsumer 延迟服务消费者
type DelayServiceConsumer struct {
	producer  sarama.SyncProducer
	delay     time.Duration
	realTopic string
}
//启动延时队列服务
// NewKafkaDelayQueueProducer 创建延迟队列生产者
// producer 生产者
// delayServiceConsumerGroup 延迟服务消费者组
// delayTime 延迟时间
// delayTopic 延迟服务主题
// realTopic 真实队列主题
func NewKafkaDelayQueueProducer(producer sarama.SyncProducer,delayServiceConsumerGroup sarama.ConsumerGroup,delayTime time.Duration, delayTopic string, realTopic string) *KafkaDelayQueueProducer {
    // var producer sarama.SyncProducer
    // var delayServiceConsumerGroup sarama.ConsumerGroup
	// var (
	// 	signals = make(chan os.Signal, 1)
	// )
	// signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT, os.Interrupt) //信号通道
	// 启动延迟服务
	consumer := NewDelayServiceConsumer(producer, delayTime, realTopic)
	log.Info("[NewKafkaDelayQueueProducer] delay queue consumer start")
	go func() {
		for {
			if err := delayServiceConsumerGroup.Consume(context.Background(),
				[]string{delayTopic}, consumer); err != nil {
				log.Error("[NewKafkaDelayQueueProducer] delay queue consumer failed,err: ", zap.Error(err))
				break
			}
			time.Sleep(2 * time.Second)
			log.Info("[NewKafkaDelayQueueProducer] 检测消费函数是否一直执行")
			// 检查是否接收到中断信号，如果是则退出循环
			// select {
			// case sin := <-signals:
			// 	log.Info("[NewKafkaDelayQueueProducer]get signal,", zap.Any("signal", sin))
			// 	return
			// default:
			// }
		}
		log.Info("[NewKafkaDelayQueueProducer] consumer func exit")
	}()
	log.Info("[NewKafkaDelayQueueProducer] return KafkaDelayQueueProducer")

	return &KafkaDelayQueueProducer{
		producer:   producer,
		delayTopic: delayTopic,
	}
}


func NewDelayServiceConsumer(producer sarama.SyncProducer, delay time.Duration,
	realTopic string) *DelayServiceConsumer {
	return &DelayServiceConsumer{
		producer:  producer,
		delay:     delay,
		realTopic: realTopic,
	}
}

// SendMessage 发送消息到延时队列里面
func (q *KafkaDelayQueueProducer) SendMessage(msg *sarama.ProducerMessage) {
	msg.Topic = q.delayTopic
    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForAll          // ACK
    config.Producer.Partitioner = sarama.NewRandomPartitioner // 分区
    // 异步回调(两个channel, 分别是成功和错误)
    config.Producer.Return.Successes = true // 确认
    config.Producer.Return.Errors = true

    // sarama.Logger = log.New(os.Stdout, "[Sarama]", log.LstdFlags)

    // 连接kafka
    // 同步
    client, err := sarama.NewSyncProducer(addrs, config)
    // 异步
    // client, err := sarama.NewAsyncProducer(addrs, config)
    // if err != nil {
    //     fmt.Println("producer error", err)
    //     return
    // }

    defer func() {
        _ = client.Close()
    }()

    //封装消息
    pid, offset, err := client.SendMessage(msg)
    if err != nil {
        fmt.Println("send failed", err)
        return
    }

    fmt.Printf("pid:%v offset:%v \n", pid, offset)
	//client.Input() <- msg   //异步发送消息
}
// 实现 sarama.ConsumerGroup.Consume 接口的Setup
func (c *DelayServiceConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}
// 实现 sarama.ConsumerGroup.Consume 接口的Cleanup
func (c *DelayServiceConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

//延时队列结束


//消费者组消费
func (c *DelayServiceConsumer) ConsumeClaim(session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
    log.Info("[delaye ConsumerClaim] cc")
	for message := range claim.Messages() {
		fmt.Println("延迟消费启动了")
		fmt.Println(string(message.Value))
		now := time.Now()
		// topic := message.Topic
		if now.Sub(message.Timestamp) > c.delay {
			//加入真实队列
			// topic = c.realTopic
			log.Info("延时队列判断成功开始加入真实队列")
			msg := sarama.ProducerMessage{
				Topic:     c.realTopic,
				Timestamp: message.Timestamp,
				// Key:       sarama.ByteEncoder(message.Key),
				Value:     sarama.ByteEncoder(message.Value),
			}
			err := services.SyncKakfaService.ProducerSync(msg)
			if err != nil {
				log.Info("[delay ConsumeClaim] delay already send to real topic failed", zap.Error(err))
				return nil
			}
			if err == nil {
				session.MarkMessage(message, "")
				log.Info("[delay ConsumeClaim] delay already send to real topic success")
				continue
			}
		}
		// fmt.Println("结束了")
		// session.MarkMessage(message, "")
		// 否则休眠一秒
		time.Sleep(time.Second)
		return nil
	}

	log.Info("[delay ConsumeClaim] ph",
		zap.Any("partitiion", claim.Partition()),
		zap.Any("HighWaterMarkOffset", claim.HighWaterMarkOffset()))
    log.Info("[delay ConsumeClaim] delay consumer end")
	return nil
}

type ConsumerRta struct {
}

func ConsumerToRequestRta(consumerGroup sarama.ConsumerGroup) {
	var (
		// signals = make(chan os.Signal, 1)
		wg = &sync.WaitGroup{}
	)
	// signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	wg.Add(1)
	// 启动消费者协程
	go func() {
		defer wg.Done()
		consumer := NewConsumerRta()
		log.Info("[ConsumerToRequestRta] consumer group start")
		// 执行消费者组消费
		for {
			log.Info("[ConsumerToRequestRta] consumer group start1")
			// err := consumerGroup.Consume(context.Background(), []string{"test"}, consumer)
			// fmt.Println(err)
			if err := consumerGroup.Consume(context.Background(), []string{RealTopic}, consumer); err != nil {
				log.Error("[ConsumerToRequestRta] Error from consumer group:", zap.Error(err))
				break
			}
			log.Info("[ConsumerToRequestRta] consumer group start2")
			time.Sleep(2 * time.Second) // 等待一段时间后重试

			//检查是否接收到中断信号，如果是则退出循环
			// select {
			// case sin := <-signals:
			// 	log.Info("get signal,", zap.Any("signal", sin))
			// 	return
			// }
		}
	}()
	// select {
	// case sin := <-signals:
	// 	log.Info("get signal,", zap.Any("signal", sin))
	// 	return
	// }	
	wg.Wait()
	log.Info("[ConsumerToRequestRta] consumer end & exit")
}

func NewConsumerRta() *ConsumerRta {
	return &ConsumerRta{
	}
}



func (c *ConsumerRta) ConsumeClaim(session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	fmt.Println("真实队列启动")
	fmt.Println(claim)
	for message := range claim.Messages() {
		// 消费逻辑
		fmt.Println(message)
		fmt.Println("真实队列开始消费")
		setConsumerFunction(message)
		// fmt.Println(message.Topic)
		session.MarkMessage(message, "")
		return nil
	}
	fmt.Println("真实队列结束")
	return nil
}

func (c *ConsumerRta) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerRta) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
