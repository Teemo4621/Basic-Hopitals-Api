package middlewares

import (
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		duration := time.Since(start)
		method := color.New(color.FgCyan, color.Bold).Sprint(c.Request.Method)
		url := color.New(color.FgGreen).Sprint(c.Request.URL.String())
		query := color.New(color.FgHiMagenta).Sprint(c.Request.URL.Query().Encode())

		log.Printf("📥 %s | %s | Query: %s | ⏱️ %v ",
			method,
			url,
			query,
			duration,
		)

		c.Next()
	}
}
