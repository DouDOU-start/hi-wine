package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/internal/service"
)

func (c *ControllerV1) AdminPackageStats(ctx context.Context, req *v1.AdminPackageStatsReq) (res *v1.AdminPackageStatsRes, err error) {
	// 调用服务获取套餐使用统计
	stats, err := service.Package().GetPackageStats(ctx, req.PackageID)
	if err != nil {
		return nil, err
	}

	// 构建响应
	res = &v1.AdminPackageStatsRes{}
	res.Code = 200
	res.Message = "获取套餐使用统计成功"
	res.Data = *stats

	return res, nil
}
