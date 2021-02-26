/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/23 13:52
 **/
package xemail

import (
	"fmt"
	"github.com/coder2m/component/xinvoker"
	"sync"
)

var email *emailInvoker

func Register(k string) xinvoker.Invoker {
	email = &emailInvoker{key: k}
	return email
}

func Invoker(key string) *Email {
	if val, ok := email.instances.Load(key); ok {
		return val.(*Email)
	}
	panic(fmt.Sprintf("no email(%s) invoker found", key))
}

type emailInvoker struct {
	xinvoker.Base
	instances sync.Map
	key       string
}

func (i *emailInvoker) Init(opts ...xinvoker.Option) error {
	i.instances = sync.Map{}
	for name, cfg := range i.loadConfig() {
		i.instances.Store(name, i.newEmail(cfg))
	}
	return nil
}

func (i *emailInvoker) Reload(opts ...xinvoker.Option) error {
	for name, cfg := range i.loadConfig() {
		i.instances.Store(name, i.newEmail(cfg))
	}
	return nil
}

func (i *emailInvoker) Close(opts ...xinvoker.Option) error {
	return nil
}
