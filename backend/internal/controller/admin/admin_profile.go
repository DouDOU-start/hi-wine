package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/internal/service"
	"backend/internal/utility/jwt"

	"github.com/gogf/gf/v2/errors/gerror"
)

// AdminProfile 获取当前管理员信息
func (c *ControllerV1) AdminProfile(ctx context.Context, req *v1.AdminProfileReq) (res *v1.AdminProfileRes, err error) {
	// 从上下文中获取管理员ID
	adminId := jwt.GetAdminIdFromCtx(ctx)
	if adminId <= 0 {
		return nil, gerror.New("未登录或登录已过期")
	}

	// 调用服务获取管理员信息
	return service.Admin().GetProfile(ctx, adminId)
}
