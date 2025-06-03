package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"backend/api/admin/v1"
)

func (c *ControllerV1) AdminPackageDelete(ctx context.Context, req *v1.AdminPackageDeleteReq) (res *v1.AdminPackageDeleteRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
