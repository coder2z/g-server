/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2020/11/4 11:18
 */
package xlog

import (
	"fmt"
	"github.com/myxy99/component/pkg/xcolor"
	cfg "github.com/myxy99/component/xcfg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type options struct {
	// Dir 日志输出目录
	Dir string `mapStructure:"dir"`
	// Name 日志文件名称
	Name string `mapStructure:"name"`
	// Level 日志初始等级
	Level string `mapStructure:"level"`
	// 日志初始化字段
	Fields []zap.Field `mapStructure:"fields"`
	// 是否添加调用者信息
	AddCaller bool `mapStructure:"add_caller"`
	// 日志前缀
	Prefix string `mapStructure:"prefix"`
	// 日志输出文件最大长度，超过改值则截断
	MaxSize   int `mapStructure:"max_size"`
	MaxAge    int `mapStructure:"max_age"`
	MaxBackup int `mapStructure:"max_backup"`
	// 日志磁盘刷盘间隔
	Interval      time.Duration          `mapStructure:"interval"`
	CallerSkip    int                    `mapStructure:"caller_skip"`
	Async         bool                   `mapStructure:"async"`
	Queue         bool                   `mapStructure:"queue"`
	QueueSleep    time.Duration          `mapStructure:"queue_sleep"`
	Core          zapcore.Core           `mapStructure:"core"`
	Debug         bool                   `mapStructure:"debug"`
	EncoderConfig *zapcore.EncoderConfig `mapStructure:"encoder_config"`
	ConfigKey     string                 `mapStructure:"config_key"`
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
	config.ConfigKey = key
	return config
}

// StdConfig Jupiter Standard logger xcfg
func StdConfig(name string) *options {
	return RawConfig("jupiter.logger." + name)
}

// DefaultConfig ...
func defaultConfig() *options {
	return &options{
		Name:          "xlog.xlog",
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
	if o.ConfigKey != "" {
		logger.AutoLevel(o.ConfigKey + ".level")
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
