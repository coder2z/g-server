/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2021/1/1 22:20
 */
package xetcd

import (
	"context"
	"fmt"
	"github.com/myxy99/component/pkg/xjson"
	"github.com/myxy99/component/xlog"
	"github.com/myxy99/component/xregistry"
	"go.etcd.io/etcd/clientv3"
)

type etcdDiscovery struct {
	client *clientv3.Client
}

func newDiscovery(conf EtcdV3Cfg) (xregistry.Discovery, error) {
	d := &etcdDiscovery{}
	c, err := clientv3.New(conf)
	if err != nil {
		return nil, err
	}
	d.client = c
	return d, nil
}

// target: serviceName
func (d *etcdDiscovery) Discover(target string) (<-chan []xregistry.Instance, error) {
	ch := make(chan []xregistry.Instance)
	go d.watch(ch, target)
	return ch, nil
}

func (d *etcdDiscovery) watch(ch chan<- []xregistry.Instance, serviceName string) {
	prefix := fmt.Sprintf("/%s/%s/", etcdPrefix, serviceName)

	update := func() []xregistry.Instance {
		resp, err := d.client.Get(context.Background(), prefix, clientv3.WithPrefix())
		if err != nil {
			xlog.Warn("etcd discovery watch", xlog.FieldErr(err), xlog.Any("servicename", serviceName))
			return nil
		}
		var i []xregistry.Instance
		for _, kv := range resp.Kvs {
			ins := xregistry.Instance{}
			if err = xjson.Unmarshal(kv.Value, &ins); err == nil {
				i = append(i, ins)
			} else {
				xlog.Warn("etcd discovery watch unmarshal servicename", xlog.FieldErr(err), xlog.Any("servicename", serviceName))
			}
		}
		return i
	}
	if i := update(); len(i) > 0 {
		ch <- i
	}

	eventCh := d.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for range eventCh {
		ch <- update()
	}
	return
}

func (d *etcdDiscovery) Close() {
	_ = d.client.Close()
}
