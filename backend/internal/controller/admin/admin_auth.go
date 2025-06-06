package admin

import (
	"context"

	"backend/api/admin"
	v1 "backend/api/admin/v1"
	"backend/api/common"
	"backend/internal/service"
)

// AuthController 管理员认证控制器 - 处理管理员登录
type AuthController struct{}

// NewAuth 创建管理员认证控制器
func NewAuth() admin.IAdminAuth {
	return &AuthController{}
}

// AdminLogin 管理员登录
func (c *AuthController) AdminLogin(ctx context.Context, req *v1.AdminLoginReq) (res *v1.AdminLoginRes, err error) {
	// 调用管理员服务进行登录验证
	result, err := service.Admin().Login(ctx, req.Username, req.Password)
	if err != nil {
		return &v1.AdminLoginRes{
			Response: common.Response[struct {
				Token     string       `json:"token"`
				AdminUser v1.AdminUser `json:"admin_user"`
			}]{
				Code:    common.CodeServerError,
				Message: err.Error(),
				Data: struct {
					Token     string       `json:"token"`
					AdminUser v1.AdminUser `json:"admin_user"`
				}{
					Token: "",
					AdminUser: v1.AdminUser{
						ID:       0,
						Username: "",
						Role:     "",
					},
				},
			},
		}, err
	}

	return result, nil
}
