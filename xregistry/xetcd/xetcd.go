package xetcd

import (
	"context"
	"fmt"
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/g-saber/xstring"
	"github.com/coder2z/g-server/xapp"
	"github.com/coder2z/g-server/xregistry"
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
		uid:       xapp.AppId(),
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
		xlog.Panic("Application Starting",
			xlog.FieldComponentName("XRegistry"),
			xlog.FieldMethod("XRegistry.XEtcd.Register"),
			xlog.FieldDescription("Service name required"),
		)
	}
	if r.options.Namespaces == "" {
		xlog.Panic("Application Starting",
			xlog.FieldComponentName("XRegistry"),
			xlog.FieldMethod("XRegistry.XEtcd.Register"),
			xlog.FieldDescription("Service namespaces required"),
		)
	}
	if r.options.Address == "" {
		xlog.Panic("Application Starting",
			xlog.FieldComponentName("XRegistry"),
			xlog.FieldMethod("XRegistry.XEtcd.Register"),
			xlog.FieldDescription("Service address required"),
		)
	}
	if r.options.RegisterTTL == 0 {
		r.options.RegisterTTL = time.Second * 30
	}
	if r.options.RegisterInterval > r.options.RegisterTTL || r.options.RegisterInterval < r.options.RegisterTTL/3 {
		r.options.RegisterInterval = r.options.RegisterTTL / 3
	}

	xlog.Info("Application Starting",
		xlog.FieldComponentName("XRegistry"),
		xlog.FieldMethod("XRegistry.XEtcd.Register"),
		xlog.FieldDescription(fmt.Sprintf("Service registration to Etcd Endpoints:%v", r.conf.Endpoints)),
		xlog.Any("service name", r.options.ServiceName),
		xlog.FieldValueAny(r.options),
	)

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
		err error
		one = sync.Once{}
	)
	defer func() {
		if err != nil {
			xlog.Warn("Application Running",
				xlog.FieldComponentName("XRegistry"),
				xlog.FieldMethod("XRegistry.XEtcd.Register"),
				xlog.FieldDescription("Etcd register error"),
				xlog.Any("service name", r.options.ServiceName),
				xlog.FieldErr(err),
				xlog.FieldValueAny(r.options),
			)
		} else {
			one.Do(func() {
				xlog.Info("Application Running",
					xlog.FieldComponentName("XRegistry"),
					xlog.FieldMethod("XRegistry.XEtcd.Register"),
					xlog.FieldDescription("Etcd register success"),
					xlog.Any("service name", r.options.ServiceName),
					xlog.FieldValueAny(r.options),
				)
			})
		}
	}()
	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	ttl, err := r.client.Grant(timeout, int64(r.options.RegisterTTL/time.Second))
	if err != nil {
		return err
	}
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
		close(r.closeCh)
	})
	r.Wait()
}

func (r *etcdReg) getKey() string {
	key := fmt.Sprintf("%s/%s/%s", etcdPrefix, fmt.Sprintf("%v/%v", r.options.Namespaces, r.options.ServiceName), r.uid)
	return key
}

func (r *etcdReg) unregister() {
	key := r.getKey()
	if _, err := r.client.Delete(context.Background(), key); err != nil {
		xlog.Warn("Application Stopping",
			xlog.FieldComponentName("XRegistry"),
			xlog.FieldMethod("XRegistry.XEtcd.unregister"),
			xlog.FieldDescription("Etcd register unregister error"),
			xlog.Any("service name", r.options.ServiceName),
			xlog.FieldErr(err),
			xlog.FieldValueAny(r.options),
			xlog.String("app id", r.uid),
		)
	}
	_, _ = r.client.Revoke(context.Background(), r.leaseId) // 回收租约
	xlog.Info("Application Stopping",
		xlog.FieldComponentName("XRegistry"),
		xlog.FieldMethod("XRegistry.XEtcd.Unregister"),
		xlog.FieldDescription("Service stopping,Registration cancellation"),
	)
	//_ = r.client.Close()
}
