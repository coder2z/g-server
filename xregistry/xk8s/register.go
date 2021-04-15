package xk8s

import (
	"github.com/coder2z/g-saber/xconsole"
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/g-server/xregistry"
	"google.golang.org/grpc/resolver"
	"sync"
	"time"
)

// grpc.Dial("k8s://namespace/servicename:portname")
// grpc.Dial("k8s://namespace/servicename:port")
// grpc.Dial("k8s:///servicename:portname") // namespace=default
// grpc.Dial("k8s:///servicename:port")
type k8sBuilder struct {
	sync.Mutex
	rs map[string]xregistry.Discovery
}

func RegisterBuilder() error {
	b := &k8sBuilder{
		rs: map[string]xregistry.Discovery{},
	}
	resolver.Register(b)
	xconsole.Greenf("Service registration discovery init:", "k8s")
	return nil
}

func (b *k8sBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var (
		err        error
		discovery  xregistry.Discovery
		namespaces = target.Authority
	)
	if namespaces == "" {
		namespaces = "default"
	}
	if discovery, err = b.getDiscovery(namespaces); err != nil {
		return nil, err
	}
	ch, err := discovery.Discover(target.Endpoint)
	if err != nil {
		return nil, err
	}

	select {
	case i := <-ch:
		xregistry.UpdateAddress(i, cc)
	case <-time.After(time.Minute):
		xlog.Warn("Application Starting",
			xlog.FieldComponentName("XRegistry"),
			xlog.FieldMethod("XRegistry.XK8S.Build"),
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

func (b *k8sBuilder) getDiscovery(namespace string) (r xregistry.Discovery, err error) {
	b.Lock()
	defer b.Unlock()
	if r = b.rs[namespace]; r != nil {
		return
	}
	if r, err = newDiscovery(namespace); err != nil {
		return
	}
	b.rs[namespace] = r
	return
}

func (b *k8sBuilder) Scheme() string {
	return "k8s"
}
