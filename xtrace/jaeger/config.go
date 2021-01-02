/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2021/1/2 14:49
 */
package jaeger

import (
	xapp "github.com/myxy99/component"
	"github.com/myxy99/component/pkg/xdefer"
	"github.com/myxy99/component/xcfg"
	"github.com/myxy99/component/xlog"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jCfg "github.com/uber/jaeger-client-go/config"
	"os"
	"time"
)

type Config struct {
	ServiceName      string
	Sampler          *jCfg.SamplerConfig
	Reporter         *jCfg.ReporterConfig
	Headers          *jaeger.HeadersConfig
	EnableRPCMetrics bool
	tags             []opentracing.Tag
	options          []jCfg.Option
	PanicOnError     bool
}

func DefaultConfig() *Config {
	agentAddr := "127.0.0.1:6831"
	headerName := "x-trace-id"
	if addr := os.Getenv("JAEGER_AGENT_ADDR"); addr != "" {
		agentAddr = addr
	}
	return &Config{
		ServiceName: xapp.Name(),
		Sampler: &jCfg.SamplerConfig{
			Type:  "const",
			Param: 0.001,
		},
		Reporter: &jCfg.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  agentAddr,
		},
		EnableRPCMetrics: true,
		Headers: &jaeger.HeadersConfig{
			TraceBaggageHeaderPrefix: "ctx-",
			TraceContextHeaderName:   headerName,
		},
		tags: []opentracing.Tag{
			{Key: "hostname", Value: xapp.HostName()},
		},
		PanicOnError: true,
	}
}

func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := xcfg.UnmarshalKey(key, config); err != nil {
		xlog.Panic("unmarshal key", xlog.Any("err", err))
	}
	return config
}

func (config *Config) WithTag(tags ...opentracing.Tag) *Config {
	if config.tags == nil {
		config.tags = make([]opentracing.Tag, 0)
	}
	config.tags = append(config.tags, tags...)
	return config
}

func (config *Config) WithOption(options ...jCfg.Option) *Config {
	if config.options == nil {
		config.options = make([]jCfg.Option, 0)
	}
	config.options = append(config.options, options...)
	return config
}

func (config *Config) Build() opentracing.Tracer {
	var configuration = jCfg.Configuration{
		ServiceName: config.ServiceName,
		Sampler:     config.Sampler,
		Reporter:    config.Reporter,
		RPCMetrics:  config.EnableRPCMetrics,
		Headers:     config.Headers,
		Tags:        config.tags,
	}
	tracer, closer, err := configuration.NewTracer(config.options...)
	if err != nil {
		if config.PanicOnError {
			xlog.Panic("new jaeger", xlog.String("mod", "jaeger"), xlog.FieldErr(err))
		} else {
			xlog.Error("new jaeger", xlog.String("mod", "jaeger"), xlog.FieldErr(err))
		}
	}
	xdefer.Register(closer.Close)
	return tracer
}
