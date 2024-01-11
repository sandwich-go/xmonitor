// Code generated by optiongen. DO NOT EDIT.
// optiongen: github.com/timestee/optiongen

package mongo

import (
	"sync/atomic"
	"time"
	"unsafe"
)

// Conf should use NewCollectorConf to initialize it
type Conf struct {
	// annotation@Collector(comment="metrics收集器，不能为空")
	Collector StoreCollector `xconf:"collector" usage:"metrics收集器，不能为空"`
	// annotation@Skipper(comment="跳过监控")
	Skipper Skipper `xconf:"skipper" usage:"跳过监控"`
	// annotation@SlowLogThreshold(comment="慢日志阈值")
	SlowLogThreshold time.Duration `xconf:"slow_log_threshold" usage:"慢日志阈值"`
	// annotation@OnSlowCommand(comment="慢日志回调")
	OnSlowCommand OnSlowCommand `xconf:"on_slow_command" usage:"慢日志回调"`
}

// NewCollectorConf new Conf
func NewCollectorConf(opts ...ConfOption) *Conf {
	cc := newDefaultConf()
	for _, opt := range opts {
		opt(cc)
	}
	if watchDogConf != nil {
		watchDogConf(cc)
	}
	return cc
}

// ApplyOption apply multiple new option and return the old ones
// sample:
// old := cc.ApplyOption(WithTimeout(time.Second))
// defer cc.ApplyOption(old...)
func (cc *Conf) ApplyOption(opts ...ConfOption) []ConfOption {
	var previous []ConfOption
	for _, opt := range opts {
		previous = append(previous, opt(cc))
	}
	return previous
}

// ConfOption option func
type ConfOption func(cc *Conf) ConfOption

// WithCollector metrics收集器，不能为空
func WithCollector(v StoreCollector) ConfOption {
	return func(cc *Conf) ConfOption {
		previous := cc.Collector
		cc.Collector = v
		return WithCollector(previous)
	}
}

// WithSkipper 跳过监控
func WithSkipper(v Skipper) ConfOption {
	return func(cc *Conf) ConfOption {
		previous := cc.Skipper
		cc.Skipper = v
		return WithSkipper(previous)
	}
}

// WithSlowLogThreshold 慢日志阈值
func WithSlowLogThreshold(v time.Duration) ConfOption {
	return func(cc *Conf) ConfOption {
		previous := cc.SlowLogThreshold
		cc.SlowLogThreshold = v
		return WithSlowLogThreshold(previous)
	}
}

// WithOnSlowCommand 慢日志回调
func WithOnSlowCommand(v OnSlowCommand) ConfOption {
	return func(cc *Conf) ConfOption {
		previous := cc.OnSlowCommand
		cc.OnSlowCommand = v
		return WithOnSlowCommand(previous)
	}
}

// InstallConfWatchDog the installed func will called when NewCollectorConf  called
func InstallConfWatchDog(dog func(cc *Conf)) { watchDogConf = dog }

// watchDogConf global watch dog
var watchDogConf func(cc *Conf)

// newDefaultConf new default Conf
func newDefaultConf() *Conf {
	cc := &Conf{}

	for _, opt := range [...]ConfOption{
		WithCollector(&noopCollector{}),
		WithSkipper(ignorePing),
		WithSlowLogThreshold(0),
		WithOnSlowCommand(nil),
	} {
		opt(cc)
	}

	return cc
}

// AtomicSetFunc used for XConf
func (cc *Conf) AtomicSetFunc() func(interface{}) { return AtomicConfSet }

// atomicConf global *Conf holder
var atomicConf unsafe.Pointer

// onAtomicConfSet global call back when  AtomicConfSet called by XConf.
// use ConfInterface.ApplyOption to modify the updated cc
// if passed in cc not valid, then return false, cc will not set to atomicConf
var onAtomicConfSet func(cc ConfInterface) bool

// InstallCallbackOnAtomicConfSet install callback
func InstallCallbackOnAtomicConfSet(callback func(cc ConfInterface) bool) { onAtomicConfSet = callback }

// AtomicConfSet atomic setter for *Conf
func AtomicConfSet(update interface{}) {
	cc := update.(*Conf)
	if onAtomicConfSet != nil && !onAtomicConfSet(cc) {
		return
	}
	atomic.StorePointer(&atomicConf, (unsafe.Pointer)(cc))
}

// AtomicConf return atomic *ConfVisitor
func AtomicConf() ConfVisitor {
	current := (*Conf)(atomic.LoadPointer(&atomicConf))
	if current == nil {
		defaultOne := newDefaultConf()
		if watchDogConf != nil {
			watchDogConf(defaultOne)
		}
		atomic.CompareAndSwapPointer(&atomicConf, nil, (unsafe.Pointer)(defaultOne))
		return (*Conf)(atomic.LoadPointer(&atomicConf))
	}
	return current
}

// all getter func
func (cc *Conf) GetCollector() StoreCollector       { return cc.Collector }
func (cc *Conf) GetSkipper() Skipper                { return cc.Skipper }
func (cc *Conf) GetSlowLogThreshold() time.Duration { return cc.SlowLogThreshold }
func (cc *Conf) GetOnSlowCommand() OnSlowCommand    { return cc.OnSlowCommand }

// ConfVisitor visitor interface for Conf
type ConfVisitor interface {
	GetCollector() StoreCollector
	GetSkipper() Skipper
	GetSlowLogThreshold() time.Duration
	GetOnSlowCommand() OnSlowCommand
}

// ConfInterface visitor + ApplyOption interface for Conf
type ConfInterface interface {
	ConfVisitor
	ApplyOption(...ConfOption) []ConfOption
}
