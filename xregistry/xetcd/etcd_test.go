package xetcd

import (
	"github.com/coder2z/g-saber/xtime"
	"github.com/coder2z/component/xregistry"
	"google.golang.org/grpc"
	"testing"
)

func TestEtcd(t *testing.T) {

	conf := EtcdV3Cfg{
		Endpoints: []string{"127.0.0.1:2379"},
	}

	etcdR, err := NewRegistry(conf) //注册
	if err != nil {
		t.Failed()
		return
	}

	etcdR.Register(
		xregistry.ServiceName("servicename1"),
		xregistry.ServiceNamespaces("namespaces1"),
		xregistry.Address("127.0.0.1:8888"),
		xregistry.RegisterTTL(xtime.Duration("30s")),      //服务失活一段时间后自动从注册中心删除
		xregistry.RegisterInterval(xtime.Duration("15s")), //15s 注册一次
	)

	etcdR.Close() //注销注册

	err = RegisterBuilder(conf) //服务发现

	if err != nil {
		t.Failed()
		return
	}

	conn, err := grpc.Dial("etcd://namespaces1/servicename1")
	t.Log(conn)
}
