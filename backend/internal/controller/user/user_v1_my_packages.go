package user

import (
	"backend/api/common"
	v1 "backend/api/user/v1"
	"backend/internal/service"
	"context"

	"github.com/gogf/gf/v2/util/gconv"
)

// UserMyPackages 获取当前登录用户的套餐列表
func (c *ControllerV1) UserMyPackages(ctx context.Context, req *v1.UserMyPackagesReq) (res *v1.UserMyPackagesRes, err error) {
	// 从上下文中获取当前登录用户ID
	userId := gconv.Int64(ctx.Value("userId"))
	if userId <= 0 {
		return nil, common.ErrNotAuthenticated
	}

	// 调用服务层获取用户套餐列表
	packages, err := service.UserPackage().GetUserMyPackages(ctx, userId, req.Status)
	if err != nil {
		return nil, err
	}

	// 构建响应
	res = &v1.UserMyPackagesRes{}
	res.Code = 200
	res.Message = "获取成功"
	res.Data.List = packages

	return res, nil
}
