package kafka

// KafkaConfig Данные для соединения с Кафкой
type KafkaConfig struct {
	Host string `mapstructure:"KAFKA_HOST"`
}
