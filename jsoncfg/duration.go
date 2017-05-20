package jsoncfg

import (
	"fmt"
	"strings"
	"time"
)

type Duration time.Duration

func (d Duration) String() string {
	return time.Duration(d).String()
}

func (d Duration) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%q", d)
	return []byte(s), nil
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	du, err := time.ParseDuration(strings.Trim(string(b), `"`))
	if err != nil {
		return err
	}
	*d = Duration(du)
	return nil
}
