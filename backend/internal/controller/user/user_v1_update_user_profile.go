package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"backend/api/user/v1"
)

func (c *ControllerV1) UpdateUserProfile(ctx context.Context, req *v1.UpdateUserProfileReq) (res *v1.UpdateUserProfileRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
