/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/23 13:52
 **/
package email

import (
	"fmt"
	invoker "github.com/myxy99/component"
	"sync"
)

var email *emailInvoker

func Register(k string) invoker.Invoker {
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
	invoker.Base
	instances sync.Map
	key       string
}

func (i *emailInvoker) Init(opts ...invoker.Option) error {
	i.instances = sync.Map{}
	for name, cfg := range i.loadConfig() {
		i.instances.Store(name, i.newEmail(cfg))
	}
	return nil
}

func (i *emailInvoker) Reload(opts ...invoker.Option) error {
	for name, cfg := range i.loadConfig() {
		i.instances.Store(name, i.newEmail(cfg))
	}
	return nil
}

func (i *emailInvoker) Close(opts ...invoker.Option) error {
	return nil
}
