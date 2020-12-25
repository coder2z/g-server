/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/23 13:18
 **/
package xinvoker

import (
	"errors"
	"github.com/myxy99/component/xcfg"
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
		_ = invoker.Init(opts...)
	}

	return nil
}

// Reload invoker执行热更新具体实现
func Reload(opts ...Option) error {
	for _, invoker := range invokers {
		_ = invoker.Reload(opts...)
	}

	return nil
}

// Close invoker执行退出具体实现
func Close(opts ...Option) error {
	for i := len(invokers) - 1; i >= 0; i-- {
		_ = invokers[i].Close(opts...)
	}

	return nil
}
