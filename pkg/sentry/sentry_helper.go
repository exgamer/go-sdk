package sentry

import (
	"github.com/getsentry/sentry-go"
	"time"
)

// SendError - отправить ошибку в sentry
func SendError(message string, tags map[string]string, extra map[string]interface{}) {
	defer sentry.Flush(2 * time.Second)
	event := sentry.NewEvent()
	event.Message = message
	event.Tags = tags
	event.Extra = extra
	event.Level = sentry.LevelError

	sentry.CaptureEvent(event)
}
