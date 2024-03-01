package httpMiddleware

import (
	"fmt"
	"github.com/exgamer/go-sdk/pkg/config"
	"github.com/exgamer/go-sdk/pkg/exception"
	"github.com/exgamer/go-sdk/pkg/logger"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"net/http"
)

// ResponseMiddleware Middleware для обработки ответа
func ResponseMiddleware(appInfo *config.AppInfo) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			sentry.CaptureException(err)
			logger.FormattedErrorWithAppInfo(appInfo, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

			return
		}

		appExceptionObject, exists := c.Get("exception")
		fmt.Printf("%+v\n", appExceptionObject)

		if !exists {
			data, _ := c.Get("data")
			//logInfo("", c, config)
			c.JSON(http.StatusOK, gin.H{"success": true, "data": data})

			return
		}

		appException := exception.AppException{}
		mapstructure.Decode(appExceptionObject, &appException)
		sentry.CaptureException(appException.Error)
		fmt.Printf("%+v\n", appException)
		logger.FormattedErrorWithAppInfo(appInfo, appException.Error.Error())
		c.JSON(appException.Code, gin.H{"message": appException.Error.Error(), "details": appException.Context})
	}
}
