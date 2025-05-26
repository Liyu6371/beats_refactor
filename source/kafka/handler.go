package kafka

import (
	"github.com/IBM/sarama"
)

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
				// logger.Infof("Source->%s message chan claim is closed.", Name)
				return nil
			}
			c.ch <- msg.Value
		case <-session.Context().Done():
			// logger.Info("session exit")
			return nil
		}
	}
}
