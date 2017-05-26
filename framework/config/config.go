package config

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ironzhang/matrix/jsoncfg"
)

type ConfigSet struct {
	configs map[string]interface{}
}

func (c *ConfigSet) Register(name string, config interface{}) error {
	if c.configs == nil {
		c.configs = make(map[string]interface{})
	}

	rv := reflect.ValueOf(config)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("config(%s) invalid", name)
	}
	if _, ok := c.configs[name]; ok {
		return fmt.Errorf("config(%s) duplicate", name)
	}
	c.configs[name] = config
	return nil
}

func (c *ConfigSet) LoadFromFile(file string) (err error) {
	if c.configs == nil {
		c.configs = make(map[string]interface{})
	}

	var m map[string]byteSlice
	if err = jsoncfg.LoadFromFile(file, &m); err != nil {
		return err
	}
	for name, config := range c.configs {
		data, ok := m[name]
		if !ok {
			continue
		}
		if err = json.Unmarshal(data, config); err != nil {
			return err
		}
	}
	return nil
}

func (c *ConfigSet) WriteToFile(file string) error {
	if c.configs == nil {
		c.configs = make(map[string]interface{})
	}
	return jsoncfg.WriteToFile(file, c.configs)
}

type byteSlice []byte

func (bs *byteSlice) UnmarshalJSON(b []byte) error {
	*bs = b
	return nil
}
