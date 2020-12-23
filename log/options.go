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
	Level string `json:"level" mapStructure:"level"`

	Development bool `json:"development" mapStructure:"development"`

	DisableCaller bool `json:"disableCaller" mapStructure:"disableCaller"`

	DisableStacktrace bool `json:"disableStacktrace" mapStructure:"disableStacktrace"`

	Sampling *zap.SamplingConfig `json:"sampling" mapStructure:"sampling"`

	Encoding string `json:"encoding" mapStructure:"encoding"`

	EncoderConfig zapcore.EncoderConfig `json:"encoderConfig" mapStructure:"encoderConfig"`

	OutputPaths []string `json:"outputPaths" mapStructure:"outputPaths"`

	ErrorOutputPaths []string               `json:"errorOutputPaths" mapStructure:"errorOutputPaths"`
	InitialFields    map[string]interface{} `json:"initialFields" mapStructure:"initialFields"`
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
