package tlog_test

import (
	"testing"

	"go.uber.org/zap"

	"github.com/ironzhang/matrix/tlog"
)

func TestTlog(t *testing.T) {
	cfg := tlog.Config{
		Level: zap.DebugLevel,
		//Development: true,
		//DisableCaller:     true,
		DisableStacktrace: true,
		DisableStderr:     true,
		EnableFile:        true,
	}
	if err := tlog.Init(cfg); err != nil {
		t.Fatal(err)
	}

	log := tlog.Std()
	defer log.Sync()

	log.Debug("debug message", zap.String("function", "TestTlog"))
	log.Info("info message", zap.String("function", "TestTlog"))
	log.Warn("warn message", zap.String("function", "TestTlog"))
	log.Error("error message", zap.String("function", "TestTlog"))
}
