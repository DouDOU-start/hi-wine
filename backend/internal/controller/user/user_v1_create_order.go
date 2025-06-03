package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"backend/api/user/v1"
)

func (c *ControllerV1) CreateOrder(ctx context.Context, req *v1.CreateOrderReq) (res *v1.CreateOrderRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
