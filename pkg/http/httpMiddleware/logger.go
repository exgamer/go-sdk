package httpMiddleware

import (
	"github.com/exgamer/go-sdk/pkg/config"
	"github.com/exgamer/go-sdk/pkg/exception"
	"github.com/exgamer/go-sdk/pkg/logger"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

// LoggerMiddleware Middleware для логирования ответа и отправки ошибок в сентри
func LoggerMiddleware(appInfo *config.AppInfo) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			sentry.CaptureException(err)
			logger.FormattedErrorWithAppInfo(appInfo, err.Error())
		}

		appExceptionObject, exists := c.Get("exception")

		if exists {
			appException := exception.AppException{}
			mapstructure.Decode(appExceptionObject, &appException)
			sentry.CaptureException(appException.Error)
			logger.FormattedErrorWithAppInfo(appInfo, appException.Error.Error())
		}
	}
}
