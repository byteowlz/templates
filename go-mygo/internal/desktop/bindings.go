package desktop

import (
	"context"
	"github.com/byteowlz/{{project_name}}/internal/service"
)

// Desktop exposes only semantic operations, not host lifecycle/configuration.
type Desktop struct{ Service *service.Service }

// Execute uses the exact service called by the headless CLI.
func (d Desktop) Execute(ctx context.Context, request service.Request) service.Envelope {
	return d.Service.Execute(ctx, request)
}
