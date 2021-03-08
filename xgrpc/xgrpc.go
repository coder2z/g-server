package xgrpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

func WithStreamServerInterceptors(interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	return grpc.ChainStreamInterceptor(interceptors...)
}

func WithUnaryServerInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(interceptors...)
}

func WithStreamClientInterceptors(interceptors ...grpc.StreamClientInterceptor) grpc.DialOption {
	return grpc.WithChainStreamInterceptor(interceptors...)
}

func WithUnaryClientInterceptors(interceptors ...grpc.UnaryClientInterceptor) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(interceptors...)
}

func ExtractFromCtx(ctx context.Context, key string) string {
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		return strings.Join(md.Get(key), ",")
	}
	return ""
}
