package httpMiddleware

import (
	"github.com/exgamer/go-sdk/pkg/http/httpHelpers"
	"github.com/gin-gonic/gin"
)

// FormattedResponseMiddleware Middleware для обработки ответа
func FormattedResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		httpHelpers.FormattedResponse(c)
	}
}
