package xsms

import (
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/g-server/xinvoker"
	"sync"
)

var smsI *smsInvoker

func Register(k string) xinvoker.Invoker {
	smsI = &smsInvoker{key: k}
	return smsI
}

func Invoker(key string) *Client {
	if val, ok := smsI.instances.Load(key); ok {
		return val.(*Client)
	}
	xlog.Panicf("no sms(%s) invoker found", key)
	return nil
}

type smsInvoker struct {
	xinvoker.Base
	instances sync.Map
	key       string
}

func (i *smsInvoker) Init(opts ...xinvoker.Option) error {
	i.instances = sync.Map{}
	for name, cfg := range i.loadConfig() {
		i.instances.Store(name, i.newSMSClient(cfg))
	}
	return nil
}

func (i *smsInvoker) Reload(opts ...xinvoker.Option) error {
	for name, cfg := range i.loadConfig() {
		i.instances.Store(name, i.newSMSClient(cfg))
	}
	return nil
}

func (i *smsInvoker) Close(opts ...xinvoker.Option) error {
	i.instances.Range(func(key, value interface{}) bool {
		i.instances.Delete(key)
		return true
	})
	return nil
}
