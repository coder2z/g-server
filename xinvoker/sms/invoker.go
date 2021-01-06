/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/23 13:52
 **/
package xsms

import (
	"fmt"
	"github.com/myxy99/component/xinvoker"
	"sync"
)

var smsI *smsInvoker

func Register(k string) xinvoker.Invoker {
	smsI = &smsInvoker{key: k}
	return smsI
}

func Invoker(key string) *client {
	if val, ok := smsI.instances.Load(key); ok {
		return val.(*client)
	}
	panic(fmt.Sprintf("no redis(%s) invoker found", key))
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
