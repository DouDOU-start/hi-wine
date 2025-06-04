package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/internal/service"
)

func (c *ControllerV1) AdminUserPackageUpdateStatus(ctx context.Context, req *v1.AdminUserPackageUpdateStatusReq) (res *v1.AdminUserPackageUpdateStatusRes, err error) {
	// 调用服务更新用户套餐状态
	detail, err := service.UserPackage().UpdateUserPackageStatus(ctx, req)
	if err != nil {
		return nil, err
	}

	// 构建响应
	res = &v1.AdminUserPackageUpdateStatusRes{}
	res.Code = 200
	res.Message = "更新用户套餐状态成功"
	res.Data = *detail

	return res, nil
}
