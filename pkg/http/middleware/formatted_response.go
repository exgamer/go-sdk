package middleware

import (
	"github.com/exgamer/go-sdk/pkg/http/helpers"
	"github.com/gin-gonic/gin"
)

// FormattedResponseMiddleware Middleware для обработки ответа
func FormattedResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		helpers.FormattedResponse(c)
	}
}
