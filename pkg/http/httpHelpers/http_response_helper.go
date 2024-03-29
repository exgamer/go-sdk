package httpHelpers

import (
	"fmt"
	"github.com/exgamer/go-sdk/pkg/exception"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"net/http"
)

func ErrorResponse(c *gin.Context, statusCode int, err error, context map[string]any) {
	AppExceptionResponse(c, exception.NewAppException(statusCode, err, context))
}

func AppExceptionResponse(c *gin.Context, exception *exception.AppException) {
	c.Set("exception", exception)
	c.Status(exception.Code)
}

func SuccessResponse(c *gin.Context, data any) {
	c.Set("data", data)
}

func FormattedResponse(c *gin.Context) {
	for _, err := range c.Errors {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

		return
	}

	appExceptionObject, exists := c.Get("exception")
	fmt.Printf("%+v\n", appExceptionObject)

	if !exists {
		data, _ := c.Get("data")

		c.JSON(http.StatusOK, gin.H{"success": true, "data": data})

		return
	}

	appException := exception.AppException{}
	mapstructure.Decode(appExceptionObject, &appException)
	fmt.Printf("%+v\n", appException)

	c.JSON(appException.Code, gin.H{"message": appException.Error.Error(), "details": appException.Context})
}
