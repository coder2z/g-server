/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/23 13:52
 **/
package xoss

import (
	"github.com/coder2z/component/xinvoker"
	"github.com/coder2z/component/xinvoker/oss/standard"
	"github.com/coder2z/g-saber/xlog"
	"sync"
)

var ossI *ossInvoker

func Register(k string) xinvoker.Invoker {
	ossI = &ossInvoker{key: k}
	return ossI
}

func Invoker(key string) standard.Oss {
	if val, ok := ossI.instances.Load(key); ok {
		return val.(standard.Oss)
	}
	xlog.Panicf("no oss(%s) invoker found", key)
	return nil
}

type ossInvoker struct {
	xinvoker.Base
	instances sync.Map
	key       string
}

func (i *ossInvoker) Init(opts ...xinvoker.Option) error {
	i.instances = sync.Map{}
	for name, cfg := range i.loadConfig() {
		i.instances.Store(name, i.new(cfg))
	}
	return nil
}

func (i *ossInvoker) Reload(opts ...xinvoker.Option) error {
	for name, cfg := range i.loadConfig() {
		i.instances.Store(name, i.new(cfg))
	}
	return nil
}

func (i *ossInvoker) Close(opts ...xinvoker.Option) error {
	return nil
}
