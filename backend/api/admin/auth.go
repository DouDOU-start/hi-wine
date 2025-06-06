package admin

import (
	"context"

	v1 "backend/api/admin/v1"
)

// IAdminAuth 管理员认证接口
type IAdminAuth interface {
	// AdminLogin 管理员登录
	AdminLogin(ctx context.Context, req *v1.AdminLoginReq) (res *v1.AdminLoginRes, err error)
}
