package framework

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ironzhang/matrix/jsoncfg"
)

var Config = config{}

func loadConfig(file string) error {
	if file != "" {
		return Config.load(file)
	}
	return nil
}

type config map[string]interface{}

func (c config) Register(module string, conf interface{}) {
	rv := reflect.ValueOf(conf)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic(fmt.Sprintf("module(%s) config invalid", module))
	}

	if _, ok := c[module]; ok {
		panic(fmt.Sprintf("module(%s) duplicate", module))
	}

	c[module] = conf
}

func (c config) load(file string) (err error) {
	var m map[string]byteSlice
	if err = jsoncfg.LoadFromFile(file, &m); err != nil {
		return err
	}
	for k, v := range c {
		b, ok := m[k]
		if !ok {
			continue
		}
		if err = json.Unmarshal(b, v); err != nil {
			return err
		}
	}
	return nil
}

func (c config) write(file string) error {
	return jsoncfg.WriteToFile(file, c)
}

type byteSlice []byte

func (bs *byteSlice) UnmarshalJSON(b []byte) error {
	*bs = b
	return nil
}
