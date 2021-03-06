package xlog

import (
	"go.uber.org/zap"
)

func logger() *Logger {
	if defaultLogger == nil {
		cfg := StdConfig()
		defaultLogger = cfg.Build()
	}
	return defaultLogger
}

func DefaultLogger() *Logger {
	return logger()
}


func SetDefaultLogger(o *options) {
	defaultLogger = o.Build()
	return
}

var defaultLogger *Logger

// Auto ...
func Auto(err error) Func {
	if err != nil {
		return logger().With(zap.Any("err", err.Error())).Error
	}
	return logger().Info
}

// Info ...
func Info(msg string, fields ...Field) {
	logger().Info(msg, fields...)
}

// Debug ...
func Debug(msg string, fields ...Field) {
	logger().Debug(msg, fields...)
}

// Warn ...
func Warn(msg string, fields ...Field) {
	logger().Warn(msg, fields...)
}

// Error ...
func Error(msg string, fields ...Field) {
	logger().Error(msg, fields...)
}

// Panic ...
func Panic(msg string, fields ...Field) {
	logger().Panic(msg, fields...)
}

// DPanic ...
func DPanic(msg string, fields ...Field) {
	logger().DPanic(msg, fields...)
}

// Fatal ...
func Fatal(msg string, fields ...Field) {
	logger().Fatal(msg, fields...)
}

// Debugw ...
func Debugw(msg string, keysAndValues ...interface{}) {
	logger().Debugw(msg, keysAndValues...)
}

// Infow ...
func Infow(msg string, keysAndValues ...interface{}) {
	logger().Infow(msg, keysAndValues...)
}

// Warnw ...
func Warnw(msg string, keysAndValues ...interface{}) {
	logger().Warnw(msg, keysAndValues...)
}

// Errorw ...
func Errorw(msg string, keysAndValues ...interface{}) {
	logger().Errorw(msg, keysAndValues...)
}

// Panicw ...
func Panicw(msg string, keysAndValues ...interface{}) {
	logger().Panicw(msg, keysAndValues...)
}

// DPanicw ...
func DPanicw(msg string, keysAndValues ...interface{}) {
	logger().DPanicw(msg, keysAndValues...)
}

// Fatalw ...
func Fatalw(msg string, keysAndValues ...interface{}) {
	logger().Fatalw(msg, keysAndValues...)
}

// Debugf ...
func Debugf(msg string, args ...interface{}) {
	logger().Debugf(msg, args...)
}

// Infof ...
func Infof(msg string, args ...interface{}) {
	logger().Infof(msg, args...)
}

// Warnf ...
func Warnf(msg string, args ...interface{}) {
	logger().Warnf(msg, args...)
}

// Errorf ...
func Errorf(msg string, args ...interface{}) {
	logger().Errorf(msg, args...)
}

// Panicf ...
func Panicf(msg string, args ...interface{}) {
	logger().Panicf(msg, args...)
}

// DPanicf ...
func DPanicf(msg string, args ...interface{}) {
	logger().DPanicf(msg, args...)
}

// Fatalf ...
func Fatalf(msg string, args ...interface{}) {
	logger().Fatalf(msg, args...)
}

// Log ...
func (fn Func) Log(msg string, fields ...Field) {
	fn(msg, fields...)
}

// With ...
func With(fields ...Field) *Logger {
	return logger().With(fields...)
}
