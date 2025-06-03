package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"backend/api/user/v1"
)

func (c *ControllerV1) WechatLogin(ctx context.Context, req *v1.WechatLoginReq) (res *v1.WechatLoginRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
