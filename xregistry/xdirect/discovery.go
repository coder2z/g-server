package xdirect

import (
	"fmt"
	"github.com/coder2z/g-saber/xlog"
	xbalancer "github.com/coder2z/g-server/xgrpc/balancer"
	"github.com/coder2z/g-server/xregistry"
	"google.golang.org/grpc/metadata"
	"net/url"
	"strings"
)

type directDiscovery struct{}

func newDiscovery() xregistry.Discovery {
	return &directDiscovery{}
}

// target 格式： "127.0.0.1:8000?w=1,127.0.0.1:8001?w=2"
func (d *directDiscovery) Discover(target string) (<-chan []xregistry.Instance, error) {
	endpoints := strings.Split(target, ",")
	if len(endpoints) == 0 {
		return nil, fmt.Errorf("no endpoint")
	}
	var i []xregistry.Instance
	for _, addr := range endpoints {
		urlObj, err := url.Parse(addr)
		if err != nil {
			xlog.Panic("Application Starting",
				xlog.FieldComponentName("XRegistry"),
				xlog.FieldMethod("XRegistry.XDirect.Discover"),
				xlog.FieldDescription("endpoints format error"),
				xlog.FieldErr(err),
			)
		}
		ins := xregistry.Instance{Address: urlObj.Host, Metadata: metadata.Pairs(xbalancer.WeightKey, urlObj.Query().Get("w"))}
		i = append(i, ins)
	}

	ch := make(chan []xregistry.Instance)
	go func() {
		ch <- i
	}()
	return ch, nil
}

func (d *directDiscovery) Close() {}
