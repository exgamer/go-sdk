package messenger

import (
	"encoding/json"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/exgamer/go-sdk/pkg/config"
	"github.com/exgamer/go-sdk/pkg/kafka/messenger/structures"
	"github.com/exgamer/go-sdk/pkg/logger"
	"github.com/exgamer/go-sdk/pkg/validation"
)

func NewMessageSender(appInfo *config.AppInfo, configMap *kafka.ConfigMap) *MessageSender {

	if configMap == nil {
		configMap = &kafka.ConfigMap{
			"bootstrap.servers": "localhost:9092",
		}
	}

	return &MessageSender{
		appInfo:   appInfo,
		configMap: configMap,
	}
}

// MessageSender - отсылка сообщений через кафку
type MessageSender struct {
	appInfo   *config.AppInfo
	configMap *kafka.ConfigMap
}

func (s *MessageSender) SendSms(phone string, text string) error {
	producer, _ := kafka.NewProducer(s.configMap)

	smsMessage := structures.SmsMessage{
		Phone:       phone,
		Text:        text,
		ServiceName: s.appInfo.ServiceName,
	}

	if !validation.CheckValidPhone(smsMessage.Phone) {
		message := "Sms message send error: invalid phone number:" + phone
		logger.FormattedError(s.appInfo.ServiceName, s.appInfo.RequestMethod, s.appInfo.RequestUrl, 0, s.appInfo.RequestId, message)

		return errors.New("invalid phone number")
	}

	topic := s.appInfo.AppEnv + "." + "messenger-service.command.sms" // хард код потому что по идее никогда не изменится
	jsonValue, _ := json.Marshal(smsMessage)
	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          jsonValue,
	}, nil)
	message := "Sms message send: phone:" + phone + "; text:" + text
	logger.FormattedInfo(s.appInfo.ServiceName, s.appInfo.RequestMethod, s.appInfo.RequestUrl, 0, s.appInfo.RequestId, message)

	return nil
}
