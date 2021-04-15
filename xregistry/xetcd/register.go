package xetcd

import (
	"fmt"
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/g-server/xregistry"
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
	xlog.Info("Application Starting",
		xlog.FieldComponentName("XRegistry"),
		xlog.FieldMethod("XRegistry.XEtcd.RegisterBuilder"),
		xlog.FieldDescription("Service registration discovery initialization"),
		xlog.FieldAddr(fmt.Sprintf("%v", conf.Endpoints)),
	)
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
		xlog.Warn("Application Starting",
			xlog.FieldComponentName("XRegistry"),
			xlog.FieldMethod("XRegistry.XEtcd.Build"),
			xlog.FieldDescription("Server discover not resolve success in one minute"),
			xlog.Any("target", target),
		)
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
