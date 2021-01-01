/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2021/1/1 22:44
 */
package xdirect

import (
	"github.com/myxy99/component/xlog"
	"github.com/myxy99/component/xregistry"
	"google.golang.org/grpc/resolver"
	"time"
)

func init() {
	resolver.Register(&directBuilder{
		discovery: NewDiscovery(),
	})
}

func RegisterBuilder() error {
	return nil
}

// grpc.Dial("direct://namespaces/127.0.0.1:8000,127.0.0.1:8001")
type directBuilder struct {
	discovery xregistry.Discovery
}

func (b *directBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	ch, err := b.discovery.Discover(xregistry.KeyFormat(target))
	if err != nil {
		return nil, err
	}

	select {
	case i := <-ch:
		xregistry.UpdateAddress(i, cc)
	case <-time.After(time.Minute):
		xlog.Warnw("not resolve succuss in one minute, target:%v", target)
	}
	return &xregistry.NoopResolver{}, nil
}

func (b *directBuilder) Scheme() string {
	return "direct"
}
