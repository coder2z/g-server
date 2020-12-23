package log

import (
	"github.com/myxy99/component/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	level = map[string]zapcore.Level{
		"info":   zapcore.InfoLevel,
		"debug":  zapcore.DebugLevel,
		"warn":   zapcore.WarnLevel,
		"error":  zapcore.ErrorLevel,
		"dPanic": zapcore.DPanicLevel,
		"panic":  zapcore.PanicLevel,
		"fatal":  zapcore.FatalLevel,
	}
)

func (i *logInvoker) loadConfig() map[string]*options {
	conf := make(map[string]*options)

	prefix := i.key
	for name := range config.GetStringMap(prefix) {
		cfg := config.UnmarshalWithExpect(prefix+"."+name, newLogOptions()).(*options)
		conf[name] = cfg
	}

	return conf
}

func (i *logInvoker) new(o *options) *zap.SugaredLogger {
	z := zap.NewProductionConfig()
	if v, ok := level[o.Level]; ok {
		z.Level = zap.NewAtomicLevelAt(v)
	} else {
		z.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
	z.Development = o.Development
	z.DisableCaller = o.DisableCaller
	z.DisableStacktrace = o.DisableStacktrace
	z.Sampling = o.Sampling
	z.Encoding = o.Encoding
	z.EncoderConfig = o.EncoderConfig
	z.OutputPaths = o.OutputPaths
	z.ErrorOutputPaths = o.ErrorOutputPaths
	z.InitialFields = o.InitialFields
	logger, err := z.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil
	}
	return logger.Sugar()
}
