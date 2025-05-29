package controller

import (
	"backend/api"
	"backend/internal/service"
	"context"
)

type UserController struct{}

func NewUser() *UserController {
	return &UserController{}
}

// Update 更新用户信息
func (c *UserController) Update(ctx context.Context, req *api.UserUpdateReq) (res *api.UserUpdateRes, err error) {
	return service.User().Update(ctx, req)
}

// Info 获取用户信息
func (c *UserController) Info(ctx context.Context, req *api.UserInfoReq) (res *api.UserInfoRes, err error) {
	return service.User().GetInfo(ctx, req)
}
