package product

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"backend/api/product/v1"
)

func (c *ControllerV1) UserProductList(ctx context.Context, req *v1.UserProductListReq) (res *v1.UserProductListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
