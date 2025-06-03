package print

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"backend/api/print/v1"
)

func (c *ControllerV1) PrintOrder(ctx context.Context, req *v1.PrintOrderReq) (res *v1.PrintOrderRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
