package restful

import "fmt"

type Interface struct {
	Method  string
	Path    string
	Handler interface{}
}

func Register(m *ServeMux, interfaces []Interface) (err error) {
	for _, i := range interfaces {
		if err = m.Add(i.Method, i.Path, i.Handler); err != nil {
			return fmt.Errorf("add interface(%s, %s, %T): %v", i.Method, i.Path, i.Handler, err)
		}
	}
	return nil
}
