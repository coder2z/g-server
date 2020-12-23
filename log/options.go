/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2020/11/4 11:18
 */
package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Server struct {
	Open bool   `json:"open" yaml:"open"`
	Addr string `json:"addr" yaml:"addr"`
}

type options struct {
	Level string `json:"level" yaml:"level"`

	Development bool `json:"development" yaml:"development"`

	DisableCaller bool `json:"disableCaller" yaml:"disableCaller"`

	DisableStacktrace bool `json:"disableStacktrace" yaml:"disableStacktrace"`

	Sampling *zap.SamplingConfig `json:"sampling" yaml:"sampling"`

	Encoding string `json:"encoding" yaml:"encoding"`

	EncoderConfig zapcore.EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`

	OutputPaths []string `json:"outputPaths" yaml:"outputPaths"`

	ErrorOutputPaths []string               `json:"errorOutputPaths" yaml:"errorOutputPaths"`
	InitialFields    map[string]interface{} `json:"initialFields" yaml:"initialFields"`
}

func newLogOptions() *options {
	return &options{
		Level: "info",
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:05:05"),
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}
}
