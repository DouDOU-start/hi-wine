package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/internal/service"
)

func (c *ControllerV1) AdminPackageDetail(ctx context.Context, req *v1.AdminPackageDetailReq) (res *v1.AdminPackageDetailRes, err error) {
	// 调用服务获取套餐详情
	detail, err := service.Package().GetPackageDetail(ctx, req.PackageID)
	if err != nil {
		return nil, err
	}

	// 构建响应
	res = &v1.AdminPackageDetailRes{}
	res.Code = 200
	res.Message = "获取套餐详情成功"
	res.Data = *detail

	return res, nil
}
