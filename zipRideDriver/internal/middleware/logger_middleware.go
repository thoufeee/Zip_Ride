package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()
		ua := c.Request.UserAgent()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Error("request",
					zap.Int("status", status),
					zap.String("method", method),
					zap.String("path", path),
					zap.String("ip", clientIP),
					zap.String("ua", ua),
					zap.Duration("latency", latency),
					zap.String("error", e.Error()),
				)
			}
			return
		}

		log.Info("request",
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", clientIP),
			zap.String("ua", ua),
			zap.Duration("latency", latency),
		)
	}
}
