// Code generated by optiongen. DO NOT EDIT.
// optiongen: github.com/timestee/optiongen

package xmonitor

import (
	"net/http"
	"sync/atomic"
	"unsafe"
)

// CollectorConf should use NewCollectorConf to initialize it
type CollectorConf struct {
	// annotation@ConstLabels(comment="监控默认添加的静态labels")
	ConstLabels map[string]string `xconf:"const_labels" usage:"监控默认添加的静态labels"`
	// annotation@MonitorRegister(xconf="-"，comment="监控注册函数")
	MonitorRegister MonitorRegisterFunc `xconf:"监控注册函数"`
	// annotation@Buckets(comment="histogram buckets 监控耗时桶(单位秒)，参考 https://cloud.tencent.com/developer/article/1495303")
	Buckets       []float64         `xconf:"buckets" usage:"histogram buckets 监控耗时桶(单位秒)，参考 https://cloud.tencent.com/developer/article/1495303"`
	IgnoreLatency IgnoreLatencyFunc `xconf:"ignore_latency"`
}

// NewCollectorConf new CollectorConf
func NewCollectorConf(opts ...CollectorConfOption) *CollectorConf {
	cc := newDefaultCollectorConf()
	for _, opt := range opts {
		opt(cc)
	}
	if watchDogCollectorConf != nil {
		watchDogCollectorConf(cc)
	}
	return cc
}

// ApplyOption apply multiple new option and return the old ones
// sample:
// old := cc.ApplyOption(WithTimeout(time.Second))
// defer cc.ApplyOption(old...)
func (cc *CollectorConf) ApplyOption(opts ...CollectorConfOption) []CollectorConfOption {
	var previous []CollectorConfOption
	for _, opt := range opts {
		previous = append(previous, opt(cc))
	}
	return previous
}

// CollectorConfOption option func
type CollectorConfOption func(cc *CollectorConf) CollectorConfOption

// WithConstLabels 监控默认添加的静态labels
func WithConstLabels(v map[string]string) CollectorConfOption {
	return func(cc *CollectorConf) CollectorConfOption {
		previous := cc.ConstLabels
		cc.ConstLabels = v
		return WithConstLabels(previous)
	}
}

// WithMonitorRegister option func for filed MonitorRegister
func WithMonitorRegister(v MonitorRegisterFunc) CollectorConfOption {
	return func(cc *CollectorConf) CollectorConfOption {
		previous := cc.MonitorRegister
		cc.MonitorRegister = v
		return WithMonitorRegister(previous)
	}
}

// WithBuckets histogram buckets 监控耗时桶(单位秒)，参考 https://cloud.tencent.com/developer/article/1495303
func WithBuckets(v ...float64) CollectorConfOption {
	return func(cc *CollectorConf) CollectorConfOption {
		previous := cc.Buckets
		cc.Buckets = v
		return WithBuckets(previous...)
	}
}

// WithIgnoreLatency option func for filed IgnoreLatency
func WithIgnoreLatency(v IgnoreLatencyFunc) CollectorConfOption {
	return func(cc *CollectorConf) CollectorConfOption {
		previous := cc.IgnoreLatency
		cc.IgnoreLatency = v
		return WithIgnoreLatency(previous)
	}
}

// InstallCollectorConfWatchDog the installed func will called when NewCollectorConf  called
func InstallCollectorConfWatchDog(dog func(cc *CollectorConf)) { watchDogCollectorConf = dog }

// watchDogCollectorConf global watch dog
var watchDogCollectorConf func(cc *CollectorConf)

// newDefaultCollectorConf new default CollectorConf
func newDefaultCollectorConf() *CollectorConf {
	cc := &CollectorConf{}

	for _, opt := range [...]CollectorConfOption{
		WithConstLabels(nil),
		WithMonitorRegister(nil),
		WithBuckets([]float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}...),
		WithIgnoreLatency(func(code int) bool {
			return code != http.StatusOK
		}),
	} {
		opt(cc)
	}

	return cc
}

// AtomicSetFunc used for XConf
func (cc *CollectorConf) AtomicSetFunc() func(interface{}) { return AtomicCollectorConfSet }

// atomicCollectorConf global *CollectorConf holder
var atomicCollectorConf unsafe.Pointer

// onAtomicCollectorConfSet global call back when  AtomicCollectorConfSet called by XConf.
// use CollectorConfInterface.ApplyOption to modify the updated cc
// if passed in cc not valid, then return false, cc will not set to atomicCollectorConf
var onAtomicCollectorConfSet func(cc CollectorConfInterface) bool

// InstallCallbackOnAtomicCollectorConfSet install callback
func InstallCallbackOnAtomicCollectorConfSet(callback func(cc CollectorConfInterface) bool) {
	onAtomicCollectorConfSet = callback
}

// AtomicCollectorConfSet atomic setter for *CollectorConf
func AtomicCollectorConfSet(update interface{}) {
	cc := update.(*CollectorConf)
	if onAtomicCollectorConfSet != nil && !onAtomicCollectorConfSet(cc) {
		return
	}
	atomic.StorePointer(&atomicCollectorConf, (unsafe.Pointer)(cc))
}

// AtomicCollectorConf return atomic *CollectorConfVisitor
func AtomicCollectorConf() CollectorConfVisitor {
	current := (*CollectorConf)(atomic.LoadPointer(&atomicCollectorConf))
	if current == nil {
		defaultOne := newDefaultCollectorConf()
		if watchDogCollectorConf != nil {
			watchDogCollectorConf(defaultOne)
		}
		atomic.CompareAndSwapPointer(&atomicCollectorConf, nil, (unsafe.Pointer)(defaultOne))
		return (*CollectorConf)(atomic.LoadPointer(&atomicCollectorConf))
	}
	return current
}

// all getter func
func (cc *CollectorConf) GetConstLabels() map[string]string       { return cc.ConstLabels }
func (cc *CollectorConf) GetMonitorRegister() MonitorRegisterFunc { return cc.MonitorRegister }
func (cc *CollectorConf) GetBuckets() []float64                   { return cc.Buckets }
func (cc *CollectorConf) GetIgnoreLatency() IgnoreLatencyFunc     { return cc.IgnoreLatency }

// CollectorConfVisitor visitor interface for CollectorConf
type CollectorConfVisitor interface {
	GetConstLabels() map[string]string
	GetMonitorRegister() MonitorRegisterFunc
	GetBuckets() []float64
	GetIgnoreLatency() IgnoreLatencyFunc
}

// CollectorConfInterface visitor + ApplyOption interface for CollectorConf
type CollectorConfInterface interface {
	CollectorConfVisitor
	ApplyOption(...CollectorConfOption) []CollectorConfOption
}
