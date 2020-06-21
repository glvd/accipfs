package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const module = "log"

// Level ...
var Level = "info"

// Output ...
var Output = "stderr"

var _log *zap.SugaredLogger

func logLvToAtomicLv(lv string) zap.AtomicLevel {
	a := zap.NewAtomicLevel()

	level := zapcore.InfoLevel
	switch lv {
	case "debug":
		level = zapcore.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "dpanic":
		level = zapcore.DPanicLevel
	case "panic":
		level = zapcore.PanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	}

	a.SetLevel(level)

	return a
}

// InitLog ...
func InitLog() {
	cfg := zap.NewProductionConfig()
	cfg.Level = logLvToAtomicLv(Level)
	cfg.OutputPaths = []string{Output}
	cfg.ErrorOutputPaths = []string{Output}
	logger, e := cfg.Build(
		zap.AddCaller(),
		//zap.AddCallerSkip(1),
	)
	if e != nil {
		panic(e)
	}
	_log = logger.Sugar()
	_log.Debugw("log init", "module", module, "level", Level, "output", Output)
}

// Module ...
func Module(m string) *zap.SugaredLogger {
	if _log == nil {
		InitLog()
	}
	return _log.With("module", m)
}
