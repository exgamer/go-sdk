package middleware

import (
	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

// RateLimiterMiddleware Middleware для rate limit
func RateLimiterMiddleware(redis *redis.Client, limit uint) gin.HandlerFunc {
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
