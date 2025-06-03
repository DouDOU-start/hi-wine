package product

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"backend/api/product/v1"
)

func (c *ControllerV1) UserCategoryList(ctx context.Context, req *v1.UserCategoryListReq) (res *v1.UserCategoryListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
