package accipfs

import (
	"fmt"
	"github.com/goextension/log"
	zap "go.uber.org/zap"
)

// InitLog ...
func InitLog() {
	fmt.Println("log init info:", LogLevel, LogOutput)
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
}
