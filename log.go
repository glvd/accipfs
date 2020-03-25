package accipfs

import (
	"github.com/goextension/log"
	"go.uber.org/zap"
)

// InitLog ...
func InitLog() {
	cfg := zap.NewProductionConfig()
	cfg.Level = logLvToAtomicLv(LogLevel)
	cfg.OutputPaths = []string{LogOutput}
	cfg.ErrorOutputPaths = []string{LogOutput}
	logger, e := cfg.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	if e != nil {
		panic(e)
	}
	log.Register(logger.Sugar())

	log.Debugw("log init", "level", LogLevel, "output", LogOutput)
}
