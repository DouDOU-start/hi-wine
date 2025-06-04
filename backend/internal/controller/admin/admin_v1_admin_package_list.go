package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/internal/service"
)

func (c *ControllerV1) AdminPackageList(ctx context.Context, req *v1.AdminPackageListReq) (res *v1.AdminPackageListRes, err error) {
	// 调用服务获取套餐列表
	list, total, err := service.Package().GetPackageList(ctx, req)
	if err != nil {
		return nil, err
	}

	// 构建响应
	res = &v1.AdminPackageListRes{}
	res.Code = 200
	res.Message = "获取套餐列表成功"
	res.Data.List = list
	res.Data.Total = total

	return res, nil
}
