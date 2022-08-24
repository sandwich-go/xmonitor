package xmonitor

import (
	"github.com/gin-gonic/gin"
)

func NewGinMonitor(collector Collector) gin.HandlerFunc {
	return func(c *gin.Context) {
		after := collector.MonitorRequest(c.Request)
		c.Next()
		after(c.Writer.Status(), c.Writer.Size())
	}
}
