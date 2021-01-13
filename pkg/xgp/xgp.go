package xgp

import (
	"bytes"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

var goPool = NewWaitPool(50)

const (
	start int32 = iota
	running
	shutdown
	stop
)

type simpleFSM struct {
	status int32
}

func newSimpleFSM() *simpleFSM {
	return &simpleFSM{
		status: start,
	}
}

func (s *simpleFSM) actEvent(stat int32) bool {
	if s.Current() > stat {
		return false
	} else {
		return atomic.CompareAndSwapInt32(&s.status, s.status, stat)
	}
}

func (s *simpleFSM) isRunning() bool {
	return s.Current() == running
}

func (s *simpleFSM) Current() int32 {
	return atomic.LoadInt32(&s.status)
}

func (s *simpleFSM) Action(stat int32) bool {
	switch stat {
	case start:
	case running:
	case shutdown:
	case stop:
		return s.actEvent(stat)
	default:
		return false
	}
	return false
}

type basePoolExecutor struct {
	corePoolSize chan int
	fsm          *simpleFSM
	oc           sync.Once
}

type poolExecutor struct {
	basePoolExecutor
}

type waitPoolExecutor struct {
	basePoolExecutor
	gp *sync.WaitGroup
}

func newBasePoolExecutor(cap int) basePoolExecutor {
	return basePoolExecutor{
		corePoolSize: make(chan int, cap),
		fsm:          newSimpleFSM(),
		oc:           sync.Once{},
	}
}

func NewPool(cap int) *poolExecutor {
	checkCap(cap)
	return &poolExecutor{
		basePoolExecutor: newBasePoolExecutor(cap),
	}
}

func NewWaitPool(cap int) *waitPoolExecutor {
	checkCap(cap)
	return &waitPoolExecutor{
		basePoolExecutor: newBasePoolExecutor(cap),
		gp:               new(sync.WaitGroup),
	}
}

func checkCap(cap int) {
	if cap < 0 {
		panic("The pool cap cannot lower zero")
	}
}

func (b *basePoolExecutor) checkSubmit(f func()) {
	if f == nil {
		panic("The submit func is nil")
	}
	b.oc.Do(func() {
		b.fsm.actEvent(running)
	})
	if !b.fsm.isRunning() {
		panic("The pool is not running")
	}
}

func (b *basePoolExecutor) ShutDown() {
	b.fsm.actEvent(shutdown)
}

func (b *basePoolExecutor) IsShutDown() bool {
	return b.fsm.Current() >= shutdown
}

func (b *basePoolExecutor) IsTerminated() bool {
	return b.fsm.Current() >= stop
}

func (t *poolExecutor) Submit(f func(), ef func(error)) {
	t.checkSubmit(f)
	t.corePoolSize <- 1
	go func() {
		defer func() {
			if err := recover(); err != nil {
				var buf bytes.Buffer
				stack := debug.Stack()
				buf.Write(stack)
				err := fmt.Errorf("gp: panic recovered: %s \n %s", err, buf.String())
				if ef == nil {
					ef(err)
				} else {
					panic(err)
				}
			}
			<-t.corePoolSize
		}()
		if t.IsTerminated() {
			return
		}
		f()
	}()
}

func (t *waitPoolExecutor) Submit(f func(), ef func(error)) {
	t.checkSubmit(f)
	t.gp.Add(1)
	t.corePoolSize <- 1
	go func() {
		defer func() {
			if err := recover(); err != nil {
				var buf bytes.Buffer
				stack := debug.Stack()
				buf.Write(stack)
				err := fmt.Errorf("gp: panic recovered: %s \n %s", err, buf.String())
				if ef == nil {
					ef(err)
				} else {
					panic(err)
				}
			}
			<-t.corePoolSize
			t.gp.Done()
		}()
		if t.IsTerminated() {
			return
		}
		f()
	}()
}

func (t *waitPoolExecutor) Wait() {
	t.gp.Wait()
	t.fsm.actEvent(stop)
	close(t.corePoolSize)
}

func SafeGo(fn func(), rec func(error)) {
	goPool.Submit(fn, rec)
}

func Go(fn func()) {
	goPool.Submit(fn, nil)
}

func Wait() {
	goPool.Wait()
}
