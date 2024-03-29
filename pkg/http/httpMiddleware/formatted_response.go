package httpMiddleware

import (
	"fmt"
	"github.com/exgamer/go-sdk/pkg/exception"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"net/http"
)

// FormattedResponseMiddleware Middleware для обработки ответа
func FormattedResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
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
		fmt.Printf("%+v\n", appException)

		c.JSON(appException.Code, gin.H{"message": appException.Error.Error(), "details": appException.Context})
	}
}
