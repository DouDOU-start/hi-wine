package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"backend/api/admin/v1"
)

func (c *ControllerV1) AdminOrderUpdateStatus(ctx context.Context, req *v1.AdminOrderUpdateStatusReq) (res *v1.AdminOrderUpdateStatusRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
