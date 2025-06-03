package product

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"backend/api/product/v1"
)

func (c *ControllerV1) ProductList(ctx context.Context, req *v1.ProductListReq) (res *v1.ProductListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
