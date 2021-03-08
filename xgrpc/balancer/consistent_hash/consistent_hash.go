/**
 * @Author: yangon
 * @Description
 * @Date: 2021/3/8 14:39
 **/
package consistent_hash

import (
	"fmt"
	"github.com/coder2z/component/xgrpc"
	"github.com/coder2z/component/xgrpc/balancer"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

const ConsistentHash = "consistent_hash_x"

var defaultConsistentHashKey = "consistent-hash"

func init() {
	balancer.Register(newConsistentHashBuilder(defaultConsistentHashKey))
}

// new Consistent HashBuilder creates a new Consistent Hash balancer builder.
func newConsistentHashBuilder(consistentHashKey string) balancer.Builder {
	return base.NewBalancerBuilderV2(ConsistentHash, &consistentHashPickerBuilder{consistentHashKey}, base.Config{HealthCheck: true})
}

type consistentHashPickerBuilder struct {
	consistentHashKey string
}

func (b *consistentHashPickerBuilder) Build(buildInfo base.PickerBuildInfo) balancer.V2Picker {
	if len(buildInfo.ReadySCs) == 0 {
		return base.NewErrPickerV2(balancer.ErrNoSubConnAvailable)
	}

	picker := &consistentHashPicker{
		subConns:          make(map[string]balancer.SubConn),
		hash:              newKetama(10, nil),
		consistentHashKey: b.consistentHashKey,
	}

	for sc, conInfo := range buildInfo.ReadySCs {
		weight := xbalancer.GetWeight(conInfo.Address)
		for i := 0; i < weight; i++ {
			node := wrapAddr(conInfo.Address.Addr, i)
			picker.hash.Add(node)
			picker.subConns[node] = sc
		}
	}
	return picker
}

type consistentHashPicker struct {
	subConns          map[string]balancer.SubConn
	hash              *ketama
	consistentHashKey string
}

func (p *consistentHashPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	var ret balancer.PickResult
	targetAddr, ok := p.hash.Get(
		defaultConsistentHashKey +
			info.FullMethodName +
			`/` +
			xgrpc.ExtractFromCtx(info.Ctx, "ip"),
	)
	if ok {
		ret.SubConn = p.subConns[targetAddr]
	}
	return ret, nil
}

func wrapAddr(addr string, idx int) string {
	return fmt.Sprintf("%s/%d", addr, idx)
}
