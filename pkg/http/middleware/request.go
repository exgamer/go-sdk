package middleware

import (
	"github.com/exgamer/go-sdk/pkg/config"
	"github.com/exgamer/go-sdk/pkg/constants"
	"github.com/gin-gonic/gin"
)

// RequestMiddleware Middleware заполняющий данные запроса
func RequestMiddleware(appInfo *config.AppInfo) gin.HandlerFunc {
	return func(c *gin.Context) {
		appInfo.RequestId = c.GetHeader(constants.RequestIdHeaderName)
		// если request id не пришел с заголовком, генерим его, чтобы прокидывать дальше при http запросах
		if appInfo.RequestId == "" {
			appInfo.GenerateRequestId()
			c.Request.Header.Add(constants.RequestIdHeaderName, appInfo.RequestId)
		}

		appInfo.LanguageCode = c.GetHeader(constants.LanguageHeaderName)
		appInfo.RequestUrl = c.Request.URL.Path
		appInfo.RequestMethod = c.Request.Method
		appInfo.RequestScheme = c.Request.URL.Scheme
		appInfo.RequestHost = c.Request.Host

		c.Next()
	}
}
