/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/23 13:52
 **/
package xsms

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/myxy99/component/xinvoker"
	"sync"
)

var redisI *smsInvoker

func Register(k string) xinvoker.Invoker {
	redisI = &smsInvoker{key: k}
	return redisI
}

func Invoker(key string) *redis.Client {
	if val, ok := redisI.instances.Load(key); ok {
		return val.(*redis.Client)
	}
	panic(fmt.Sprintf("no redis(%s) invoker found", key))
}

type smsInvoker struct {
	xinvoker.Base
	instances sync.Map
	key       string
}

func (i *smsInvoker) Init(opts ...xinvoker.Option) error {
	i.instances = sync.Map{}
	for name, cfg := range i.loadConfig() {
		i.instances.Store(name, i.newSMSClient(cfg))
	}
	return nil
}

func (i *smsInvoker) Reload(opts ...xinvoker.Option) error {
	for name, cfg := range i.loadConfig() {
		i.instances.Store(name, i.newSMSClient(cfg))
	}
	return nil
}

func (i *smsInvoker) Close(opts ...xinvoker.Option) error {
	i.instances.Range(func(key, value interface{}) bool {
		_ = value.(*redis.Client).Close()
		i.instances.Delete(key)
		return true
	})
	return nil
}
