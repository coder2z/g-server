/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2020/11/4 11:18
 */
package xlog

import (
	"fmt"
	cfg "github.com/myxy99/component/config"
	"github.com/myxy99/component/pkg/xcolor"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type options struct {
	// Dir 日志输出目录
	Dir string
	// Name 日志文件名称
	Name string
	// Level 日志初始等级
	Level string
	// 日志初始化字段
	Fields []zap.Field
	// 是否添加调用者信息
	AddCaller bool
	// 日志前缀
	Prefix string
	// 日志输出文件最大长度，超过改值则截断
	MaxSize   int
	MaxAge    int
	MaxBackup int
	// 日志磁盘刷盘间隔
	Interval      time.Duration
	CallerSkip    int
	Async         bool
	Queue         bool
	QueueSleep    time.Duration
	Core          zapcore.Core
	Debug         bool
	EncoderConfig *zapcore.EncoderConfig
	configKey     string
}

// Filename ...
func (o *options) Filename() string {
	return fmt.Sprintf("%s/%s", o.Dir, o.Name)
}

// RawConfig ...
func RawConfig(key string) *options {
	var config = defaultConfig()
	if err := cfg.UnmarshalKey(key, &config); err != nil {
		panic(err)
	}
	config.configKey = key
	return config
}

// StdConfig Jupiter Standard logger config
func StdConfig(name string) *options {
	return RawConfig("jupiter.logger." + name)
}

// DefaultConfig ...
func defaultConfig() *options {
	return &options{
		Name:          "log.log",
		Dir:           ".",
		Level:         "info",
		MaxSize:       500, // 500M
		MaxAge:        1,   // 1 day
		MaxBackup:     10,  // 10 backup
		Interval:      24 * time.Hour,
		CallerSkip:    1,
		AddCaller:     false,
		Async:         true,
		Queue:         false,
		QueueSleep:    100 * time.Millisecond,
		EncoderConfig: DefaultZapConfig(),
	}
}

// Build ...
func (o options) Build() *Logger {
	if o.EncoderConfig == nil {
		o.EncoderConfig = DefaultZapConfig()
	}
	if o.Debug {
		o.EncoderConfig.EncodeLevel = DebugEncodeLevel
	}
	logger := newLogger(&o)
	if o.configKey != "" {
		logger.AutoLevel(o.configKey + ".level")
	}
	return logger
}

func DefaultZapConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "lv",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func DebugEncodeLevel(lv zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var colorize = xcolor.Red
	switch lv {
	case zapcore.DebugLevel:
		colorize = xcolor.Blue
	case zapcore.InfoLevel:
		colorize = xcolor.Green
	case zapcore.WarnLevel:
		colorize = xcolor.Yellow
	case zapcore.ErrorLevel, zap.PanicLevel, zap.DPanicLevel, zap.FatalLevel:
		colorize = xcolor.Red
	default:
	}
	enc.AppendString(colorize(lv.CapitalString()))
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.Unix())
}
