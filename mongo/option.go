package mongo

import (
	"go.mongodb.org/mongo-driver/event"
	"time"

	"github.com/sandwich-go/xmonitor"
)

//go:generate optiongen --new_func=NewCollectorConf --xconf=true --empty_composite_nil=true --usage_tag_name=usage
func ConfOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// annotation@Collector(comment="metrics收集器，不能为空")
		"Collector": xmonitor.StoreCollector(&noopCollector{}),
		// annotation@SlowLogThreshold(comment="慢日志阈值")
		"SlowLogThreshold": time.Duration(0),
		// annotation@OnSlowCommand(comment="慢日志回调")
		"OnSlowCommand": OnSlowCommand(nil),
	}
}

type OnSlowCommand func(event event.CommandFinishedEvent)

type noopCollector struct{}

func (n noopCollector) RequestCost(tableName, commandName string, cost time.Duration) {}

func (n noopCollector) RequestCount(tableName, commandName, status string) {}

var _ xmonitor.StoreCollector = (*noopCollector)(nil)
