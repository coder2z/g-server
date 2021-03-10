package xgrpc

import (
	"fmt"
	"github.com/coder2z/g-server/xgrpc/balancer/least_connection"
	clientinterceptors "github.com/coder2z/g-server/xgrpc/client"
	serverinterceptors "github.com/coder2z/g-server/xgrpc/server"
	"google.golang.org/grpc"
	"net"
	"testing"
	"time"
)

func TestXGrpcServer(t *testing.T) {
	lis, err := net.Listen("tcp", ":8888")
	if err != nil {
		t.Fatal(err.Error())
		return
	}

	options := []grpc.ServerOption{
		WithUnaryServerInterceptors(
			serverinterceptors.CrashUnaryServerInterceptor(),
			serverinterceptors.PrometheusUnaryServerInterceptor(),
			serverinterceptors.XTimeoutUnaryServerInterceptor(5*time.Second),
			serverinterceptors.TraceUnaryServerInterceptor(),
		),
		WithStreamServerInterceptors(
			serverinterceptors.CrashStreamServerInterceptor(),
			serverinterceptors.PrometheusStreamServerInterceptor(),
		),
	}
	serve := grpc.NewServer(options...)
	_ = serve.Serve(lis)
}

func TestXGrpcClint(t *testing.T) {
	dialOptions := []grpc.DialOption{
		WithStreamClientInterceptors(
			clientinterceptors.PrometheusStreamClientInterceptor("servername"),
		),
		WithUnaryClientInterceptors(
			clientinterceptors.XLoggerUnaryClientInterceptor("servername"),
			clientinterceptors.PrometheusUnaryClientInterceptor("servername"),
			clientinterceptors.XTraceUnaryClientInterceptor(),
			clientinterceptors.XTimeoutUnaryClientInterceptor(time.Minute, time.Second),
			clientinterceptors.XAidUnaryClientInterceptor(),

		),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, least_connection.LeastConnection)),
		grpc.WithBlock(),
	}
	_, _ = grpc.Dial(":8888", dialOptions...)
}
