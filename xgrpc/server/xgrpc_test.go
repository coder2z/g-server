package serverinterceptors

import (
	"github.com/coder2z/g-server/xgrpc"
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
		xgrpc.WithUnaryServerInterceptors(
			CrashUnaryServerInterceptor(),
			PrometheusUnaryServerInterceptor(),
			XTimeoutUnaryServerInterceptor(5*time.Second),
			TraceUnaryServerInterceptor(),
		),
		xgrpc.WithStreamServerInterceptors(
			CrashStreamServerInterceptor(),
			PrometheusStreamServerInterceptor(),
		),
	}
	serve := grpc.NewServer(options...)
	_ = serve.Serve(lis)
}
