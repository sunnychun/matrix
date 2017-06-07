package dashboard_module

import (
	"context"
	"net/url"

	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/restful"
)

type handlers struct {
}

func (h *handlers) Register(m *restful.ServeMux) error {
	apis := []restful.API{
		{"GET", "/dashboard/configs", h.GetConfigs},
		{"GET", "/dashboard/configs/:module", h.GetModuleConfig},
		{"PUT", "/dashboard/configs/:module", h.PutModuleConfig},
	}
	return restful.Register(m, apis)
}

func (h *handlers) GetConfigs(ctx context.Context, values url.Values, req interface{}, resp *interface{}) error {
	m := make(map[string]interface{})
	for k, c := range framework.Configs() {
		m[k] = c.Get()
	}
	*resp = m
	return nil
}

func (h *handlers) GetModuleConfig(ctx context.Context, values url.Values, req interface{}, resp *interface{}) error {
	if c, ok := framework.Configs()[values.Get(":module")]; ok {
		*resp = c.Get()
	}
	return nil
}

func (h *handlers) PutModuleConfig(ctx context.Context, values url.Values, req map[string]interface{}, resp *interface{}) error {
	if c, ok := framework.Configs()[values.Get(":module")]; ok {
		if err := c.Set(req); err != nil {
			return err
		}
		*resp = c.Get()
	}
	return nil
}
