Пример использования функции для отправки смс через кафку

```go
package app

"github.com/exgamer/go-sdk/pkg/kafka/messenger"
kafkaLib "github.com/confluentinc/confluent-kafka-go/v2/kafka"

func main() {
	messageSender := messenger.NewMessageSender(app.appInfo, &kafkaLib.ConfigMap{
		"bootstrap.servers": strings.Join([]string{app.kafkaConfig.Host}, ","),
	})

	err := messageSender.SendSms("+7 (778)-300-00-44", "test message")

	if err != nil {

		println(err)
	}
}

```