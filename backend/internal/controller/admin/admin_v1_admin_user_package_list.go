package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/internal/service"
)

func (c *ControllerV1) AdminUserPackageList(ctx context.Context, req *v1.AdminUserPackageListReq) (res *v1.AdminUserPackageListRes, err error) {
	// 调用服务获取用户套餐列表
	list, total, err := service.UserPackage().GetUserPackageList(ctx, req)
	if err != nil {
		return nil, err
	}

	// 构建响应
	res = &v1.AdminUserPackageListRes{}
	res.Code = 200
	res.Message = "获取用户套餐列表成功"
	res.Data.List = list
	res.Data.Total = total

	return res, nil
}
