package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/internal/service"
)

func (c *ControllerV1) AdminUserActivePackages(ctx context.Context, req *v1.AdminUserActivePackagesReq) (res *v1.AdminUserActivePackagesRes, err error) {
	// 调用服务获取用户有效套餐
	list, err := service.UserPackage().GetUserActivePackages(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// 构建响应
	res = &v1.AdminUserActivePackagesRes{}
	res.Code = 200
	res.Message = "获取用户有效套餐成功"
	res.Data.List = list

	return res, nil
}
