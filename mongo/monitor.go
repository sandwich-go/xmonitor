package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/event"
)

func NewCommandMonitor(opts ...ConfOption) *event.CommandMonitor {
	cc := NewCollectorConf(opts...)
	if cc.Collector == nil {
		panic("new mongo monitor fail, must set collector")
	}
	m := &commandMonitor{
		cc: cc,
	}
	return &event.CommandMonitor{
		Started:   m.Started,
		Succeeded: m.Succeeded,
		Failed:    m.Failed,
	}
}

type commandMonitor struct {
	cc *Conf
}

func (c *commandMonitor) finished(_ context.Context, e event.CommandFinishedEvent, status string) {
	if c.cc.Skipper != nil && c.cc.Skipper(e) {
		return
	}
	c.cc.Collector.RequestCost(e.DatabaseName, e.CommandName, e.Duration)
	c.cc.Collector.RequestCount(e.DatabaseName, e.CommandName, status)
	if c.cc.SlowLogThreshold > 0 && e.Duration >= c.cc.SlowLogThreshold {
		if c.cc.OnSlowCommand != nil {
			c.cc.OnSlowCommand(e)
		}
	}
}

func (c *commandMonitor) Started(context.Context, *event.CommandStartedEvent) {}

func (c *commandMonitor) Succeeded(ctx context.Context, e *event.CommandSucceededEvent) {
	c.finished(ctx, e.CommandFinishedEvent, "Succeeded")
}
func (c *commandMonitor) Failed(ctx context.Context, e *event.CommandFailedEvent) {
	c.finished(ctx, e.CommandFinishedEvent, e.Failure)
}
