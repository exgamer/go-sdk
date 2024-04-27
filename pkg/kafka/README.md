Пример использования консьюмера для кафки

```go
package app

func (app *App) RunConsumers() {
	//Указываем массив с хостами кафки
	brokerList := []string{app.kafkaConfig.Host}
	//создаем мапу с обработчиками
	handlers := map[string]kafka.IKafkaHandler{
		"messenger-service.command.sms": DefaultKafkaHandler{},
	}
    //инициализация консьюмера
	consumer := kafka.NewConsumer(app.appInfo, handlers, &kafkaLib.ConfigMap{
		"bootstrap.servers":     strings.Join(brokerList, ","),
		"broker.address.family": "v4",
		"group.id":              app.baseConfig.Name,
		"auto.offset.reset":     "earliest",
		"enable.auto.commit":    "false",
	})
	err := consumer.Init()

	if err != nil {
		log.Println(err)
		panic(err)
	}
    //старт консьюмера
	consumer.StartConsume()

	done := make(chan struct{})
	<-done
}

//обработчик сообщения кафки
type DefaultKafkaHandler struct{}

func (d DefaultKafkaHandler) Handle(consumer *kafka.Consumer, message *kafka.Message) error {
	fmt.Printf("Topic: %s, Message: %+v\n", *message.TopicPartition.Topic, string(message.Value))
	consumer.CommitMessage(message)

	return nil
}

```