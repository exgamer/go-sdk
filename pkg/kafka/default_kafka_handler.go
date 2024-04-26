package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type DefaultKafkaHandler struct{}

func (d DefaultKafkaHandler) Handle(consumer *kafka.Consumer, message *kafka.Message) error {
	fmt.Printf("Topic: %s, Message: %+v\n", *message.TopicPartition.Topic, string(message.Value))
	consumer.CommitMessage(message)

	return nil
}
