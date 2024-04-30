package structures

// SmsMessage - модель пэйлоада входящего сообщения из кафки для отправки смс
type SmsMessage struct {
	ServiceName string `json:"service_name"  validate:"required"`
	Text        string `json:"text"  validate:"required"`
	Phone       string `json:"phone"  validate:"required"`
}
