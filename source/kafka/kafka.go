package kafka

import (
	"beats_refactor/logger"
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
	c  KafkaInstanceConfig
	wg sync.WaitGroup

	alive bool
	mu    sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc
}

func New(conf interface{}) (source.SourceInstance, error) {
	c := KafkaInstanceConfig{}
	if err := mapstructure.Decode(conf, &c); err != nil {
		return nil, fmt.Errorf("kafka decode config error: %w", err)
	}
	if !c.Enabled {
		return nil, fmt.Errorf("kafka source is not enabled")
	}
	return &KafkaSourceInstance{
		c: c,
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
	k.turnOn()
	k.ctx, k.cancel = context.WithCancel(ctx)
	// 拉起一定数量的 goroutine 来进行 kafka 的数据消费
	for i := range k.c.WorkerNum {
		k.wg.Add(1)
		go func(workId int, ch chan<- []byte) {
			defer k.wg.Done()
			client, err := createConsumeClient(k.c)
			if err != nil {
				logger.Errorf("kafka source instance worker: %d create consumer client error: %v", workId, err)
				return
			}
			defer func() {
				if err := client.Close(); err != nil {
					logger.Errorf("kafka source instance worker: %d error closing consumer client: %v", workId, err)
				}
			}()
			for {
				if err := client.Consume(k.ctx, k.c.Topics, &consumerGroupHandler{ch: ch}); err != nil {
					if errors.Is(err, sarama.ErrClosedConsumerGroup) {
						logger.Errorf("kafka source instance worker: %d error is : %s", workId, err)
						return
					}
					logger.Errorf("kafka source instance: error from consumer: %s", err)
					return
				}
				// 检查上下文是否被取消
				if k.ctx.Err() != nil {
					logger.Infof("kafka source instance worker: %d context cancelled, exiting consume loop", workId)
					return
				}
			}
		}(i, ch)
	}
}

// Stop implements source.SourceInstance.
func (k *KafkaSourceInstance) Stop() {
	k.cancel()
	k.wg.Wait() // 等待所有 goroutine 完成
	k.turnOff()
}
