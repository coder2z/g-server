package xredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/coder2m/component/xcfg"
)

func (i *redisInvoker) newRedisClient(o *options) (c *redis.Client) {
	rdb := redis.NewClient(&redis.Options{
		Network:            o.Network,
		Addr:               o.Addr,
		Username:           o.Username,
		Password:           o.Password,
		DB:                 o.DB,
		MaxRetries:         o.MaxRetries,
		MinRetryBackoff:    o.MinRetryBackoff,
		MaxRetryBackoff:    o.MaxRetryBackoff,
		DialTimeout:        o.DialTimeout,
		ReadTimeout:        o.ReadTimeout,
		WriteTimeout:       o.WriteTimeout,
		PoolSize:           o.PoolSize,
		MinIdleConns:       o.MinIdleConns,
		MaxConnAge:         o.MaxConnAge,
		PoolTimeout:        o.PoolTimeout,
		IdleTimeout:        o.IdleTimeout,
		IdleCheckFrequency: o.IdleCheckFrequency,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	return rdb
}

func (i *redisInvoker) loadConfig() map[string]*options {
	conf := make(map[string]*options)
	prefix := i.key
	for name := range xcfg.GetStringMap(prefix) {
		cfg := xcfg.UnmarshalWithExpect(prefix+"."+name, newRedisOptions()).(*options)
		conf[name] = cfg
	}
	return conf
}
