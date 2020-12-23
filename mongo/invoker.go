/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/23 13:52
 **/
package mongo

import (
	"fmt"
	invoker "github.com/myxy99/component"
	"sync"
)

var mongoI *mongoInvoker

func Register(k string) invoker.Invoker {
	mongoI = &mongoInvoker{key: k}
	return mongoI
}

func Invoker(key string) MongoImp {
	if val, ok := mongoI.instances.Load(key); ok {
		return val.(MongoImp)
	}
	panic(fmt.Sprintf("no log(%s) invoker found", key))
}

type mongoInvoker struct {
	invoker.Base
	instances sync.Map
	key       string
}

func (i *mongoInvoker) Init(opts ...invoker.Option) error {
	i.instances = sync.Map{}
	for name, cfg := range i.loadConfig() {
		log := i.new(cfg)
		i.instances.Store(name, log)
	}
	return nil
}

func (i *mongoInvoker) Reload(opts ...invoker.Option) error {
	for name, cfg := range i.loadConfig() {
		log := i.new(cfg)
		i.instances.Store(name, log)
	}
	return nil
}

func (i *mongoInvoker) Close(opts ...invoker.Option) error {
	return nil
}
