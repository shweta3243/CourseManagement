package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware captures URL, Method, Time Taken, and Status Code for every request
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// After request
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		log.Printf("| %3d | %13v | %-7s | %s",
			statusCode,
			duration,
			method,
			path,
		)
	}
}
