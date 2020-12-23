/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/23 13:52
 **/
package log

import (
	"fmt"
	invoker "github.com/myxy99/component"
	"go.uber.org/zap"
	"sync"
)

var log *logInvoker

func Register(k string) invoker.Invoker {
	log = &logInvoker{key: k}
	return log
}

func Invoker(key string) *zap.SugaredLogger {
	if val, ok := log.instances.Load(key); ok {
		return val.(*zap.SugaredLogger)
	}
	panic(fmt.Sprintf("no log(%s) invoker found", key))
}

type logInvoker struct {
	invoker.Base
	instances sync.Map
	key       string
}

func (i *logInvoker) Init(opts ...invoker.Option) error {
	i.instances = sync.Map{}
	for name, cfg := range i.loadConfig() {
		log := i.new(cfg)
		i.instances.Store(name, log)
	}
	return nil
}

func (i *logInvoker) Reload(opts ...invoker.Option) error {
	for name, cfg := range i.loadConfig() {
		log := i.new(cfg)
		i.instances.Store(name, log)
	}
	return nil
}

func (i *logInvoker) Close(opts ...invoker.Option) error {
	return nil
}
