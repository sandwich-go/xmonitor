package xmonitor

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type HTTPCollector interface {
	MonitorRequest(request *http.Request) func(statusCode, size int)
}

func NewHttpCollector(opts ...CollectorConfOption) HTTPCollector {
	conf := NewCollectorConf(opts...)
	if conf.MonitorRegister == nil {
		panic("must set MonitorRegister")
	}
	initMetrics(conf)
	return &httpCollector{
		conf: conf,
	}
}

const (
	tagMethod = "method"
	tagPath   = "path"
	tagStatus = "http_status"
)

var (
	latencyHistogram *prometheus.HistogramVec
	inFlowCounter    *prometheus.CounterVec
	outFlowCounter   *prometheus.CounterVec
	requestCounter   *prometheus.CounterVec
	initOnce         sync.Once
)

func initMetrics(cc *CollectorConf) {
	initOnce.Do(func() {
		latencyHistogram = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "http_server_latency",
				Help:        "The HTTP request latencies in seconds.",
				Buckets:     cc.Buckets,
				ConstLabels: cc.ConstLabels,
			}, []string{tagMethod, tagPath})

		inFlowCounter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_server_request_size_bytes",
				Help:        "The HTTP request sizes in bytes.",
				ConstLabels: cc.ConstLabels,
			}, []string{tagMethod, tagPath})

		outFlowCounter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_server_response_size_byte",
				Help:        "The HTTP response sizes in bytes.",
				ConstLabels: cc.ConstLabels,
			}, []string{tagMethod, tagPath})

		requestCounter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_server_request_count",
				Help:        "Total number of HTTP requests made.",
				ConstLabels: cc.ConstLabels,
			}, []string{tagMethod, tagPath, tagStatus})

		cc.MonitorRegister(latencyHistogram)
		cc.MonitorRegister(inFlowCounter)
		cc.MonitorRegister(outFlowCounter)
		cc.MonitorRegister(requestCounter)
	})
}

type httpCollector struct {
	conf *CollectorConf
}

func (h httpCollector) MonitorRequest(request *http.Request) func(statusCode, size int) {
	if h.conf.Skip != nil && h.conf.Skip(request) {
		return func(_, _ int) {}
	}
	path := request.URL.EscapedPath()
	method := strings.ToLower(request.Method)

	inFlowCounter.WithLabelValues(method, path).Add(float64(calcRequestSize(request)))
	start := time.Now()

	return func(statusCode, size int) {
		//size must >= 0
		if size < 0 {
			size = 0
		}
		if h.conf.IgnoreLatency != nil && !h.conf.IgnoreLatency(statusCode) {
			latencyHistogram.WithLabelValues(method, path).Observe(time.Since(start).Seconds())
		}
		requestCounter.WithLabelValues(method, path, strconv.Itoa(statusCode)).Inc()
		outFlowCounter.WithLabelValues(method, path).Add(float64(size))
	}
}
