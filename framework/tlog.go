package framework

import (
	"github.com/ironzhang/matrix/jsoncfg"
	"github.com/ironzhang/matrix/tlog"
	"go.uber.org/zap"
)

func tlogLoadFromFile(file string) (*zap.Logger, error) {
	if file != "" {
		var cfg tlog.Config
		if err := jsoncfg.LoadFromFile(file, &cfg); err != nil {
			return nil, err
		}
		return tlog.Init(cfg)
	}
	return tlog.Std(), nil
}
