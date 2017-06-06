package dashboard_module

import (
	"context"
	"net/url"

	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/framework/pkg/assign"
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
	*resp = framework.Configs()
	return nil
}

func (h *handlers) GetModuleConfig(ctx context.Context, values url.Values, req interface{}, resp *interface{}) error {
	configs := framework.Configs()
	if cfg, ok := configs[values.Get(":module")]; ok {
		*resp = cfg
	}
	return nil
}

func (h *handlers) PutModuleConfig(ctx context.Context, values url.Values, req map[string]interface{}, resp *interface{}) error {
	configs := framework.Configs()
	if cfg, ok := configs[values.Get(":module")]; ok {
		if err := assign.Assign(cfg, req); err != nil {
			return err
		}
		*resp = cfg
	}
	return nil
}
