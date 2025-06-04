package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/internal/service"
)

func (c *ControllerV1) AdminPackageCreate(ctx context.Context, req *v1.AdminPackageCreateReq) (res *v1.AdminPackageCreateRes, err error) {
	// 调用服务创建套餐
	detail, err := service.Package().CreatePackage(ctx, req)
	if err != nil {
		return nil, err
	}

	// 构建响应
	res = &v1.AdminPackageCreateRes{}
	res.Code = 200
	res.Message = "创建套餐成功"
	res.Data = *detail

	return res, nil
}
