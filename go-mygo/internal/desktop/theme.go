package desktop

import (
	"context"
	"fmt"

	design "github.com/byteowlz/design-system/impls/mygo-go"
	"github.com/byteowlz/{{project_name}}/internal/config"
	"github.com/byteowlz/{{project_name}}/internal/service"
)

// ResolveTheme uses the canonical native adapter, never copies MyGo's default colors.
func ResolveTheme(svc *service.Service, cfg config.Config) (design.Resolved, error) {
	if cfg.Theme == "omarchy" {
		reply := svc.Inspect(context.Background(), service.Request{Operation: "appearance"})
		if !reply.OK {
			return design.Resolved{}, fmt.Errorf("%s: %s", reply.Error.Code, reply.Error.Message)
		}
		scheme, err := design.OmarchyColors(reply.Result.Appearance.Colors)
		if err != nil {
			return design.Resolved{}, err
		}
		return design.Resolve(scheme, float32(cfg.Radius))
	}
	scheme, err := design.Builtin("oqto-" + cfg.Theme)
	if err != nil {
		return design.Resolved{}, err
	}
	return design.Resolve(scheme, float32(cfg.Radius))
}
