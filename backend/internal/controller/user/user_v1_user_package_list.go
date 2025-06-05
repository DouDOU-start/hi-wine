package user

import (
	"context"

	v1 "backend/api/user/v1"
	"backend/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

// UserPackageList 获取所有可购买的套餐列表 - 公开API
func (c *ControllerV1) UserPackageList(ctx context.Context, req *v1.UserPackageListReq) (res *v1.UserPackageListRes, err error) {
	// 1. 调用服务获取套餐列表
	list, total, err := service.UserPackageForUser().GetPackageList(ctx, req)
	if err != nil {
		g.Log().Error(ctx, "获取套餐列表失败:", err)
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
