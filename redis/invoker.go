/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/23 13:52
 **/
package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	invoker "github.com/myxy99/component"
	"sync"
)

var redisI *redisInvoker

func Register(k string) invoker.Invoker {
	redisI = &redisInvoker{key: k}
	return redisI
}

func Invoker(key string) *redis.Client {
	if val, ok := redisI.instances.Load(key); ok {
		return val.(*redis.Client)
	}
	panic(fmt.Sprintf("no redis(%s) invoker found", key))
}

type redisInvoker struct {
	invoker.Base
	instances sync.Map
	key       string
}

func (i *redisInvoker) Init(opts ...invoker.Option) error {
	i.instances = sync.Map{}
	for name, cfg := range i.loadConfig() {
		i.instances.Store(name, i.newRedisClient(cfg))
	}
	return nil
}

func (i *redisInvoker) Reload(opts ...invoker.Option) error {
	for name, cfg := range i.loadConfig() {
		i.instances.Store(name, i.newRedisClient(cfg))
	}
	return nil
}

func (i *redisInvoker) Close(opts ...invoker.Option) error {
	i.instances.Range(func(key, value interface{}) bool {
		_ = value.(*redis.Client).Close()
		i.instances.Delete(key)
		return true
	})
	return nil
}
