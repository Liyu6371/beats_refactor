package kafka

import (
	"beats_refactor/source"
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/mitchellh/mapstructure"
)

func init() {
	source.RegisterSourceInstance(Name, New)
}

const Name = "kafka"

type KafkaSourceInstance struct {
	c KafkaInstanceConfig

	alive bool
	mu    sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc

	client  sarama.ConsumerGroup
	handler sarama.ConsumerGroupHandler
}

func New(conf interface{}) (source.SourceInstance, error) {
	c := KafkaInstanceConfig{}
	if err := mapstructure.Decode(conf, &c); err != nil {
		return nil, fmt.Errorf("kafka decode config error: %w", err)
	}
	if !c.Enabled {
		return nil, fmt.Errorf("kafka source is not enabled")
	}
	saramaConf := sarama.NewConfig()
	// Set the kafka version
	if c.Version != "" {
		if v, err := sarama.ParseKafkaVersion(c.Version); err == nil {
			saramaConf.Version = v
		}
	}
	// kafka 认证设置
	if c.UserName != "" && c.Password != "" {
		saramaConf.Net.SASL.Enable = true
		saramaConf.Net.SASL.User = c.UserName
		saramaConf.Net.SASL.Password = c.Password
	}
	// 设置消费者组消费位置
	if c.ConsumeOldest {
		saramaConf.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		saramaConf.Consumer.Offsets.Initial = sarama.OffsetNewest
	}
	// 配置消费者组分区平衡策略
	if v, ok := kafkaRebalanceMap[c.Assignor]; ok {
		saramaConf.Consumer.Group.Rebalance.Strategy = v
	} else {
		saramaConf.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	}
	// 配置消费者组名称
	consumerGroupName := "default_beats_kafka_consumer_group"
	if c.ConsumerGroup != "" {
		consumerGroupName = c.ConsumerGroup
	}

	client, err := sarama.NewConsumerGroup(c.Hosts, consumerGroupName, saramaConf)
	if err != nil {
		return nil, fmt.Errorf("kafka create consumer group error: %w", err)
	}
	return &KafkaSourceInstance{
		c:      c,
		client: client,
	}, nil
}

// Alive implements source.SourceInstance.
func (k *KafkaSourceInstance) Alive() bool {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.alive
}

func (k *KafkaSourceInstance) turnOn() {
	k.mu.Lock()
	defer k.mu.Unlock()
	if !k.alive {
		k.alive = true
	}
}

func (k *KafkaSourceInstance) turnOff() {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.alive {
		k.alive = false
	}
}

// GetName implements source.SourceInstance.
func (k *KafkaSourceInstance) GetName() string {
	return Name
}

// Start implements source.SourceInstance.
func (k *KafkaSourceInstance) Start(ctx context.Context, ch chan<- []byte) {
	k.ctx, k.cancel = context.WithCancel(ctx)
	k.handler = &consumerGroupHandler{ch: ch}

	k.turnOn()
	defer func() {
		k.turnOff()
		// 确保在退出时取消上下文
		// 以及关闭消费者组
		if k.cancel != nil {
			k.cancel()
		}
		if err := k.client.Close(); err != nil {
			fmt.Printf("kafka consumer group close error: %v\n", err)
		}
	}()

	for {
		if err := k.client.Consume(k.ctx, k.c.Topics, k.handler); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				// fmt.Errorf("kafka source instance: error is : %s", err)
				return
			}
			// logger.Errorf("kafka source instance: Error from consumer: %s", err)
			return
		}
		// 检查上下文是否被取消
		if k.ctx.Err() != nil {
			// logger.Infof("kafka source instance: context cancelled, exiting consume loop")
			return
		}
	}
}

// Stop implements source.SourceInstance.
func (k *KafkaSourceInstance) Stop() {
	k.cancel()
}
