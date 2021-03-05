/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2021/1/1 22:21
 */
package xetcd

import (
	"context"
	"fmt"
	"github.com/coder2m/component/pkg/xconsole"
	"github.com/coder2m/component/pkg/xstring"
	"github.com/coder2m/component/xlog"
	"github.com/coder2m/component/xregistry"
	"go.etcd.io/etcd/clientv3"
	"sync"
	"time"
)

var (
	etcdPrefix = "registry.etcd" // 不需包含 '/'
)

type etcdReg struct {
	conf    clientv3.Config
	client  *clientv3.Client
	options *xregistry.Options

	*sync.WaitGroup
	isOk      bool
	closeCh   chan struct{}
	closeOnce sync.Once
	uid       string
	leaseId   clientv3.LeaseID
}

func NewRegistry(conf EtcdV3Cfg) (xregistry.Registry, error) {
	r := &etcdReg{
		conf:      conf,
		options:   &xregistry.Options{},
		closeCh:   make(chan struct{}),
		uid:       xstring.GenerateID(),
		WaitGroup: new(sync.WaitGroup),
	}
	c, err := clientv3.New(conf)
	if err != nil {
		return nil, err
	}
	r.client = c
	return r, nil
}

func (r *etcdReg) Register(ops ...xregistry.Option) {
	for _, o := range ops {
		o(r.options)
	}
	if r.options.ServiceName == "" {
		panic("service name required")
	}
	if r.options.Namespaces == "" {
		panic("service namespaces required")
	}
	if r.options.Address == "" {
		panic("service address required")
	}
	if r.options.RegisterTTL == 0 {
		r.options.RegisterTTL = time.Second * 30
	}
	if r.options.RegisterInterval > r.options.RegisterTTL || r.options.RegisterInterval < r.options.RegisterTTL/3 {
		r.options.RegisterInterval = r.options.RegisterTTL / 3
	}

	xconsole.Greenf("Service registration to:", fmt.Sprintf("etcd:%v", r.conf.Endpoints))

	go func() {
		r.Add(1)
		defer r.Done()
		var err error
		err = r.register() // 先注册一次
		ticker := time.NewTicker(r.options.RegisterInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err == nil { // 注册成功则续租
					err = r.keepAliveOnce()
				}
				if err != nil { // 注册/续租失败则重新注册
					err = r.register()
				}
			case <-r.closeCh:
				if r.isOk {
					r.unregister()
				}
				return
			}
		}
	}()
}

func (r *etcdReg) register() error {
	var (
		err  error
		step int
		one  = sync.Once{}
	)
	defer func() {
		if err != nil {
			xlog.Warn("etcd register error", xlog.FieldErr(err), xlog.Any("step", step), xlog.Any("options", r.options))
		} else {
			one.Do(func() {
				xconsole.Greenf("etcd register success to:", xstring.Json(r.options))
			})
		}
	}()
	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	ttl, err := r.client.Grant(timeout, int64(r.options.RegisterTTL/time.Second))
	if err != nil {
		return err
	}

	step += 1
	_, err = r.client.Put(timeout, r.getKey(), xstring.Json(r.options), clientv3.WithLease(ttl.ID))
	if err == nil {
		r.leaseId = ttl.ID
		r.isOk = true
	}
	return err
}

// 续租一次
func (r *etcdReg) keepAliveOnce() error {
	_, err := r.client.KeepAliveOnce(context.Background(), r.leaseId)
	return err
}

func (r *etcdReg) Close() {
	r.closeOnce.Do(func() {
		xconsole.Red("Service registration shutdown")
		close(r.closeCh)
	})
	r.Wait()
}

func (r *etcdReg) getKey() string {
	key := fmt.Sprintf("/%s/%s/%s", etcdPrefix, fmt.Sprintf("%v.%v", r.options.Namespaces, r.options.ServiceName), r.uid)
	return key
}

func (r *etcdReg) unregister() {
	key := r.getKey()
	if _, err := r.client.Delete(context.Background(), key); err != nil {
		xlog.Warn("unregister error", xlog.FieldErr(err), xlog.Any("uid", r.uid), xlog.Any("options", r.options))
	}
	_, _ = r.client.Revoke(context.Background(), r.leaseId) // 回收租约
	xconsole.Redf("unregister success", xstring.Json(r.options))
	//_ = r.client.Close()
}
