package internal

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start).Seconds()
		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		method := c.Request.Method

		HTTPRequestsTotal.WithLabelValues(route, method).Inc()
		RequestDuration.WithLabelValues(route, method).Observe(duration)
		StatusCodesTotal.WithLabelValues(strconv.Itoa(c.Writer.Status())).Inc()
	}
}
