Пример использования консьюмера для кафки

```go
package app

func (app *App) RunConsumers() {
	brokerList := []string{app.kafkaConfig.Host}
	topicList := []string{
		"messenger-service.command.sms", //@TODO константы
	}

	consumer := NewConsumer(app.appInfo, topicList, &kafkaLib.ConfigMap{
		"bootstrap.servers":     strings.Join(brokerList, ","),
		"broker.address.family": "v4",
		"group.id":              app.baseConfig.Name,
		"auto.offset.reset":     "earliest",
		"enable.auto.commit":    "false",
	}, nil)
	err := consumer.Init()

	if err != nil {
		log.Println(err)
		panic(err)
	}

	//go ProduceTestMessages(brokerList, "group1", "local.messenger-service.command.message")
	consumer.StartConsume()

	done := make(chan struct{})
	<-done
}

```