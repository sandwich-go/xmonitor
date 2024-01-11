package xmonitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"time"
)

type StoreCollector interface {
	RequestCost(tableName, commandName string, cost time.Duration)
	RequestCount(tableName, commandName, status string)
}

type storeCollector struct {
	conf *CollectorConf
}

func (s storeCollector) RequestCost(tableName, commandName string, cost time.Duration) {
	storeRequestCost.WithLabelValues(tableName, commandName).Observe(cost.Seconds())
}

func (s storeCollector) RequestCount(tableName, commandName, status string) {
	storeRequestCount.WithLabelValues(tableName, commandName, status).Inc()
}

func NewStoreCollector(opts ...CollectorConfOption) StoreCollector {
	conf := NewCollectorConf(opts...)
	if conf.MonitorRegister == nil {
		panic("must set MonitorRegister")
	}
	initStoreMetrics(conf)
	return &storeCollector{
		conf: conf,
	}
}

const (
	tagTable   = "table"
	tagCommand = "command"
)

var (
	storeRequestCost  *prometheus.HistogramVec
	storeRequestCount *prometheus.CounterVec

	initStoreOnce sync.Once
)

func initStoreMetrics(cc *CollectorConf) {
	initStoreOnce.Do(func() {
		storeRequestCost = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "store_client_request_cost",
				Help:        "The store client request cost latencies in seconds.",
				Buckets:     cc.Buckets,
				ConstLabels: cc.ConstLabels,
			}, []string{tagTable, tagCommand})
		storeRequestCount = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "store_client_request_count",
				Help:        "The store client request count.",
				ConstLabels: cc.ConstLabels,
			}, []string{tagTable, tagCommand, tagStatus})

		cc.MonitorRegister(storeRequestCost)
		cc.MonitorRegister(storeRequestCount)
	})
}
