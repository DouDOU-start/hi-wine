package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"backend/api/admin/v1"
)

func (c *ControllerV1) AdminProductList(ctx context.Context, req *v1.AdminProductListReq) (res *v1.AdminProductListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
