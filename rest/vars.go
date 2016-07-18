package rest

import "net/url"

type Vars map[string][]string

func (v Vars) Get(key string) string {
	return url.Values(v).Get(key)
}
