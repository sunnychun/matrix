package debug_module

import "github.com/ironzhang/matrix/restful"

type handler struct {
	*restful.ServeMux
}

func (h *handler) Init() error {
	h.ServeMux = restful.NewServeMux(nil)
	apis := []restful.API{}
	if err := restful.Register(h.ServeMux, apis); err != nil {
		return err
	}
	return nil
}
