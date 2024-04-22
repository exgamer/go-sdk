package structures

// HttpResponse Модель описывающая ответ от rest запроса
type HttpResponse[E interface{}] struct {
	Status     string
	Body       []byte
	StatusCode int
	Url        string
	Method     string
	Headers    map[string]string
	Result     E
}
