package gin

import (
	"github.com/JGLTechnologies/gin-rate-limit"
	"github.com/exgamer/go-sdk/pkg/config"
	"github.com/exgamer/go-sdk/pkg/exception"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/redis/go-redis/v9"
	timeout "github.com/vearne/gin-timeout"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"net/http"
	"time"
)

// InitRouter Базовая инициализация gin
func InitRouter(baseConfig *config.BaseConfig) *gin.Engine {
	if baseConfig.AppEnv == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Options
	router := gin.New()
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "404 page not found"})
	})
	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"code": "METHOD_NOT_ALLOWED", "message": "405 method not allowed"})
	})
	p := ginprometheus.NewPrometheus("ginHelpers")
	p.Use(router)
	router.Use(sentrygin.New(sentrygin.Options{}))
	router.Use(gin.Logger())
	router.Use(timeout.Timeout(timeout.WithTimeout(time.Duration(baseConfig.HandlerTimeout) * time.Second)))
	router.Use(gin.CustomRecovery(ErrorHandler))

	return router
}

func getGinRedisRateLimiter(redis *redis.Client, limit uint) gin.HandlerFunc {
	// This makes it so each ip can only make n requests per second
	store := ratelimit.RedisStore(&ratelimit.RedisOptions{
		RedisClient: redis,
		Rate:        time.Second,
		Limit:       limit,
	})

	return ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: func(context *gin.Context, info ratelimit.Info) {
			context.String(http.StatusTooManyRequests, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
		},
		KeyFunc: func(context *gin.Context) string {
			return context.ClientIP()
		},
	})
}

// ErrorHandler Обработчик ошибок gin
func ErrorHandler(c *gin.Context, err any) {
	goErr := errors.Wrap(err, 2)
	details := make([]string, 0)

	for _, frame := range goErr.StackFrames() {
		details = append(details, frame.String())
	}

	sentry.CaptureException(goErr)
	c.JSON(http.StatusInternalServerError, gin.H{"message": goErr.Error(), "details": details})
}

func Error(c *gin.Context, exception *exception.AppException) {
	c.Set("exception", exception)
	c.Status(exception.Code)
}

func Success(c *gin.Context, data any) {
	c.Set("data", data)
}
