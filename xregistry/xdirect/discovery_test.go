package xdirect

import (
	"context"
	"fmt"
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/g-saber/xstring"
	"github.com/coder2z/g-server/xgrpc"
	"github.com/coder2z/g-server/xgrpc/balancer/least_connection"
	clientinterceptors "github.com/coder2z/g-server/xgrpc/client"
	serverinterceptors "github.com/coder2z/g-server/xgrpc/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"net"
	"testing"
	"time"
)

func TestClientDirect(t *testing.T) {

	err := RegisterBuilder() //发现
	if err != nil {
		t.Failed()
		return
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

	conn, err := grpc.Dial("direct://namespaces///127.0.0.1:5656?w=1", dialOptions...)

	client := helloworld.NewGreeterClient(conn)
	for {
		hello, err := client.SayHello(context.Background(), &helloworld.HelloRequest{
			Name: "TestClientDirect",
		})
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		xlog.Info(xstring.Json(hello))
		time.Sleep(1 * time.Second)
	}
}

type Hello struct {
}

func (h Hello) SayHello(ctx context.Context, request *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{
		Message: request.GetName() + "direct",
	}, nil
}

func TestServerDirect(t *testing.T) {
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

	helloworld.RegisterGreeterServer(serve, new(Hello))

	listener, _ := net.Listen("tcp", "127.0.0.1:5656")
	_ = serve.Serve(listener)
}
