package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/davecgh/go-spew/spew"
	"github.com/exgamer/go-sdk/pkg/config"
	"github.com/exgamer/go-sdk/pkg/logger"
	"github.com/exgamer/go-sdk/pkg/sentry"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func NewConsumer(appInfo *config.AppInfo, topicList []string, configMap *kafka.ConfigMap, writer *DefaultKafkaHandler) *KafkaConsumer {
	topics := make([]string, len(topicList))

	for i, topic := range topicList {
		topics[i] = appInfo.AppEnv + "." + topic
	}

	return &KafkaConsumer{
		appInfo:   appInfo,
		topicList: topics,
		writer:    writer,
		run:       true,
		consume:   false,
		configMap: configMap,
	}
}

type KafkaConsumer struct {
	consumer  *kafka.Consumer
	appInfo   *config.AppInfo
	topicList []string
	writer    *DefaultKafkaHandler
	run       bool
	consume   bool
	configMap *kafka.ConfigMap
}

func (kc *KafkaConsumer) SetConfig(configMap *kafka.ConfigMap) {
	kc.configMap = configMap
}

func (kc *KafkaConsumer) Init() error {
	spew.Dump(kc.configMap)

	if kc.configMap == nil {
		kc.configMap = &kafka.ConfigMap{
			"bootstrap.servers":     "localhost:9092",
			"broker.address.family": "v4",
			"group.id":              "default-group",
			"auto.offset.reset":     "earliest",
		}
	}

	c, err := kafka.NewConsumer(kc.configMap)

	if err != nil {
		return err
	}

	kc.consumer = c

	kc.consumer.SubscribeTopics(kc.topicList, nil)

	if kc.writer == nil {
		kc.writer = &DefaultKafkaHandler{}
	}

	go kc.initConsume()

	return nil
}

func (kc *KafkaConsumer) initConsume() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for kc.run {
		for kc.consume {
			select {
			case sig := <-sigchan:
				fmt.Printf("Caught signal %v: terminating\n", sig)
				kc.run = false
				kc.consume = false
			default:
				ev := kc.consumer.Poll(100)

				if ev == nil {
					continue
				}

				configValue, _ := kc.configMap.Get("group.id", "-")
				groupId := configValue.(string)

				switch e := ev.(type) {
				case *kafka.Message:
					message := "key:" + string(e.Key) + "; value:" + string(e.Value)
					logger.FormattedInfo(kc.appInfo.ServiceName, "", *e.TopicPartition.Topic, 0, groupId, message)
					err := kc.writer.Handle(kc.consumer, e)

					if err != nil {
						log.Println(err)
						sentry.SendError("Kafka Consumer Error: "+err.Error(),
							map[string]string{
								"service_name": kc.appInfo.ServiceName,
								"env":          kc.appInfo.AppEnv,
								"kafka_group":  groupId,
							},
							map[string]interface{}{
								"key":             string(e.Key),
								"value":           string(e.Value),
								"topic_partition": e.TopicPartition,
								"timestamp":       e.Timestamp,
							},
						)
					}
				case kafka.Error:
					// Errors should generally be considered
					// informational, the client will try to
					// automatically recover.
					fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)

					sentry.SendError("Kafka Error: "+e.Error(),
						map[string]string{
							"service_name": kc.appInfo.ServiceName,
							"env":          kc.appInfo.AppEnv,
							"kafka_group":  configValue.(string),
						},
						map[string]interface{}{
							"error": e,
						},
					)

					if e.Code() == kafka.ErrAllBrokersDown {
						kc.run = false
						kc.consume = false
					}

				default:
					fmt.Printf("Ignored %v\n", e)
				}
			}
		}
	}
}

func (kc *KafkaConsumer) StartConsume() {
	fmt.Println("Starting consumer")
	kc.consume = true
}

func (kc *KafkaConsumer) StopConsume() {
	fmt.Println("Stopping consumer")
	kc.consume = false
}

func (kc *KafkaConsumer) Close() {
	fmt.Println("Closing consumer")
	kc.run = false
	kc.consume = false
	kc.consumer.Close()
}
