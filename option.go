package xmonitor

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

//go:generate optiongen --new_func=NewCollectorConf --xconf=true --empty_composite_nil=true --usage_tag_name=usage
func CollectorConfOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// annotation@ConstLabels(comment="监控默认添加的静态labels")
		"ConstLabels": map[string]string{},
		// annotation@MonitorRegister(xconf="-"，comment="监控注册函数")
		"MonitorRegister": MonitorRegisterFunc(nil),
		// annotation@Buckets(comment="histogram buckets 监控耗时桶(单位秒)，参考 https://cloud.tencent.com/developer/article/1495303")
		"Buckets": []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		"IgnoreLatency": IgnoreLatencyFunc(func(code int) bool {
			return code != http.StatusOK
		}),
	}
}

// MonitorRegisterFunc monitor collector register
type MonitorRegisterFunc func(prometheus.Collector)

// IgnoreLatencyFunc ignore latency func
type IgnoreLatencyFunc func(code int) bool
