package bootstrap

import (
    "context"
    "fmt"
    "github.com/IBM/sarama"
	"testproject/app/services"
    // "log"
    "go.uber.org/zap"
    "time"
    "os/signal"
    // "log"
    // "math"
    "os"
    // "os/signal"
    "sync"
    "syscall"
    log "github.com/sirupsen/logrus"
)

var addrs = []string{"192.168.31.232:9092"}
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
                    fmt.Println(msg.Topic);
                    //执行消费代码
                    setConsumerFunction(msg.Topic,string(msg.Value))
                    fmt.Printf("Partition:%d Offset:%d Key:%s Value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
                }
            }(pc)
        }
    }

    //time.Sleep(5 * time.Second)
    select {}
}

//根据topic而选择不同的消费
func setConsumerFunction(topic string,value string){
    switch topic {
    case "sun":
        services.ElasticSearchHander.Sync(value);
    case "redisKeyDel"://用于删除缓存的真实队列
        services.ElasticSearchHander.Sync(value);
    case "redisKeyDelDelay": //延迟队列
        services.ElasticSearchHander.Sync(value);
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
// var KafkaDelayQueue    *KafkaDelayQueueProducer
func InitializekafkaGroup(){

	// 创建配置
	config := sarama.NewConfig()
	// 设置消费者组名称
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
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
// func GetKafkaDelayQueue() {
// 	KafkaDelayQueue = NewKafkaDelayQueueProducer(Producer, ConsumerGroupDelay, DelayTime, DelayTopic, RealTopic, log)
// }

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
	var (
		signals = make(chan os.Signal, 1)
	)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT, os.Interrupt) //信号通道
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
			select {
			case sin := <-signals:
				log.Info("[NewKafkaDelayQueueProducer]get signal,", zap.Any("signal", sin))
				return
			default:
			}
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
func (q *KafkaDelayQueueProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	msg.Topic = q.delayTopic
    //client.Input() <- msg   //异步发送消息
	return q.producer.SendMessage(msg)
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
        setConsumerFunction(message.Topic,string(message.Value)) //根据不同选择加入不同的真实队列处理
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
		signals = make(chan os.Signal, 1)
		wg = &sync.WaitGroup{}
	)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	wg.Add(1)
	// 启动消费者协程
	go func() {
		defer wg.Done()
		consumer := NewConsumerRta()
		log.Info("[ConsumerToRequestRta] consumer group start")
		// 执行消费者组消费
		for {
			if err := consumerGroup.Consume(context.Background(), []string{"test","test2","test3"}, consumer); err != nil {
				log.Error("[ConsumerToRequestRta] Error from consumer group:", zap.Error(err))
				break
			}
			time.Sleep(2 * time.Second) // 等待一段时间后重试

			// 检查是否接收到中断信号，如果是则退出循环
			select {
			case sin := <-signals:
				log.Info("get signal,", zap.Any("signal", sin))
				return
			default:
			}
		}
	}()
	wg.Wait()
	log.Info("[ConsumerToRequestRta] consumer end & exit")
}

func NewConsumerRta() *ConsumerRta {
	return &ConsumerRta{
	}
}



func (c *ConsumerRta) ConsumeClaim(session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		// 消费逻辑
		session.MarkMessage(message, "")
		return nil
	}

	return nil
}

func (c *ConsumerRta) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerRta) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
