package server

import (
	"fmt"
	"strings"
)

func pathKey(path string) string {
	i := len(path) - 1
	if path[i] != '/' {
		return path + "/"
	}
	return path
}

func parsePath(rawpath string) (op string, key string, err error) {
	if len(rawpath) <= 0 || rawpath[0] != '/' {
		return "", "", fmt.Errorf("invalid path: %q", rawpath)
	}

	path := rawpath[1:]
	i := strings.IndexRune(path, '/')
	if i == -1 {
		op = path
		key = "/"
	} else {
		op = path[:i]
		key = pathKey(path[i:])
	}

	return
}
