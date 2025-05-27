package kafka

import (
	"beats_refactor/logger"

	"github.com/IBM/sarama"
)

func createConsumeClient(c KafkaInstanceConfig) (sarama.ConsumerGroup, error) {
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
		return nil, err
	}
	return client, nil
}

type consumerGroupHandler struct {
	ch chan<- []byte
}

func (c *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				logger.Infof("Source->%s message chan claim is closed.", Name)
				return nil
			}
			c.ch <- msg.Value
		case <-session.Context().Done():
			logger.Infof("Source->%s context done, exiting consume claim loop.", Name)
			return nil
		}
	}
}
