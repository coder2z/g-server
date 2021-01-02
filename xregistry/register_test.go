/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2021/1/2 11:11
 */
package xregistry

import (
	"github.com/myxy99/component/xregistry/xdirect"
	"github.com/myxy99/component/xregistry/xetcd"
	"github.com/myxy99/component/xregistry/xk8s"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestEtcd(t *testing.T) {

	conf := xetcd.EtcdV3Cfg{
		Endpoints: []string{"127.0.0.1:2379"},
	}

	etcdR, err := xetcd.NewRegistry(conf) //注册
	if err != nil {
		t.Failed()
		return
	}

	etcdR.Register(
		ServiceName("servicename1"),
		ServiceNamespaces("namespaces1"),
		Address("127.0.0.1:8888"),
		RegisterTTL(time.Second*30),      //服务失活一段时间后自动从注册中心删除
		RegisterInterval(time.Second*15), //15s 注册一次
	)

	etcdR.Close() //注销注册

	err = xetcd.RegisterBuilder(conf) //服务发现

	if err != nil {
		t.Failed()
		return
	}

	conn, err := grpc.Dial("etcd://namespaces1/servicename1")
	t.Log(conn)
}

func TestK8s(t *testing.T) {

	err := xk8s.RegisterBuilder() //发现
	if err != nil {
		t.Failed()
		return
	}

	conn, err := grpc.Dial("k8s://namespaces/servicename:portname")
	t.Log(conn)
}

func TestDirect(t *testing.T) {

	err := xdirect.RegisterBuilder() //发现
	if err != nil {
		t.Failed()
		return
	}

	conn, err := grpc.Dial("direct://namespaces/127.0.0.1:8000,127.0.0.1:8001")

	t.Log(conn)
}
