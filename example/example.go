package main

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sandwich-go/logbus/monitor"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/sandwich-go/xmonitor/fibermid"
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

	// new collector
	collector := xmonitor.NewCollector(
		xmonitor.WithConstLabels(map[string]string{
			"some": "globalValue",
		}),
		xmonitor.WithMonitorRegister(func(c prometheus.Collector) {
			monitor.RegisterCollector(c)
		}),
	)

	wg.Add(1)
	go startGinExample(collector)
	wg.Add(1)
	go startFiberExample(collector)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")
	close(stop)
	wg.Wait()

}

func startFiberExample(collector xmonitor.Collector) {
	app := fiber.New()
	go func() {
		defer wg.Done()
		<-stop
		_ = app.Shutdown()
		log.Println("fiber server exiting")
	}()

	skipper := func(ctx *fiber.Ctx) bool {
		if ctx.Path() == "/health" {
			return true
		}
		return false
	}
	app.Use(
		fibermid.NewMonitorMid(skipper, collector),
	)
	app.Get("/fiber/:id", func(ctx *fiber.Ctx) error {
		log.Println("handler", ctx.Route().Path, ctx.Path())
		return ctx.SendString("Hello fiber!")
	})

	if err := app.Listen(":29080"); err != nil {
		log.Panic(err)
	}
}

func startGinExample(collector xmonitor.Collector) {
	r := gin.New()

	srv := &http.Server{
		Addr:    ":29090",
		Handler: r,
	}

	skipper := func(ctx *gin.Context) bool {
		if ctx.Request.URL.Path == "/health" {
			return true
		}
		return false
	}
	r.Use(ginmid.NewMonitorMid(skipper, collector))

	r.GET("/gin/:id", func(c *gin.Context) {
		c.JSON(200, "Hello gin!")
	})
	go func() {
		defer wg.Done()
		<-stop
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("gin server Shutdown:", err)
		}
		log.Println("gin server exiting")
	}()
	// 服务连接
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}
}
