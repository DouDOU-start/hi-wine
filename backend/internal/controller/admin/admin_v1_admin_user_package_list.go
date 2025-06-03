package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"backend/api/admin/v1"
)

func (c *ControllerV1) AdminUserPackageList(ctx context.Context, req *v1.AdminUserPackageListReq) (res *v1.AdminUserPackageListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
