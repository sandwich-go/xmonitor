package ginmid

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// SkipFunc if skip monitor return true else return false
type SkipFunc func(*gin.Context) bool

type Collector interface {
	MonitorRequest(method, path string, reqSize int) func(statusCode, respSize int)
	MonitorLogic(uri string) func(status string)
}

func NewMonitorMid(skipper SkipFunc, collector Collector) gin.HandlerFunc {
	return func(c *gin.Context) {
		if skipper != nil && skipper(c) {
			c.Next()
			return
		}
		path := c.FullPath()
		method := strings.ToLower(c.Request.Method)
		after := collector.MonitorRequest(method, path, CalcRequestSize(c.Request))
		c.Next()
		after(c.Writer.Status(), c.Writer.Size())
	}
}

// From https://github.com/zsais/go-gin-prometheus/blob/2199a42d96c1d40f249909ed2f27d42449c7fc94/middleware.go#L397
func CalcRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.String())
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}
