package user

import (
	"context"

	v1 "backend/api/user/v1"
	"backend/internal/service"
	"backend/internal/utility"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// UserPackageDetail 获取套餐详情
func (c *ControllerV1) UserPackageDetail(ctx context.Context, req *v1.UserPackageDetailReq) (res *v1.UserPackageDetailRes, err error) {
	// 0. 从上下文中获取用户ID并验证
	userId, err := utility.GetUserID(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取套餐详情失败: 未找到用户ID", err)
		return nil, gerror.New("未登录或登录已过期")
	}

	g.Log().Debug(ctx, "获取套餐详情，用户ID:", userId, "套餐ID:", req.PackageID)

	// 1. 调用服务获取套餐详情
	detail, err := service.UserPackageForUser().GetPackageDetail(ctx, req.PackageID)
	if err != nil {
		return nil, err
	}

	// 2. 构建响应
	res = &v1.UserPackageDetailRes{}
	res.Code = 200
	res.Message = "获取套餐详情成功"
	res.Data = *detail

	return res, nil
}
