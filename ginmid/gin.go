package ginmid

import (
	"github.com/gin-gonic/gin"

	"github.com/sandwich-go/xmonitor"
)

func NewMonitorMid(collector xmonitor.Collector) gin.HandlerFunc {
	return func(c *gin.Context) {
		after := collector.MonitorRequest(c.Request)
		c.Next()
		after(c.Writer.Status(), c.Writer.Size())
	}
}
