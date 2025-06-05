package user

import (
	"context"

	v1 "backend/api/user/v1"
	"backend/internal/service"
	"backend/internal/utility"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// UserPackageList 获取所有可购买的套餐列表
func (c *ControllerV1) UserPackageList(ctx context.Context, req *v1.UserPackageListReq) (res *v1.UserPackageListRes, err error) {
	// 0. 从上下文中获取用户ID并验证
	userId, err := utility.GetUserID(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取套餐列表失败: 未找到用户ID", err)
		return nil, gerror.New("未登录或登录已过期")
	}

	g.Log().Debug(ctx, "获取套餐列表，用户ID:", userId)

	// 1. 调用服务获取套餐列表
	list, total, err := service.UserPackageForUser().GetPackageList(ctx, req)
	if err != nil {
		return nil, err
	}

	// 2. 构建响应
	res = &v1.UserPackageListRes{}
	res.Code = 200
	res.Message = "获取套餐列表成功"
	res.Data.List = list
	res.Data.Total = total

	return res, nil
}
