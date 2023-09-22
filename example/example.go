package main

import (
	"context"
	"errors"
	"github.com/sandwich-go/logbus/monitor"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/sandwich-go/xmonitor/ginmid"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sandwich-go/logbus"

	"github.com/sandwich-go/xmonitor"
)

var (
	metricsPath = "/metrics"
)

var (
	stop = make(chan struct{})
	wg   sync.WaitGroup
)

func main() {
	// close logger before exit
	defer logbus.Close()

	// 主线程中使用 非线程安全
	logbus.Init(logbus.NewConf(logbus.WithMonitorOutput(logbus.Prometheus)))

	wg.Add(1)
	go startGinExample()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")
	close(stop)
	wg.Wait()

}

func startGinExample() {
	r := gin.New()

	srv := &http.Server{
		Addr:    ":29090",
		Handler: r,
	}

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
	skipper := func(ctx *gin.Context) bool {
		if ctx.Request.URL.Path == "/health" {
			return true
		}
		return false
	}
	r.Use(ginmid.NewMonitorMid(skipper, collector))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "Hello world!")
	})
	go func() {
		defer wg.Done()
		select {
		case <-stop:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				log.Fatal("Server Shutdown:", err)
			}
			log.Println("Server exiting")
		}
	}()
	// 服务连接
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}
}
