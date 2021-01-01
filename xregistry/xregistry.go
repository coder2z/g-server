/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2021/1/1 22:04
 */
package xregistry

import (
	"fmt"
	"google.golang.org/grpc/resolver"
	"time"
)

// Registry Options
type Options struct {
	ServiceName string            `json:"servicename"`
	Namespaces  string            `json:"namespaces"`
	Address     string            `json:"address"`
	Metadata    map[string]string `json:"metadata"`
	// 服务有效时长
	RegisterTTL      time.Duration `json:"-"` // time to live, 服务失活一段时间后自动从注册中心删除
	RegisterInterval time.Duration `json:"-"` // 注册间隔时长，也可不要 默认为RegisterTTL/3
}

type Option func(*Options)

type Registry interface {
	Register(ops ...Option)
	Close()
}

func ServiceName(name string) Option {
	return func(o *Options) {
		o.ServiceName = name
	}
}

func ServiceNamespaces(namespaces string) Option {
	return func(o *Options) {
		o.Namespaces = namespaces
	}
}

func Address(address string) Option {
	return func(o *Options) {
		o.Address = address
	}
}

func Metadata(m map[string]string) Option {
	return func(o *Options) {
		o.Metadata = m
	}
}

func RegisterTTL(ttl time.Duration) Option {
	return func(o *Options) {
		o.RegisterTTL = ttl
	}
}

func RegisterInterval(interval time.Duration) Option {
	return func(o *Options) {
		o.RegisterInterval = interval
	}
}

type Instance struct {
	ServiceName string            `json:"servicename"`
	Address     string            `json:"address"`
	Metadata    map[string]string `json:"metadata"`
}

// 服务发现接口
// target的具体格式由其实现决定
type Discovery interface {
	Discover(target string) (<-chan []Instance, error)
	Close()
}

func UpdateAddress(inss []Instance, conn resolver.ClientConn) {
	var address []resolver.Address
	for _, ins := range inss {
		address = append(address, resolver.Address{Addr: ins.Address})
	}
	conn.UpdateState(resolver.State{
		Addresses: address,
	})
}

func KeyFormat(target resolver.Target) string {
	return fmt.Sprintf("%v.%v", target.Authority, target.Endpoint)
}

type NoopResolver struct{}

func (r *NoopResolver) ResolveNow(resolver.ResolveNowOptions) {}
func (r *NoopResolver) Close()                                {}
