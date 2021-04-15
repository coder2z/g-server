package clientinterceptors

import (
	"fmt"
	"github.com/coder2z/g-server/xgrpc"
	"github.com/coder2z/g-server/xgrpc/balancer/least_connection"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestXGrpcClint(t *testing.T) {
	dialOptions := []grpc.DialOption{
		xgrpc.WithStreamClientInterceptors(
			PrometheusStreamClientInterceptor("servername"),
		),
		xgrpc.WithUnaryClientInterceptors(
			XLoggerUnaryClientInterceptor("servername"),
			PrometheusUnaryClientInterceptor("servername"),
			XTraceUnaryClientInterceptor(),
			XTimeoutUnaryClientInterceptor(time.Minute, time.Second),
			XAidUnaryClientInterceptor(),

		),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, least_connection.LeastConnection)),
		grpc.WithBlock(),
	}
	_, _ = grpc.Dial(":8888", dialOptions...)
}
