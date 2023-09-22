package main

import (
	"github.com/sandwich-go/logbus/monitor"
	"github.com/sandwich-go/xmonitor/ginmid"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sandwich-go/logbus"

	"github.com/sandwich-go/xmonitor"
)

var (
	metricsPath = "/metrics"
)

func main() {
	// close logger before exit
	defer logbus.Close()

	// 主线程中使用 非线程安全
	logbus.Init(logbus.NewConf(logbus.WithMonitorOutput(logbus.Prometheus))

	r := gin.New()
	// new collector
	collector := xmonitor.NewCollector(
		xmonitor.WithConstLabels(map[string]string{
			"project":  "wcc",
			"env_name": "prod",
			"service":  "game",
		}),
		xmonitor.WithSkip(func(r *http.Request) bool {
			return r.URL.EscapedPath() == metricsPath
		}),
		xmonitor.WithMonitorRegister(func(c prometheus.Collector) {
			monitor.RegisterCollector(c)
		}),
	)
	r.Use(ginmid.NewMonitorMid(collector))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "Hello world!")
	})

	_ = r.Run(":29090")
}
