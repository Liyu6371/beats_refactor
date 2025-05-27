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
	Enabled       bool     `yaml:"enabled" mapstructure:"enabled"`
	WorkerNum     int      `yaml:"worker_num" mapstructure:"worker_num"`
	UserName      string   `yaml:"user_name,omitempty" mapstructure:"user_name,omitempty"`
	Password      string   `yaml:"password,omitempty" mapstructure:"password,omitempty"`
	Version       string   `yaml:"version" mapstructure:"version"`
	ConsumerGroup string   `yaml:"kafka_consumer_group" mapstructure:"kafka_consumer_group"`
	ConsumeOldest bool     `yaml:"kafka_consume_oldest" mapstructure:"kafka_consume_oldest"`
	Assignor      string   `yaml:"kafka_assignor" mapstructure:"kafka_assignor"`
	Hosts         []string `yaml:"hosts" mapstructure:"hosts"`
	Topics        []string `yaml:"topics" mapstructure:"topics"`
}
