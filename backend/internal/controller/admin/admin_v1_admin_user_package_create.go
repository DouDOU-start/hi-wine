package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/internal/service"
)

func (c *ControllerV1) AdminUserPackageCreate(ctx context.Context, req *v1.AdminUserPackageCreateReq) (res *v1.AdminUserPackageCreateRes, err error) {
	// 调用服务创建用户套餐
	detail, err := service.UserPackage().CreateUserPackage(ctx, req)
	if err != nil {
		return nil, err
	}

	// 构建响应
	res = &v1.AdminUserPackageCreateRes{}
	res.Code = 200
	res.Message = "创建用户套餐成功"
	res.Data = *detail

	return res, nil
}
