package config

import "github.com/google/uuid"

// AppInfo Данные приложения
type AppInfo struct {
	RequestId     string
	LanguageCode  string
	RequestScheme string
	RequestHost   string
	RequestMethod string
	RequestUrl    string
	ServiceName   string
}

func (s *AppInfo) GenerateRequestId() {
	s.RequestId = uuid.New().String()
}

func (s *AppInfo) SetConsoleMode(name string) {
	s.RequestId = uuid.New().String()
	s.RequestMethod = "console"
	s.RequestUrl = name
	s.GenerateRequestId()
}
