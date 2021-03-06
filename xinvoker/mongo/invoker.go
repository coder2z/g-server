/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/23 13:52
 **/
package xmongo

import (
	"fmt"
	"github.com/coder2z/component/xinvoker"
	"sync"
)

var mongoI *mongoInvoker

func Register(k string) xinvoker.Invoker {
	mongoI = &mongoInvoker{key: k}
	return mongoI
}

func Invoker(key string) MongoImp {
	if val, ok := mongoI.instances.Load(key); ok {
		return val.(MongoImp)
	}
	panic(fmt.Sprintf("no xlog(%s) invoker found", key))
}

type mongoInvoker struct {
	xinvoker.Base
	instances sync.Map
	key       string
}

func (i *mongoInvoker) Init(opts ...xinvoker.Option) error {
	i.instances = sync.Map{}
	for name, cfg := range i.loadConfig() {
		log := i.new(cfg)
		i.instances.Store(name, log)
	}
	return nil
}

func (i *mongoInvoker) Reload(opts ...xinvoker.Option) error {
	for name, cfg := range i.loadConfig() {
		log := i.new(cfg)
		i.instances.Store(name, log)
	}
	return nil
}

func (i *mongoInvoker) Close(opts ...xinvoker.Option) error {
	return nil
}
