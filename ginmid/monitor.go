package ginmid

import (
	"github.com/gin-gonic/gin"
	"strings"

	"github.com/sandwich-go/xmonitor"
)

// SkipFunc if skip monitor return true else return false
type SkipFunc func(*gin.Context) bool

func NewMonitorMid(skipper SkipFunc, collector xmonitor.Collector) gin.HandlerFunc {
	return func(c *gin.Context) {
		if skipper != nil && skipper(c) {
			c.Next()
			return
		}
		path := c.FullPath()
		method := strings.ToLower(c.Request.Method)
		after := collector.MonitorRequest(method, path, xmonitor.CalcRequestSize(c.Request))
		c.Next()
		after(c.Writer.Status(), c.Writer.Size())
	}
}
