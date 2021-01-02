/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2021/1/2 12:50
 */
package xgrpc

import (
	"context"
	"fmt"
	"github.com/myxy99/component/pkg/xcode"
	"github.com/myxy99/component/xlog"
	"github.com/myxy99/component/xmonitor"
	"github.com/myxy99/component/xtrace"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"runtime/debug"
	"strings"
	"time"
)

func extractAID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return strings.Join(md.Get("aid"), ",")
	}
	return "unknown"
}

func handleCrash(handler func(interface{})) {
	if r := recover(); r != nil {
		handler(r)
	}
}

func toPanicError(r interface{}) error {
	xlog.Errorf(fmt.Sprintf("%+v", r), xlog.Any("Stack", debug.Stack()))
	return status.Errorf(codes.Internal, "panic: %v", r)
}

func CrashUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer handleCrash(func(r interface{}) {
			err = toPanicError(r)
		})
		return handler(ctx, req)
	}
}

func CrashStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) (err error) {
		defer handleCrash(func(r interface{}) {
			err = toPanicError(r)
		})
		return handler(srv, stream)
	}
}

func PrometheusUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		resp, err := handler(ctx, req)
		code := xcode.ExtractCodes(err)
		xmonitor.ServerHandleHistogram.WithLabelValues(xmonitor.TypeGRPCUnary, info.FullMethod, extractAID(ctx)).Observe(time.Since(startTime).Seconds())
		xmonitor.ServerHandleCounter.WithLabelValues(xmonitor.TypeGRPCUnary, info.FullMethod, extractAID(ctx), code.GetMessage()).Inc()
		return resp, err
	}
}

func PrometheusStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()
		err := handler(srv, ss)
		code := xcode.ExtractCodes(err)
		xmonitor.ServerHandleHistogram.WithLabelValues(xmonitor.TypeGRPCStream, info.FullMethod, extractAID(ss.Context())).Observe(time.Since(startTime).Seconds())
		xmonitor.ServerHandleCounter.WithLabelValues(xmonitor.TypeGRPCStream, info.FullMethod, extractAID(ss.Context()), code.GetMessage()).Inc()
		return err
	}
}

func TraceUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		span, ctx := xtrace.StartSpanFromContext(
			ctx,
			info.FullMethod,
			xtrace.FromIncomingContext(ctx),
			xtrace.TagComponent("gRPC"),
			xtrace.TagSpanKind("server.unary"),
		)
		defer span.Finish()
		resp, err := handler(ctx, req)
		if err != nil {
			code := codes.Unknown
			if s, ok := status.FromError(err); ok {
				code = s.Code()
			}
			span.SetTag("code", code)
			ext.Error.Set(span, true)
			span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
		}
		return resp, err
	}
}

func TimeoutUnaryServerInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		if deadline, ok := ctx.Deadline(); ok {
			leftTime := time.Until(deadline)
			if leftTime < timeout {
				timeout = leftTime
			}
		}
		ctx, cancel := context.WithDeadline(ctx, time.Now().Add(timeout))
		defer cancel()
		return handler(ctx, req)
	}
}
