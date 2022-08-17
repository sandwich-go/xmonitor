package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sandwich-go/xmonitor"
)

var (
	metricsPath = "/metrics"
)

func main() {
	r := gin.New()
	// new collector
	collector := xmonitor.NewHttpCollector(
		xmonitor.WithConstLabels(map[string]string{
			"project":  "wcc",
			"env_name": "prod",
			"service":  "game",
		}),
		xmonitor.WithSkip(func(r *http.Request) bool {
			return r.URL.EscapedPath() == metricsPath
		}),
		xmonitor.WithMonitorRegister(func(c prometheus.Collector) {
			prometheus.MustRegister(c)
		}),
	)
	r.Use(xmonitor.NewGinMonitor(collector))

	r.GET(metricsPath, prometheusHandler())
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "Hello world!")
	})

	_ = r.Run(":29090")
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
