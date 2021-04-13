package xinvoker

import (
	"errors"
	"github.com/coder2z/g-saber/xcfg"
	"github.com/coder2z/g-saber/xconsole"
	"reflect"
)

func init() {
	xcfg.OnChange(func(*xcfg.Configuration) {
		_ = Reload()
	})
}

type Option func()

type Invoker interface {
	Init(opts ...Option) error
	Reload(opts ...Option) error
	Close(opts ...Option) error
}

type Base struct{}

func (*Base) Init(opts ...Option) error {
	return errors.New("init not implemented")
}

func (*Base) Reload(opts ...Option) error {
	return errors.New("reload not implemented")
}

func (*Base) Close(opts ...Option) error {
	return errors.New("close not implemented")
}

var invokers []Invoker

// Register invoker注册具体实现
func Register(ivk ...Invoker) {
	invokers = append(invokers, ivk...)
}

// Init invoker执行初始化具体实现
func Init(opts ...Option) error {
	for _, invoker := range invokers {
		key := reflect.ValueOf(invoker).Elem().FieldByName("key").String()
		xconsole.Bluef("invoker running start init:", key)
		_ = invoker.Init(opts...)
	}

	return nil
}

// Reload invoker执行热更新具体实现
func Reload(opts ...Option) error {
	for _, invoker := range invokers {
		key := reflect.ValueOf(invoker).Elem().FieldByName("key").String()
		xconsole.Redf("invoker running start Reload:", key)
		_ = invoker.Reload(opts...)
	}

	return nil
}

// Close invoker执行退出具体实现
func Close(opts ...Option) error {
	for i := len(invokers) - 1; i >= 0; i-- {
		key := reflect.ValueOf(invokers[i]).Elem().FieldByName("key").String()
		xconsole.Redf("invoker running start close:", key)
		_ = invokers[i].Close(opts...)
	}

	return nil
}
