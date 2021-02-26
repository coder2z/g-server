package xlog

import (
	"github.com/coder2m/component/xcfg"
	"go.uber.org/zap"
)

func GetDefaultLogger() *Logger {
	if DefaultLogger == nil {
		prefix := `xlog`
		cfg := xcfg.UnmarshalWithExpect(prefix, defaultConfig()).(*options)
		DefaultLogger = newLogger(cfg)

	}
	return DefaultLogger
}

var DefaultLogger *Logger

// Auto ...
func Auto(err error) Func {
	if err != nil {
		return GetDefaultLogger().With(zap.Any("err", err.Error())).Error
	}

	return GetDefaultLogger().Info
}

// Info ...
func Info(msg string, fields ...Field) {
	GetDefaultLogger().Info(msg, fields...)
}

// Debug ...
func Debug(msg string, fields ...Field) {
	GetDefaultLogger().Debug(msg, fields...)
}

// Warn ...
func Warn(msg string, fields ...Field) {
	GetDefaultLogger().Warn(msg, fields...)
}

// Error ...
func Error(msg string, fields ...Field) {
	GetDefaultLogger().Error(msg, fields...)
}

// Panic ...
func Panic(msg string, fields ...Field) {
	GetDefaultLogger().Panic(msg, fields...)
}

// DPanic ...
func DPanic(msg string, fields ...Field) {
	GetDefaultLogger().DPanic(msg, fields...)
}

// Fatal ...
func Fatal(msg string, fields ...Field) {
	GetDefaultLogger().Fatal(msg, fields...)
}

// Debugw ...
func Debugw(msg string, keysAndValues ...interface{}) {
	GetDefaultLogger().Debugw(msg, keysAndValues...)
}

// Infow ...
func Infow(msg string, keysAndValues ...interface{}) {
	GetDefaultLogger().Infow(msg, keysAndValues...)
}

// Warnw ...
func Warnw(msg string, keysAndValues ...interface{}) {
	GetDefaultLogger().Warnw(msg, keysAndValues...)
}

// Errorw ...
func Errorw(msg string, keysAndValues ...interface{}) {
	GetDefaultLogger().Errorw(msg, keysAndValues...)
}

// Panicw ...
func Panicw(msg string, keysAndValues ...interface{}) {
	GetDefaultLogger().Panicw(msg, keysAndValues...)
}

// DPanicw ...
func DPanicw(msg string, keysAndValues ...interface{}) {
	GetDefaultLogger().DPanicw(msg, keysAndValues...)
}

// Fatalw ...
func Fatalw(msg string, keysAndValues ...interface{}) {
	GetDefaultLogger().Fatalw(msg, keysAndValues...)
}

// Debugf ...
func Debugf(msg string, args ...interface{}) {
	GetDefaultLogger().Debugf(msg, args...)
}

// Infof ...
func Infof(msg string, args ...interface{}) {
	GetDefaultLogger().Infof(msg, args...)
}

// Warnf ...
func Warnf(msg string, args ...interface{}) {
	GetDefaultLogger().Warnf(msg, args...)
}

// Errorf ...
func Errorf(msg string, args ...interface{}) {
	GetDefaultLogger().Errorf(msg, args...)
}

// Panicf ...
func Panicf(msg string, args ...interface{}) {
	GetDefaultLogger().Panicf(msg, args...)
}

// DPanicf ...
func DPanicf(msg string, args ...interface{}) {
	GetDefaultLogger().DPanicf(msg, args...)
}

// Fatalf ...
func Fatalf(msg string, args ...interface{}) {
	GetDefaultLogger().Fatalf(msg, args...)
}

// Log ...
func (fn Func) Log(msg string, fields ...Field) {
	fn(msg, fields...)
}

// With ...
func With(fields ...Field) *Logger {
	return GetDefaultLogger().With(fields...)
}
