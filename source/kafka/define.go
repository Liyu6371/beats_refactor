package kafka

import (
	"github.com/IBM/sarama"
)

var (
	kafkaRebalanceMap = map[string]sarama.BalanceStrategy{
		"sticky":     sarama.NewBalanceStrategySticky(),
		"roundrobin": sarama.NewBalanceStrategyRoundRobin(),
		"range":      sarama.NewBalanceStrategyRange(),
	}
)

type KafkaInstanceConfig struct {
	Enabled       bool     `yaml:"enabled"`
	UserName      string   `yaml:"user_name,omitempty"`
	Password      string   `yaml:"password,omitempty"`
	Version       string   `yaml:"version"`
	ConsumerGroup string   `yaml:"kafka_consumer_group"`
	ConsumeOldest bool     `yaml:"kafka_consume_oldest"`
	Assignor      string   `yaml:"kafka_assignor"`
	Hosts         []string `yaml:"hosts"`
	Topics        []string `yaml:"topics"`
}
