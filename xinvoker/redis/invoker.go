package xredis

import (
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/g-server/xinvoker"
	"github.com/go-redis/redis/v8"
	"sync"
)

var redisI *redisInvoker

func Register(k string) xinvoker.Invoker {
	redisI = &redisInvoker{key: k}
	return redisI
}

func Invoker(key string) *redis.Client {
	if val, ok := redisI.instances.Load(key); ok {
		return val.(*redis.Client)
	}
	xlog.Panicf("no redis(%s) invoker found", key)
	return nil
}

type redisInvoker struct {
	xinvoker.Base
	instances sync.Map
	key       string
}

func (i *redisInvoker) Init(opts ...xinvoker.Option) error {
	i.instances = sync.Map{}
	for name, cfg := range i.loadConfig() {
		i.instances.Store(name, i.newRedisClient(cfg))
	}
	return nil
}

func (i *redisInvoker) Reload(opts ...xinvoker.Option) error {
	for name, cfg := range i.loadConfig() {
		i.instances.Store(name, i.newRedisClient(cfg))
	}
	return nil
}

func (i *redisInvoker) Close(opts ...xinvoker.Option) error {
	i.instances.Range(func(key, value interface{}) bool {
		_ = value.(*redis.Client).Close()
		i.instances.Delete(key)
		return true
	})
	return nil
}
