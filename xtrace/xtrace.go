/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2021/1/2 13:44
 */
package xtrace

import (
	"context"
	"github.com/coder2m/component/pkg/xconsole"
	"github.com/coder2m/component/xtrace/jaeger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc/metadata"
)

func Init(cfgKey string) {
	cfg := jaeger.RawConfig(cfgKey)
	SetGlobalTracer(
		cfg.Build(),
	)
	xconsole.Greenf("trace init:", cfg.Reporter.LocalAgentHostPort)
}

// SetGlobalTracer ...
func SetGlobalTracer(tracer opentracing.Tracer) {
	opentracing.SetGlobalTracer(tracer)
}

// Start ...
func StartSpanFromContext(ctx context.Context, op string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, op, opts...)
}

// SpanFromContext ...
func SpanFromContext(ctx context.Context) opentracing.Span {
	return opentracing.SpanFromContext(ctx)
}

// CustomTag ...
func CustomTag(key string, val interface{}) opentracing.Tag {
	return opentracing.Tag{
		Key:   key,
		Value: val,
	}
}

// TagComponent ...
func TagComponent(component string) opentracing.Tag {
	return opentracing.Tag{
		Key:   "component",
		Value: component,
	}
}

// TagSpanKind ...
func TagSpanKind(kind string) opentracing.Tag {
	return opentracing.Tag{
		Key:   "span.kind",
		Value: kind,
	}
}

// TagSpanURL ...
func TagSpanURL(url string) opentracing.Tag {
	return opentracing.Tag{
		Key:   "span.url",
		Value: url,
	}
}

// FromIncomingContext ...
func FromIncomingContext(ctx context.Context) opentracing.StartSpanOption {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	sc, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, MetadataReaderWriter{MD: md})
	if err != nil {
		return NullStartSpanOption{}
	}
	return ext.RPCServerOption(sc)
}

// HeaderExtractor ...
func HeaderExtractor(hdr map[string][]string) opentracing.StartSpanOption {
	sc, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, MetadataReaderWriter{MD: hdr})
	if err != nil {
		return NullStartSpanOption{}
	}
	return opentracing.ChildOf(sc)
}

type hdrRequestKey struct{}

// HeaderInjector ...
func HeaderInjector(ctx context.Context, hdr map[string][]string) context.Context {
	span := opentracing.SpanFromContext(ctx)
	err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, MetadataReaderWriter{MD: hdr})
	if err != nil {
		span.LogFields(log.String("event", "inject failed"), log.Error(err))
		return ctx
	}
	return context.WithValue(ctx, hdrRequestKey{}, hdr)
}

// MetadataExtractor ...
func MetadataExtractor(md map[string][]string) opentracing.StartSpanOption {
	sc, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, MetadataReaderWriter{MD: md})
	if err != nil {
		return NullStartSpanOption{}
	}
	return opentracing.ChildOf(sc)
}

// MetadataInjector ...
func MetadataInjector(ctx context.Context, md metadata.MD) context.Context {
	span := opentracing.SpanFromContext(ctx)
	err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, MetadataReaderWriter{MD: md})
	if err != nil {
		span.LogFields(log.String("event", "inject failed"), log.Error(err))
		return ctx
	}
	return metadata.NewOutgoingContext(ctx, md)
}

// NullStartSpanOption ...
type NullStartSpanOption struct{}

// Apply ...
func (sso NullStartSpanOption) Apply(options *opentracing.StartSpanOptions) {}
