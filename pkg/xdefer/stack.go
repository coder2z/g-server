package xdefer

import (
	"sync"
)

func newStack() *DeferStack {
	return &DeferStack{
		fns: make([]func() error, 0),
		mu:  sync.RWMutex{},
	}
}

type DeferStack struct {
	fns []func() error
	mu  sync.RWMutex
}

func (ds *DeferStack) Push(fns ...func() error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.fns = append(ds.fns, fns...)
}

func (ds *DeferStack) Clean() {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	for i := len(ds.fns) - 1; i >= 0; i-- {
		_ = ds.fns[i]()
	}
}

var (
	globalDefers = newStack()
)

// Register 注册一个defer函数
func Register(fns ...func() error) {
	globalDefers.Push(fns...)
}

// Clean 清除
func Clean() {
	globalDefers.Clean()
}
