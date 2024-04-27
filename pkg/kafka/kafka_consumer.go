package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/davecgh/go-spew/spew"
	"github.com/exgamer/go-sdk/pkg/config"
	"github.com/exgamer/go-sdk/pkg/logger"
	"github.com/exgamer/go-sdk/pkg/sentry"
	"os"
	"os/signal"
	"syscall"
)

func NewConsumer(appInfo *config.AppInfo, handlers map[string]IKafkaHandler, configMap *kafka.ConfigMap) *KafkaConsumer {
	return &KafkaConsumer{
		appInfo:   appInfo,
		handlers:  handlers,
		run:       true,
		consume:   false,
		configMap: configMap,
	}
}

type KafkaConsumer struct {
	consumer  *kafka.Consumer
	appInfo   *config.AppInfo
	handlers  map[string]IKafkaHandler
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

	h := map[string]IKafkaHandler{}
	var topicList []string

	for topic, handler := range kc.handlers {
		topicList = append(topicList, kc.appInfo.AppEnv+"."+topic)
		h[kc.appInfo.AppEnv+"."+topic] = handler
	}

	kc.handlers = h

	if kc.handlers == nil {
		kc.handlers = map[string]IKafkaHandler{
			"default": DefaultKafkaHandler{},
		}
	}

	spew.Dump(kc.handlers)
	spew.Dump(topicList)
	kc.consumer.SubscribeTopics(topicList, nil)
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
					topic := e.TopicPartition.Topic
					handler := kc.handlers[*topic].(IKafkaHandler)

					if handler == nil {
						handler = kc.handlers["default"].(IKafkaHandler)
					}

					message := "groupId:" + groupId + ";key:" + string(e.Key) + "; value:" + string(e.Value)
					logger.FormattedInfo(kc.appInfo.ServiceName, "consumer", *e.TopicPartition.Topic, 0, kc.appInfo.RequestId, message)
					err := handler.Handle(kc.consumer, e)

					if err != nil {
						message = message + "; error_text:" + err.Error()
						logger.FormattedError(kc.appInfo.ServiceName, "consumer", *e.TopicPartition.Topic, 0, kc.appInfo.RequestId, message)
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
						kc.consumer.CommitMessage(e) // кафка не имеет механизма удаления сообщения, поэтому при ошибке просто комитим его
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
