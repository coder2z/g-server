package p2c

import (
	"github.com/coder2z/g-saber/xtime"
	xbalancer "github.com/coder2z/g-server/xgrpc/balancer"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	P2C             = "p2c_x"
	decayTime       = int64(time.Second * 10)
	forcePick       = int64(time.Second)
	initSuccess     = 1000
	throttleSuccess = initSuccess / 2
	penalty         = int64(math.MaxInt32)
	pickTimes       = 3
)

func init() {
	balancer.Register(newBuilder())
}

type p2cPickerBuilder struct{}

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilderV2(P2C, new(p2cPickerBuilder), base.Config{HealthCheck: true})
}

func (b *p2cPickerBuilder) Build(buildInfo base.PickerBuildInfo) balancer.V2Picker {
	if len(buildInfo.ReadySCs) == 0 {
		return base.NewErrPickerV2(balancer.ErrNoSubConnAvailable)
	}

	var conn []*subConn
	for subCon, subConnInfo := range buildInfo.ReadySCs {
		weight := xbalancer.GetWeight(subConnInfo.Address)
		for i := 0; i < weight; i++ {
			conn = append(conn, &subConn{
				addr:    subConnInfo.Address,
				conn:    subCon,
				success: initSuccess,
			})
		}
	}

	return &p2cPicker{
		conn: conn,
		r:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type p2cPicker struct {
	conn []*subConn
	r    *rand.Rand
	lock sync.Mutex
}

func (p *p2cPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	var (
		chosen *subConn
		ret    balancer.PickResult
	)
	switch len(p.conn) {
	case 0:
		return ret, balancer.ErrNoSubConnAvailable
	case 1:
		chosen = p.choose(p.conn[0], nil)
	case 2:
		chosen = p.choose(p.conn[0], p.conn[1])
	default:
		var node1, node2 *subConn
		for i := 0; i < pickTimes; i++ {
			a := p.r.Intn(len(p.conn))
			b := p.r.Intn(len(p.conn) - 1)
			if b >= a {
				b++
			}
			node1 = p.conn[a]
			node2 = p.conn[b]
			if node1.healthy() && node2.healthy() {
				break
			}
		}

		chosen = p.choose(node1, node2)
	}

	atomic.AddInt64(&chosen.inflight, 1)
	atomic.AddInt64(&chosen.requests, 1)
	ret.SubConn = chosen.conn
	ret.Done = p.buildDoneFunc(chosen)
	return ret, nil
}

func (p *p2cPicker) buildDoneFunc(c *subConn) func(info balancer.DoneInfo) {
	start := xtime.Now().Unix()
	return func(info balancer.DoneInfo) {
		atomic.AddInt64(&c.inflight, -1)
		now := xtime.Now().Unix()
		last := atomic.SwapInt64(&c.last, now)
		td := now - last
		if td < 0 {
			td = 0
		}
		w := math.Exp(float64(-td) / float64(decayTime))
		lag := now - start
		if lag < 0 {
			lag = 0
		}
		olag := atomic.LoadUint64(&c.lag)
		if olag == 0 {
			w = 0
		}
		atomic.StoreUint64(&c.lag, uint64(float64(olag)*w+float64(lag)*(1-w)))
		success := initSuccess
		if info.Err != nil && !acceptable(info.Err) {
			success = 0
		}
		s := atomic.LoadUint64(&c.success)
		atomic.StoreUint64(&c.success, uint64(float64(s)*w+float64(success)*(1-w)))
	}
}

func (p *p2cPicker) choose(c1, c2 *subConn) *subConn {
	start := xtime.Now().Unix()
	if c2 == nil {
		atomic.StoreInt64(&c1.pick, start)
		return c1
	}

	if c1.load() > c2.load() {
		c1, c2 = c2, c1
	}

	pick := atomic.LoadInt64(&c2.pick)
	if start-pick > forcePick && atomic.CompareAndSwapInt64(&c2.pick, pick, start) {
		return c2
	}

	atomic.StoreInt64(&c1.pick, start)
	return c1
}

type subConn struct {
	addr     resolver.Address
	conn     balancer.SubConn
	lag      uint64
	inflight int64
	success  uint64
	requests int64
	last     int64
	pick     int64
}

func (c *subConn) healthy() bool {
	return atomic.LoadUint64(&c.success) > throttleSuccess
}

func (c *subConn) load() int64 {
	lag := int64(math.Sqrt(float64(atomic.LoadUint64(&c.lag) + 1)))
	load := lag * (atomic.LoadInt64(&c.inflight) + 1)
	if load == 0 {
		return penalty
	}

	return load
}

func acceptable(err error) bool {
	switch status.Code(err) {
	case codes.DeadlineExceeded, codes.Internal, codes.Unavailable, codes.DataLoss:
		return false
	default:
		return true
	}
}
