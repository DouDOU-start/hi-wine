package order

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"backend/api/order/v1"
)

func (c *ControllerV1) OrderDetail(ctx context.Context, req *v1.OrderDetailReq) (res *v1.OrderDetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
