package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/internal/service"
)

func (c *ControllerV1) AdminPackageUpdate(ctx context.Context, req *v1.AdminPackageUpdateReq) (res *v1.AdminPackageUpdateRes, err error) {
	// 调用服务更新套餐
	detail, err := service.Package().UpdatePackage(ctx, req)
	if err != nil {
		return nil, err
	}

	// 构建响应
	res = &v1.AdminPackageUpdateRes{}
	res.Code = 200
	res.Message = "更新套餐成功"
	res.Data = *detail

	return res, nil
}
