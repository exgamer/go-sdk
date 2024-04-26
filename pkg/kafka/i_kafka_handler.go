package kafka

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type IKafkaHandler interface {
	Handle(consumer *kafka.Consumer, message *kafka.Message) error
}
