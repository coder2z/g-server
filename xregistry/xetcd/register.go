/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2021/1/1 22:34
 */
package xetcd

import (
	"github.com/myxy99/component/xregistry"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

// grpc.Dial("etcd://default/servicename")
type etcdBuilder struct {
	discovery xregistry.Discovery
}

func RegisterBuilder(conf clientv3.Config) error {
	d, err := NewDiscovery(conf)
	if err != nil {
		return err
	}
	b := &etcdBuilder{
		discovery: d,
	}

	resolver.Register(b)
	return nil
}

func (b *etcdBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	ch, err := b.discovery.Discover(xregistry.KeyFormat(target))
	if err != nil {
		return nil, err
	}

	select {
	case inss := <-ch:
		xregistry.UpdateAddress(inss, cc)
	case <-time.After(time.Minute):
		log.Printf("not resolve succuss in one minute, target:%v", target)
	}
	go func() {
		for inss := range ch {
			xregistry.UpdateAddress(inss, cc)
		}
	}()
	return &xregistry.NoopResolver{}, nil
}

func (b *etcdBuilder) Scheme() string {
	return "etcd"
}
