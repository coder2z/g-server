package gp

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var goPool = NewGoPool(50, 3*time.Second)

type Pool struct {
	head        goroutine
	tail        *goroutine
	count       int
	idleTimeout time.Duration
	maxNum      int
	sync.Mutex
}

type goroutine struct {
	ch     chan handle
	next   *goroutine
	status int32
}

type handle struct {
	f    func()
	errF func(err error)
}

const (
	statusIdle  int32 = 0
	statusInUse int32 = 1
	statusDead  int32 = 2
)

func NewGoPool(c int, idleTimeout time.Duration) *Pool {
	pool := &Pool{
		idleTimeout: idleTimeout,
		maxNum:      c,
	}
	pool.tail = &pool.head
	return pool
}

func (pool *Pool) Go(f func(), ef func(err error)) {
	for {
		g := pool.get()
		if atomic.CompareAndSwapInt32(&g.status, statusIdle, statusInUse) {
			g.ch <- handle{
				f:    f,
				errF: ef,
			}
			return
		}
	}
}

func (pool *Pool) get() *goroutine {
	pool.Lock()
	defer pool.Unlock()
	head := &pool.head
	if head.next == nil {
		if pool.maxNum <= 0 {
			return &pool.head
		} else {
			return pool.alloc()
		}
	}

	ret := head.next
	head.next = ret.next
	if ret == pool.tail {
		pool.tail = head
	}
	pool.count--
	ret.next = nil
	return ret
}

func (pool *Pool) alloc() *goroutine {
	g := &goroutine{
		ch: make(chan handle),
	}
	go g.workLoop(pool)
	pool.Lock()
	pool.maxNum--
	pool.Unlock()
	return g
}

func (g *goroutine) put(pool *Pool) {
	g.status = statusIdle
	pool.Lock()
	pool.tail.next = g
	pool.tail = g
	pool.count++
	pool.Unlock()
}

func (g *goroutine) workLoop(pool *Pool) {
	timer := time.NewTimer(pool.idleTimeout)
	for {
		select {
		case <-timer.C:
			if atomic.CompareAndSwapInt32(&g.status, statusIdle, statusDead) {
				return
			}
		case handle := <-g.ch:
			func() {
				defer panicRecover(handle.errF)
				handle.f()
			}()
			g.put(pool)
		}
		timer.Reset(pool.idleTimeout)
	}
}

func panicRecover(ef func(err error)) {
	if r := recover(); r != nil {
		buf := make([]byte, 64<<10)
		buf = buf[:runtime.Stack(buf, false)]
		err := fmt.Errorf("gp: panic recovered: %s\n%s", r, buf)
		if ef != nil {
			ef(err)
		} else {
			fmt.Println(err)
		}
	}
}

func SafeGo(fn func(), rec func(error)) {
	goPool.Go(fn, rec)
}

func Go(fn func()) {
	goPool.Go(fn, nil)
}
