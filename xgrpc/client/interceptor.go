/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2021/1/2 12:50
 */
package xgrpc

import (
	"context"
	"errors"
	"github.com/myxy99/component/pkg/xcode"
	"github.com/myxy99/component/xlog"
	"github.com/myxy99/component/xmonitor"
	"github.com/myxy99/component/xtrace"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"time"
)

func XMonitorUnaryClientInterceptor(name string) func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		t := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		spbStatus := xcode.ExtractCodes(err)
		if spbStatus.Code < xcode.EcodeNum {
			//系统错误
			xmonitor.ClientHandleCounter.WithLabelValues(xmonitor.TypeGRPCUnary, name, method, cc.Target(), spbStatus.GetMessage()).Inc()
			xmonitor.ClientHandleHistogram.WithLabelValues(xmonitor.TypeGRPCUnary, name, method, cc.Target()).Observe(time.Since(t).Seconds())
		} else {
			xmonitor.ClientHandleCounter.WithLabelValues(xmonitor.TypeGRPCUnary, name, method, cc.Target(), "biz error").Inc()
			xmonitor.ClientHandleHistogram.WithLabelValues(xmonitor.TypeGRPCUnary, name, method, cc.Target()).Observe(time.Since(t).Seconds())
		}
		return err
	}
}

func XMonitorStreamClientInterceptor(name string) func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		t := time.Now()
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		spbStatus := xcode.ExtractCodes(err)
		if spbStatus.Code < xcode.EcodeNum {
			//系统错误
			xmonitor.ClientHandleCounter.WithLabelValues(xmonitor.TypeGRPCUnary, name, method, cc.Target(), spbStatus.GetMessage()).Inc()
			xmonitor.ClientHandleHistogram.WithLabelValues(xmonitor.TypeGRPCUnary, name, method, cc.Target()).Observe(time.Since(t).Seconds())
		} else {
			xmonitor.ClientHandleCounter.WithLabelValues(xmonitor.TypeGRPCUnary, name, method, cc.Target(), "biz error").Inc()
			xmonitor.ClientHandleHistogram.WithLabelValues(xmonitor.TypeGRPCUnary, name, method, cc.Target()).Observe(time.Since(t).Seconds())
		}
		return clientStream, err
	}
}

func XTraceUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}

		span, ctx := xtrace.StartSpanFromContext(
			ctx,
			method,
			xtrace.TagSpanKind("client"),
			xtrace.TagComponent("grpc"),
		)
		defer span.Finish()

		err := invoker(xtrace.MetadataInjector(ctx, md), method, req, reply, cc, opts...)
		if err != nil {
			code := codes.Unknown
			if s, ok := status.FromError(err); ok {
				code = s.Code()
			}
			span.SetTag("response_code", code)
			ext.Error.Set(span, true)

			span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
		}
		return err
	}
}

func XTimeoutUnaryClientInterceptor(timeout time.Duration, slowThreshold time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		now := time.Now()
		_, ok := ctx.Deadline()
		if !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		err := invoker(ctx, method, req, reply, cc, opts...)
		du := time.Since(now)
		remoteIP := "unknown"
		if remote, ok := peer.FromContext(ctx); ok && remote.Addr != nil {
			remoteIP = remote.Addr.String()
		}

		if slowThreshold > time.Duration(0) && du > slowThreshold {
			xlog.Error("slow",
				xlog.FieldErr(errors.New("grpc unary slow command")),
				xlog.FieldMethod(method),
				xlog.FieldName(cc.Target()),
				xlog.FieldCost(du),
				xlog.FieldAddr(remoteIP),
			)
		}
		return err
	}
}

func XLoggerUnaryClientInterceptor(name string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		beg := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)

		spbStatus := xcode.ExtractCodes(err)
		if err != nil {
			// 只记录系统级别错误
			if spbStatus.Code < xcode.EcodeNum {
				xlog.Error(
					"access",
					xlog.FieldType("unary"),
					xlog.FieldCode(spbStatus.Code),
					xlog.FieldErrKind(spbStatus.Message),
					xlog.FieldName(name),
					xlog.FieldMethod(method),
					xlog.FieldCost(time.Since(beg)),
					xlog.Any("req", req),
					xlog.Any("reply", reply),
				)
			} else {
				xlog.Warn(
					"access",
					xlog.FieldType("unary"),
					xlog.FieldCode(spbStatus.Code),
					xlog.FieldErrKind(spbStatus.Message),
					xlog.FieldName(name),
					xlog.FieldMethod(method),
					xlog.FieldCost(time.Since(beg)),
					xlog.Any("req", req),
					xlog.Any("reply", reply),
				)
			}
			return err
		}
		return nil
	}
}
