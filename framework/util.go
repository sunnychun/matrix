package framework

import (
	"encoding/json"

	"go.uber.org/zap"

	"github.com/ironzhang/matrix/framework/pkg/values"
	"github.com/ironzhang/matrix/jsoncfg"
	"github.com/ironzhang/matrix/tlog"
)

type byteSlice []byte

func (bs *byteSlice) UnmarshalJSON(b []byte) error {
	*bs = b
	return nil
}

func loadAppConfig(configs values.Values, file string) (err error) {
	if file == "" {
		return nil
	}
	var m map[string]byteSlice
	if err = jsoncfg.LoadFromFile(file, &m); err != nil {
		return err
	}
	for k, v := range m {
		if cfg, ok := configs[k]; ok {
			if err = json.Unmarshal(v, cfg); err != nil {
				return err
			}
		}
	}
	return nil
}

func loadLogConfig(file string) (*zap.Logger, error) {
	if file == "" {
		return tlog.Std(), nil
	}
	var cfg tlog.Config
	if err := jsoncfg.LoadFromFile(file, &cfg); err != nil {
		return nil, err
	}
	return tlog.Init(cfg)
}
