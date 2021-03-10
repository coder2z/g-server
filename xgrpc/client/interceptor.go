/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2021/1/2 12:50
 */
package clientinterceptors

import (
	"context"
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/coder2z/g-saber/xcast"
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/g-server/xapp"
	"github.com/coder2z/g-server/xcode"
	"github.com/coder2z/g-server/xmonitor"
	"github.com/coder2z/g-server/xtrace"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"time"
)

func PrometheusUnaryClientInterceptor(name string) func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		t := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		spbStatus := xcode.ExtractCodes(err)
		if spbStatus.Code < xcast.ToInt32(xcode.CodeBreakUp) {
			//系统错误
			xmonitor.ClientHandleCounter.WithLabelValues(xmonitor.TypeGRPCUnary, name, method, cc.Target(), xcast.ToString(spbStatus.GetCode())).Inc()
		} else {
			xmonitor.ClientHandleCounter.WithLabelValues(xmonitor.TypeGRPCUnary, name, method, cc.Target(), "biz error").Inc()
		}
		xmonitor.ClientHandleHistogram.WithLabelValues(xmonitor.TypeGRPCUnary, name, method, cc.Target()).Observe(time.Since(t).Seconds())
		return err
	}
}

func PrometheusStreamClientInterceptor(name string) func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		t := time.Now()
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		spbStatus := xcode.ExtractCodes(err)
		if spbStatus.Code < xcast.ToInt32(xcode.CodeBreakUp) {
			//系统错误
			xmonitor.ClientHandleCounter.WithLabelValues(xmonitor.TypeGRPCUnary, name, method, cc.Target(), xcast.ToString(spbStatus.GetCode())).Inc()
		} else {
			xmonitor.ClientHandleCounter.WithLabelValues(xmonitor.TypeGRPCUnary, name, method, cc.Target(), "biz error").Inc()
		}
		xmonitor.ClientHandleHistogram.WithLabelValues(xmonitor.TypeGRPCUnary, name, method, cc.Target()).Observe(time.Since(t).Seconds())
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
			spbStatus := xcode.ExtractCodes(err)
			span.SetTag("response_code", spbStatus.GetCode())
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
			if spbStatus.Code < xcast.ToInt32(xcode.CodeBreakUp) {
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
				err = spbStatus.SetMsg("server internal error") //吃掉内部错误
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

func XAidUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		clientAidMD := metadata.Pairs(
			"info", fmt.Sprintf("%s/%s/%s",
				xapp.Name(),
				xapp.HostName(),
				xapp.AppId(),
			),
			"ip", xapp.BuildHost(),
			"app_id", xapp.AppId(),
			"app_name", xapp.Name(),
			"host_name", xapp.HostName(),
		)
		if ok {
			md = metadata.Join(md, clientAidMD)
		} else {
			md = clientAidMD
		}
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func HystrixUnaryClientIntercept(timeout, maxConcurrentRequests, requestVolumeThreshold, errorPercentThreshold, sleepWindow int,
	fallback func(context.Context, error) error) grpc.UnaryClientInterceptor {
	//		Timeout：				1000	// 超时时间设置  单位毫秒
	//		MaxConcurrentRequests:  1,		// 最大请求数
	//		RequestVolumeThreshold: 2,		// 默认20，如果错误超过该次数，才开始计算错误百分比
	//		ErrorPercentThreshold:  50, 	// 错误百分比，默认50，即50%
	//		SleepWindow:            5000, 	// 过多长时间，熔断器再次检测是否开启。单位毫秒
	config := hystrix.CommandConfig{
		Timeout:                timeout,
		MaxConcurrentRequests:  maxConcurrentRequests,
		RequestVolumeThreshold: requestVolumeThreshold,
		ErrorPercentThreshold:  errorPercentThreshold,
		SleepWindow:            sleepWindow,
	}
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		name := "GUC" + method
		hystrix.ConfigureCommand(name, config)
		return hystrix.DoC(ctx, name, func(ctx context.Context) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}, fallback)
	}
}

func HystrixStreamClientInterceptor(timeout, maxConcurrentRequests, requestVolumeThreshold, errorPercentThreshold, sleepWindow int,
	fallback func(context.Context, error) error) grpc.StreamClientInterceptor {
	//		Timeout：				1000	// 超时时间设置  单位毫秒
	//		MaxConcurrentRequests:  1,		// 最大请求数
	//		RequestVolumeThreshold: 2,		// 默认20，如果错误超过该次数，才开始计算错误百分比
	//		ErrorPercentThreshold:  50, 	// 错误百分比，默认50，即50%
	//		SleepWindow:            5000, 	// 过多长时间，熔断器再次检测是否开启。单位毫秒
	config := hystrix.CommandConfig{
		Timeout:                timeout,
		MaxConcurrentRequests:  maxConcurrentRequests,
		RequestVolumeThreshold: requestVolumeThreshold,
		ErrorPercentThreshold:  errorPercentThreshold,
		SleepWindow:            sleepWindow,
	}
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (cs grpc.ClientStream, err error) {
		name := "GSC" + method
		hystrix.ConfigureCommand(name, config)
		err = hystrix.DoC(ctx, name, func(ctx context.Context) error {
			cs, err = streamer(ctx, desc, cc, method, opts...)
			return err
		}, fallback)
		return
	}
}
