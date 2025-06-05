package user

import (
	"context"

	v1 "backend/api/user/v1"
	"backend/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

// UserPackageDetail 获取套餐详情 - 公开API
func (c *ControllerV1) UserPackageDetail(ctx context.Context, req *v1.UserPackageDetailReq) (res *v1.UserPackageDetailRes, err error) {
	// 1. 调用服务获取套餐详情
	detail, err := service.UserPackageForUser().GetPackageDetail(ctx, req.PackageID)
	if err != nil {
		g.Log().Error(ctx, "获取套餐详情失败:", err, "套餐ID:", req.PackageID)
		return nil, err
	}

	// 2. 构建响应
	res = &v1.UserPackageDetailRes{}
	res.Code = 200
	res.Message = "获取套餐详情成功"
	res.Data = *detail

	return res, nil
}
