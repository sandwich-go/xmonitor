package xmonitor

import (
	"github.com/gin-gonic/gin"
)

func NewGinMonitor(collector HTTPCollector) gin.HandlerFunc {
	return func(c *gin.Context) {
		after := collector.MonitorRequest(c.Request)
		c.Next()
		after(c.Writer.Status(), c.Writer.Size())
	}
}
