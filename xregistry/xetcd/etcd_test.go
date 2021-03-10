package xetcd

import (
	"context"
	"fmt"
	"github.com/coder2z/g-server/xgrpc"
	xbalancer "github.com/coder2z/g-server/xgrpc/balancer"
	"github.com/coder2z/g-server/xgrpc/balancer/least_connection"
	clientinterceptors "github.com/coder2z/g-server/xgrpc/client"
	serverinterceptors "github.com/coder2z/g-server/xgrpc/server"
	"github.com/coder2z/g-server/xregistry"
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/g-saber/xtime"
	proto "github.com/grpc-ecosystem/go-grpc-prometheus/examples/grpc-server-with-prometheus/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	"testing"
	"time"
)

type server struct {
	add string
}

func (s server) SayHello(c context.Context, request *proto.HelloRequest) (*proto.HelloResponse, error) {
	text := "Hello " + request.Name + ", I am " + s.add
	log.Println(text)

	return &proto.HelloResponse{Message: text}, nil
}

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
		xregistry.ServiceName("servicename"),
		xregistry.ServiceNamespaces("namespaces"),
		xregistry.Address("127.0.0.1:8888"),
		xregistry.RegisterTTL(xtime.Duration("30s")),      //服务失活一段时间后自动从注册中心删除
		xregistry.RegisterInterval(xtime.Duration("15s")), //15s 注册一次
		xregistry.Metadata(metadata.Pairs(xbalancer.WeightKey, "2")),
	)
	options := []grpc.ServerOption{
		xgrpc.WithUnaryServerInterceptors(
			serverinterceptors.CrashUnaryServerInterceptor(),
			serverinterceptors.PrometheusUnaryServerInterceptor(),
			serverinterceptors.XTimeoutUnaryServerInterceptor(5*time.Second),
			serverinterceptors.TraceUnaryServerInterceptor(),
		),
		xgrpc.WithStreamServerInterceptors(
			serverinterceptors.CrashStreamServerInterceptor(),
			serverinterceptors.PrometheusStreamServerInterceptor(),
		),
	}
	serve := grpc.NewServer(options...)
	listener, err := net.Listen("tcp", ":8888")
	proto.RegisterDemoServiceServer(serve, server{add: ":8888"})
	_ = serve.Serve(listener)
}

func TestEtcd2(t *testing.T) {

	conf := EtcdV3Cfg{
		Endpoints: []string{"127.0.0.1:2379"},
	}

	etcdR, err := NewRegistry(conf) //注册
	if err != nil {
		t.Failed()
		return
	}

	options := []grpc.ServerOption{
		xgrpc.WithUnaryServerInterceptors(
			serverinterceptors.CrashUnaryServerInterceptor(),
			serverinterceptors.PrometheusUnaryServerInterceptor(),
			serverinterceptors.XTimeoutUnaryServerInterceptor(5*time.Second),
			serverinterceptors.TraceUnaryServerInterceptor(),
		),
		xgrpc.WithStreamServerInterceptors(
			serverinterceptors.CrashStreamServerInterceptor(),
			serverinterceptors.PrometheusStreamServerInterceptor(),
		),
	}
	serve := grpc.NewServer(options...)

	etcdR.Register(
		xregistry.ServiceName("servicename"),
		xregistry.ServiceNamespaces("namespaces"),
		xregistry.Address("127.0.0.1:7777"),
		xregistry.RegisterTTL(xtime.Duration("30s")),      //服务失活一段时间后自动从注册中心删除
		xregistry.RegisterInterval(xtime.Duration("15s")), //15s 注册一次
		xregistry.Metadata(metadata.Pairs(xbalancer.WeightKey, "2")),
	)

	listener, err := net.Listen("tcp", ":7777")
	proto.RegisterDemoServiceServer(serve, server{add: ":7777"})
	_ = serve.Serve(listener)
}

func TestEtcdDiscovery(t *testing.T) {
	conf := EtcdV3Cfg{
		Endpoints: []string{"127.0.0.1:2379"},
	}
	dialOptions := []grpc.DialOption{
		xgrpc.WithStreamClientInterceptors(
			clientinterceptors.PrometheusStreamClientInterceptor("servername"),
		),
		xgrpc.WithUnaryClientInterceptors(
			clientinterceptors.XAidUnaryClientInterceptor(),
			clientinterceptors.HystrixUnaryClientIntercept(1000,
				30,
				20,
				30,
				20,
				func(ctx context.Context, err error) error {
					return err
				}),
			clientinterceptors.XTimeoutUnaryClientInterceptor(time.Minute, time.Second),
			clientinterceptors.XLoggerUnaryClientInterceptor("servername"),
			clientinterceptors.PrometheusUnaryClientInterceptor("servername"),
			clientinterceptors.XTraceUnaryClientInterceptor(),
		),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, least_connection.LeastConnection)),
		grpc.WithInsecure(),
	}
	_ = RegisterBuilder(conf) //服务发现
	conn, err := grpc.Dial("etcd://namespaces/servicename", dialOptions...)
	if err != nil {
		t.Logf("grpc dial: %s", err)
		return
	}
	defer conn.Close()
	client := proto.NewDemoServiceClient(conn)
	for i := 0; i < 100; i++ {
		ctx := context.Background()
		resp, err := client.SayHello(ctx,
			&proto.HelloRequest{
				Name: "123",
			})
		if err != nil {
			xlog.Error("error", xlog.FieldErr(err))
			time.Sleep(time.Second)
			continue
		}
		xlog.Info("success", xlog.String("msg", resp.GetMessage()))
		time.Sleep(time.Second)
	}
}
