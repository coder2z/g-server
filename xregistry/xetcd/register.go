/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2021/1/1 22:34
 */
package xetcd

import (
	"fmt"
	"github.com/coder2m/g-saber/xconsole"
	"github.com/coder2m/component/xlog"
	"github.com/coder2m/component/xregistry"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/resolver"
	"time"
)

// grpc.Dial("etcd://namespaces/servicename")
type etcdBuilder struct {
	discovery xregistry.Discovery
}

type (
	EtcdV3Cfg = clientv3.Config
)

func RegisterBuilder(conf EtcdV3Cfg) error {
	d, err := newDiscovery(conf)
	if err != nil {
		return err
	}
	b := &etcdBuilder{
		discovery: d,
	}
	resolver.Register(b)
	xconsole.Greenf("Service registration discovery init:", fmt.Sprintf("etcd:%v", conf.Endpoints))
	return nil
}

func (b *etcdBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	ch, err := b.discovery.Discover(xregistry.KeyFormat(target))
	if err != nil {
		return nil, err
	}

	select {
	case x := <-ch:
		xregistry.UpdateAddress(x, cc)
	case <-time.After(time.Minute):
		xlog.Warn("not resolve succuss in one minute", xlog.Any("target", target))
	}
	go func() {
		for i := range ch {
			xregistry.UpdateAddress(i, cc)
		}
	}()
	return &xregistry.NoopResolver{}, nil
}

func (b *etcdBuilder) Scheme() string {
	return "etcd"
}
