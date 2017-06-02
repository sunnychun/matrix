package dashboard_module

import (
	"context"

	"github.com/ironzhang/matrix/restful"
)

type handlers struct {
}

func (h *handlers) Register(m *restful.ServeMux) error {
	apis := []restful.API{
		{"POST", "/config/", h.GetConfig},
		{"PUT", "/config/", h.PutConfig},
	}
	return restful.Register(m, apis)
}

func (h *handlers) GetConfig(ctx context.Context, req interface{}, resp *interface{}) error {
	return nil
}

func (h *handlers) PutConfig(ctx context.Context, req interface{}, resp interface{}) error {
	return nil
}
