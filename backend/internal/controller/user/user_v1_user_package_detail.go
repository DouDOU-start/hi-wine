package user

import (
	"context"

	v1 "backend/api/user/v1"
	"backend/internal/service"
)

// UserPackageDetail 获取套餐详情
func (c *ControllerV1) UserPackageDetail(ctx context.Context, req *v1.UserPackageDetailReq) (res *v1.UserPackageDetailRes, err error) {
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
