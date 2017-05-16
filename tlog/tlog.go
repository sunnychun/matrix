package tlog

import (
	"os"

	"github.com/ironzhang/matrix/tlog/writers/file"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var std *zap.Logger

func init() {
	var err error
	std, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
}

type Config struct {
	Level             zapcore.Level
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	DisableStderr     bool
	EnableFile        bool
	FileOptions       file.Options
}

func (cfg Config) openSinks() (zapcore.WriteSyncer, error) {
	var writers []zapcore.WriteSyncer

	if !cfg.DisableStderr {
		writers = append(writers, os.Stderr)
	}

	if cfg.EnableFile {
		f, err := file.Open(cfg.FileOptions)
		if err != nil {
			return nil, err
		}
		writers = append(writers, f)
	}

	return zap.CombineWriteSyncers(writers...), nil
}

func (cfg Config) buildOptions(sink zapcore.WriteSyncer) []zap.Option {
	opts := []zap.Option{zap.ErrorOutput(sink)}

	if cfg.Development {
		opts = append(opts, zap.Development())
	}

	if !cfg.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	stackLevel := zapcore.ErrorLevel
	if cfg.Development {
		stackLevel = zapcore.WarnLevel
	}
	if !cfg.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	return opts
}

func Init(cfg Config) error {
	sink, err := cfg.openSinks()
	if err != nil {
		return err
	}
	opts := cfg.buildOptions(sink)
	enc := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	core := zapcore.NewCore(enc, sink, zap.NewAtomicLevelAt(cfg.Level))
	std = zap.New(core, opts...)
	return nil
}

func Std() *zap.Logger {
	return std
}
