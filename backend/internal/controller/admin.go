package controller

import (
	"backend/api"
	"backend/internal/service"
	"context"
)

type AdminController struct{}

func NewAdmin() *AdminController {
	return &AdminController{}
}

// 管理员登录
func (c *AdminController) Login(ctx context.Context, req *api.AdminLoginReq) (res *api.AdminLoginRes, err error) {
	return service.Admin().Login(ctx, req)
}

// 获取管理员信息
func (c *AdminController) Info(ctx context.Context, req *api.AdminInfoReq) (res *api.AdminInfoRes, err error) {
	return service.Admin().Info(ctx, req)
}

// 获取用户列表
func (c *AdminController) UserList(ctx context.Context, req *api.UserListReq) (res *api.UserListRes, err error) {
	return service.Admin().UserList(ctx, req)
}

// 获取用户详情
func (c *AdminController) UserDetail(ctx context.Context, req *api.UserDetailReq) (res *api.UserDetailRes, err error) {
	return service.Admin().UserDetail(ctx, req)
}

// 更新用户状态
func (c *AdminController) UserStatusUpdate(ctx context.Context, req *api.UserStatusUpdateReq) (res *api.UserStatusUpdateRes, err error) {
	return service.Admin().UserStatusUpdate(ctx, req)
}
