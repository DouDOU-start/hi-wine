package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/internal/service"
)

// AdminLogin 管理员登录
func (c *ControllerV1) AdminLogin(ctx context.Context, req *v1.AdminLoginReq) (res *v1.AdminLoginRes, err error) {
	// 调用管理员服务进行登录验证
	result, err := service.Admin().Login(ctx, req.Username, req.Password)
	if err != nil {
		return &v1.AdminLoginRes{
			Token: "",
			AdminUser: v1.AdminUser{
				ID:       0,
				Username: "",
				Role:     "",
			},
		}, err
	}

	return result, nil
}
