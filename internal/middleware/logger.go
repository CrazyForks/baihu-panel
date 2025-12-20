package middleware

import (
	"fmt"
	"time"

	"baihu/internal/logger"

	"github.com/gin-gonic/gin"
)

// GinLogger 返回使用 logrus 的 Gin 日志中间件
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		if query != "" {
			path = path + "?" + query
		}

		msg := fmt.Sprintf("%3d | %13v | %15s | %-7s %s",
			status, latency, clientIP, method, path)

		if status >= 500 {
			logger.Error(msg)
		} else if status >= 400 {
			logger.Warn(msg)
		} else {
			logger.Info(msg)
		}
	}
}

// GinRecovery 返回使用 logrus 的 Gin 恢复中间件
func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("Panic recovered: %v | path: %s", err, c.Request.URL.Path)
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
