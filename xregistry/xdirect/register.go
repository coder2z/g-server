package xdirect

import (
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/g-server/xregistry"
	"google.golang.org/grpc/resolver"
	"time"
)

func RegisterBuilder() error {
	resolver.Register(&directBuilder{
		discovery: newDiscovery(),
	})
	xlog.Info("Application Starting",
		xlog.FieldComponentName("XRegistry"),
		xlog.FieldMethod("XRegistry.XDiscovery.init"),
		xlog.FieldDescription("Service use discovery registration discovery initialization"),
	)
	return nil
}

// 通过 ip:port 设置需要在前面加上//
// grpc.Dial("direct://namespaces///127.0.0.1:8000?w=1,//127.0.0.1:8001?w=1")
type directBuilder struct {
	discovery xregistry.Discovery
}

func (b *directBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	ch, err := b.discovery.Discover(target.Endpoint)
	if err != nil {
		return nil, err
	}

	select {
	case i := <-ch:
		xregistry.UpdateAddress(i, cc)
	case <-time.After(time.Minute):
		xlog.Warn("Application Starting",
			xlog.FieldComponentName("XRegistry"),
			xlog.FieldMethod("XRegistry.XDirect.Build"),
			xlog.FieldDescription("Server discover not resolve success in one minute"),
			xlog.Any("target", target),
		)
	}
	return &xregistry.NoopResolver{}, nil
}

func (b *directBuilder) Scheme() string {
	return "direct"
}
