/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/23 13:52
 **/
package oss

import (
	"fmt"
	invoker "github.com/myxy99/component"
	"github.com/myxy99/component/oss/standard"
	"sync"
)

var ossI *ossInvoker

func Register(k string) invoker.Invoker {
	ossI = &ossInvoker{key: k}
	return ossI
}

func Invoker(key string) standard.Oss {
	if val, ok := ossI.instances.Load(key); ok {
		return val.(standard.Oss)
	}
	panic(fmt.Sprintf("no oss(%s) invoker found", key))
}

type ossInvoker struct {
	invoker.Base
	instances sync.Map
	key       string
}

func (i *ossInvoker) Init(opts ...invoker.Option) error {
	i.instances = sync.Map{}
	for name, cfg := range i.loadConfig() {
		i.instances.Store(name, i.new(cfg))
	}
	return nil
}

func (i *ossInvoker) Reload(opts ...invoker.Option) error {
	for name, cfg := range i.loadConfig() {
		i.instances.Store(name, i.new(cfg))
	}
	return nil
}

func (i *ossInvoker) Close(opts ...invoker.Option) error {
	return nil
}
