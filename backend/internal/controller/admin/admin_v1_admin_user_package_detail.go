package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/internal/service"
)

func (c *ControllerV1) AdminUserPackageDetail(ctx context.Context, req *v1.AdminUserPackageDetailReq) (res *v1.AdminUserPackageDetailRes, err error) {
	// 调用服务获取用户套餐详情
	detail, err := service.UserPackage().GetUserPackageDetail(ctx, req.UserPackageID)
	if err != nil {
		return nil, err
	}

	// 构建响应
	res = &v1.AdminUserPackageDetailRes{}
	res.Code = 200
	res.Message = "获取用户套餐详情成功"
	res.Data = *detail

	return res, nil
}
