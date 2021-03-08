/**
 * @Author: yangon
 * @Description
 * @Date: 2021/3/8 14:58
 **/
package random

import (
	"github.com/coder2z/component/xgrpc/balancer"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math/rand"
	"sync"
	"time"
)

const Random = "random_x"

// newRandomBuilder creates a new random balancer builder.
func newRandomBuilder() balancer.Builder {
	return base.NewBalancerBuilderV2(Random, &randomPickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newRandomBuilder())
}

type randomPickerBuilder struct{}

func (*randomPickerBuilder) Build(buildInfo base.PickerBuildInfo) balancer.V2Picker {
	if len(buildInfo.ReadySCs) == 0 {
		return base.NewErrPickerV2(balancer.ErrNoSubConnAvailable)
	}
	var scs []balancer.SubConn

	for subCon, subConnInfo := range buildInfo.ReadySCs {
		weight := xbalancer.GetWeight(subConnInfo.Address)
		for i := 0; i < weight; i++ {
			scs = append(scs, subCon)
		}
	}
	return &randomPicker{
		subConns: scs,
		rand:     rand.New(rand.NewSource(time.Now().Unix())),
	}
}

type randomPicker struct {
	subConns []balancer.SubConn
	mu       sync.Mutex
	rand     *rand.Rand
}

func (p *randomPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	ret := balancer.PickResult{}
	p.mu.Lock()
	ret.SubConn = p.subConns[p.rand.Intn(len(p.subConns))]
	p.mu.Unlock()
	return ret, nil
}
