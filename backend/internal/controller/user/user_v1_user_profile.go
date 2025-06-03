package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"backend/api/user/v1"
)

func (c *ControllerV1) UserProfile(ctx context.Context, req *v1.UserProfileReq) (res *v1.UserProfileRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
